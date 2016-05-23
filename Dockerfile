FROM alpine

WORKDIR /tmp

RUN apk add --update bash python py-pip && pip install awscli

ADD . /tmp

ENTRYPOINT ["./tmp/jon-backup-service"]
