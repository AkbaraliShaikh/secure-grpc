
# Secure-gRPC

[![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)](https://travis-ci.org/joemccann/dillinger)

Go-lang: Secure Bi-Directional GRPC streaming server & client example using TLS to find max number from the stream.

# How to run

Run Server
```sh
$ cd server
$ go run server.go
```
Run Client
``` sh
$ cd client
$ go run client.go
```

## Output

<a target="_blank" href="https://github.com/AkbaraliShaikh/AspNetCore2Docker/blob/master/img/sgrpc.png" class="rich-diff-level-one"><img src="https://github.com/AkbaraliShaikh/AspNetCore2Docker/blob/master/img/sgrpc.png" alt="text" width=85%  height=500px></a>

# Guide to generate proto and certs
#### Step 1: Proto Generation:

``` sh
$ cd proto
$ protoc --go_out=plugins=grpc:. *.proto
```
#### Step 2: Cert Generation:
```sh
$ git clone https://github.com/square/certstrap
$ cd certstrap
$ ./build
```

```sh
$ ./bin/certstrap-master-linux-amd64 init --common-name "akbar.com"
Created out/akbar.com.key
Created out/akbar.com.crt
Created out/akbar.com.crl
```
```sh
$ ./bin/certstrap-master-linux-amd64 request-cert --common-name client
$ ./bin/certstrap-master-linux-amd64 sign out/client --CA akbar.com 

$ ./bin/certstrap-master-linux-amd64 request-cert --common-name server
$ ./bin/certstrap-master-linux-amd64 sign out/server --CA akbar.com
```
The above files will be generated inside the `/out folder`

- `Copy akbar.com.crt, server.crt and server.key to secure-grpc/server/cert folder`
- `Copy akbar.com.crt, client.crt and client.key to secure-grpc/client/cert folder`

# Happy Coding!
