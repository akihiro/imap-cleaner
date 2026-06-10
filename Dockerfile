FROM golang:1.26.4 AS build
WORKDIR /work
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
RUN make

FROM quay.io/prometheus/busybox:glibc
COPY --from=build --link ./work/imap-cleaner /
USER daemon
CMD ["/imap-cleaner"]
