#!/bin/bash

set -Eeuo pipefail

# Set base variables with their default values
script_dir=$(cd "$(dirname ${BASH_SOURCE[0]})" &>/dev/null && pwd -P)
ARGS=""
DATA_DIR_OWNER=$(stat -c '%U:%G' "${script_dir}/data" )
BACKUP_REPO="none"
BACKUP_MODE="full"
BACKUP_ID="none"
required_role_type="SqlStorage"

if [[ $(nproc) -ge 4 ]] ; then
  BACKUP_PROC=4
else
  BACKUP_PROC=$(nproc)
fi

# -----------------------------------------------------------------
#  Helper functions
# -----------------------------------------------------------------
usage() {
  cat <<EOF
Usage: $(basename "${BASH_SOURCE[0]}") [-h] [-v]
       $(basename "${BASH_SOURCE[0]}") checkdb -B <backup repo dir> [-j num] [-v] [-j <N>]
       $(basename "${BASH_SOURCE[0]}") backup -B <backup repo dir> -b <full|delta> [-j num] [-n] [-v] [-j <N>]
       $(basename "${BASH_SOURCE[0]}") restore -B <backup repo dir> [-i <backup ID>] [-j num] [-n] [-v] [-j <N>]
       $(basename "${BASH_SOURCE[0]}") show -B <backup repo dir> [-i <backup ID>][-v]
       $(basename "${BASH_SOURCE[0]}") merge -B <backup repo dir> -i <backup ID> [-j num] [-v] [-j <N>]
       $(basename "${BASH_SOURCE[0]}") delete -B <backup repo dir> -i <backup ID> [-v]
       $(basename "${BASH_SOURCE[0]}") cleanup -B <backup repo dir> [-v]

This script is designed for preforming backup and restore of SqlStorage data.

Available commands:
  backup                          Perform backup
  checkdb                         Perform check for corrupted pages in a database (not in backups)
  restore                         Restore a backup
  delete                          Delete a backup from backup repository
  show                            List backups stored in backup repository
  merge                           Merge incremental (delta) backups into the closest full backup
  cleanup                         Delete backups of old role versions data

Available options:
  -B <dir>       Path to the backup repository.
  -b <mode>      Backup mode. Only "full" and "delta" modes are currentry supported. Applicable for "backup" command.
  -i backup_id   Specifies the unique identifier of the backup. Available for "merge", "show" and "delete" commands
  -n             Skip backup validation. Available for "backup", "merge" and "restore" commands. 
  -j <N>         Run operation in N parallel threads
  --no-color     Don't use colorized output
  -h, --help     Print this help and exit
  -v, --verbose  Print script debug info
EOF
  exit
}

setup_colors() {
  if [[ -t 2 && -z "${NO_COLOR-}" && "${TERM-}" != "dumb" ]]; then
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
  local code="${2-1}" # default exit status 1
  msg "$msg"
  exit "$code"
}

# -----------------------------------------------------------------
#  Argument parsing
# -----------------------------------------------------------------
declare -a ARGS      # array for options that will be forwarded to pg_probackup

parse_params() {
  if [[ "${1-}" == "" || "${1-}" == "-h" || "${1-}" == "--help" ]] ; then
    usage
    exit 0
  fi

  case "$1" in
      backup|restore|checkdb|merge|show|delete|cleanup)
        ACTIVE_COMMAND=$1
        shift
      ;;
      *) die "Unknown command: $1" ;;
  esac

  while true ; do
    case "${1-}" in
        -B)
            shift
            if [[ "${1-}" != "" ]] ; then
            if [[ -d "${1}" ]] ; then
                BACKUP_REPO=${1};
            else
                die "Error: Backup repository directory does not exist!"
            fi
            else
                die "Error: Backup repository path is not set!"
            fi
        ;;
        -b)
            shift
            if [[ "${1-}" == "full" || "${1-}" == "delta" ]] ; then
                BACKUP_MODE=${1};
            else
                die "Error: Backup mode is incorrect or not set!"
            fi
        ;;
        -n)
          ARGS+=(--no-validate)
        ;;
        -j)
            shift
            if [[ "${1-}" == "" ]] ; then
              die "Error: \"-j\" parameter is set but no number provided!"
            else
              CORES=$(nproc)
              if [[ "${1-}" =~ ^[0-9]+$ && "${1-}" -le $CORES && "${1-}" -gt 0 ]] ; then
                BACKUP_PROC="${1}"
              else
                die "Error: Number of processes set by \"-j\" parameters is incorrect or higher than CPU cores you have"
              fi
            fi
       ;;
        -i)
            shift
            if [[ "${1-}" != "" ]] ; then
              if [[ "${1-}" =~ ^[A-Z0-9]{6}$ ]] ; then
                ARGS+=(-i ${1});
                BACKUP_ID=${1}
              else
                die "Error: Backup ID is incorrect or not set!"
              fi
            else
              die "Error: Backup ID is incorrect or not set!"
            fi
        ;;
        -h | --help) usage ;;
        -v | --verbose) set -x ;;
        --no-color) NO_COLOR=1 ;;
        -?*) die "Error: Unknown option: \"$1\"" ;;
        *) break ;;
    esac
    shift
  done

  # Checking if all necessary parameters are passed to the script
  if [[ "${BACKUP_REPO}" == "none" ]] ; then
    die "Error: Backup repository is not set!"
  fi

  if [[ (${ACTIVE_COMMAND} == "merge" || "${ACTIVE_COMMAND}" == "delete" || "${ACTIVE_COMMAND}" == "show"} ) && "${BACKUP_ID}" == "none" ]] ; then
    die "Error: Backup ID is not set! Please use \"-i <backup ID>\" parameter to specify backup ID."
  fi

  return 0
}

