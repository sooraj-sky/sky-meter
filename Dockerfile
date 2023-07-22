# syntax=docker/dockerfile:1
FROM golang:1.18-alpine AS builder
WORKDIR /go/src/app
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .
WORKDIR /go/src/app/cmd
RUN go build -o /main

FROM alpine

COPY --from=builder /main .
COPY --from=builder /go/src/app/templates .
ARG USER=skyuser
ENV HOME /home/$USER

# install sudo as root
RUN apk add --update sudo

RUN adduser -D $USER \
        && echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER \
        && chmod 0440 /etc/sudoers.d/$USER

USER $USER
WORKDIR $HOME
USER skyuser

EXPOSE 8080

CMD ["/main"]
