FROM ubuntu:20.04

# Install prerequisite
RUN apt update && \
    apt-get install software-properties-common -y
RUN apt-add-repository ppa:git-core/ppa && \
    apt update && \
    apt install -y wget curl build-essential git git-lfs unzip

# Use Bash instead of Dash
RUN ln -sf bash /bin/sh

# Install python3
RUN apt install -y python3 python3-pip

# Install azure-cli
RUN curl -sL https://aka.ms/InstallAzureCLIDeb | bash

# Install docker
RUN apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common lsb-release && \
    curl -fsSL https://download.docker.com/linux/$(lsb_release -is | tr '[:upper:]' '[:lower:]')/gpg | apt-key add - 2>/dev/null && \
    add-apt-repository "deb [arch=$(dpkg --print-architecture)] https://download.docker.com/linux/$(lsb_release -is | tr '[:upper:]' '[:lower:]') $(lsb_release -cs) stable" && \
    apt-get update &&\
    apt-get install -y docker-ce-cli

# Install golang
ARG GO_VERSION
RUN curl -LO https://golang.org/dl/go${GO_VERSION}.linux-$(dpkg --print-architecture).tar.gz && \
    tar -C /usr/local -xvzf go${GO_VERSION}.linux-$(dpkg --print-architecture).tar.gz && \
    rm -rf go${GO_VERSION}.linux-$(dpkg --print-architecture).tar.gz

# Install kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/$(dpkg --print-architecture)/kubectl" && \
    curl -LO "https://dl.k8s.io/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/$(dpkg --print-architecture)/kubectl.sha256" && \
    echo "$(cat kubectl.sha256)  kubectl" | sha256sum --check && \
    chmod +x ./kubectl && mv ./kubectl /usr/bin/kubectl && \
    rm kubectl.sha256

# Install operator-sdk
RUN apt-get update && \
    apt-get install -y jq && \
    ARCH="$(case $(uname -m) in x86_64) echo -n amd64 ;; aarch64) echo -n arm64 ;; *) echo -n $(uname -m) ;; esac)" && \
    OS=$(uname | awk '{print tolower($0)}') && \
    RELEASE_VERSION="$(curl -s https://api.github.com/repos/operator-framework/operator-sdk/releases/latest | jq -r .tag_name)" && \
    curl -LO https://github.com/operator-framework/operator-sdk/releases/download/${RELEASE_VERSION}/operator-sdk_${OS}_${ARCH} && \
    curl -LO https://github.com/operator-framework/operator-sdk/releases/download/${RELEASE_VERSION}/checksums.txt && \
    curl -LO https://github.com/operator-framework/operator-sdk/releases/download/${RELEASE_VERSION}/checksums.txt.asc && \
    gpg --keyserver keyserver.ubuntu.com --recv-key 052996E2A20B5C7E && \
    gpg --verify checksums.txt.asc && \
    sha256sum --ignore-missing --check checksums.txt && \
    chmod +x operator-sdk_${OS}_${ARCH} && \
    mkdir -p /usr/local/bin/ && \
    mv operator-sdk_${OS}_${ARCH} /usr/local/bin/operator-sdk

ENV PATH=${PATH}:/usr/local/go/bin
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=${PATH}:${GOPATH}/bin

# Install gh
RUN curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list > /dev/null && \
    apt update && \
    apt install -y gh

# Protocol Buffer Compiler
RUN PROTOC_VERSION=21.9 \
    && if [ $(dpkg --print-architecture) = "amd64" ]; then PROTOC_ARCH="x86_64"; else PROTOC_ARCH="aarch_64" ; fi \
    && curl -LO "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-$PROTOC_ARCH.zip" \
    && unzip "protoc-${PROTOC_VERSION}-linux-$PROTOC_ARCH.zip" -d $HOME/.local \
    && mv $HOME/.local/bin/protoc /usr/local/bin/protoc \
    && mv $HOME/.local/include/ /usr/local/bin/include/
