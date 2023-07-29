# toll-calculator

### run kafka docker server

```
docker-compose up
```


### Installing gRPC and Protobuffer for Golang
1. CLI
```
brew install protobuf
```

2. protoc-gen-go pkg
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
```

3. protoc-gen-go-grpc pkg
```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

IMPORTANT: try to install latest if 1.3 is not working(was not working for me first time)
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

4. install the package dependencies
```
go get google.golang.org/protobuf
go get google.golang.org/grpc
```

5. Set up env vars
```
export GOPATH=$HOME/go
export PATH="$PATH:$(go env GOPATH)/bin"
```