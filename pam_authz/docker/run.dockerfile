FROM phusion/baseimage

ARG identifier
ENV identifier=$identifier

# Enable the sshd service
RUN rm -f /etc/service/sshd/down

# Regenerate SSH host keys. baseimage-docker does not contain any, so you
# have to do that yourself. You may also comment out this instruction; the
# init system will auto-generate one during boot.
RUN /etc/my_init.d/00_regen_ssh_host_keys.sh

RUN apt-get update && apt-get install -y \
	vim \
	sudo && \
	rm -rf /var/lib/apt/lists/*

# Copy the JSON library contents from the build container and install them.
COPY jansson_lib /jansson_lib
RUN tar -xvf /jansson_lib -C /usr/lib && \
    rm -f /jansson_lib

# Do not cache sudo authorization (that is, check the PAM auth stack every invocation)
RUN sed -i 's/env_reset/env_reset,timestamp_timeout=0/g' /etc/sudoers
# Create a link where processes can write to the container's stdout
RUN ln -sf /proc/1/fd/1 /var/log/stdout.log
RUN mkdir -p /lib/security

# Copy the demonstration key into the repo
COPY id_rsa.pub /tmp/your_key.pub
RUN cat /tmp/your_key.pub >> /root/.ssh/authorized_keys

# Create some user accounts. All of these users use the same ssh key for authentication
COPY create_user.sh /create_user.sh
RUN /create_user.sh frontend-dev backend-dev ops alice bob charlie dave eve frank sally 

# Replace the default ssh and sudo PAM configs. Our config requires the PAM authz
# authorization module and disables standard linux authorization
COPY /pam.d/* /etc/pam.d/

# Replace the default sshd config with our config. These enables PAM in sshd
COPY /sshd_config /etc/ssh/sshd_config
COPY pam_authz.so /lib/security/pam_authz.so

# Delete the ssh key common to all users
RUN rm -f /tmp/your_key.pub
