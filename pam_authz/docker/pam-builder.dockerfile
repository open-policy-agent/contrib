FROM golang:1.8

RUN apt-get update && apt-get install -y \
    libpam0g-dev && \
    rm -rf /var/lib/apt/lists/*

COPY pam /pam
COPY docker/run.dockerfile /run.dockerfile
COPY docker/keys/id_rsa.pub /id_rsa.pub
COPY docker/etc/pam.d /pam.d
COPY docker/etc/sshd_config /sshd_config
COPY docker/create_user.sh /create_user.sh

WORKDIR /pam

RUN make clean && make

CMD tar -cf - \
        -C /pam pam_authz.so \
        -C / run.dockerfile id_rsa.pub pam.d sshd_config create_user.sh
