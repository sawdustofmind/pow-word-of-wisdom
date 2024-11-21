# "Word of Wisdom" TCP-server with protection from DDOS based on Proof of Work

## 1. Description

- Server replies with a random wisdom if PoW passed
- Communication between client and server features with JSON messages
- Client implemented 
- Challenge response implementation [wiki](https://en.wikipedia.org/wiki/Proof_of_work)

## 2 Getting started

### 2.1 Requirements

- Go 1.22+
- Docker 27.3+
- Docker compose 2.27+

### 2.2 Start server and client by docker-compose

```
make up
```

### 2.3 Launch tests

```
go test ./...
```

### 3 Benchmarks

| Name              | Count | Avg Time (ms/op) |             
|-------------------|-------|------------------|
| l = 32, zeros = 3 | 1033  | 	    0.984       |
| l = 32, zeros = 4 | 72    | 	  14.143        |
| l = 32, zeros = 5 | 6     | 	 203.555        |
| l = 32, zeros = 6 | 1     | 	2644.084        |
| l = 64, zeros = 3 | 1080  | 	   1.133        |
| l = 64, zeros = 4 | 79    | 	  14.556        |
| l = 64, zeros = 5 | 5     | 	 352.181        |
| l = 64, zeros = 6 | 1     | 	3158.835        |

### 4 Notes
- necessary to implement rate limiting by ip, and throttling for requests
- server relies on well configured firewall still, like cloudflare or iptables to avoid SYN-Flood and other attacks
- monitoring server with pprof and prometheus needed for server
- think of changing JSON to protobuff, messagepack
- client reconnect
- code needs some refactoring