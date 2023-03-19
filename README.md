

## Usage

### Start a redis server using docker

```sh
docker run -d -p 6379:6379 --name redis-server redis/redis-stack
```

### Verify redis server is running:

```sh
docker ps
```

### Install redis-cli 

On a Mac:

```sh
brew install redis
```

```sh
redis-cli
```

## Init golang project

How the project was initialized

* Change for your github username;

```sh
go mod init github.com/mitchallen/go-redis-gin
```

Get the redis package

```sh
go get -u github.com/redis/go-redis/v9
```

Get the gin package

```sh
go get -u github.com/gin-gonic/gin
```

How the main file was created (see actual source code):

```sh
touch server.go
```

## Run the program

The redis server must be running on port 6379:

```sh
docker ps 
```

```sh
go run server.go
```

Test with curl:

```sh
curl http://localhost:8080/lock/alpha/admin
{"resource":"alpha","userId":"admin","duration":"5s"}
```

### Check the results in redis-cli

```sh
% redis-cli
127.0.0.1:6379> GET id123
"{\"name\":\"Otto\",\"age\":45}"
```

## Cleanup

```sh
docker stop redis-server
docker rm redis-server
docker rmi redis/redis-stack
```

## References

* https://redis.io/docs/getting-started/installation/install-redis-on-mac-os/
* https://developer.redis.com/develop/golang/
* https://redis.io/docs/ui/cli/
* https://tutorialedge.net/golang/go-redis-tutorial/