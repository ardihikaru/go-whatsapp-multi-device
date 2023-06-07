# build: `docker build -t sea-cucumber-api-service:develop .`
# run: `docker run --name sea-cucumber-api-service --network host --rm -it sea-cucumber-api-service:develop`
# in local, after you've built you image, run: `docker image prune --filter label=stage=gobuild`
FROM golang:1.18 as gobuild
LABEL stage=gobuild
ARG VERSION=latest

WORKDIR /go/src/github.com/satumedis/go-template-api-service
ADD go.mod go.sum ./
ADD cmd ./cmd
ADD .env ./cmd
ADD internal ./internal
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build  \
    -o api-service -ldflags "-X main.Version=$VERSION" \
    github.com/satumedishub/go-template-api-service/cmd

FROM gcr.io/distroless/base

COPY --from=gobuild /go/src/github.com/satumedis/go-template-api-service/api-service /bin

ENTRYPOINT ["/bin/api-service"]
