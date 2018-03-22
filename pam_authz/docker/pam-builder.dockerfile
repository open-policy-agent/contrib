FROM golang:1.8

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    libpam0g-dev \
    apt-utils \
    libcurl4-gnutls-dev && \
    rm -rf /var/lib/apt/lists/*

# Fetch and install the JSON library.
# Source is at: https://github.com/akheron/jansson
RUN wget http://www.digip.org/jansson/releases/jansson-2.11.tar.gz && \
    tar -xvf jansson-2.11.tar.gz && \
    cd jansson-2.11 && \
    ./configure --prefix=/usr && \
    make && \
    make check && \
    make install

# Archive the JSON library contents for installation in the running container.
RUN cd /usr/lib && \
    tar -cf jansson_lib libjansson* && \
    mv jansson_lib /

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
        -C / run.dockerfile id_rsa.pub pam.d sshd_config create_user.sh jansson_lib