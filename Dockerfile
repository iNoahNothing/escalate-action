FROM golang:1.12.9-alpine

COPY main.go /usr/src/escalate.go

ARG GOBIN=/usr/local/bin

RUN go install /usr/src/escalate.go

ENTRYPOINT [ "escalate" ]