parse_params "$@"
setup_colors

msg "Using backup repository at ${BACKUP_REPO}"

if [[ $(stat -c '%U:%G' "${BACKUP_REPO}") != "${DATA_DIR_OWNER}" ]] ; then
  die "${RED}Error${NOFORMAT}: Backup repository owner must be set to ${DATA_DIR_OWNER}!"
else
  msg "Checking backup repository directory owner: ${GREEN}OK${NOFORMAT}"
fi

role_id=${script_dir##*/}
app_id=$(basename $(dirname $script_dir))
params=$(salt-call pillar.item local pillarenv=$app_id --out=json)

role_type=$(echo $params | $script_dir/jq -r --arg role_id "$role_id" '.local.local[] | select(.Id==$role_id) | .RoleTypeId')
version=$(echo $params | $script_dir/jq -r --arg role_id "$role_id" '.local.local[].Instance | select(.Id==$role_id) | .Version')
PG_USER=$(echo $params | $script_dir/jq -r --arg role_id "$role_id" '.local.local[] | select(.Id==$role_id) | .Params.PgUser')

[[ $role_type != $required_role_type ]] && die "${RED}Script file is outside of ${NOFORMAT}${required_role_type}${RED} role directory${NOFORMAT}"

[[ -z $version || $version = null ]] && die "${RED}Can't detect role version${NOFORMAT}"

postgres_container=$(docker ps -a | grep $role_id | grep postgres | awk 'NR==1{print $1}')
postgres_image=$(docker inspect $postgres_container --format "{{ .Image }}" | cut -d ":" -f 2)

[[ -z $postgres_image ]] && die "${RED}Storage docker image not found${NOFORMAT}"
[[ -z $postgres_container ]] && die "${RED}Storage docker container not found${NOFORMAT}"

# Check backup repo target directory if it's empty or contains any data
if [[ -z $(ls -A "${BACKUP_REPO}") ]] ; then
  echo "Backup repository is empty, performing initialization..."
  docker run --rm -v "${BACKUP_REPO}":/pg_backup "${postgres_image}" pg_probackup-16 init -B /pg_backup
else
  # Let's check if the repo is initialized  proreply
  # We expect that the repo contains two subdirs: "backups" and "wal"
  expected_backup_repo_content=(backups wal) # array of directories we want in BACKUP_REPO
  IFS=$'\n' sorted_expected=($(printf '%s\n' "${expected_backup_repo_content[@]}" | sort))
  unset IFS

  mapfile -t actual_backup_repo_content < <(
      find "$BACKUP_REPO" -mindepth 1 -maxdepth 1 -printf '%f\n' | sort
  )

  if [[ "${actual_backup_repo_content[*]}" != "${sorted_expected[*]}" ]]; then
    die "${RED}Error${NOFORMAT}: Looks like your backup repository is corrupted or contains some other data!"
  fi
fi

# Checking if we have SqlStorage instance in the backup repo exists and created proreply
if [[ -d "${BACKUP_REPO}/backups/${required_role_type}_${version}" ]] && [[ -d "${BACKUP_REPO}/wal/${required_role_type}_${version}" ]] ;
then
  msg "SqlStorage instance found: ${GREEN}OK${NOFORMAT}"
else
  echo "No instance configuration found in the backup repository, trying to set one..."
  docker run --rm -v "${BACKUP_REPO}":"${BACKUP_REPO}" \
                  -v "${script_dir}"/data:/var/lib/postgresql/data "${postgres_image}" \
                  pg_probackup-16 add-instance -B "${BACKUP_REPO}" --instance ${required_role_type}_$version
  docker run --rm -v "${BACKUP_REPO}":"${BACKUP_REPO}" \
                  -v "${script_dir}"/data:/var/lib/postgresql/data "${postgres_image}" \
                  pg_probackup-16 set-config -B "${BACKUP_REPO}" --instance ${required_role_type}_$version -h localhost -U pt_system
fi

echo "Role version: $version"

case "${ACTIVE_COMMAND}" in
  backup)
    if [[ $(docker inspect --format "{{ .State.Status }}" "${postgres_container}") != "running" ]] ; then
      die "${RED}Storage docker image is not running${NOFORMAT}"
    fi
    docker run -ti --rm \
      -v "${BACKUP_REPO}":"${BACKUP_REPO}" \
      -v "${script_dir}"/data:/var/lib/postgresql/data \
      --network=container:${postgres_container} \
      "${postgres_image}" \
      pg_probackup-16 backup -B "${BACKUP_REPO}" \
                             --instance "${required_role_type}_${version}" \
                             -b "${BACKUP_MODE}" --stream --temp-slot --compress --progress \
                             -j "${BACKUP_PROC}" "${ARGS[@]}"
  ;;
  show)
    docker run -ti --rm \
      -v "${BACKUP_REPO}":"${BACKUP_REPO}" \
      "${postgres_image}" \
      pg_probackup-16 show \
        -B "${BACKUP_REPO}" \
        --instance "${required_role_type}_${version}" "${ARGS[@]}"
  ;;
  merge)
    docker run -ti --rm \
      -v ${BACKUP_REPO}:"${BACKUP_REPO}" \
      "${postgres_image}" \
      pg_probackup-16 merge \
          -B "${BACKUP_REPO}" \
          --instance "${required_role_type}_${version}" "${ARGS[@]}"
  ;;
  delete)
    docker run -ti --rm \
      -v ${BACKUP_REPO}:"${BACKUP_REPO}" \
      "${postgres_image}" \
      pg_probackup-16 delete \
          -B "${BACKUP_REPO}" \
          --instance "${required_role_type}_${version}" "${ARGS[@]}"
  ;;
  checkdb)
    if [[ $(docker inspect --format "{{ .State.Status }}" "${postgres_container}") != "running" ]] ; then
      die "${RED}Storage docker image is not running${NOFORMAT}"
    fi
    echo "Starting database check..."
    docker run --rm \
      -v ${BACKUP_REPO}:"${BACKUP_REPO}" \
      -v ${script_dir}/data:/var/lib/postgresql/data \
      --network=container:${postgres_container} \
      "${postgres_image}" \
      pg_probackup-16 checkdb \
          -B "${BACKUP_REPO}" \
          --instance "${required_role_type}_${version}" "${ARGS[@]}"
  ;;
  restore)
    backup_count=$(docker run --rm -v "${BACKUP_REPO}":"${BACKUP_REPO}" ${postgres_image} pg_probackup-16 show -B "${BACKUP_REPO}" \
                              --instance "${required_role_type}_${version}" --format json "${ARGS[@]}" | \
                  grep -v Started | ${script_dir}/jq -r '.[0].backups | length')
    if [[ ${backup_count:-0} -eq "0" ]] ; then
      die "${RED}No valid backups to restore. Exiting. Use \"./probackup.sh show\" command to check existing backups' list.${NOFORMAT}"
    fi
    if [[ $(docker inspect --format "{{ .State.Status }}" "${postgres_container}") == "running" ]] ; then
      echo "Stopping container ${postgres_container}..."
      docker stop ${postgres_container}
    fi
    echo "Deleting old data..."
    find "${script_dir}"/data -maxdepth 1 -mindepth 1 -exec rm -r "{}" \;
    if [[ "${BACKUP_ID}" == "none" ]] ; then
      echo "Restoring data from latest backup..."
    else
      echo "Restoring data from backup with ID ${BACKUP_ID}..."
    fi
    docker run --rm \
      -v "${BACKUP_REPO}":"${BACKUP_REPO}" \
      -v "${script_dir}/data":/var/lib/postgresql/data \
      "${postgres_image}" \
      pg_probackup-16 restore \
        -B "${BACKUP_REPO}"  \
        --instance "${required_role_type}_${version}" \
        -j "${BACKUP_PROC}" "${ARGS[@]}"
    echo "Starting container ${postgres_container}..."
    docker start "${postgres_container}"
  ;;
  cleanup)
    old_backups=$(find "${BACKUP_REPO}" -mindepth 2 -maxdepth 2 -type d -not -name "${required_role_type}_${version}")
    if [[ -z ${old_backups} ]] ; then
      die "No old backups found. There is nothing to cleanup or wrong backup repository path passed with \"-B\" parameter."
    else
    echo -e "\n${RED}You are going to delete old role versions backups. This operation can not be undone!${NOFORMAT}"
    echo -e "${RED}The following directories will be deleted:${NOFORMAT}"
    echo "${old_backups}"
      echo -ne "\n\n${RED}Please confirm the deletion by typing \"Yes\" and pressing \"Enter\": ${NOFORMAT}"
      read -p "" confirm
      if [[ "${confirm}" == "Yes" ]] ; then
        find "${BACKUP_REPO}" -mindepth 2 -maxdepth 2 -type d -not -name "${required_role_type}_${version}" -exec rm -rf {} \;
        echo -e "${GREEN}Old role version backup data deleted${NOFORMAT}"
      else
        echo "Exiting..."
      fi
    fi
  ;;
esac
