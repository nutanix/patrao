FROM golang:1.12.4


WORKDIR /upgrade_agent

COPY . ./upgrade_agent

# temporary value
EXPOSE 1533

# default poll interval is 30m
ENV UPGRADE_AGENT_POLL_INTERVAL_S 1800

CMD ["./upgrade_agent/upgrade_agent", "--run-once"]