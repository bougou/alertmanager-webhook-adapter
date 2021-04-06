# alertmanager-webhook-adapter

## Run

### Build and Run

```bash
$ cd cmd/alertmanager-webhook-adapter
$ go build -v -x

$ ./alertmanager-webhook-adapter

# see help
$ ./alertmanager-webhook-adapter --listen-address=:8060

# Add signature for messages
$ ./alertmanager-webhook-adapter --listen-address=:8060 --signature "Anything-You-Like"
# the signature normally will be added to the begining of the messsage:
# 【Anything-You-Like】this-is-the-the-the-the-the-xxxxxxxxxx-message

```

### Start as systemd service

```bash
# Install the binary alertmanager-webhook-adpater file to some directory
# like /usr/local/bin/alertmanager-webhook-adapater
# and chmod +x /usr/local/bin/alertmanager-webhook-adapater

$ cp deploy/alertmanager-webhook-adapter.service /etc/systemd/system/

# make sure the bin path be consistent
$ vim /usr/local/bin/alertmanager-webhook-adapater

$ systemctl daemon-reload
$ systemctl start
```
