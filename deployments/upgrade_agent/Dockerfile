FROM centos:latest


WORKDIR /upgrade_agent

COPY . ./upgrade_agent
RUN curl -L https://github.com/docker/compose/releases/download/1.23.2/docker-compose-`uname -s`-`uname -m` -o ./upgrade_agent/docker-compose
RUN chmod +x ./upgrade_agent/docker-compose
# temporary value
EXPOSE 1533

# default poll interval is 30m
ENV UPGRADE_AGENT_POLL_INTERVAL_S 1800

CMD ["./upgrade_agent/upgrade_agent"]