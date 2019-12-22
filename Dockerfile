# docker build -t webvision .
FROM golang:1.13-buster as builder
COPY . /app
WORKDIR /app
RUN go mod vendor
RUN go build -o webvision .

FROM debian:buster
RUN apt-get update
RUN apt-get install -y ffmpeg

COPY --from=builder /app /app

RUN addgroup --gid 1000 webvision
RUN adduser --disabled-password --gecos "" --force-badname --ingroup webvision -u 1000 webvision
RUN mkdir /data && chown :1000 /data && chmod 775 /data && chmod g+s /data

USER webvision
VOLUME /data
WORKDIR /app

ENTRYPOINT /app/webvision -log info -config /data/config.yaml
