#!/bin/bash
# https://github.com/complexorganizations/unbound-manager

# Require script to be run as root
function super-user-check() {
  if [ "${EUID}" -ne 0 ]; then
    echo "You need to run this script as super user."
    exit
  fi
}

# Check for root
super-user-check

# Detect Operating System
function dist-check() {
  if [ -e /etc/os-release ]; then
    # shellcheck disable=SC1091
    source /etc/os-release
    DISTRO=${ID}
  fi
}

# Check Operating System
dist-check

# Pre-Checks system requirements
function installing-system-requirements() {
  if { [ "${DISTRO}" == "ubuntu" ] || [ "${DISTRO}" == "debian" ] || [ "${DISTRO}" == "raspbian" ] || [ "${DISTRO}" == "pop" ] || [ "${DISTRO}" == "kali" ] || [ "${DISTRO}" == "linuxmint" ] || [ "${DISTRO}" == "fedora" ] || [ "${DISTRO}" == "centos" ] || [ "${DISTRO}" == "rhel" ] || [ "${DISTRO}" == "arch" ] || [ "${DISTRO}" == "manjaro" ] || [ "${DISTRO}" == "alpine" ] || [ "${DISTRO}" == "freebsd" ]; }; then
    if { [ ! -x "$(command -v curl)" ] || [ ! -x "$(command -v cron)" ]; }; then
      if { [ "${DISTRO}" == "ubuntu" ] || [ "${DISTRO}" == "debian" ] || [ "${DISTRO}" == "raspbian" ] || [ "${DISTRO}" == "pop" ] || [ "${DISTRO}" == "kali" ] || [ "${DISTRO}" == "linuxmint" ]; }; then
        apt-get update && apt-get install curl cron -y
      elif { [ "${DISTRO}" == "fedora" ] || [ "${DISTRO}" == "centos" ] || [ "${DISTRO}" == "rhel" ]; }; then
        yum update -y && yum install curl cron -y
      elif { [ "${DISTRO}" == "arch" ] || [ "${DISTRO}" == "manjaro" ]; }; then
        pacman -Syu && pacman -Syu --noconfirm --needed curl cronie
      elif [ "${DISTRO}" == "alpine" ]; then
        apk update && apk add curl cron
      elif [ "${DISTRO}" == "freebsd" ]; then
        pkg update && pkg install curl cron
      fi
    fi
  else
    echo "Error: ${DISTRO} not supported."
    exit
  fi
}

# Run the function and check for requirements
installing-system-requirements

# Global variables
RESOLV_CONFIG="/etc/resolv.conf"
RESOLV_CONFIG_OLD="/etc/resolv.conf.old"
UNBOUND_ROOT="/etc/unbound"
UNBOUND_MANAGER="${UNBOUND_ROOT}/unbound-manager"
UNBOUND_CONFIG="${UNBOUND_ROOT}/unbound.conf"
UNBOUND_ROOT_HINTS="${UNBOUND_ROOT}/root.hints"
UNBOUND_ANCHOR="/var/lib/unbound/root.key"
UNBOUND_ROOT_SERVER_CONFIG_URL="https://www.internic.net/domain/named.cache"
UNBOUND_CONFIG_HOST_URL="https://raw.githubusercontent.com/complexorganizations/unbound-manager/main/configs/host"
UNBOUND_CONFIG_HOST="/etc/unbound/unbound.conf.d/host.conf"
UNBOUND_CONFIG_HOST_TMP="/tmp/host"
UNBOUND_MANAGER_UPDATE_URL="https://raw.githubusercontent.com/complexorganizations/unbound-manager/main/unbound-manager.sh"

function usage-help() {
  echo "usage: ./$(basename "$0") <command>"
  echo "  --install     Install Unbound"
  echo "  --start       Start Unbound"
  echo "  --stop        Stop Unbound"
  echo "  --restart     Restart Unbound"
  echo "  --uninstall   Uninstall Unbound"
  echo "  --update      Update Unbound Manager"
  echo "  --help        Show Usage Guide"
}

