FROM golang:1.15.1-alpine3.12 AS build-env
ENV CGO_ENABLED 0
ENV COOS linux
ADD . /go/src/kubetest2
WORKDIR /go/src/kubetest2
RUN apk add --no-cache git && \
go build cmd/apiserver/main.go

# Final stage
FROM alpine:3.12
LABEL Author="Artur Kryukov <artur@kryukov.biz>"
EXPOSE 8080
WORKDIR /kubetest2
ADD templates templates
ADD .env ./
COPY --from=build-env /go/src/kubetest2/main /kubetest2
ENTRYPOINT /kubetest2/main