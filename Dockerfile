# Builder
FROM golang as builder

ENV GOROOT /usr/local/go
ENV ETHTOOL "${GOPATH}/src/blbl.cr/ethtool"
ENV GO111MODULE on

WORKDIR "${ETHTOOL}"
COPY . "${ETHTOOL}"

RUN apt update && apt -y dist-upgrade
# build statically so that we only need go binary in the final image
RUN go test
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ethtool-exporter .

# Final app
FROM scratch

COPY --from=builder ${ETHTOOL}/ethtool-exporter /

CMD ["/ethtool-exporter", "-ifaceregexp", "'ens.*'", "-sleep", "1", "-output", "/tmp/ethtool"]
