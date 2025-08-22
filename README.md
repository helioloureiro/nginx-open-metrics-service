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
🚚 fetching data from: http://localhost:8080/api
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
# HELP nginx_active_connections The number of active connections
# TYPE nginx_active_connections gauge
nginx_active_connections 1
# HELP nginx_reading_connections The number of active reading connections
# TYPE nginx_reading_connections gauge
nginx_reading_connections 0
# HELP nginx_server_accepts_total The total number of server accepted connections
# TYPE nginx_server_accepts_total counter
nginx_server_accepts_total 221
# HELP nginx_server_handled_total The total number of server handled connections
# TYPE nginx_server_handled_total counter
nginx_server_handled_total 221
# HELP nginx_server_requests_total The total number of server requests
# TYPE nginx_server_requests_total counter
nginx_server_requests_total 521
# HELP nginx_waiting_connections The number of waiting connections
# TYPE nginx_waiting_connections gauge
nginx_waiting_connections 0
# HELP nginx_writing_connections The number of active writing connections
# TYPE nginx_writing_connections gauge
nginx_writing_connections 1
```

__Note:__ numbers might divert a bit since it updates every 15s.

## Build

Just have a Go compiler and make.

```shell
❯ make
go mod tidy
go mod vendor
cd nginx-openmetrics
go test -v ./...
?       nginx-openmetrics/v/nginx-openmetrics   [no test files]
cd nginx-openmetrics
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o nginx-openmetrics -modcacherw -ldflags="-w -X 'main.Version=$(git tag -l --sort taggerdate | tail -1)'" -buildmode=pie -tags netgo,osusergo -trimpath ./...
```

Binary will be available at `nginx-openmetrics/nginx-openmetrics`.

## systemd

There is a template for systemd service.

It is available on `nginx-openmetrics.service`.

You can change the parameters to fit your settings:

```systemd
[Service]
Restart=always
User=www-data
Group=www-data
ExecStart=/usr/sbin/nginx-openmetrics --service=http://localhost:8080/api
```

User `www-data` and group `www-data` are the default for debian based systems
(debian, ubuntu, mint, etc).

Be sure you placed the binary into `/usr/sbin`.

Move the `nginx-openmetrics.service` to `/etc/systemd/system` and reload systemd.
Then enable the service.

```shell
❯ mv nginx-openmetrics.service /etc/systemd/system
❯ sytemctl daemon-reload
❯ sytemctl enable --now nginx-openmetrics
```

## debian package

Now it is possible to generate a debian package to make easier to install.

You need to have `debuild` in place, which is parte of the package `devscripts`.

Then it is just matter to run:

```shell
❯ make debian
```