# The usage of the script
function usage() {
  while [ $# -ne 0 ]; do
    case ${1} in
    --install)
      shift
      HEADLESS_INSTALL=${HEADLESS_INSTALL:-y}
      ;;
    --start)
      shift
      USER_OPTIONS=${USER_OPTIONS:-1}
      ;;
    --stop)
      shift
      USER_OPTIONS=${USER_OPTIONS:-2}
      ;;
    --restart)
      shift
      USER_OPTIONS=${USER_OPTIONS:-3}
      ;;
    --uninstall)
      shift
      USER_OPTIONS=${USER_OPTIONS:-4}
      ;;
    --update)
      shift
      USER_OPTIONS=${USER_OPTIONS:-5}
      ;;
    --help)
      shift
      usage-help
      ;;
    *)
      echo "Invalid argument: ${1}"
      usage-help
      exit
      ;;
    esac
  done
}

usage "$@"

function headless-install() {
  if [[ ${HEADLESS_INSTALL} =~ ^[Yy]$ ]]; then
    LIST_CHOICE_SETTINGS=${LIST_CHOICE_SETTINGS:-1}
    AUTOMATIC_UPDATES_SETTINGS=${AUTOMATIC_UPDATES_SETTINGS:-1}
  fi
}

# No GUI
headless-install

