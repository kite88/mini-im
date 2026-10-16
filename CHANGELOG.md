# 更新日志

本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)，版本号形如 `v主版本.次版本.修订号`。

## [v1.1.1] - 2026-10-16

### 修复

- **首屏闪烁**：聊天页与登录页在 Vue 编译前会先把原始模板画出来，刷新时能看到 `{{ }}` 原文，以及表情盘 / 设置面板 / 弹窗遮罩闪一下。现在两端根节点都挂 `v-cloak` 整块隐藏（编译完成由 Vue 移除属性），编译期间用 `#boot-mask` 占位顶住画面（`mounted` 后移除，另有 6 秒兜底淡出，脚本加载失败也不会永久遮挡）；带 `v-show` 的浮层与登录页的注册表单补了初始 `display: none`。规则统一收在共用的 `common.css`，两个页面引用它的缓存串更新为 `?v=20261016`。

### 工程

- **README 截图全部重截**（PC 端 1280×800、登录 / 注册页 900×620、手机端 390×844）：截图取自真实运行的实例，演示账号与 `config.example.yaml` 的 `seed.accounts` 对齐（Mini-IM No. 1 / No. 2 / No. 3、Xiaoha、Xiaowu、Xiaozhu、GoForIt），界面标题也已换成 `Mini-IM`；README / README.en.md 中演示账号表、`seed.accounts` 片段与接口示例里的昵称同步更新。
- **`config.yaml` 不再入库**：`.gitignore` 新增 `/config.yaml` 与 `/config.yaml.bak`，仓库只带模板 `config.example.yaml`，本地配置由初始化配置向导写出或 `cp config.example.yaml config.yaml` 生成；此前入库的 `config.yaml` 带着开发用的 JWT 密钥与本地库名（如 `mini-im-2`），每次改配置都得在提交里挑拣。README / README.en.md 的快速开始、配置、项目结构与 Redis key 前缀说明已同步改为指向模板。

## [v1.1.0] - 2026-10-16

新增首次启动的初始化配置向导，发布包改为「不带配置文件、首次启动现场生成」的形态。

### 新增

- **初始化配置向导**：目录下没有 `config.yaml` 时不再用内置默认值硬跑，而是先起一个临时服务（`/setup`）引导填写 PostgreSQL / Redis：可真实点一次「测试连接」，保存时按 `config.example.yaml` 的注释与排版渲染写回（权限 0600，旧配置备份为 `.bak`），写回后业务服务自动接管同一端口；已配置后访问 `/setup` 是服务状态页，`SETUP_FORCE=1` 可重新配置。
- **监听失败的可操作提示**（`pkg/port`）：Windows 的 `winnat` 会随机把一段端口划为保留区间，落在区间里的端口即使没有任何进程占用也 bind 不了，而报错原文只是一句「权限不足」。现在这类错误会补上排查命令与换端口的做法，端口被占用的情况同理。

### 工程

- 配置文件缺失时返回 `ErrConfigNotFound` 并进入向导（此前是打印一句提示后使用内置默认值继续启动）；向导写回只替换表单覆盖到的标量，`app.name` / `app.web_dir` / `postgres.timezone` / `snowflake.node_id` / `seed.accounts` 按模板值保留，中文注释与排版原样不动。
- `config.example.yaml` 入库，兼作初始化向导写回配置时的底稿。
- 发布归档里**不再带 `config.yaml`**（避免把开发用的默认 JWT 密钥发出去），改为带 `config.example.yaml`；首次双击启动即进入向导。
- 归档名与包内顶层目录都带版本号（如 `mini-im-v1.1.0-windows-amd64.zip`），包内新增 `version.txt`，启动脚本启动时打印版本；`build.sh` 新增 `-v` 选项（省略时取最近的 git 标签，且必须是 `vX.Y.Z` 形式，避免把 `v1.0.0-3-gabc123` 打进归档名），打包前校验 `web/index.html`、`config.example.yaml`、`LICENSE` 齐全，缺 `zip` 时退回 Python `zipfile`。
- GitHub Actions 在上传前自检产物：归档名含标签、`version.txt` 与标签一致、包内有页面与配置模板、包内没有 `config.yaml`。
- 默认 HTTP 端口由 `2580` 改为 `8280`：配置文件与模板、初始化向导默认监听端口、`tools/e2e` 默认地址、启动脚本与文档里的示例地址同步更新。改端口前请先确认目标端口没有被系统保留（Windows 上 `netsh int ipv4 show excludedportrange protocol=tcp`），必要时用 `IM_HTTP_PORT` / `SETUP_PORT` 覆盖。

## [v1.0.0] - 2026-10-01

首个正式版本：一个可运行、可阅读的即时通讯 / 客服系统。

### 新增

