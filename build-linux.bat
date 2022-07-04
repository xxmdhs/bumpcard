SET CGO_ENABLED=1
SET GOOS=linux
SET GOARCH=amd64
set CC=zig cc -target x86_64-linux-musl
set CXX=zig c++ -target x86_64-linux-musl
go build -trimpath -ldflags "-w -s -linkmode \"external\" -extldflags \"-static\"" 
