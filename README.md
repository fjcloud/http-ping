# HTTP Ping

## Deploy
```shell
oc new-project http-ping
oc new-app https://github.com/fjcloud/http-ping.git --strategy docker
```

## Expose
### HTTP
```shell
oc expose svc/http-ping
```
### HTTPS
```shell
oc create route edge --service http-ping
```

For minify version (under 800 bytes) go on http://<url>/min

This app purpose is to debug MTU Path issue between client and server
