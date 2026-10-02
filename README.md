 **简体中文** | [English](README.en.md)

# mini-im

Go + Gin + GORM(PostgreSQL) + Redis 实现的即时通讯 / 客服系统。
前端是原生 Vue 2 + Axios 静态页，由后端直接托管，不需要单独构建。
同一套页面响应式适配 PC 与手机，无构建流程、无外部图片资源。

<img src="docs/images/pc-04-chat.png" width="820" alt="PC 端聊天界面">

---

## 目录

- [界面预览](#界面预览)
- [快速开始](#快速开始)
- [功能演示](#功能演示)
- [接口](#接口)
- [端到端回归](#端到端回归)
- [项目结构](#项目结构)
- [数据表](#数据表)
- [Redis key](#redis-key)
- [设计要点](#设计要点)
- [部署提示](#部署提示)
- [更新日志](#更新日志)
- [许可](#许可)

---

## 界面预览

> 下面的截图都取自 `config.yaml` 中 `seed` 段写入的演示账号（客服小哈 / 客服小吴 / 客服小朱 / 客服505），
> 把 `seed.enabled` 改成 `true` 重启即可得到完全一致的初始数据。

### PC 端

| 会话列表（含未读角标） | 聊天窗口 |
| --- | --- |
| <img src="docs/images/pc-03-sessions.png" width="460"> | <img src="docs/images/pc-04-chat.png" width="460"> |

| 搜索用户 / 发起好友申请 | 黑名单 |
| --- | --- |
| <img src="docs/images/pc-06-add-friend.png" width="460"> | <img src="docs/images/pc-08-blacklist.png" width="460"> |

### 手机端

| 登录 | 会话列表 | 聊天 | 设置 |
| --- | --- | --- | --- |
| <img src="docs/images/mob-01-login.png" width="180"> | <img src="docs/images/mob-02-sessions.png" width="180"> | <img src="docs/images/mob-03-chat.png" width="180"> | <img src="docs/images/mob-04-settings.png" width="180"> |

### 深色主题 / 多语言

| PC 端深色 | 手机端深色 | English |
| --- | --- | --- |
| <img src="docs/images/pc-11-dark.png" width="400"> | <img src="docs/images/mob-05-dark.png" width="170"> | <img src="docs/images/pc-13-english.png" width="400"> |

---

## 快速开始

```bash
docker compose up -d          # 启动 PostgreSQL + Redis（本机已有实例时跳过）
go run .                      # 启动服务
```

打开 <http://127.0.0.1:2580/>，会自动跳转到 `/web/index.html`；
未登录时前端会引导到 `/web/login.html`，注册或登录后进入聊天页。

新库是空的，直接在登录页注册账号即可（用户名 6-20 位字母、数字、下划线或减号，密码至少 6 位），
再通过左侧「新的朋友 → 添加好友」按账号或昵称搜索、发起申请，对方同意后即可互发消息。

### 演示账号（默认不创建）

`config.yaml` 里 `seed.enabled` 默认是 `false`。改成 `true` 后重启，会写入下面这组账号，
并**两两互加好友**，开箱即可互发消息，密码取 `seed.password`（默认 `123123`）：

| 账号 | 昵称 |
| --- | --- |
| `ha1234` | 客服小哈 |
| `wu1234` | 客服小吴 |
| `zhu333` | 客服小朱 |
| `for505` | 客服505 |

账号名与昵称都在 `config.yaml` 的 `seed.accounts` 里，按需增删即可：

```yaml
seed:
  enabled: false            # 改成 true 后重启生效
  password: "123123"
  accounts:
    - { username: ha1234,  nickname: 客服小哈 }
    - { username: wu1234,  nickname: 客服小吴 }
    - { username: zhu333,  nickname: 客服小朱 }
    - { username: for505,  nickname: 客服505 }
```

初始化是幂等的：账号已存在时只补齐昵称与头像，不覆盖密码，重复启动安全。
演示账号由 seed 直接写库，绕过了注册接口的校验，因此不受注册规则（用户名 6-20 位）限制。

### 前置条件

- Go 1.23+
- PostgreSQL（库需先存在，表由 GORM 自动迁移）
- Redis

不使用 `docker compose` 时，手工建库即可：

```sql
CREATE DATABASE "mini-im" ENCODING 'UTF8';
```

### 配置

`config.yaml` 默认为：PostgreSQL `127.0.0.1:5432` / `postgres` / `123456` / 库 `mini-im`，
Redis `127.0.0.1:6379` DB `1`，HTTP 端口 `2580`。
`app.name`（默认 `mini-im`）同时作为 Redis key 的前缀，改它等于换一套缓存命名空间。

除 `seed.*` 外的配置都能用环境变量覆盖，容器部署时无需改文件：

```bash
IM_PG_HOST=pg IM_PG_PASSWORD=secret IM_REDIS_ADDR=redis:6379 IM_JWT_SECRET=xxx go run .
```

<details>
<summary>全部环境变量</summary>

| 变量 | 说明 |
| --- | --- |
| `IM_APP_NAME` | 应用名，同时作为 Redis key 前缀 |
| `IM_HTTP_PORT` / `IM_RUN_MODE` / `IM_WEB_DIR` | 端口、模式、静态资源目录 |
| `IM_PG_HOST` `IM_PG_PORT` `IM_PG_USER` `IM_PG_PASSWORD` `IM_PG_DBNAME` `IM_PG_SSLMODE` | PostgreSQL 连接 |
| `IM_REDIS_ADDR` `IM_REDIS_PASSWORD` `IM_REDIS_DB` | Redis 连接 |
| `IM_JWT_SECRET` `IM_JWT_EXPIRE` | JWT 密钥与有效期（秒） |
| `IM_SNOWFLAKE_NODE_ID` | 雪花 ID 节点号（0~1023，默认 -1 按本机地址自动推导），只影响用户与好友申请 |

`seed.*` 没有环境变量覆盖，只有配置文件能开关。

</details>

生产环境请务必通过 `IM_JWT_SECRET` 注入密钥，不要沿用默认值。

`snowflake.node_id` 默认 `-1`，即按本机网卡地址自动推导节点号；**多实例部署时必须为每个实例
显式配置不同值（0~1023）**，否则两个实例在同一毫秒可能生成相同的用户 / 好友申请 ID。
消息主键是 UUID v7，不含节点号，多实例无需协调。

---

## 功能演示

下面按一次完整的使用流程串起各个功能，截图均为实际运行时的界面（PC 端 1280×800、手机端 390×844）。

### 1. 注册 / 登录

| 登录 | 注册 |
| --- | --- |
| <img src="docs/images/pc-01-login.png" width="440"> | <img src="docs/images/pc-02-register.png" width="440"> |

登录与注册在同一个页面切换：用户名 6-20 位（字母、数字、下划线或减号），密码至少 6 位，
昵称选填、留空时回退为账号名。勾选「请记住我」会把登录态写进 `localStorage`，
下次打开直接进入聊天页；token 剩余有效期不足 12 小时时，前端会自动续签。

同一账号在别处登录后，旧连接的 token 立即失效，前端会弹出
「账号在别处登录，尝试重新登录」并自动清理本地态跳回登录页（单点登录，见[设计要点](#设计要点)）。

### 2. 搜索并添加好友

点左上角 `+` 或「新的朋友」里的按钮打开「添加好友」，按账号 / 昵称模糊搜索，
结果里会标注你与对方的关系（可添加 / 申请中 / 待你处理 / 已是好友 / 已拉黑），
可附带一句 100 字以内的验证消息。

<img src="docs/images/pc-06-add-friend.png" width="820">

### 3. 新的朋友（处理申请）

对方申请后会实时收到 `friend_request` 事件，左侧「新的朋友」出现角标，
列表里直接「同意 / 拒绝」；同意后双方即时成为好友，申请方会收到 `friend_accepted` 事件。

<img src="docs/images/pc-05-friend-requests.png" width="820">

### 4. 聊天

聊天窗口支持文本、emoji 与图片表情（图片以 HTML 片段作为消息内容，`content_type=image`）。
对方不在线时消息进入离线队列，上线后按原顺序补发；
未读数由服务端的已读游标算出，打开会话即标记已读。

| PC 端聊天 | 手机端聊天 |
| --- | --- |
| <img src="docs/images/pc-04-chat.png" width="460"> | <img src="docs/images/mob-03-chat.png" width="190"> |

发送前服务端会校验**双方仍互为好友**、**接收方没有拉黑我**，任一不通过都不转发也不落库，
只给发送方回一个事件并撤回本地乐观渲染的那条消息：

<img src="docs/images/pc-12-blacklist-blocked.png" width="560">

| event | 含义 |
| --- | --- |
| `blacklisted` | 对方已把我拉入黑名单（上图） |
| `unfriended` | 对方已把我从好友中删除 |
| `not_friend` | 双方本就不是好友，或我已删除对方 |

### 5. 好友管理：删除 / 拉黑

鼠标悬停好友条目会出现「拉黑」「删除」两个操作，都会二次确认后才执行。
删除好友是**单向软删除**：只清我这一侧的列表与聊天记录，对方不受影响，但双方都发不出消息；
再次添加是全新会话。

| 好友列表（悬停出现操作） | 删除好友确认 |
| --- | --- |
| <img src="docs/images/pc-07-friends.png" width="460"> | <img src="docs/images/pc-09-confirm.png" width="460"> |

拉黑同样是单向的：被拉黑的人从我的好友列表与会话列表中消失（好友关系与聊天记录都不删除），
只拦截「被拉黑方 → 拉黑方」方向的消息，在「黑名单」页可以随时移出、移出后自动恢复。

<img src="docs/images/pc-08-blacklist.png" width="560">

### 6. 设置：主题 / 字号 / 字重 / 语言

右上角齿轮打开设置抽屉，改动实时生效并落到 `localStorage`。

| 浅色 | 深色 |
| --- | --- |
| <img src="docs/images/pc-10-settings.png" width="460"> | <img src="docs/images/pc-11-dark.png" width="460"> |

### 7. 手机端

手机端是同一套页面的响应式布局（断点 `991px`）：左侧导航收成窄栏，
点开会话后聊天面板全屏滑入，顶部出现返回箭头；输入区是微信风格的单行输入条，
输入框字号 ≥16px 以避免 iOS 聚焦时页面缩放，并监听 `visualViewport` 做键盘避让。

| 登录 | 会话列表 | 聊天（全屏） | 设置 |
| --- | --- | --- | --- |
| <img src="docs/images/mob-01-login.png" width="180"> | <img src="docs/images/mob-02-sessions.png" width="180"> | <img src="docs/images/mob-03-chat.png" width="180"> | <img src="docs/images/mob-04-settings.png" width="180"> |

---

## 接口

鉴权：请求头 `Authorization: Bearer <token>`。

| 方法 | 路径 | 鉴权 | 说明 |
| --- | --- | --- | --- |
| POST | `/api/register` | 否 | 注册新账号（用户名 6-20 位字母、数字、下划线或减号，密码 ≥ 6 位） |
| POST | `/api/login` | 否 | 登录，返回 `user` 与 `tokens` |
| PUT | `/api/token` | 是 | 续签 token |
| POST | `/api/logout` | 是 | 退出登录 |
| GET | `/api/friends` | 是 | 好友列表（不含我拉黑的人） |
| DELETE | `/api/friends/:FID` | 是 | 删除好友（单向：我的列表与会话清空，对方不受影响） |
| GET | `/api/messages` | 是 | 会话列表，每个对象只取最新一条，带未读数（不含我拉黑的人） |
| GET | `/api/chats/:FID` | 是 | 与指定好友的聊天记录 |
| POST | `/api/read/:FID` | 是 | 把与指定好友的会话标记为已读 |
| GET | `/api/users?keyword=` | 是 | 按账号 / 昵称搜索用户，带与我的关系 |
| POST | `/api/friend-requests` | 是 | 发起好友申请（`{to_id, message}`） |
| GET | `/api/friend-requests` | 是 | 待我处理的好友申请 |
| POST | `/api/friend-requests/:ID/accept` | 是 | 同意好友申请 |
| POST | `/api/friend-requests/:ID/reject` | 是 | 拒绝好友申请 |
| GET | `/api/blacklist` | 是 | 我拉黑的用户列表（按拉黑时间倒序） |
| POST | `/api/blacklist` | 是 | 拉入黑名单（`{user_id}`，幂等） |
| DELETE | `/api/blacklist/:UID` | 是 | 移出黑名单 |
| GET | `/ws?token=xxx` | 是 | WebSocket 长连接，也兼容 Authorization 头 |
| GET | `/avatar?name=&label=` | 否 | 动态生成 SVG 头像 |
| GET | `/healthz` | 否 | 健康检查，返回在线连接数 |

### 响应格式

```json
{ "StatusCode": 0, "Message": "请求成功", "Data": {} }
```

| 码 | 含义 |
| --- | --- |
| `0` | 成功 |
| `-1` | 业务失败，原因见 `Message` |
| `-10` | 未登录 / 登录态失效 / 账号在别处登录，前端会自动清理本地态并跳登录页 |

HTTP 状态码固定 200（WebSocket 鉴权失败除外，返回 401），业务结果只看 `StatusCode`。
无数据时列表接口返回 `[]`，其余返回 `{}`，不会是 `null`。

### 登录示例

```http
POST /api/login
Content-Type: application/json

{ "username": "ha1234", "password": "123123" }
```

```json
{
  "StatusCode": 0,
  "Message": "登录成功",
  "Data": {
    "user": {
      "id": "364066726138089472",
      "username": "ha1234",
      "nickname": "客服小哈",
      "avatar": "/avatar?label=%E5%AE%A2%E6%9C%8D%E5%B0%8F%E5%93%88&name=ha1234"
    },
    "tokens": { "token": "eyJhbGciOi...", "exp": 1790843147 }
  }
}
```

> `im_user` / `im_friend_request` 的 ID 是雪花 ID（18~19 位），超出 JavaScript 安全整数范围（2^53），
> 因此 JSON 中一律按**字符串**输出；消息主键是 UUID v7，本身就是字符串。
> `create_time`、`msg_count` 等小整数仍是数字；`content_type` 是字符串标识（见下）。

聊天记录 `GET /api/chats/364066726138089472` 返回数组，元素在实体字段之外额外带两个只读字段：

```json
{
  "id": "01a0f691-d2f8-7a3d-9e2f-4b6a8c1d5e70",
  "content_type": "text",
  "content": "你好 👋",
  "user_id": "364066726138089472",
  "friend_id": "1",
  "create_time": 1790843147,
  "type": 1,
  "user": { "id": "1", "nickname": "客服小哈", "avatar": "/avatar?..." }
}
```

`type` 为 `1` 表示自己发出、`2` 表示对方发来，仅在接口返回时计算。
列表按 `create_time` 倒序，前端会 reverse 成正序渲染，单会话最多返回 200 条。

会话列表 `GET /api/messages` 的元素形如：

```json
{
  "time": 1790843147,
  "message": "你好 👋",
  "content_type": "text",
  "uid": "364066726138089472",
  "msg_count": 2,
  "user": { "id": "364066726138089472", "nickname": "客服小哈", "avatar": "/avatar?..." }
}
```

`uid` 是对方的 ID，`msg_count` 是服务端按已读游标算出的未读数，调用 `POST /api/read/:FID` 后清零。

### WebSocket

连接 `ws://127.0.0.1:2580/ws?token=<token>`，建立后会先收到一条
`content` 为 `socket服务连接成功` 的消息，随后依次补发积压的离线消息。

上行报文：

```json
{
  "sender": "",
  "recipient": "2",
  "content": "你好",
  "content_type": "text",
  "subjoin": { "avatar": "", "nickname": "" }
}
```

`sender`、`time`、`subjoin` 由服务端覆盖填写，客户端传什么都无效；
`content_type` 是字符串标识：`text` 文本、`image` 图片（内容为 HTML 片段），
未知值由服务端归一化为 `text`。
单条报文上限 8KB，超过会被断开连接。下行报文结构相同。

发送前服务端会依次校验：**双方仍互为好友**（任一方删除好友后消息都不再可达）、
**接收方没有拉黑我**；任一不通过都不转发、不落库，并给发送方回一个
`not_friend` / `unfriended` / `blacklisted` 事件。

此外服务端会主动推送**业务事件**，报文只有 `event` 与 `data` 两个字段，
前端收到后刷新对应列表即可（离线时丢弃，登录后拉接口依旧能看到）：

```json
{
  "event": "friend_request",
  "data": { "id": "364070295146860544", "nickname": "AddTest", "avatar": "/avatar?...", "message": "你好" }
}
```

| event | 触发时机 | data 携带 |
| --- | --- | --- |
| `friend_request` | 有人申请加我为好友 | 申请人资料 + `message` |
| `friend_accepted` | 我发出的申请被同意（已互为好友） | 同意者资料 |
| `not_friend` | 我发的消息因我已删除对方（或双方本就不是好友）而未送达 | 对方资料（消息已被服务端丢弃） |
| `unfriended` | 我发的消息因对方已把我从好友中删除而未送达 | 对方资料（消息已被服务端丢弃） |
| `blacklisted` | 我发的消息因被对方拉黑而未送达 | 拉黑方资料（消息已被服务端丢弃） |

---

## 端到端回归

`tools/e2e` 是一个只依赖被测接口的回归脚本（HTTP + WebSocket），
覆盖「注册 → 登录 → 搜索 → 加好友 → 在线消息 → 聊天记录 → 已读 → 离线消息 → 登录态」全链路，
以及拒绝重申请、互相申请直接成为好友、并发消息、超长消息、伪造 sender 等边界。

```bash
go run ./tools/e2e                  # 全部场景
go run ./tools/e2e -scenario=core   # 只跑主链路
go run ./tools/e2e -scenario=edge   # 只跑边界 / 并发
go run ./tools/e2e -addr=http://127.0.0.1:2580
```

脚本会真实注册固定账号（`carol_01` `bob_01` `dave_01` `erin_01` `frank_01` `nick0_01`，密码统一 `123456`）并写入消息，
断言依赖「账号此前不存在、彼此尚不是好友」，因此**请在干净库上运行**，重复运行前先清库：

```bash
psql -h 127.0.0.1 -U postgres -d "mini-im" -c \
  "TRUNCATE im_message, im_friend, im_friend_request, im_read_state, im_blacklist, im_contact_state, im_user;"
redis-cli -n 1 flushdb
```

全部通过时退出码为 `0`，任一项失败为 `1`，可直接接入 CI。
当前基线：**94 项断言全部通过**。

---

## 项目结构

```
main.go                 入口：初始化组件 → 启动服务 → 优雅退出
config.yaml             配置
docker-compose.yml      PostgreSQL + Redis
conf/                   配置加载、Redis key 生成
db/                     GORM 连接与自动迁移、Redis 客户端、演示数据
model/                  实体：im_user / im_friend / im_friend_request / im_message / im_read_state / im_blacklist / im_contact_state
dao/                    数据访问
service/                业务：鉴权、会话、已读、在线状态、消息队列
api/handler             HTTP 处理器
api/middleware          鉴权、跨域
ws/                     WebSocket 连接中心
task/                   后台消息落库任务
route/                  路由注册
pkg/                    统一响应、JWT、HTTP 工具、雪花 ID、UUID v7
web/                    前端静态页
docs/images/            README 截图
tools/                  维护脚本：e2e 回归、favicon 生成、Bootstrap CSS 精简
```

依赖方向单向向下：`api / ws / task → service → dao → db`，不反向引用。

### 前端资源

前端没有构建流程，页面只依赖 **Vue 2 + Axios**，没有引入 jQuery 和 Bootstrap JS
（Bootstrap 仅用其 CSS 的网格 / 按钮 / 表单部分）。

右上角齿轮打开的设置抽屉里，各项设置都存在 `localStorage`，并落到 `<html>` 的属性上：

| 设置 | key | 落地方式 |
| --- | --- | --- |
| 主题（浅色 / 跟随系统 / 深色） | `mini-im-theme` | `data-theme` |
| 语言（中文 / English） | `mini-im-lang` | `lang` |
| 字号（滑杆 10~30px） | `mini-im-font-size` | 直接写根字号内联样式（`rem` 基准），`data-font-size` 留记号 |
| 字重（常规 / 加粗） | `mini-im-font-weight` | `data-font-weight`，切 `--fw-base` / `--fw-medium` / `--fw-bold` |

CSS 的字号统一写成 `rem`（基准 14px，即默认字号与改造前像素级一致），
`common.css` 开头把 Bootstrap 里几个固定 `px` 的基础字号（`body` / `.btn` / `.form-control`）
也改写为 `rem`，因此切字号只缩放文字，间距与图片尺寸不变；
字重则通过 `--fw-*` 变量控制，加粗档位下标题层级同步上浮。

设置抽屉本身刻意保留固定 `px` 字号：面板就是「调字号」的地方，拖动滑杆时它自己不能跟着缩放，
否则滑杆会上下跳动；面板外的聊天界面是实时预览，拖到哪文字就变到哪。

`bootstrap.min.css` 虽是完整包，但已按「当前页面实际用到的 class」裁剪到约 12KB
（原始 121KB）。裁剪脚本在 `tools/`，重新下载完整版 Bootstrap 后跑一次即可：

```bash
node tools/trim-bootstrap.cjs            # 空跑，打印统计与 5 项校验
node tools/trim-bootstrap.cjs --write    # 校验通过后写入
```

裁剪按选择器分支进行：分支里含有当前 DOM 不存在的 class 才会被删除，
且 `@media` 递归、`@keyframes` 整块保留，因此对现有页面渲染等价。

导航栏与设置面板的图标用内联 SVG 实现，不再需要 Glyphicons 字体；
`web/favicon.ico` 由 `go run ./tools/genfavicon.go` 生成。

### 数据表

| 表 | 字段 |
| --- | --- |
| `im_user` | `id`(雪花) `username`(唯一) `password`(bcrypt) `nickname` `avatar` `create_time` `update_time` |
| `im_friend` | `id`(自增) `user_id` `friend_id` `add_time`，两列均有索引（删除好友不删行） |
| `im_friend_request` | `id`(雪花) `from_id` `to_id` `message` `status` `create_time` `update_time` |
| `im_message` | `id`(UUID v7) `content_type`(字符串：`text` / `image`) `content` `user_id` `friend_id` `create_time`（只增不改） |
| `im_contact_state` | `id`(自增) `user_id` `peer_id` `removed` `clear_time`，`(user_id, peer_id)` 唯一 |
| `im_read_state` | `id`(自增) `user_id` `peer_id` `read_time` `update_time`，`(user_id, peer_id)` 唯一 |
| `im_blacklist` | `id`(自增) `user_id` `blocked_id` `add_time`，`(user_id, blocked_id)` 唯一，`blocked_id` 建立索引 |

`im_friend_request` 上 `(from_id, to_id)` 唯一、`(to_id, status)` 建立索引：
同一对用户只保留一行，重复申请走 upsert；`status` 为 `1` 待处理、`2` 已同意、`3` 已拒绝。

`im_blacklist` 是**单向**关系：`user_id` 拉黑了 `blocked_id`。拉黑后：

- 被拉黑的人从我的**好友列表**与**会话列表**中消失，只保留在黑名单列表里
  （好友关系与聊天记录都不删除，只是列表接口按黑名单过滤，移出后自动恢复）；
- 只拦截「被拉黑方 → 拉黑方」这个方向的消息投递，发送方会收到 `blacklisted` 事件；
- 两人之间待处理的好友申请会被清理，被拉黑方无法再向我发起好友申请。

`im_contact_state` 记录「我这一侧」的联系人状态，每个 `(user_id, peer_id)` 一行：

- `removed`：我已把对方从好友中移除（我的好友列表 / 会话列表过滤掉他，且双方都发不出消息）；
- `clear_time`：我清空与该用户会话的时间点，早于它的消息不再出现在我的会话列表 / 聊天记录里
  （消息行本身不动，对方仍看得到自己的记录）。

删除好友、重新加回来、清空会话都只写这张表，`im_friend` 行始终保留。

时间字段统一为 Unix 秒。好友关系按 `(user_id = ? OR friend_id = ?)` 双向匹配，
一条记录即可表达双向关系；「谁删了谁」这类单侧状态在 `im_contact_state`。

### Redis key

前缀取 `app.name`（`config.yaml` 中为 `mini-im`），生成逻辑集中在 `conf/config.go`。

| Key | 类型 | 用途 |
| --- | --- | --- |
| `<prefix>:auth_token:<uid>` | string | 当前有效 token，TTL = `jwt.expire`，用于单点登录 |
| `<prefix>:online_users` | set | 在线用户 ID |
| `<prefix>:msg_queue` | list | 待落库消息，`LPUSH` 入队 / `BRPOP` 出队 |
| `<prefix>:msg_queue_error` | list | 落库失败的消息，便于补偿 |
| `<prefix>:offline:<uid>` | list | 离线消息，`LPUSH` + `RPOP` 保证原顺序 |

---

## 设计要点

**主键选型**
`im_user` / `im_friend_request` 用雪花 ID（`pkg/snowflake`，41 位毫秒时间戳 + 10 位节点号 + 12 位序列号，
起始时间 2024-01-01 UTC）；`im_message` 用 UUID v7（`pkg/uuidv7`，48 位毫秒时间戳 + 4 位版本号 +
12 位序列号 + 2 位变体 + 62 位随机数，遵循 RFC 9562）。
两类主键都由应用侧生成（雪花字段标注 `autoIncrement:false`，UUID 字段是字符串主键），
插入前由模型的 `BeforeCreate` 钩子补齐，
因此 ID 在落库前就已确定、且整体按时间递增（B 树索引表现为顺序追加）；
发生时钟回拨时一律沿用上次时间戳继续自增，保证同一进程内不重复、不回退，单机每毫秒最多产出 4096 个。
对外一律以字符串形式序列化：雪花 ID 避免前端 JS 大整数精度丢失（见上文登录示例说明），UUID 本身就是字符串。

消息是三张表里增长最快、且只增不改的一张，主键换用 UUID v7 后不再需要节点号协调
（`snowflake.node_id` 只影响用户与好友申请），多实例部署可直接水平扩容，也不存在 64 位容量问题；
代价是主键宽一倍（16 字节）。同毫秒内序列号递增（而非纯随机）保留了「按 ID 排序即按时间排序」的性质，
`im_message` 的排序索引因此不会因随机主键而频繁页分裂。
`im_friend` / `im_read_state` 仍是数据库自增，无需全局唯一。

**加好友**
搜索（`username` / `nickname` 模糊匹配，排除自己，`ILIKE` 且转义 `%` `_`）→ 发起申请 →
对方同意后写入一条 `im_friend` 记录（一条记录即双向关系）。
同一对用户只有一行申请，重复申请走 `ON CONFLICT` upsert 并重新置为待处理，验证消息上限 100 字；
若「对方已申请我」时我再申请，则等价于双方互相确认，直接成为好友。
申请与同意都会通过 `service.Notify` 推给在线用户，推送实现（`*ws.Hub`）在 `main` 中注入：
`service.SetNotifier`，因此业务层不必反向依赖 ws 包。
被拒绝的申请不会推送事件，申请方需要自行拉接口或再次申请。

**删好友（单向，不删数据）**
`DELETE /api/friends/:FID` 不动 `im_friend` 行，只在 `im_contact_state` 写入
`(user_id = 我, peer_id = 对方, removed = true, clear_time = 现在)`：

- 我这一侧：好友列表 / 会话列表不再有对方，`GET /api/chats/:FID` 只返回 `clear_time`
  之后的消息（此刻为空，所以「重新加回来是空的」）；
- 对方那一侧：好友列表、会话、聊天记录都不受影响，只是双方都发不出消息
  （`IsFriend` 双向校验），对方发送时会收到 `unfriended` 提示；
- 重新加为好友时清掉双方的 `removed`，但保留各自的 `clear_time`：
  删过好友的一方是全新空会话，另一方原本的记录不会被翻出来；
- 和拉黑一样，删除与清空都只影响「我这一侧」，消息行本身始终保留。

**单点登录**
登录签发 token 后写入 `<prefix>:auth_token:<uid>` 覆盖旧值，校验时比对 Redis 中的值。
旧 token 再访问返回 `-10` 与「账号在别处登录」，前端自动跳登录页。
同一账号的 WebSocket 长连接不受影响（建连时校验一次）。

**消息链路**
WebSocket 收到上行消息后，服务端补全发送者资料 → 扇出给接收方的全部连接 →
对端无活跃连接则写入离线队列 → 同时投递到待落库队列，由后台任务阻塞消费写入 PostgreSQL。
落库是异步的，正常毫秒级完成；失败的消息转入错误队列而非丢弃。
`recipient` 非法或等于自己时直接丢弃该报文。

**未读数**
`im_read_state` 记录每个会话的已读游标（`read_time`），未读数 = 对方发来且晚于游标的消息条数，
`POST /api/read/:FID` 推进游标（`GREATEST` 保证并发下只前进不后退）。
之所以记游标而不是给每条消息打已读标记：消息是异步落库的，用游标可以避免
「标记已读发生在消息落库之前」而产生的假未读。游标缺失时回退到「我在该会话中最后一次发言时间」，
避免把早已看过的历史消息全算成未读。

**在线状态**
Hub 维护「用户 → 连接集合」，同一账号可多端在线。首个连接建立时标记在线，
最后一个连接断开才标记离线。

**前端布局**
桌面端与移动端是同一套 DOM，靠 `@media (max-width: 991px)` 切换：
桌面端左侧固定 650px 高的卡片布局，移动端改成 `100svh` 纵向 flex，
左侧导航收成 56px 窄栏、聊天面板绝对定位全屏滑入（`.app-mobile-chat`）。
移动端 `.right-box` 必须显式 `height: auto` 复位桌面端的 `--Height`，
否则 `top/bottom` 同时生效时 `bottom` 被忽略，聊天面板只占 650px、底部会露出导航栏。

**连接保活**
服务端 54s 发一次 ping，客户端 60s 内无响应即判定断线；单次写超时 10s，
发送缓冲写满时断开该连接而不阻塞其他用户。

**头像**
`/avatar?name=&label=` 按用户名哈希出固定配色，取 `label` 首字符渲染 SVG 并缓存 24h，
因此服务不依赖任何图片资源，注册账号的头像即指向该接口。

**emoji**
PostgreSQL 原生支持四字节 UTF-8，直接明文存储。旧 MySQL 实现里的 `[\uXXXX]`
转义与还原逻辑已整体移除。

---

## 部署提示

- 关闭 debug：`IM_RUN_MODE=release`，同时 gorm 日志降为 warn 级别
- 反向代理需要透传 WebSocket 升级头 `Upgrade` / `Connection`，且不要设置过短的读超时
- 应用监听 `SIGINT` / `SIGTERM`，退出时等待 10s 让在途请求结束并关闭数据库连接

---

## 更新日志

版本历史见 [CHANGELOG.md](CHANGELOG.md)，遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## 许可

[MIT](LICENSE) © 2026 kite88
