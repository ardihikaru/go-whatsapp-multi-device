# build: `docker build -t go-whatsapp-multi-device:1.0 .`
# run: `docker run --name go-whatsapp-multi-device --network host --rm -it go-whatsapp-multi-device:1.0`
# in local, after you've built you image, run: `docker image prune --filter label=stage=gobuild`
FROM golang:1.18 as gobuild
LABEL stage=gobuild
ARG VERSION=latest

WORKDIR /go/src/github.com/ardihikaru/go-whatsapp-multi-device
ADD go.mod go.sum ./
ADD cmd ./cmd
ADD .env ./cmd
ADD internal ./internal
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build  \
    -o api-service -ldflags "-X main.Version=$VERSION" \
    github.com/ardihikaru/go-whatsapp-multi-device/cmd

FROM gcr.io/distroless/base

COPY --from=gobuild /go/src/github.com/ardihikaru/go-whatsapp-multi-device/api-service /bin

ENTRYPOINT ["/bin/api-service"]
