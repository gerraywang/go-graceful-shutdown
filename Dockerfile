FROM golang:1.17-buster AS debug
WORKDIR /app
COPY . .
RUN go install github.com/cosmtrek/air@v1.27.3
CMD air

FROM golang:1.17-buster AS build
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates
WORKDIR /app
ARG UID
RUN adduser test -u ${UID:-1000} && chown -R test /app
USER test
ENV CGO_ENABLED=0
ENV GIN_MODE=release
COPY go.* ./
RUN go mod download
COPY . ./
RUN GOOS=linux go build -tags timetzdata -mod=readonly -v -o server

# k8s、cloud-run用ビルド --target deploy
FROM scratch AS deploy
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/server /server
CMD ["/server"]
