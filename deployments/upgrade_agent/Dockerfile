FROM alpine:3.9.2

WORKDIR /upgrade_agent
ADD bin/upgrade_agent .
RUN apk add ca-certificates && update-ca-certificates

# temporary value
EXPOSE 1533

# default poll interval is 30m
ENV UPGRADE_AGENT_POLL_INTERVAL_S 1800

CMD ["./upgrade_agent"]
