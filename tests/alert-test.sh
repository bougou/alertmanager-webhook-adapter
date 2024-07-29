#!/bin/bash

dingtalk() {
  token="${DINGTALK_TOKEN}"
  channel_type="dingtalk"
  msg_type="markdown"
  payload=$(cat ./alert.json)

  echo "$payload" | curl -s -H "Content-Type: application/json" -v -XPOST "http://127.0.0.1:8090/webhook/send?channel_type=${channel_type}&token=${token}&msg_type=${msg_type}" -d @-
}

feishu() {
  token="${FEISHU_TOKEN}"
  channel_type="feishu"
  msg_type="markdown"

  payload=$(cat ./alert.json)

  echo "$payload" | curl -s -H "Content-Type: application/json" -v -XPOST "http://127.0.0.1:8090/webhook/send?channel_type=${channel_type}&token=${token}&msg_type=${msg_type}" -d @-

  # curl -X POST -H "Content-Type: application/json" -d '{"msg_type":"interactive","card":{"elements":[{"tag":"div","text":{"tag":"lark_md","content":"Hello"}}]}}'
}

slack() {
  token="${SLACK_APP_TOKEN}"
  channel_type="slack"
  msg_type="markdown"
  channel="jenkins-ci"

  payload=$(cat ./alert.json)

  echo "$payload" | curl -s -H "Content-Type: application/json" -v -XPOST "http://127.0.0.1:8090/webhook/send?channel_type=${channel_type}&token=${token}&channel=${channel}" -d @-

  ## Invite the slack app to the channel, then the slack app can send messages to this channel.
  # /invite @BOT_NAME
}

weixin() {
  token="${WEIXIN_TOKEN}"
  channel_type="weixin"
  msg_type="markdown"

  payload=$(cat ./alert.json)

  echo "$payload" | curl -s -H "Content-Type: application/json" -v -XPOST "http://127.0.0.1:8090/webhook/send?channel_type=${channel_type}&token=${token}&msg_type=${msg_type}" -d @-
}

weixinapp() {
  corpID="${WEIXIN_APP_CORP_ID}"
  agentID=${WEIXIN_APP_AGENT_ID}
  agentSecret="${WEIXIN_APP_SECRET}"

  toParty=2

  channel_type="weixinapp"
  msg_type="markdown"

  payload=$(cat ./alert.json)

  echo "$payload" | curl -s -H "Content-Type: application/json" -v -XPOST "http://127.0.0.1:8090/webhook/send?channel_type=${channel_type}&msg_type=${msg_type}&corp_id=${corpID}&agent_id=${agentID}&agent_secret=${agentSecret}&to_party=${toParty}" -d @-
}

discord-webhook() {
  id="${DISCORD_WEBHOOK_ID}"
  token="${DISCORD_WEBHOOK_TOKEN}"

  channel_type="discord-webhook"
  msg_type="markdown"
  payload=$(cat ./alert.json)

  echo "$payload" | curl -s -H "Content-Type: application/json" -v -XPOST "http://127.0.0.1:8090/webhook/send?channel_type=${channel_type}&msg_type=${msg_type}&id=${id}&token=${token}" -d @-
}


failed-test-1() {
  corpID="${WEIXIN_APP_CORP_ID}"
  agentID=${WEIXIN_APP_AGENT_ID}
  agentSecret="${WEIXIN_APP_SECRET}"

  toParty=2

  channel_type="notsupported"
  msg_type="markdown"

  payload=$(cat ./alert.json)

  echo "$payload" | curl -s -H "Content-Type: application/json" -v -XPOST "http://127.0.0.1:8090/webhook/send?channel_type=${channel_type}&msg_type=${msg_type}&corp_id=${corpID}&agent_id=${agentID}&agent_secret=${agentSecret}&to_party=${toParty}" -d @-
}

weixin_fail_msg_type() {
  token="${WEIXIN_TOKEN}"
  channel_type="weixin"
  msg_type="type-not-exist"

  payload=$(cat ./alert.json)

  echo "$payload" | curl -s -H "Content-Type: application/json" -v -XPOST "http://127.0.0.1:8090/webhook/send?channel_type=${channel_type}&token=${token}&msg_type=${msg_type}" -d @-
}
