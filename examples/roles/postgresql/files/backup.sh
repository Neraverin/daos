#!/bin/bash

set -Eeuo pipefail
trap cleanup SIGINT SIGTERM ERR EXIT

script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd -P)
tmp_dir=$script_dir/tmp

usage() {
  cat <<EOF
Usage: $(basename "${BASH_SOURCE[0]}") [-h] [-v] [--tmpdir|-t <temp directory>] output_dir

Script description here.

Available options:

-t, --tmpdir <temp directory>   Set temporary directory path to <temp directory> instead of default path
-h, --help                      Print this help and exit
-v, --verbose                   Print script debug info
EOF
  exit
}

cleanup() {
  trap - SIGINT SIGTERM ERR EXIT
  rv=$?
  rm -f $script_dir/version
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
    -t | --tmpdir)
      shift
      if [[ ${1-} != "" ]] ;
      then
        if ! tmp_dir=$(mktemp -d -p ${1} -t pgbkp.XXXX) ;
        then
          die "Error creating directory for temp files in ${1}"
        fi
      else
        die "Error: --tmpdir parameter is used but no temp directory path set!"
      fi
    ;;
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
  [[ ${#args[@]} -eq 0 ]] && die "Target directory for saving backup files is not set"

  return 0
}

parse_params "$@"
setup_colors

# script logic here

required_role_type="SqlStorage"

msg "Backing up role ${GREEN}${required_role_type}${NOFORMAT}"

output_dir=${args[0]-}
role_id=${script_dir##*/}
app_id=$(basename $(dirname $script_dir))
params=$(salt-call pillar.item local pillarenv=$app_id --out=json)

role_type=$(echo $params | $script_dir/jq -r --arg role_id "$role_id" '.local.local[] | select(.Id==$role_id) | .RoleTypeId')
version=$(echo $params | $script_dir/jq -r --arg role_id "$role_id" '.local.local[].Instance | select(.Id==$role_id) | .Version')
PG_USER=$(echo $params | $script_dir/jq -r --arg role_id "$role_id" '.local.local[] | select(.Id==$role_id) | .Params.PgUser')

[[ $role_type != $required_role_type ]] && die "${RED}Script file is outside of ${NOFORMAT}${required_role_type}${RED} role directory${NOFORMAT}"

[[ -z $version || $version = null ]] && die "${RED}Can't detect role version${NOFORMAT}"

echo "Role version: $version"

TMP_DIR_OWNER=$(ls -la ${script_dir} | awk '$9 == "data" {print $3":"$4}')
mkdir -p $tmp_dir
chown ${TMP_DIR_OWNER} $tmp_dir

output_dir=${output_dir%/}/$app_id/$role_id

[[ -d "$output_dir" ]] && die "${RED}Directory ${NOFORMAT}${output_dir}${RED} already exists${NOFORMAT}"

mkdir -p $output_dir

backup_archive=$output_dir/backup.tar

postgres_container=$(docker ps | grep $role_id | grep postgres | awk 'NR==1{print $1}')
postgres_image=$(docker inspect $postgres_container --format "{{ .Image }}" | cut -d ":" -f 2)

[[ -z $postgres_container ]] && die "${RED}Storage docker container not found${NOFORMAT}"
[[ -z $postgres_image ]] && die "${RED}Storage docker image not found${NOFORMAT}"

msg "Exec pg_basebackup..."

docker run --rm --name postgres-backup \
           --network=container:$postgres_container \
           -v $tmp_dir:/postgres_backup/ \
           $postgres_image \
           pg_basebackup -U $PG_USER -h localhost -D /postgres_backup -Ft -P -c fast -v -z -X s

echo "role_type=${role_type}" > $tmp_dir/version
echo "version=${version}" >> $tmp_dir/version

msg "Archiving to $backup_archive..."

tar -cvf $backup_archive --directory=$tmp_dir .

if [ $? -eq 0 ]
then
  rm -rf $tmp_dir
  msg "${GREEN}Backup is done${NOFORMAT}"
else
  echo "${RED}Something went wrong with creating role backup archive. Please check temp directory in${NOFORMAT} ${tmp_dir}"
fi
