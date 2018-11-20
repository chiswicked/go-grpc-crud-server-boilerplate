FROM golang:1.11.2-alpine3.8 AS build
ARG SERVICE
ARG GITHUB_TOKEN
RUN apk update && apk add git gcc make musl-dev
ADD . /go/src/github.com/${ORG}/${SERVICE}
WORKDIR /go/src/github.com/${ORG}/${SERVICE}
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"
RUN make clean install test build
RUN mv build/${SERVICE} /${SERVICE} 

FROM alpine:3.8
ARG SERVICE
ENV APP=${SERVICE}
RUN mkdir /app
COPY --from=build /${SERVICE} /app/${SERVICE}
EXPOSE 8080 8090
ENTRYPOINT exec /app/${APP}
