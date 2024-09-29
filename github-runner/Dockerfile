FROM ubuntu:24.04

ARG RUNNER_VERSION="2.319.1"
ARG DEBIAN_FRONTEND=noninteractive

RUN apt update -y && apt upgrade -y && useradd -m docker
RUN apt install -y --no-install-recommends \
    curl unzip jq build-essential libssl-dev libffi-dev python3 python3-venv python3-dev python3-pip openssh-client sshpass nodejs npm golang
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
RUN unzip awscliv2.zip
RUN ./aws/install
RUN aws --version
RUN npm install -g yarn
RUN apt install docker.io -y
RUN docker --version
RUN cd /home/docker && mkdir actions-runner && cd actions-runner \
    && curl -O -L https://github.com/actions/runner/releases/download/v${RUNNER_VERSION}/actions-runner-linux-x64-${RUNNER_VERSION}.tar.gz \
    && tar xzf ./actions-runner-linux-x64-${RUNNER_VERSION}.tar.gz
RUN chown -R docker ~docker && /home/docker/actions-runner/bin/installdependencies.sh
COPY start.sh start.sh
RUN chmod +x start.sh
COPY gcp_rsa /home/docker/.ssh/gcp_rsa
RUN chmod 400 /home/docker/.ssh/gcp_rsa
RUN chown docker:docker -R /home/docker/.ssh/

USER docker
ENTRYPOINT ["./start.sh"]