# Templates

## Reference

### feishu

- https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/using-markdown-tags#abc9b025

### weixin (bot) 企业微信群机器人

- https://developer.work.weixin.qq.com/document/path/91770#markdown%E7%B1%BB%E5%9E%8B

### weixinapp 企业微信应用

- https://developer.work.weixin.qq.com/document/path/90250#markdown%E6%B6%88%E6%81%AF

### dingtalk

Custom group robot (`oapi.dingtalk.com/robot/send`) markdown:

- https://open.dingtalk.com/document/orgapp/enterprise-internal-robots-send-markdown-messages
- https://open.dingtalk.com/document/development/message-types-and-data-format#title-afc-2nh-5kk
- Official syntax: headings `#`–`######`, quote `>`, `**bold**`, `*italic*`, links, images, lists
- Line break: `\n` with two spaces around it. A lone `\n` does not start a new paragraph.
- Empty lines next to `#` headings are collapsed by the renderer.
- Do not use `<br>` in custom-robot markdown: DingTalk mobile often hides the following section (commonly everything until the next `---`).
- Color: `<font color="#RRGGBB">text</font>` with double quotes (PC and mobile). Named colors like `red` are PC-only.
- Do not use Feishu `<text_tag>` or card `colorTokenV2` / `sizeToken` (those apply to interactive-card markdown, not custom robots)

### slack

- https://api.slack.com/reference/surfaces/formatting#basic-formatting
