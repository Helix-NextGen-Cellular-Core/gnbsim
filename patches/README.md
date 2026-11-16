# sctp-nodelay.patch

The Helix gnbsim image sets `SCTP_NODELAY` on its SCTP connections. The Go SCTP library used
(`git.cs.nctu.edu.tw/calee/sctp@v1.1.0`) has no option for it, so the vendored copy is patched
(`sctp_linux.go`, `dialSCTPExtConfig`). The Dockerfile applies it after `go mod vendor`; by hand:

```bash
go mod vendor
git apply patches/sctp-nodelay.patch
go build -mod=vendor -o bin/gnbsim .
```
