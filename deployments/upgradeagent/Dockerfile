FROM alpine:3.9.2 


WORKDIR /upgrade_agent

COPY . ./
RUN apk add ca-certificates && update-ca-certificates && apk add py-pip && pip install 'docker-compose==1.23'

# upstream service port number
EXPOSE 1080

# default poll interval is 1 hour
ENV PATRAO_UPGRADE_INTERVAL_S 60

CMD ["./upgradeagent", "--upstreamHost=http://192.168.1.51:1080"]