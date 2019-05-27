FROM centos:latest


WORKDIR /upgrade_agent

COPY . ./
RUN curl -L https://github.com/docker/compose/releases/download/1.23.2/docker-compose-`uname -s`-`uname -m` -o ./docker-compose && chmod +x ./docker-compose && ln -s /upgrade_agent/docker-compose /usr/local/bin/docker-compose

# upstream service port number
EXPOSE 1080

# default poll interval is 1 hour
ENV PATRAO_UPGRADE_INTERVAL_S 60

CMD ["./upgrade_agent", "--upstreamHost=http://192.168.1.58:1080"]