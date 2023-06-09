# build: docker build -t go-whatsapp-multi-device:1.0 . --build-arg GIT_COMMIT=$(git rev-parse HEAD) --build-arg VERSION=develop
# run: docker run --name go-whatsapp-multi-device --network host --rm -it --env-file ./.env go-whatsapp-multi-device:1.0
# in local, after you've built you image, run: `docker image prune --filter label=stage=gobuild`
FROM debian:buster-slim AS base
LABEL maintainer="Muhammad Febrian Ardiansyah <mfardiansyah@outlook.com>"
WORKDIR /app

FROM golang:1.20 as gobuild
LABEL stage=gobuild

# captures argument
ARG GIT_COMMIT
# e.g. latest, development, production
ARG VERSION=latest

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
RUN mkdir ./data/qrcode
RUN mkdir ./data/sqlitedb
RUN go mod download

RUN go build -o api-service -ldflags "-X main.Version=$VERSION" \
    /go/src/github.com/ardihikaru/go-whatsapp-multi-device/cmd

FROM base AS release

COPY --from=gobuild /go/src/github.com/ardihikaru/go-whatsapp-multi-device/api-service .
COPY --from=gobuild /go/src/github.com/ardihikaru/go-whatsapp-multi-device/data .
COPY --from=gobuild /go/src/github.com/ardihikaru/go-whatsapp-multi-device/BUILD.txt ./BUILD.txt

EXPOSE 80 443
ENTRYPOINT ["./api-service"]
