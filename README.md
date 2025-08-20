# nginx-open-metrics-service

A Go server that reads the data from nginx and expose for prometheus or grafana alloy to collect.

It runs every 15s to collect the data from nginx.

## nginx configuration

Create the file `statistics.conf` with the following content:

```
❯ cat statistics.conf
server {
    listen 127.0.0.1:8080;
    location /api {
        stub_status;
        allow 127.0.0.1;
        deny all;
    }
}

```

Place the file on `/etc/nginx/conf.d` and reload the service.

```shell
❯ mv statistics.conf /etc/nginx/conf.d
❯ systemctl reload nginx
```

Check the interface is working with curl:

```shell
❯ curl http://localhost:8080/api
Active connections: 1 
server accepts handled requests
 219 219 515 
Reading: 0 Writing: 1 Waiting: 0 
```

## starting service

Start the service pointing to the same location you tested with curl:

```shell
❯ ./nginx-openmetrics --service=http://localhost:8080/api
🚚 fetching data from: http://localhost:8081/api
🎬 starting service at port: 9090
```

## testing

```shell
❯ curl localhost:8080/api
Active connections: 2 
server accepts handled requests
 221 221 521 
Reading: 0 Writing: 1 Waiting: 1 
❯ curl localhost:9090/metrics
# HELP active_connections The number of active connections
# TYPE active_connections gauge
active_connections 1
# HELP reading_connections The number of active reading connections
# TYPE reading_connections gauge
reading_connections 0
# HELP server_accepts_total The total number of server accepted connections
# TYPE server_accepts_total counter
server_accepts_total 221
# HELP server_handled_total The total number of server handled connections
# TYPE server_handled_total counter
server_handled_total 221
# HELP server_requests_total The total number of server requests
# TYPE server_requests_total counter
server_requests_total 521
# HELP waiting_connections The number of waiting connections
# TYPE waiting_connections gauge
waiting_connections 0
# HELP writing_connections The number of active writing connections
# TYPE writing_connections gauge
writing_connections 1
```

__Note:__ numbers might divert a bit since it updates every 15s.
