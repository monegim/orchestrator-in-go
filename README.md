# orchestrator-in-go
Build an Orchestrator in Go(From Scratch)

## Docker API
### List containers 
```shell
curl --unix-socket /var/run/docker.sock http://localhost/v1.43/containers/json
```
### Stop a container
```shell
curl --unix-socket /var/run/docker.sock http://localhost/v1.43/containers/b169f7a8780b/stop -XPOST
```
### rm a container
```shell
curl --unix-socket /var/run/docker.sock http://localhost/v1.43/containers/b169f7a8780b -XDELETE
```