- **注册 / 登录**：JWT + bcrypt；「请记住我」把登录态写进 `localStorage`，token 剩余有效期不足 12 小时时自动续签。
- **单点登录**：同一账号在别处登录后旧 token 立即失效，前端提示「账号在别处登录」并自动清理本地态跳回登录页。
- **好友关系**：按账号 / 昵称模糊搜索，发起 / 同意 / 拒绝好友申请；对方已申请我时我再申请等价于互相确认，直接成为好友。
- **实时聊天**：WebSocket 投递，支持文本、emoji 与图片表情（`content_type=image`）；对方不在线时进入离线队列，上线后按原顺序补发。
- **未读数**：服务端按已读游标算出，打开会话即标记已读。
- **好友管理**：删除好友是单向软删除（只清我这一侧的列表与记录），拉黑是单向拦截（仅拦截「被拉黑方 → 拉黑方」方向），可随时移出黑名单自动恢复。
- **业务事件**：`friend_request` / `friend_accepted` / `not_friend` / `unfriended` / `blacklisted` 主动推送。
- **聊天页面**：同一套 DOM 响应式适配 PC 与手机，亮色 / 跟随系统 / 深色三态、字号与字重可调、中英文切换，设置落 `localStorage`。
- **SVG 头像**：`GET /avatar?name=&label=` 按用户名哈希配色动态渲染并缓存，服务不依赖任何图片资源。
- **健康检查**：`GET /healthz` 返回当前在线连接数。

### 实现要点

- 主键选型：`im_user` / `im_friend_request` 用雪花 ID，`im_message` 用 UUID v7，均由应用侧在落库前生成且整体按时间递增；对外一律按字符串序列化，避免前端大整数精度丢失。
- 分层单向依赖：`api / ws / task → service → dao → db`；推送实现 `*ws.Hub` 在 `main` 中经 `service.SetNotifier` 注入，业务层不反向依赖 `ws` 包。
- 消息链路：收到上行报文 → 补全发送者资料 → 扇出给接收方全部连接 → 无活跃连接则写入离线队列 → 同时投递待落库队列，由后台任务 `BRPOP` 阻塞消费写入 PostgreSQL，失败转错误队列而非丢弃。
- 未读数用已读游标（`GREATEST` 保证并发下只前进不后退）而不是逐条打已读标记：消息异步落库，用游标可避免「标记已读早于消息落库」造成的假未读。
- 删除好友 / 拉黑 / 清空会话都只写 `im_contact_state`，`im_friend` 行始终保留。
- 发送前双向校验「仍互为好友」与「接收方未拉黑我」，不通过则既不转发也不落库，只回一个事件并让前端撤回乐观渲染的消息。
- 连接保活：服务端 54s 发一次 ping，客户端 60s 内无响应判定断线；单次写超时 10s，发送缓冲写满即断开该连接而不阻塞其他用户。

### 工程

- 前端是原生 Vue 2 + Axios 静态页，由后端直接托管，无构建流程、无外部图片资源。
- 配置集中在 `config.yaml`，除 `seed.*` 外都能用 `IM_*` 环境变量覆盖。
- 演示数据 `seed` 默认关闭，开启后幂等写入演示账号并两两互加好友，开箱即可互发消息。
- `docker-compose.yml` 提供 PostgreSQL + Redis；应用监听 `SIGINT` / `SIGTERM`，退出时等待 10s 让在途请求结束并关闭数据库连接。
- HTTP 不设写超时（WebSocket 长连接会被掐断），只设读超时与空闲超时。
- `bootstrap.min.css` 按页面实际用到的 class 从 121KB 裁剪到约 12KB，裁剪脚本在 `tools/`。
- `build.sh` 交叉编译 windows / linux / macOS（amd64 / arm64）共 6 个平台的归档（含 `config.yaml`、`web/` 与 `start.sh` / `start.bat` 启动脚本）并生成 `checksums.txt`；推送 `v*` 标签由 GitHub Actions 自动创建或更新 Release。
- Windows 启动脚本 `start.bat` 在进程退出后保留窗口并显示退出码与常见失败原因，不会闪退；脚本自身保持纯 ASCII，避免 `chcp 65001` 下 cmd 解析多字节字符出错。
- 入口包内嵌 IANA 时区数据库（`_ "time/tzdata"`）：没装 Go、系统也没有 tzdata 的机器（典型是 Windows）同样能解析 `timezone: Asia/Shanghai`，否则 pgx 建连时直接报 `unknown time zone`。
- Release 说明取自 `docs/release-notes/<tag>.md`，本文件为中文版本历史。

### 测试

- `tools/e2e` 是只依赖被测接口（HTTP + WebSocket）的端到端回归脚本，覆盖「注册 → 登录 → 搜索 → 加好友 → 在线消息 → 聊天记录 → 已读 → 离线消息 → 登录态」全链路，以及拒绝重申请、互相申请直接成为好友、并发消息、超长消息、伪造 `sender` 等边界。
- 当前基线：**94 项断言全部通过**；全部通过时退出码为 `0`，任一项失败为 `1`，可直接接入 CI。

### 已知局限

- 无 TLS：HTTPS / WSS 需由前置反向代理终结，并透传 `Upgrade` / `Connection` 升级头。
- 消息异步落库：正常情况下毫秒级完成，极端情况下接口返回与持久化之间存在极短时间差。
- 单会话最多返回 200 条历史消息，未提供向前翻页。
- 多实例部署必须为每个实例显式配置不同的 `snowflake.node_id`（0~1023），否则同毫秒内可能生成重复的用户 / 好友申请 ID。
- 在线状态由单进程内存中的连接中心维护，跨实例需共享 Redis 中的在线集合。

[v1.0.0]: https://github.com/kite88/mini-im/releases/tag/v1.0.0
