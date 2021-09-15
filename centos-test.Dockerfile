FROM centos

RUN yum install -y curl yum-utils

# install docker cli
RUN yum-config-manager \
    --add-repo \
    https://download.docker.com/linux/centos/docker-ce.repo
RUN yum install -y docker-ce-cli

# install compose v1
RUN curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
RUN chmod +x /usr/local/bin/docker-compose

# install compose v2
RUN mkdir -p /usr/local/lib/docker/cli-plugins
RUN curl -L https://github.com/docker/compose/releases/download/v2.0.0-rc.3/docker-compose-linux-amd64 > /usr/local/lib/docker/cli-plugins/docker-compose
RUN chmod +x /usr/local/lib/docker/cli-plugins/docker-compose

COPY install_on_linux.sh /tmp/install_on_linux.sh
RUN /tmp/install_on_linux.sh

RUN docker-compose --version
