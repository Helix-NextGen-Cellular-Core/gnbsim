# Builds the Helix image ghcr.io/helix-nextgen-cellular-core/gnbsim-realtrace:v1.0.0
# the same way it was built and verified:
#   docker build -t gnbsim-realtrace:local .
FROM golang:1.18.0-stretch AS builder
WORKDIR /go/src/gnbsim
COPY . .
RUN go mod vendor \
 && mkdir -p /out/gnbsim/bin && go build -mod=vendor -o /out/gnbsim/bin/gnbsim .

FROM alpine:3.16
LABEL org.opencontainers.image.source="https://github.com/Helix-NextGen-Cellular-Core/gnbsim" \
      org.opencontainers.image.licenses="Apache-2.0"
RUN apk update && apk add -U vim strace net-tools curl netcat-openbsd bind-tools gcompat bash tcpdump
COPY --from=builder /out/ /
WORKDIR /gnbsim/bin
