To Set SCTPNO delay option 
In the following file in go folder 
/home/aether/go/pkg/mod/git.cs.nctu.edu.tw/calee/sctp@v1.1.0/sctp_linux.go

before 
_, err = SCTPConnect(sock, raddr)
add 
// err = syscall.SetsockoptInt(sock, SOL_SCTP, SCTP_NODELAY, 1)
this line