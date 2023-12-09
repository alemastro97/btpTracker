FROM golang:latest

RUN apt update && apt upgrade -y && apt install -y wamerican cron

ENV PROJECT_DIR=/go/src/backend \
    GO111MODULE=on \
    CGO_ENABLED=0

ADD . /go/src/backend
WORKDIR /go/src/backend
RUN mkdir "/build"
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon
COPY . .
EXPOSE 8080

# Prepare the cronjob
COPY ./scripts/cronjob /etc/cron.d/cronjob
RUN chmod 0644 /etc/cron.d/cronjob && crontab /etc/cron.d/cronjob

ENTRYPOINT ["bash", "scripts/init.sh"]