if [ ! -f "${UNBOUND_MANAGER}" ]; then

  # Function to install unbound
  function install-unbound() {
    if [ ! -x "$(command -v unbound)" ]; then
      if [ "${DISTRO}" == "ubuntu" ]; then
        apt-get install unbound unbound-host e2fsprogs -y
        if pgrep systemd-journal; then
          systemctl stop systemd-resolved
          systemctl disable systemd-resolved
        else
          service systemd-resolved stop
          service systemd-resolved disable
        fi
      elif { [ "${DISTRO}" == "debian" ] || [ "${DISTRO}" == "raspbian" ] || [ "${DISTRO}" == "pop" ] || [ "${DISTRO}" == "kali" ] || [ "${DISTRO}" == "linuxmint" ]; }; then
        apt-get install unbound unbound-host e2fsprogs -y
      elif { [ "${DISTRO}" == "centos" ] || [ "${DISTRO}" == "rhel" ]; }; then
        yum install unbound unbound-libs -y
      elif [ "${DISTRO}" == "fedora" ]; then
        dnf install unbound -y
      elif { [ "${DISTRO}" == "arch" ] || [ "${DISTRO}" == "archarm" ] || [ "${DISTRO}" == "manjaro" ]; }; then
        pacman -Syu --noconfirm unbound
      elif [ "${DISTRO}" == "alpine" ]; then
        apk add unbound
      elif [ "${DISTRO}" == "freebsd" ]; then
        pkg install unbound
      fi
      if [ -f "${UNBOUND_ANCHOR}" ]; then
        rm -f ${UNBOUND_ANCHOR}
      fi
      if [ -f "${UNBOUND_CONFIG}" ]; then
        rm -f ${UNBOUND_CONFIG}
      fi
      if [ -f "${UNBOUND_ROOT_HINTS}" ]; then
        rm -f ${UNBOUND_ROOT_HINTS}
      fi
      if [ -d "${UNBOUND_ROOT}" ]; then
        unbound-anchor -a ${UNBOUND_ANCHOR}
        curl ${UNBOUND_ROOT_SERVER_CONFIG_URL} -o ${UNBOUND_ROOT_HINTS}
        NPROC=$(nproc)
        echo "server:
num-threads: ${NPROC}
verbosity: 1
root-hints: ${UNBOUND_ROOT_HINTS}
auto-trust-anchor-file: ${UNBOUND_ANCHOR}
interface: 0.0.0.0
interface: ::0
max-udp-size: 3072
access-control: 0.0.0.0/0 allow
access-control: ::0 allow
access-control: 127.0.0.1 allow
do-tcp: no
hide-identity: yes
hide-version: yes
harden-glue: yes
harden-dnssec-stripped: yes
harden-referral-path: yes
unwanted-reply-threshold: 10000000
val-log-level: 1
cache-min-ttl: 1800
cache-max-ttl: 14400
prefetch: yes
qname-minimisation: yes
prefetch-key: yes" >>${UNBOUND_CONFIG}
      fi
      if [ -f "${RESOLV_CONFIG_OLD}" ]; then
        rm -f ${RESOLV_CONFIG_OLD}
      fi
      if [ -f "${RESOLV_CONFIG}" ]; then
        chattr -i ${RESOLV_CONFIG}
        mv ${RESOLV_CONFIG} ${RESOLV_CONFIG_OLD}
        echo "nameserver 127.0.0.1" >>${RESOLV_CONFIG}
        echo "nameserver ::1" >>${RESOLV_CONFIG}
        chattr +i ${RESOLV_CONFIG}
      else
        echo "nameserver 127.0.0.1" >>${RESOLV_CONFIG}
        echo "nameserver ::1" >>${RESOLV_CONFIG}
      fi
      echo "Unbound: true" >>${UNBOUND_MANAGER}
      # restart unbound
      if pgrep systemd-journal; then
        systemctl reenable unbound
        systemctl restart unbound
      else
        service unbound enable
        service unbound restart
      fi
    fi
  }

  # Running Install Unbound
  install-unbound

  function choose-your-list() {
    echo "Which list do you want to use?"
    echo "1) All (Recommended)"
    echo "2) No (Advanced)"
    until [[ "${LIST_CHOICE_SETTINGS}" =~ ^[1-2]$ ]]; do
      read -rp "List Choice [1-2]: " -e -i 1 LIST_CHOICE_SETTINGS
    done
    case ${LIST_CHOICE_SETTINGS} in
    1)
      echo "include: ${UNBOUND_CONFIG_HOST}" >>${UNBOUND_CONFIG}
      curl "${UNBOUND_CONFIG_HOST_URL}" -o ${UNBOUND_CONFIG_HOST_TMP}
      sed -i -e "s_.*_0.0.0.0 &_" ${UNBOUND_CONFIG_HOST_TMP}
      grep "^0\.0\.0\.0" "${UNBOUND_CONFIG_HOST_TMP}" | awk '{print "local-data: \""$2" IN A 0.0.0.0\""}' >"${UNBOUND_CONFIG_HOST}"
      rm -f ${UNBOUND_CONFIG_HOST_TMP}
      ;;
    2)
      echo "There are no lists selected."
      ;;
    esac
  }

  choose-your-list

  # real-time updates
  function enable-automatic-updates() {
    echo "Would you like to setup real-time updates?"
    echo "1) Yes (Recommended)"
    echo "2) No (Advanced)"
    until [[ "${AUTOMATIC_UPDATES_SETTINGS}" =~ ^[1-2]$ ]]; do
      read -rp "Automatic Updates [1-2]: " -e -i 1 AUTOMATIC_UPDATES_SETTINGS
    done
    case ${AUTOMATIC_UPDATES_SETTINGS} in
    1)
      crontab -l | {
        cat
        echo "0 0 * * * $(realpath "$0") --update"
      } | crontab -
      if pgrep systemd-journal; then
        systemctl enable cron
        systemctl start cron
      else
        service cron enable
        service cron start
      fi
      ;;
    2)
      echo "Real-time Updates Disabled"
      ;;
    esac
  }

  # real-time updates
  enable-automatic-updates

  # Install unbound manager
  function install-unbound-manager-file() {
    if [ -d "${UNBOUND_ROOT}" ]; then
      if [ ! -f "${UNBOUND_MANAGER}" ]; then
        echo "Unbound Manager: true" >>${UNBOUND_MANAGER}
      fi
    fi
  }

  # Unbound unbound
  install-unbound-manager-file

