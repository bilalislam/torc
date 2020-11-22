# Torc

How to usage
---------------
1. cd $GOPATH
2. git clone https://git.hepsiburada.com/erebor/torc.git

```go
import "git.hepsiburada.com/erebor/torc"
```

go mod update version 
-------------------------
https://medium.com/faun/managing-dependency-and-module-versioning-using-go-modules-c7c6da00787a

Go package management
-------------------------
1. Enable go modules from setting go module
2. go mod init git.hepsiburada.com/erebor/torc
3. go mod tidy for optimize  all dependencies
Not: By the way this process must for dockerizing of this app

For example;
```sh
$ cd root folder
$ go mod init git.hepsiburada.com/erebor/torc
$ go mod tidy
```

Any more can sharing this package as library

Dockerizing looks like as followings 
---------------------------------
```docker
# Dockerfile Example
# https://medium.com/@petomalina/using-go-mod-download-to-speed-up-golang-docker-builds-707591336888
# Based on this image: https:/hub.docker.com/_/golang/
FROM golang:latest as builder

RUN mkdir -p /go/src/git.hepsiburada.com/erebor/basket/erebor.basket.consumers.basketupdated
WORKDIR /go/src/git.hepsiburada.com/erebor/basket/erebor.basket.consumers.basketupdated

RUN git config --global url."https://cem.basaranoglu:TzW_3JesxtbwjuyLD1KG@git.hepsiburada.com".insteadOf "https://git.hepsiburada.com"
# Force the go compiler to use modules
ENV GO111MODULE on
# <- COPY go.mod and go.sum files to the workspace
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .


# Compile application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o basket-updated

RUN chmod +x /go/src/git.hepsiburada.com/erebor/basket/erebor.basket.consumers.basketupdated

#Image Diff
#(Not Scratch) 1.23GB
#(Scratch    ) 34.3MB
# <- Second step to build minimal image
FROM scratch
WORKDIR /root/
COPY --from=builder /go/src/git.hepsiburada.com/erebor/basket/erebor.basket.consumers.basketupdated .
ENV ENV_FILE qa
# Execite application when container is started
CMD ["./basket-updated"]

EXPOSE 8080
```

Amqp wrapper package
--------------------------
let look at https://github.com/emretiryaki as github account

Generic repository for db
--------------------------
* basic implementation 
    * https://github.com/basho/taste-of-riak/blob/master/go/ch03/repositories/repository.go
    
* bridge patter implementation 
    * http://blog.ralch.com/tutorial/design-patterns/golang-bridge/
    
Http Client
---------------------------
````
client := p.Client.NewRequest()
request := client.Get("/path")
request.AppendHeader("header-key", "header-value")
var build = request.BuildRequest()
var response interface{}
err := build.Call(&response)
````
How to mocking in golang
---------------------------
1. https://github.com/vektra/mockery
2. all tests targeted %100 coverage but some cases could not be mocking 
3. how do it that ignore coverage for unnecessary classes in golang

How to using swagger in golang
---------------------------
1. https://github.com/swaggo/swag

How to loggin in golang
-------------------------
1. https://godoc.org/go.uber.org/zap
