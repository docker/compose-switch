FROM ubuntu

# install docker CLI
RUN apt update
RUN apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
RUN echo \
  "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
RUN apt-get update
RUN apt-get install -y docker-ce-cli

# install compose v1 (ubuntu managed)
RUN apt-get install -y docker-compose

# install compose v2
RUN mkdir -p /usr/local/lib/docker/cli-plugins
RUN curl -L https://github.com/docker/compose/releases/download/v2.0.0-rc.3/docker-compose-linux-amd64 > /usr/local/lib/docker/cli-plugins/docker-compose
RUN chmod +x /usr/local/lib/docker/cli-plugins/docker-compose

COPY install_on_linux.sh /tmp/install_on_linux.sh
RUN /tmp/install_on_linux.sh

RUN docker-compose --version