else

  # take user input
  function take-user-input() {
    echo "What do you want to do?"
    echo " 1) Start Unbound"
    echo " 2) Stop Unbound"
    echo " 3) Restart Unbound"
    echo " 4) Uninstall Unbound"
    echo " 5) Update Unbound Manager"
    until [[ "$USER_OPTIONS" =~ ^[0-9]+$ ]] && [ "$USER_OPTIONS" -ge 1 ] && [ "$USER_OPTIONS" -le 5 ]; do
      read -rp "Select an Option [1-5]: " -e -i 1 USER_OPTIONS
    done
    case $USER_OPTIONS in
    1)
      if [ -x "$(command -v unbound)" ]; then
        if pgrep systemd-journal; then
          systemctl start unbound
        else
          service unbound start
        fi
      fi
      ;;
    2)
      if [ -x "$(command -v unbound)" ]; then
        if pgrep systemd-journal; then
          systemctl stop unbound
        else
          service unbound stop
        fi
      fi
      ;;
    3)
      if [ -x "$(command -v unbound)" ]; then
        if pgrep systemd-journal; then
          systemctl restart unbound
        else
          service unbound restart
        fi
      fi
      ;;
    4)
      if [ -x "$(command -v unbound)" ]; then
        if [ -f "${UNBOUND_MANAGER}" ]; then
          if pgrep systemd-journal; then
            systemctl disable unbound
            systemctl stop unbound
          else
            service unbound disable
            service unbound stop
          fi
          if [ -f "${RESOLV_CONFIG_OLD}" ]; then
            rm -f ${RESOLV_CONFIG}
            mv ${RESOLV_CONFIG_OLD} ${RESOLV_CONFIG}
          fi
          if { [ "${DISTRO}" == "centos" ] || [ "${DISTRO}" == "rhel" ]; }; then
            yum remove unbound unbound-host -y
          elif { [ "${DISTRO}" == "debian" ] || [ "${DISTRO}" == "pop" ] || [ "${DISTRO}" == "ubuntu" ] || [ "${DISTRO}" == "raspbian" ] || [ "${DISTRO}" == "kali" ] || [ "${DISTRO}" == "linuxmint" ]; }; then
            apt-get remove --purge unbound unbound-host -y
          elif { [ "${DISTRO}" == "arch" ] || [ "${DISTRO}" == "manjaro" ]; }; then
            pacman -Rs unbound unbound-host -y
          elif [ "${DISTRO}" == "fedora" ]; then
            dnf remove unbound -y
          elif [ "${DISTRO}" == "alpine" ]; then
            apk del unbound
          elif [ "${DISTRO}" == "freebsd" ]; then
            pkg delete unbound
          fi
          if [ -f "${UNBOUND_MANAGER}" ]; then
            rm -f ${UNBOUND_MANAGER}
          fi
          if [ -f "${UNBOUND_CONFIG}" ]; then
            rm -f ${UNBOUND_CONFIG}
          fi
          if [ -f "${UNBOUND_ANCHOR}" ]; then
            rm -f ${UNBOUND_ANCHOR}
          fi
          if [ -f "${UNBOUND_ROOT_HINTS}" ]; then
            rm -f ${UNBOUND_ROOT_HINTS}
          fi
          if [ -f "${UNBOUND_ROOT}" ]; then
            rm -f ${UNBOUND_ROOT}
          fi
        fi
      fi
      ;;
    5)
      # Update script
      CURRENT_FILE_PATH="$(realpath "$0")"
      if [ -f "${CURRENT_FILE_PATH}" ]; then
        curl -o "${CURRENT_FILE_PATH}" ${UNBOUND_MANAGER_UPDATE_URL}
        chmod +x "${CURRENT_FILE_PATH}" || exit
      fi
      # Update root hints
      if [ -f "${UNBOUND_ROOT_HINTS}" ]; then
        curl ${UNBOUND_ROOT_SERVER_CONFIG_URL} -o ${UNBOUND_ROOT_HINTS}
      fi
      # Update Host List
      if [ -f "${UNBOUND_CONFIG_HOST}" ]; then
        rm -f ${UNBOUND_CONFIG_HOST}
        curl "${UNBOUND_CONFIG_HOST_URL}" -o ${UNBOUND_CONFIG_HOST_TMP}
        sed -i -e "s_.*_0.0.0.0 &_" ${UNBOUND_CONFIG_HOST_TMP}
        grep "^0\.0\.0\.0" "${UNBOUND_CONFIG_HOST_TMP}" | awk '{print "local-data: \""$2" IN A 0.0.0.0\""}' >"${UNBOUND_CONFIG_HOST}"
        rm -f ${UNBOUND_CONFIG_HOST_TMP}
      fi
      ;;
    esac
  }

  # run the function
  take-user-input

fi
