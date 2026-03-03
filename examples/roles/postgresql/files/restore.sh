#!/bin/bash

set -Eeuo pipefail
trap cleanup SIGINT SIGTERM ERR EXIT

script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd -P)
tmp_dir=''

usage() {
  cat <<EOF
Usage: $(basename "${BASH_SOURCE[0]}") [-h] [-v] backup_dir

Script description here.

Available options:

-h, --help          Print this help and exit
-v, --verbose       Print script debug info
EOF
  exit
}

cleanup() {
  trap - SIGINT SIGTERM ERR EXIT
  rv=$?
  rm -rf $tmp_dir
  exit $rv
}

setup_colors() {
  if [[ -t 2 ]] && [[ -z "${NO_COLOR-}" ]] && [[ "${TERM-}" != "dumb" ]]; then
    NOFORMAT='\033[0m' RED='\033[0;31m' GREEN='\033[0;32m' ORANGE='\033[0;33m' BLUE='\033[0;34m' PURPLE='\033[0;35m' CYAN='\033[0;36m' YELLOW='\033[1;33m'
  else
    NOFORMAT='' RED='' GREEN='' ORANGE='' BLUE='' PURPLE='' CYAN='' YELLOW=''
  fi
}

msg() {
  echo >&2 -e "${1-}"
}

die() {
  local msg=$1
  local code=${2-1} # default exit status 1
  msg "$msg"
  exit "$code"
}

parse_params() {
  while :; do
    case "${1-}" in
    -h | --help) usage ;;
    -v | --verbose) set -x ;;
    --no-color) NO_COLOR=1 ;;
    -?*) die "Unknown option: $1" ;;
    *) break ;;
    esac
    shift
  done

  args=("$@")

  # check arguments
  [[ ${#args[@]} -eq 0 ]] && die "Missing script arguments"

  return 0
}

parse_params "$@"
setup_colors

# script logic here

required_role_type="SqlStorage"

msg "Restoring role ${GREEN}${required_role_type}${NOFORMAT}"

backup_dir=${args[0]-}
role_id=${script_dir##*/}
app_id=$(basename $(dirname $script_dir))
params=$(salt-call pillar.item local pillarenv=$app_id --out=json)
role_type=$(echo $params | $script_dir/jq -r --arg role_id "$role_id" '.local.local[] | select(.Id==$role_id) | .RoleTypeId')
required_version=$(echo $params | $script_dir/jq -r --arg role_id "$role_id" '.local.local[].Instance | select(.Id==$role_id) | .Version')

[[ $role_type != $required_role_type ]] && die "${RED}Script file is outside of ${NOFORMAT}${required_role_type}${RED} role directory${NOFORMAT}"

[[ -z $required_version || $required_version = null ]] && die "${RED}No version${NOFORMAT}"

backup_dir=$backup_dir/$app_id/$role_id
backup_archive=$backup_dir/backup.tar

[[ ! -r "$backup_archive" ]] && die "${RED}Backup file ${NOFORMAT}${backup_archive}${RED} not found or it is not readable${NOFORMAT}"

tmp_dir=$backup_dir/tmp

rm -rf $tmp_dir
mkdir -p $tmp_dir

msg "Extracting backup data from ${backup_archive} to ${tmp_dir}"

tar --directory=$tmp_dir -xvf $backup_archive

[[ ! -r "${tmp_dir}/version" ]] && die "${RED}Version file not found or it is not readable${NOFORMAT}"

source $tmp_dir/version

[[ $role_type != $required_role_type ]] && die "${RED}Invalid role type ${NOFORMAT}${role_type}${RED} in backup. ${NOFORMAT}${required_role_type}${RED} required${NOFORMAT}"

[[ $required_version != $version ]] && die "${RED}Current role version is ${NOFORMAT}${required_version}${RED}, but archived is ${NOFORMAT}${version}"

postgres_container=$(docker ps | grep $role_id | grep postgres | awk 'NR==1{print $1}')

[[ -z $postgres_container ]] && die "${RED}Storage container not found${NOFORMAT}"

msg "Stopping postgresql container..."

docker stop $postgres_container

msg "Removing current data..."

rm -rf $script_dir/data/*

msg "Restoring archived..."

tar --directory=$script_dir/data -xvf $tmp_dir/base.tar.gz

mkdir $script_dir/data/wal_archive
tar -C $script_dir/data/wal_archive -xvf $tmp_dir/pg_wal.tar.gz
chown -R dockerns:dockerns $script_dir/data/wal_archive

touch $script_dir/data/recovery.signal
chown -R dockerns:dockerns $script_dir/data/recovery.signal

msg "Starting postgresql container..."

docker start $postgres_container

msg "Waiting for database restore"

while [[ -f $script_dir/data/recovery.signal ]]
do
  sleep 1
done

msg "${GREEN}Done${NOFORMAT}"
