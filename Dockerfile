# build: docker build -t go-whatsapp-multi-device:1.0 . --build-arg GIT_COMMIT=$(git rev-parse HEAD) --build-arg VERSION=develop
# run: docker run --name go-whatsapp-multi-device --network host --rm -it --env-file ./.env go-whatsapp-multi-device:1.0
# in local, after you've built you image, run: `docker image prune --filter label=stage=gobuild`

# global arguments
ARG TZ_ARG="Asia/Jakarta"

# FYI: in some cases, we might got issue with missing/not found `GLIBC`, using `bookworm` image can fix this issue
#FROM debian:buster-slim AS base
#FROM debian:bullseye AS base
#FROM debian:buster AS base
FROM debian:bookworm AS base
LABEL maintainer="bellatrix Developer Team <bellatrix.developer@gmail.com>"
WORKDIR /app

FROM golang:1.20 as gobuild
LABEL stage=gobuild

# captures argument
ARG GIT_COMMIT
# e.g. latest, development, production
ARG VERSION=latest
ARG TZ_ARG

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GO111MODULE=on

RUN echo "Set ARG value of [GIT_COMMIT] as $GIT_COMMIT"
RUN echo "Set ARG value of [VERSION] as $VERSION"

WORKDIR /go/src/github.com/ardihikaru/go-whatsapp-multi-device

# get current commit and create build number
# echoing only the first 7 chars of the GIT_COMMIT
RUN echo "$VERSION -> ${GIT_COMMIT}" > /go/src/github.com/ardihikaru/go-whatsapp-multi-device/BUILD.txt

ADD go.mod go.sum ./
ADD cmd ./cmd
ADD internal ./internal
RUN mkdir ./data
RUN mkdir ./data/sqlitedb
RUN mkdir ./data/images
RUN mkdir ./data/images/qrcode

# copy docker source file
COPY files/docker-latest.tgz ./data/docker-latest.tgz

# download go modules
RUN go mod download

RUN go build -o api-service -ldflags "-X main.Version=$VERSION" \
    /go/src/github.com/ardihikaru/go-whatsapp-multi-device/cmd

FROM base AS release

ARG TZ_ARG

COPY --from=gobuild /go/src/github.com/ardihikaru/go-whatsapp-multi-device/api-service .
COPY --from=gobuild /go/src/github.com/ardihikaru/go-whatsapp-multi-device/data .
COPY --from=gobuild /go/src/github.com/ardihikaru/go-whatsapp-multi-device/BUILD.txt ./BUILD.txt

# CERT PACKAGES
RUN apt-get update
RUN apt-get install -y ca-certificates

RUN apt-get update && \
    apt-get install -yq tzdata && \
    ln -fs /usr/share/zoneinfo/Asia/Jakarta /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata
ENV TZ=${TZ_ARG}
#ENV TZ="Etc/GMT+7"

# builds its own docker to fix `GLIBC_2.xx` not found issue
# all sources: https://download.docker.com/linux/static/stable/x86_64/
#RUN  apt-get update \
#  && apt-get install -y wget \
#  && rm -rf /var/lib/apt/lists/*
#RUN wget http://get.docker.com/builds/Linux/x86_64/docker-latest.tgz
#RUN tar -xvzf docker-latest.tgz
RUN tar -xvzf /app/docker-latest.tgz
RUN mv docker/* /usr/bin/

# once installed, delete the file
RUN rm /app/docker-latest.tgz

EXPOSE 80 443
ENTRYPOINT ["/app/api-service"]
