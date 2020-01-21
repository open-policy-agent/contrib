FROM ubuntu:18.04

RUN apt-get update \
    && apt-get -y install --no-install-recommends git openssl build-essential ca-certificates nano curl python3 python3-dev python3-pip python3-setuptools python3-wheel

#Imstall KIND
RUN curl -Lo ./kind https://github.com/kubernetes-sigs/kind/releases/download/v0.6.1/kind-linux-amd64 && chmod +x ./kind && mv ./kind /usr/bin/kind
#Install OPA
RUN curl -L -o opa https://openpolicyagent.org/downloads/latest/opa_linux_amd64 && mv ./opa /usr/bin/opa && chmod +x /usr/bin/opa
#Install Kubectl 
RUN curl -sSL -o /usr/local/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl \
    && chmod +x /usr/local/bin/kubectl

#Install Docker
RUN apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common lsb-release \
    && curl -fsSL https://download.docker.com/linux/$(lsb_release -is | tr '[:upper:]' '[:lower:]')/gpg | apt-key add - 2>/dev/null \
    && add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/$(lsb_release -is | tr '[:upper:]' '[:lower:]') $(lsb_release -cs) stable" \
    && apt-get update \
    && apt-get install -y docker-ce-cli

COPY ./requirements.txt ./requirements.txt
RUN pip3 install -r ./requirements.txt