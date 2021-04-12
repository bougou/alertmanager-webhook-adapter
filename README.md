# alertmanager-webhook-adapter

A general webhook server for receiving [Prometheus AlertManager](https://prometheus.io/docs/alerting/latest/configuration/#webhook_config)'s notifications and send them through different channel types.

## Features

- Support Weixin Group Bot / 企业微信群机器人
    ```
    http(s)://<this-webhook-server-addr>/webhook/send?channel_type=weixin&token=<token>
    ```
- Support Dingtalk Group Bot / 钉钉群机器人
    ```
    http(s)://<this-webhook-server-addr>/webhook/send?channel_type=dingtalk&token=<token>
    ```
- Support Feishu Group Bot / 飞书群机器人
    ```
    http(s)://<this-webhook-server-addr>/webhook/send?channel_type=feishu&token=<token>
    ```
- Support Weixin Application / 企业微信应用
    ```
    http(s)://<this-webhook-server-addr>/webhook/send?channel_type=weixinapp&corp_id=<corp_id>&agent_id=<agent_id>&agent_secret=<agent_secret>
    ```

> More is comming...

## Run

### Build and Run

```bash
$ cd cmd/alertmanager-webhook-adapter
$ go build -v -x

$ ./alertmanager-webhook-adapter

# see help
$ ./alertmanager-webhook-adapter -h

# Add signature for sent messages
$ ./alertmanager-webhook-adapter --listen-address=:8060 --signature "Anything-You-Like"
# the signature will be added to the begining of the messsage:
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

## Configure Alertmanager to send alert messages to this webhook server

See **Features** section.
