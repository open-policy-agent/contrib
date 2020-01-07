#!/bin/bash

set -xe

for username in "$@"
do
    useradd -G sudo -ms /bin/bash $username
    mkdir -p /home/$username/.ssh
    cat /tmp/your_key.pub >> /home/$username/.ssh/authorized_keys
    chown -R $username:$username /home/$username/.ssh
    chmod 700 /home/$username/.ssh
    chmod 600 /home/$username/.ssh/authorized_keys
done

