#!/bin/bash

RITA_VERSION="REPLACE_ME"

set -e 

install_target="$1"

# change working directory to directory of this script
pushd "$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")" > /dev/null

source ./scripts/helper.sh


./scripts/ansible-installer.sh


status "Installing rita via ansible on $install_target"		#================
if [ "$install_target" = "localhost" -o "$install_target" = "127.0.0.1" -o "$install_target" = "::1" ]; then
        if [ "$(uname)" == "Darwin" ]; then
            # TODO support macOS install target
            echo "${YELLOW}Installing RITA via Ansible on the local system is not yet supported on MacOS.${NORMAL}"
            exit 1
        fi
    ansible-playbook --connection=local -K -i "127.0.0.1," -e "install_hosts=127.0.0.1," install_rita.yml
else
    status "Setting up future ssh connections to $install_target .  You may be asked to provide your ssh password to this system."		#================
    ./scripts/sshprep "$install_target"
    ansible-playbook -K -i "${install_target}," -e "install_hosts=${install_target}," install_rita.yml
fi

    # ansible-playbook -i ../digitalocean_inventory.py -e "install_hosts=all" install_rita.yml

echo \
"
░█▀▀█ ▀█▀ ▀▀█▀▀ ─█▀▀█
░█▄▄▀ ░█─ ─░█── ░█▄▄█
░█─░█ ▄█▄ ─░█── ░█─░█ ${RITA_VERSION}

Brought to you by Active CounterMeasures©
"
echo "RITA was successfully installed!"

# switch back to original working directory
popd > /dev/null