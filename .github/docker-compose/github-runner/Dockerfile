FROM ubuntu:24.04

ARG RUNNER_VERSION="2.334.0"
ARG DEBIAN_FRONTEND=noninteractive

RUN apt update -y && apt upgrade -y && useradd -m docker
RUN apt install -y --no-install-recommends sudo && \
    echo "docker ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers.d/docker && \
    chmod 0440 /etc/sudoers.d/docker && \
    apt clean
RUN apt install -y --no-install-recommends \
    curl unzip jq build-essential libicu74 libssl-dev libffi-dev python3 python3-venv python3-dev python3-pip openssh-client sshpass wget git
RUN curl -sL https://api.github.com/repos/cli/cli/releases/latest \
        | grep "browser_download_url.*_linux_arm64.deb" \
        | cut -d '"' -f 4 \
        | wget -qi - -O gh_latest.deb
RUN dpkg -i gh_latest.deb
RUN curl -fsSL https://go.dev/dl/go1.26.3.linux-arm64.tar.gz -o go.tar.gz && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm go.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/bin/go
RUN curl -fsSL https://deb.nodesource.com/setup_22.x | bash - && \
    apt-get install -y nodejs
RUN npm install -g yarn
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip" -o "awscliv2.zip"
RUN unzip awscliv2.zip
RUN ./aws/install
RUN apt install docker.io -y
RUN cd /home/docker && mkdir actions-runner && cd actions-runner \
    && curl -O -L https://github.com/actions/runner/releases/download/v${RUNNER_VERSION}/actions-runner-linux-arm64-${RUNNER_VERSION}.tar.gz \
    && tar xzf ./actions-runner-linux-arm64-${RUNNER_VERSION}.tar.gz
RUN chown -R docker ~docker && /home/docker/actions-runner/bin/installdependencies.sh
COPY start.sh start.sh
RUN chmod +x start.sh
COPY gcp_rsa /home/docker/.ssh/gcp_rsa
RUN chmod 400 /home/docker/.ssh/gcp_rsa
RUN chown docker:docker -R /home/docker/.ssh/

USER docker
ENTRYPOINT ["./start.sh"]
