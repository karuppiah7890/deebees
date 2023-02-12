FROM golang:1.20 AS builder
WORKDIR /build
COPY . /build
RUN make

FROM alpine:latest
COPY --from=builder /build/counter-machine /
CMD ["/counter-machine"]
