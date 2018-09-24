FROM golang:1.11-stretch as builder

RUN apt update && apt -y install make git

WORKDIR /app
COPY go.mod .
RUN go mod download
COPY .git ./.git
COPY cmd ./cmd
COPY pkg ./pkg
COPY Makefile ./Makefile
RUN make

FROM alpine:3.8
COPY --from=builder /app/cmd/maas-api/bin/maas-api /maas-api
CMD ["/maas-api"]
