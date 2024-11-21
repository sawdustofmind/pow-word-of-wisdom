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

### 5 Docker output
```
> make up
> make logs
client-1  | {"L":"ERROR","T":"2024-11-21T12:47:48.733Z","M":"start client"}
server-1  | {"L":"INFO","T":"2024-11-21T12:47:48.637Z","M":"start listening","component":"server","address":"172.25.0.2:8888"}
client-1  | {"L":"INFO","T":"2024-11-21T12:47:48.758Z","M":"connected","service":"client","id":"1a0057e5-a18c-433a-9e5f-e028e682c59a","address":"server:8888"}
server-1  | {"L":"INFO","T":"2024-11-21T12:47:53.772Z","M":"received message","component":"client","id":"f04a1e53-1366-4fc0-ae44-759922fb2d75","address":"172.25.0.3:46372","message":"{\"type\":\"challenge\",\"challenge_counter\":0}"}
server-1  | {"L":"INFO","T":"2024-11-21T12:47:53.784Z","M":"write message to client","component":"client","id":"f04a1e53-1366-4fc0-ae44-759922fb2d75","address":"172.25.0.3:46372","message":"{\"type\":\"challenge\",\"data\":{\"challenge\":\"d6b57772d694eda650d9322c4c3202af22b07ff1571a6103032f559c54a274f1\"},\"ts\":1732193273783347597}\n"}
client-1  | {"L":"INFO","T":"2024-11-21T12:47:53.802Z","M":"Hashcash computed","service":"client","id":"1a0057e5-a18c-433a-9e5f-e028e682c59a","address":"server:8888","elapsed":"36.84575ms"}
client-1  | d6b57772d694eda650d9322c4c3202af22b07ff1571a6103032f559c54a274f1 8553 4
server-1  | {"L":"INFO","T":"2024-11-21T12:47:53.802Z","M":"received message","component":"client","id":"f04a1e53-1366-4fc0-ae44-759922fb2d75","address":"172.25.0.3:46372","message":"{\"type\":\"quote\",\"challenge_counter\":8553}"}
server-1  | d6b57772d694eda650d9322c4c3202af22b07ff1571a6103032f559c54a274f1 8553 4
server-1  | {"L":"INFO","T":"2024-11-21T12:47:53.804Z","M":"write message to client","component":"client","id":"f04a1e53-1366-4fc0-ae44-759922fb2d75","address":"172.25.0.3:46372","message":"{\"type\":\"quote\",\"data\":{\"quote\":\"A decision has value when it foresees its impact on the future and keeps in mind the learnings from the past and touches major subjects or matters of life.\"},\"ts\":1732193273804191055}\n"}
client-1  | {"L":"INFO","T":"2024-11-21T12:47:53.804Z","M":"quote result","service":"client","id":"1a0057e5-a18c-433a-9e5f-e028e682c59a","address":"server:8888","message":"A decision has value when it foresees its impact on the future and keeps in mind the learnings from the past and touches major subjects or matters of life."}
```