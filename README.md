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

### Installing Prometheus and Grafana

1. Clone the repo
```
git clone git clone https://github.com/prometheus/prometheus.git
```

2. Install
```
cd prometheus
make build
```

3. Run the prometheus deamon
```
./prometheus --config.file=prometheus.yml
```

4. In the projects case that would be (running from inside the project directory)
```
../prometheus/prometheus --config.file=prometheus.yml
```

5. Install go client for prometheus
```
go get github.com/prometheus/client_golang/prometheus
```

6. Run Prometheus deamon
```
../../../../../my-projects/prometheus/prometheus --config.file=./.config/prometheus.yml  
```