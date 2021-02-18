#!/bin/bash
# https://github.com/complexorganizations/shell-script-boilerplate

# Require script to be run as root
function super-user-check() {
  if [ "$EUID" -ne 0 ]; then
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
    DISTRO=$ID
    DISTRO_VERSION=$VERSION_ID
  fi
}

# Check Operating System
dist-check

# Pre-Checks system requirements
function installing-system-requirements() {
  if { [ "$DISTRO" == "ubuntu" ] || [ "$DISTRO" == "debian" ] || [ "$DISTRO" == "raspbian" ] || [ "$DISTRO" == "pop" ] || [ "$DISTRO" == "kali" ] || [ "$DISTRO" == "linuxmint" ] || [ "$DISTRO" == "fedora" ] || [ "$DISTRO" == "centos" ] || [ "$DISTRO" == "rhel" ] || [ "$DISTRO" == "arch" ] || [ "$DISTRO" == "manjaro" ] || [ "$DISTRO" == "alpine" ] || [ "$DISTRO" == "freebsd" ]; }; then
    if { [ ! -x "$(command -v curl)" ] || [ ! -x "$(command -v iptables)" ]; }; then
      if { [ "$DISTRO" == "ubuntu" ] || [ "$DISTRO" == "debian" ] || [ "$DISTRO" == "raspbian" ] || [ "$DISTRO" == "pop" ] || [ "$DISTRO" == "kali" ] || [ "$DISTRO" == "linuxmint" ]; }; then
        apt-get update
      elif { [ "$DISTRO" == "fedora" ] || [ "$DISTRO" == "centos" ] || [ "$DISTRO" == "rhel" ]; }; then
        yum update -y
      elif { [ "$DISTRO" == "arch" ] || [ "$DISTRO" == "manjaro" ]; }; then
        pacman -Syu
      elif [ "$DISTRO" == "alpine" ]; then
        apk update
      elif [ "$DISTRO" == "freebsd" ]; then
        pkg update
      fi
    fi
  else
    echo "Error: $DISTRO not supported."
    exit
  fi
}

# Run the function and check for requirements
installing-system-requirements

# Global variables
GLOBAL_VARIABLES="/config/file/path"

if [ ! -f "$GLOBAL_VARIABLES" ]; then

  # comments for the first question
  function first-question() {
    echo "What is the first question that u want to ask the user?"
    echo "  1) Ansewer #1 (Recommended)"
    echo "  2) Ansewer #2"
    echo "  3) Custom (Advanced)"
    until [[ "$FIRST_QUESTION_SETTINGS" =~ ^[1-3]$ ]]; do
      read -rp "Subnetwork choice [1-3]: " -e -i 1 FIRST_QUESTION_SETTINGS
    done
    case $FIRST_QUESTION_SETTINGS in
    1)
      FIRST_QUESTION="Ansewer #1"
      ;;
    2)
      FIRST_QUESTION="Ansewer #2"
      ;;
    3)
      read -rp "User text: " -e -i "Ansewer #3" FIRST_QUESTION
      ;;
    esac
  }

  # comments for the first question
  first-question

  ### use the code above to ask any questions as u want.
  function install-the-app() {
    if { [ "$DISTRO" == "ubuntu" ] || [ "$DISTRO" == "debian" ] || [ "$DISTRO" == "raspbian" ] || [ "$DISTRO" == "pop" ] || [ "$DISTRO" == "kali" ] || [ "$DISTRO" == "linuxmint" ]; }; then
      apt-get update
    elif { [ "$DISTRO" == "fedora" ] || [ "$DISTRO" == "centos" ] || [ "$DISTRO" == "rhel" ]; }; then
      yum update -y
    elif { [ "$DISTRO" == "arch" ] || [ "$DISTRO" == "manjaro" ]; }; then
      pacman -Syu
    elif [ "$DISTRO" == "alpine" ]; then
      apk update
    elif [ "$DISTRO" == "freebsd" ]; then
      pkg update
    fi
  }

  # run the function
  install-the-app

  # configure service here
  function config-service() {
    if { [ "$DISTRO" == "ubuntu" ] || [ "$DISTRO" == "debian" ] || [ "$DISTRO" == "raspbian" ] || [ "$DISTRO" == "pop" ] || [ "$DISTRO" == "kali" ] || [ "$DISTRO" == "linuxmint" ] || [ "$DISTRO" == "fedora" ] || [ "$DISTRO" == "centos" ] || [ "$DISTRO" == "rhel" ] || [ "$DISTRO" == "arch" ] || [ "$DISTRO" == "manjaro" ] || [ "$DISTRO" == "alpine" ] || [ "$DISTRO" == "freebsd" ]; }; then
      echo $GLOBAL_VARIABLES
      echo "$FIRST_QUESTION"
      echo "$DISTRO"
      echo "$DISTRO_VERSION"
    fi
  }

  # run the function
  config-service

  function service-manager() {
    if pgrep systemd-journal; then
      systemctl disable SERVICE
      systemctl stop SERVICE
    else
      service SERVICE disable
      service SERVICE stop
    fi
  }

  # restart the chrome service
  service-manager

else

  # take user input
  function take-user-input() {
    echo "What do you want to do?"
    echo "   1) Option #1"
    echo "   2) Option #2"
    echo "   3) Option #3"
    echo "   4) Option #4"
    echo "   5) Option #5"
    until [[ "$USER_OPTIONS" =~ ^[0-9]+$ ]] && [ "$USER_OPTIONS" -ge 1 ] && [ "$USER_OPTIONS" -le 5 ]; do
      read -rp "Select an Option [1-5]: " -e -i 1 USER_OPTIONS
    done
    case $USER_OPTIONS in
    1)
      echo "Hello, World!"
      ;;
    2)
      echo "Hello, World!"
      ;;
    3)
      echo "Hello, World!"
      ;;
    4)
      echo "Hello, World!"
      ;;
    5)
      echo "Hello, World!"
      ;;
    esac
  }

  # run the function
  take-user-input

fi
