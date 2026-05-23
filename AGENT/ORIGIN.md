# urlAPI 原项目功能与结构分析

本文记录当前仓库的既有功能、模块边界、运行链路和迁移注意事项。后续重构时应以本文作为功能迁移清单，避免遗漏原有行为。

## 1. 项目定位

`urlAPI` 是一个 Go + Vue 的 API 服务项目，核心能力是把多个外部资源或模型接口包装成可通过 URL 调用的图片/文本服务，并提供后台管理页面。

主要功能：

- 文本生成：调用 OpenAI、DeepSeek、阿里巴巴兼容接口或自定义接口生成文本，再把文本绘制成 PNG 图片，通过 `/download` 返回。
- 图片生成：调用 OpenAI 或阿里巴巴通义万相生成图片，保存为本地临时 PNG，通过 `/download` 返回。
- 随机图片：从后台登记的 GitHub/Gitee 仓库图片列表中随机返回一个图片 URL。
- 网页缩略图：针对 Bilibili、YouTube、arXiv、IT之家、GitHub、Gitee 等页面生成信息卡片图片。
- 下载中转：读取内部生成的临时 PNG，作为附件返回。
- 后台管理：Vue + MDUI 页面，支持登录、任务查看、接口设置、功能设置、安全设置、仓库管理、API 工作台。
- 安全防护：Referer 白名单、IP 频率限制、功能开关、API 白名单、Prompt 通配符限制、任务记录例外过滤。

## 2. 技术栈

后端：

- Go module：`module urlAPI`。
- Web 框架：`github.com/gin-gonic/gin`。
- CORS：`github.com/gin-contrib/cors`，当前允许所有来源，允许 `GET`、`POST`，Header 包含 `Content-Type` 和 `Authorization`。
- 数据库：SQLite + GORM，数据库文件固定为 `assets/database.db`。
- 图片处理：`image/*`、`github.com/golang/freetype`、`github.com/nfnt/resize`。
- 静态资源嵌入：Go `embed`，嵌入 `static/dist/*`、字体、图标、Logo、默认设置和 `empty.png`。

前端：

- Vue 3 + Vite。
- UI：MDUI Web Components。
- 路由：Vue Router，所有后台页面在 `/dash` 下。
- Cookie：`js-cookie` 存储后台 token。
- 密码哈希：`js-sha256` 在前端对密码做 SHA-256 后提交。

部署与构建：

- 后端 Dockerfile 位于 `Dockerfile.d/amd64`、`Dockerfile.d/arm64`，基于 `golang:1.26.2-alpine` 构建，运行镜像基于 `alpine`。
- 前端独立 Dockerfile 位于 `Dockerfile.d/frontend-amd64`、`Dockerfile.d/frontend-arm64`，构建 Vite 产物后由 Nginx 托管。
- 前端独立部署通过 `BACKEND_URL` 环境变量生成 Nginx 反代配置。
- GitHub Actions：`.github/workflows/build.yml` 构建多平台二进制并发布 Release；`docker.yml`、`docker-frontend.yml` 构建并推送 Docker 镜像。

## 3. 顶层目录结构

- `main.go`：程序入口，解析命令行参数，启动 Gin 服务。
- `command/`：启动参数处理，包括端口、清空任务、清空登录态、重置密码、清除后台登录 IP 限制。
- `handler/`：Gin 路由、请求组装、安全检查调用、任务前后处理、响应返回。
- `request/`：把一次请求聚合为 `DB`、`Processor`、`Security` 三部分的中间结构。
- `processor/`：业务处理器，负责调用外部 API、生成图片、任务缓存、后台 session 操作。
- `security/`：通用安全检查、功能开关检查、API 白名单检查。
- `database/`：SQLite/GORM 模型、迁移、初始化、CRUD、内存缓存。
- `util/`：外部 API 调用、图片绘制、下载、通配符检查、IP 地区、随机字符串等工具。
- `file/`：嵌入资源，包括字体、图标、Logo、默认设置、空图片。
- `static/`：Vue 后台源码、构建产物 `dist`、独立前端 Nginx 配置与 Go 静态资源嵌入。
- `guide/`：README 中使用的示例图片。
- `AGENT/doc/`：另一个 Go 参考项目/文档工程，不属于当前 urlAPI 主服务运行链路。

## 4. 程序启动链路

入口在 `main.go`：

1. `command.Arg(os.Args)` 解析命令行参数。
2. `handler.Handler(command.Port)` 初始化并启动 Gin 服务。
3. `defer database.Disconnect()` 预期在服务退出后关闭数据库。

注意：`handler.Handler()` 内部调用 `r.Run()` 阻塞运行，因此 `database.Disconnect()` 只有服务退出时才可能执行。

数据库包有 `init()`：

1. 打印 `urlAPI` Figlet。
2. `connect()` 创建 `assets/` 并连接 `assets/database.db`。
3. `migration()` 自动迁移 `Setting`、`Task`、`Session`、`Repo` 四张表。
4. 初始化 `SettingMap`、`RepoMap`、`SessionMap` 内存缓存。
5. `settingInit()` 读取嵌入的 `file/setting.json`，补齐默认配置。

处理器包有 `init()`：

1. 删除 `assets/img/`。
2. 重新创建 `assets/img/`。

因此重启服务会清空临时生成图片目录，但不会清空 SQLite 中已有任务记录。

## 5. 启动参数

由 `command.Arg()` 处理：

- `port <端口>`：设置服务端口，默认 `2233`。
- `repwd`：把后台密码重置为 `123456` 的 SHA-256 值，并清空登录态。
- `clear`：删除并重建任务表。
- `logout`：删除并重建会话表，清空 `SessionMap`。
- `clear_ip_restriction`：把 `dashallowedip` 设置为 `["*"]`。

当前代码没有参数越界保护，`port` 后缺少端口时会访问 `args[index+1]`。

## 6. HTTP 路由

`handler.Handler()` 注册的路由：

- `GET /txt`：文本生成。
- `GET /img`：图片生成。
- `GET /rand`：随机图片。
- `GET /web`：网页缩略图。
- `GET /download`：内部图片下载/中转。
- `POST /session`：后台登录、配置、任务、仓库管理。
- `GET /assets/*`：前端静态资源。
- `NoRoute`：`/dash` 和 `/dash/*` 返回后台 `index.html`；其他路径 301 跳转到一个 Bilibili 视频 URL。

`/dash` 前端由后端内嵌的 `static/dist` 提供。若重构构建流程，需要确保 `static/dist/index.html` 和 `static/dist/assets/*` 在 Go 编译前存在，否则嵌入和运行会受影响。

## 7. 通用请求处理模型

公开 API 的典型流程如下：

1. `handler/*Handler` 创建 `request.Request`。
2. `requestBuilder` 从 Gin Context 读取 Query、Referer、IP、User-Agent、Host 等信息。
3. Builder 同时填充：
   - `r.Security.*`：安全检查所需数据。
   - `r.DB.Task`：任务记录字段。
   - `r.Processor.*`：业务处理器参数。
   - `r.Processor.Filter`：任务缓存键。
4. `checker()` 执行通用安全检查和具体功能/API 检查。
5. 如果 `General.Unsafe` 为 true，返回 `403 {"error": ...}`。
6. 对 `/txt`、`/img`、`/web`：先 `beforeTask()` 处理缓存和并发限制，再执行 processor，再 `afterTask()` 写缓存/写库。
7. 对 `/rand`：直接执行随机处理器并写库，不走任务缓存队列。
8. `returner()` 根据 `format=json` 返回 JSON，否则 `302` 跳转到结果 URL。

重要结构：

- `request.Request` 包含 `DB`、`Processor`、`Security` 三个聚合字段。
- `database.Task` 是任务记录核心模型。
- `processor.TaskQueueFilter` 用作缓存键，包含 `Type`、`Size`、`Target`、`API`。

## 8. 任务缓存与并发限制

`processor` 内维护两个全局结构：

- `TaskQueue map[TaskQueueFilter]TaskQueueItem`：按同类请求缓存上一次成功任务。
- `TaskCounter map[string]int`：按 API 标识统计当前运行任务数。

`beforeTask()` 行为：

- 根据任务类型映射到配置名：`txt.gen -> txt`、`img.gen -> img`、`web.img -> web`。
- 从配置中读取缓存过期时间：`txt[3]`、`img[2]`、`web[3]`，单位分钟。
- 如果同一个 `TaskQueueFilter` 已存在且正在运行，会每秒等待直到其完成。
- 如果已有成功任务未过期，则复用结果：当前任务状态和返回值复制旧任务，`Temp` 标记为 `Yes`，但 UUID 和 Time 使用当前请求。
- 如果已有任务过期或失败，会删除旧图片文件。
- 新任务会等待同一 API 的运行数不超过 2，然后标记 `Running=true` 并计数加一。

`afterTask()` 行为：

- 当前 API 计数减一。
- 若生成的 PNG 文件无效，且不是缓存命中、任务状态为 success，则改为 failed，返回 `download?img=empty`。
- 成功任务写入 `TaskQueue` 缓存。
- 如果 `SkipDB=false`，写入 `Task` 表。

迁移时必须保留的行为：

- 同请求去重/等待逻辑。
- 缓存过期时间来自配置数组的固定位置。
- 每个 API 同时最多约 3 个任务的限制逻辑。
- `Temp=Yes/No` 的任务记录语义。
- 图片有效性检查失败时返回空图。

## 9. 对外 API 行为

### 9.1 `GET /txt`

用途：文本生成并绘制成图片。

参数：

- `prompt`：必填。可为预置值 `laugh`、`poem`、`sentence`，也可为自定义提示词。
- `format`：选填。`json` 返回 JSON，否则 302 跳转图片下载 URL。
- `api`：选填。支持 `openai`、`alibaba`、`deepseek`、`otherapi`。
- `model`：选填。为空时使用所选 API 的默认文本生成模型。
- `more`：选填，写入任务的 `MoreInfo` 并可参与例外过滤。

处理细节：

- 预置 prompt 通过 `database.PromptMap` 映射到 `SettingMap["prompt"]`：`laugh=0`、`poem=1`、`sentence=2`。
- `api` 为空时使用 `SettingMap["txt"][1]`。
- `model` 为空时使用 `SettingMap[api][1]`。
- system prompt 使用 `SettingMap["context"][0]`。
- API endpoint 通过 `getEndpoint(api)` 获取。
- 生成文本后调用 `util.DrawTxt()` 绘制 PNG 到 `assets/img/<uuid>.png`。
- `data.Return` 存 JSON 字符串：`prompt`、`response`、`url`。
- `processor.Return` 是 `/download?img=<uuid>` 完整 URL。

### 9.2 `GET /img`

用途：图片生成。

参数：

- `prompt`：必填。
- `format`：选填。`json` 返回 JSON，否则 302 跳转。
- `api`：选填。支持 `openai`、`alibaba`。
- `model`：选填。为空时使用 API 的默认图片模型。
- `size`：选填。为空时使用 API 默认尺寸。
- `more`：选填，写入任务记录。

处理细节：

- `api` 为空时使用 `SettingMap["img"][1]`。
- `model` 为空时使用 `SettingMap[api][3]`。
- `size` 为空时使用 `SettingMap[api][4]`。
- Alibaba 使用固定接口 `https://dashscope.aliyuncs.com/api/v1/services/aigc/text2image/image-synthesis`，异步创建任务后轮询 `api/v1/tasks/<id>`。
- OpenAI 使用 `SettingMap["openai"][5]` 作为图片 endpoint。
- 成功后图片保存到 `assets/img/<uuid>.png`。
- `data.Return` 存 JSON 字符串：`original_prompt`、`actual_prompt`、`url`。

### 9.3 `GET /rand`

用途：从已登记仓库图片列表中随机返回图片 URL。

参数：

- `api`：选填。支持 `github`、`gitee`，为空时使用 `SettingMap["rand"][3]`。
- `user`：必填，仓库用户名/组织名。
- `repo`：必填，仓库名。
- `format`：选填。`json` 返回 JSON，否则 302 跳转。
- `more`：选填。

处理细节：

- `target` 组装为 `user/repo`。
- 从 `database.RepoMap[api+";"+target]` 获取图片 URL 列表。
- 随机选择一个 URL。
- 成功时 `data.Return` 为 `{"url":"..."}`。
- 当前请求不会走 `beforeTask/afterTask` 缓存队列，但会根据 `SkipDB` 决定是否写任务表。
- 仓库列表通过后台 `newRepo` 或 `refreshRepo` 从 GitHub/Gitee contents API 拉取。

### 9.4 `GET /web`

用途：为特定网站页面生成缩略图/信息卡片。

参数：

- `img`：必填，目标页面 URL。当前实现只读取 `img`，不是 `url`。
- `format`：选填。`json` 返回 JSON，否则 302 跳转。
- `more`：选填。

支持域名与处理器：

- `www.bilibili.com`：解析 BV/AV，调用 Bilibili view API，绘制视频卡片。
- `www.youtube.com`：解析视频 ID，使用 YouTube Data API，绘制视频卡片。
- `arxiv.org`：抓取 HTML，解析标题、作者、摘要，绘制文章卡片。
- `www.ithome.com`：抓取 HTML，使用文本总结 API 生成摘要，绘制文章卡片。
- `github.com`、`gitee.com`：调用仓库 API 获取名称、作者、描述、star/fork，绘制仓库卡片。

处理细节：

- Builder 会解析 `img` 的 Host 作为 `API`。
- `webimgallowed` 必须允许该 Host。
- `www.ithome.com` 额外要求文本功能启用。
- 成功后保存 PNG 到 `assets/img/<uuid>.png`，返回 `/download?img=<uuid>`。

### 9.5 `GET /download`

用途：读取内部生成图片并作为附件返回。

参数：

- `img`：必填。`empty` 读取嵌入的 `empty.png`；其他值读取 `assets/img/<img>.png`。

行为：

- 只执行 `InfoChecker()` 和 `ExceptionChecker()`，并强制 `SkipDB=true`。
- 成功时设置 `Content-Type: image/png` 和 `Content-Disposition: attachment; filename="download.png"`。
- 失败时 `ReturnError` 设置为 `SettingMap["web"][4]`，最终 302 跳转到该错误图片 URL。

### 9.6 `POST /session`

用途：后台会话、任务、配置、仓库管理。

请求：

- Header：`Authorization`。
- Body：JSON，绑定到 `processor.Session`。
- `operation` 决定操作类型。

支持 operation：

- `login`：登录。
- `logout`：删除当前 token。
- `exit`：如果当前 session 非长期登录，则删除 token。
- `newRepo`：新增仓库图片源。
- `refreshRepo`：刷新仓库图片源。
- `delRepo`：删除仓库图片源。
- `fetchRepo`：读取所有仓库图片源。
- `fetchTask`：查询任务。
- `fetchSetting`：读取配置分组。
- `editSetting`：保存配置分组。

鉴权行为：

- `login` 操作要求 `Authorization` 等于 `SettingMap["dash"][0]`，即 SHA-256 后的密码。
- 登录成功生成随机 token，写入 `Session` 表和 `SessionMap`。
- `login_term=true` 时 token 过期时间为 7 天，否则为 1 天。
- 非 `login` 操作要求 `Authorization` 对应的 token 存在且未过期。

响应：

- 成功返回完整 `processor.Session` JSON，其中可能包含 `session_token`、`setting_data`、`task_data`、`repo_data` 等。
- 失败返回 `400 {"error":"..."}`。

## 10. 安全模型

### 10.1 通用安全检查

`GeneralChecker()` 执行：

1. `FrequencyChecker()`。
2. `ExceptionChecker()`。
3. `InfoChecker()`。

频率限制：

- 以 `Type + IP` 为键。
- 0.25 秒内同类请求计数达到 10 后标记 unsafe。
- 错误信息：`<ip> accessed too frequently`。

任务记录例外：

- `taskexceptdomain`：Referer 域名通配符列表。
- `taskexceptinfo`：`General.Info` 通配符列表。
- 命中后设置 `SkipDB=true`，请求仍可继续处理。

基本信息检查：

- `Target` 为空则 unsafe，错误 `Empty Target`。
- 读取 `allowedref`，用 `util.GetDomain(Referer)` 提取域名，再用通配符匹配。
- Referer 为空或不匹配时 unsafe，错误 `Referer <referer> not allowed`。

迁移注意：当前下载接口也会要求 Referer 通过白名单，因为 `downloadChecker()` 调用了 `InfoChecker()`。

### 10.2 功能开关检查

- 文本生成：`SettingMap["txt"][0] == "true"` 才启用。
- 图片生成：`SettingMap["img"][0] == "true"` 才启用。
- 随机图片：`SettingMap["rand"][0] == "true"` 才启用。
- 网页缩略图：`SettingMap["web"][1] == "true"` 才启用。

### 10.3 API 白名单检查

- 文本：`openai`、`alibaba`、`deepseek`、`otherapi`。
- 图片：`openai`、`alibaba`。
- 随机图片：`github`、`gitee`。
- 空字符串会被 `util.ListChecker()` 视为合法，后续由处理器使用默认 API。

### 10.4 Prompt 与域名限制

- 文本预置 prompt 必须在 `txtgenenabled` 中启用。
- 自定义文本 prompt 被归类为 `other`，要求 `other` 在 `txtgenenabled` 中启用，并且 prompt 匹配 `txtacceptprompt` 通配符列表。
- 图片 prompt 必须匹配 `imgacceptprompt` 通配符列表。
- 网页缩略图 Host 必须在 `webimgallowed` 中。

## 11. 数据模型与内存缓存

数据库表：

- `Setting`：`Name` 主键，`Value` 为 JSON 字符串数组。
- `Task`：任务记录。
- `Session`：后台 token。
- `Repo`：随机图片仓库源。

全局内存缓存：

- `database.SettingMap map[string][]string`：配置名到字符串数组。
- `database.PromptMap map[string]int`：预置文本 prompt 到 `prompt` 配置下标。
- `database.RepoMap map[string][]string`：`api;user/repo` 到图片 URL 列表。
- `database.SessionMap map[string]Session`：token 到 session。

`Task` 字段：

- `uuid`：主键。
- `time`：任务创建时间。
- `ip`：发起 IP。
- `type`：中文任务类型，例如 `文字生成`、`图片生成`、`随机图片`、`网页缩略图`、`文件下载`。
- `status`：`success` 或 `failed` 等。
- `target`：prompt、repo、URL 等目标。
- `return`：结果 JSON 字符串或错误文本。
- `region`：IP 地区。
- `referer`：Referer。
- `device`：`Mobile`、`Desktop`、`Bot` 或空。
- `more_info`：附加信息。
- `api`：使用的 API 标识或 Host。
- `model`：模型。
- `temp`：缓存复用标记，`Yes` 或 `No`。
- `size`：图片尺寸。

`Repo` 字段：

- `uuid`：主键。
- `api`：`github` 或 `gitee`。
- `info`：`user/repo`。
- `content`：图片 URL 数组的 JSON 字符串。

`Session` 字段：

- `token`：主键。
- `expire`：过期时间。
- `term`：是否长期登录。

## 12. 默认配置 `file/setting.json`

配置通过 `names` 与 `edits` 的同下标数组对应。重构时如果保留兼容数据库，需要严格保留每个配置数组的含义和下标。

- `openai`：`[apiKey, 默认文本模型, 默认总结模型, 默认图片模型, 默认图片尺寸, chat/image endpoint]`。
- `deepseek`：`[apiKey, 默认文本模型, 默认总结模型, chat endpoint]`。
- `alibaba`：`[apiKey, 默认文本模型, 默认总结模型, 默认图片模型, 默认图片尺寸, compatible chat endpoint]`。
- `otherapi`：`[apiKey, 默认文本模型, 默认总结模型, endpoint]`。
- `dash`：`[后台密码 sha256]`，默认是 `123456` 的 SHA-256。
- `dashallowedip`：后台允许 IP，当前后端代码中未实际用于登录检查。
- `allowedref`：公开 API Referer 白名单。
- `txt`：`[启用, 默认生成API, 默认总结API, 缓存过期分钟, 失败图片URL]`。
- `txtgenenabled`：允许的文本生成类型列表，例如 `laugh`、`poem`、`sentence`、`other`。
- `img`：`[启用, 默认API, 缓存过期分钟, 失败图片URL]`。
- `web`：`[占位/未知, 启用, 默认API, 缓存过期分钟, 失败图片URL, GitHub token, YouTube token]`。
- `webimgallowed`：允许生成网页缩略图的 Host 列表。
- `rand`：`[启用, GitHub raw 代理/替换前缀, 失败图片URL, 默认API]`。
- `context`：`[文本生成 system prompt, 文本总结 system prompt]`。
- `prompt`：预置 prompt 文案，顺序对应 `laugh`、`poem`、`sentence`。
- `taskexceptdomain`：不记录任务的 Referer 域名通配符。
- `txtacceptprompt`：允许的自定义文本 prompt 通配符。
- `imgacceptprompt`：允许的图片 prompt 通配符。
- `taskexceptinfo`：不记录任务的 MoreInfo/Info 通配符。

后台配置分组由 `processor.PartMap` 定义：

- `openai` -> `openai`
- `deepseek` -> `deepseek`
- `alibaba` -> `alibaba`
- `otherapi` -> `otherapi`
- `security` -> `dash`、`dashallowedip`、`allowedref`
- `txt` -> `txt`、`txtgenenabled`
- `img` -> `img`
- `web` -> `web`、`webimgallowed`
- `rand` -> `rand`
- `contxt` -> `context`、`prompt`，注意 key 拼写是 `contxt`。
- `taskBehavior` -> `taskexceptdomain`、`taskexceptinfo`
- `imgSecurity` -> `imgacceptprompt`
- `txtSecurity` -> `txtacceptprompt`

## 13. 外部服务依赖

- IP 地区：`https://api.vore.top/api/IPdata?ip=<ip>`。
- OpenAI/兼容文本：`POST <endpoint>`，Bearer Token，OpenAI Chat Completions 格式。
- DeepSeek：通过兼容 chat endpoint 调用。
- 阿里巴巴文本：通过 DashScope compatible-mode endpoint 调用。
- 阿里巴巴图片：固定 DashScope image-synthesis endpoint，异步任务轮询。
- OpenAI 图片：使用配置中的 OpenAI endpoint。
- GitHub contents：`https://api.github.com/repos/<user/repo>/contents`。
- Gitee contents：`https://gitee.com/api/v5/repos/<user/repo>/contents`。
- Bilibili：`https://api.bilibili.com/x/web-interface/view?aid=...` 或 `?bvid=...`。
- YouTube：`https://www.googleapis.com/youtube/v3/videos?part=snippet,statistics&id=<id>&key=<token>`。
- arXiv/IT之家：直接抓取页面 HTML。

迁移时需要为所有外部 HTTP 调用设置合理超时。原项目全局 `http.Client` timeout 为 30 秒。

## 14. 图片生成与绘制

字体：

- 嵌入 `file/ssfonts.ttf`，README 表示使用“得意黑”。

图标与 Logo：

- `file/icon/*`：播放、收藏、star、like、coin、fork 等。
- `file/logo/*`：GitHub、Gitee、arXiv、IT之家等。

绘制函数：

- `DrawTxt()`：把文本按 20 个 rune 一行拆分，生成透明/文字阴影风格 PNG。
- `DrawRepo()`：生成仓库信息卡片。
- `DrawVideo()`：下载视频封面，生成视频信息卡片。
- `DrawArticle()`：生成文章信息卡片。
- `DrawRoundedRect()`：绘制圆角背景或边框。

临时图片保存：

- 路径固定为 `assets/img/<uuid>.png`。
- `/download?img=<uuid>` 读取该文件。
- 服务重启会清空 `assets/img/`，历史任务中的下载 URL 可能失效。

## 15. 前端后台结构

入口：

- `static/src/main.js` 创建 Vue app，提供全局注入：`url=/session`、`title`、`login`、`emitter`、`page`、`maxPage`。
- `static/src/App.vue` 挂载 Header、Sidebar、router-view，启动时用 Cookie 中的 token 试图登录。

路由：

- `/dash` 和 `/dash/task`：任务查看。
- `/dash/tool`：功能设置。
- `/dash/security`：安全设置。
- `/dash/backend`：接口设置。
- `/dash/workshop`：工作台。
- `/dash/login`：登录页。

通用请求：

- `static/src/js/fetch.js` 的 `Post()` 向 `/session` 发送 POST，body 是 `data.Send`，Header `Authorization` 是 `data.Token`。
- `static/src/js/util.js` 封装 `Login`、`Logout`、`Repo`、`Setting`、`Task`。

登录页：

- 用户输入密码后前端用 SHA-256 处理。
- 调 `Login(sha256(pwd), term)`。
- 成功后 Cookie 写入后端返回的 `session_token`，固定 `{expires: 7}`。

任务页：

- `Showcase.vue` 调用 `Task("fetchTask", catagory, by, page)`。
- 默认每页 100 条。
- Header 上的翻页/刷新通过 `emitter` 控制。
- `Filter.vue` 及各字段组件按任务字段筛选。

接口设置页：

- OpenAI、Alibaba、DeepSeek、OtherAPI、Context、WebAPI 等组件读取/保存不同 `SettingPart`。

功能设置页：

- Text、Image、Rand、Web 组件管理功能开关、默认 API、缓存过期、失败图片、允许项等。

安全设置页：

- DashSecurity、TaskBehavior、ImgSecurity、TxtSecurity 分别管理后台密码/白名单、任务记录例外、图片 Prompt 限制、文本 Prompt 限制。

工作台：

- 组合 Host、Operation、Query 参数生成 URL。
- 可用 `img` 或 `iframe` 展示调用结果。
- 页面提示需要把 API 域名加入 Referer 白名单。

## 16. 已知实现细节与潜在问题

这些不是立即要修复的内容，但重构迁移时需要明确是否兼容或修正。

- `main.go` 中 `handler.Handler()` 阻塞，`defer database.Disconnect()` 实际只在服务停止后执行。
- `command.Arg()` 处理 `port` 时未检查下一个参数是否存在。
- `GeneralChecker()` 先执行 `ExceptionChecker()` 再执行 `InfoChecker()`，而 `taskexceptinfo` 匹配的是当时的 `Info`，通常尚未被设置为错误信息。
- `dashallowedip` 有默认配置和命令行清除逻辑，但当前登录鉴权没有检查该配置。
- `util.ListChecker()` 对空 API 返回 true，所以 API 检查允许空值，依赖 processor 填默认值。
- `TxtGen.APIChecker()` 如果 API 是空字符串会通过，这属于当前行为。
- `returner()` 对 `data.Return` 做 `json.Unmarshal`，若失败会返回空结构，但非 JSON 格式仍会按 `Processor.Return` 跳转。
- 多处 `fmt.Sprintf` 手写 JSON 字符串，prompt/response 中包含引号、换行等字符时可能产生非法 JSON。
- `/web` 的参数名是 `img`，不是更直观的 `url`。
- `getBiliABV()`、`getYtbID()` 使用固定下标解析 URL，对不同 URL 形式较脆弱。
- `util.ITHome()` 打开 Logo 使用了 `file.Logos.Open("assets/logo/ithome_logo.png")`，但嵌入路径是 `logo/ithome_logo.png`，该路径可能导致运行错误。
- `util.Downloader()` 在 `GlobalHTTPClient.Get()` 返回错误且 `resp == nil` 时访问 `resp.Status` 可能 panic。
- `processor.TaskCounter.Counter[r.Processor.Filter.API]--` 对不存在 key 或异常路径可能产生负数。
- `beforeTask()` 等待已有任务时读取的 `task` 是进入循环前的副本，没有在循环中重新从 map 读取，若 Running 状态变化可能存在等待逻辑问题。
- `Repo.Create()` 不更新 `RepoMap`，新增仓库后当前进程内随机接口可能无法立即使用，除非前端刷新后有其他路径更新或重启；`Repo.Update()` 会更新 `RepoMap`。
- `fetchTask()` 分页切片 `taskList[(page-1)*100 : (page*100 - 1)]` 每页可能少一条。
- 前端 `Login.vue` 中 `if (Cookies.get("token") && await (Cookies.get("token"), false))` 少调用 `Login()`，这个判断实际不会校验 token。
- 前端 `Header.vue` 对 inject 的 `emitter` 使用 `(emitter=2)` 这类赋值，可能不是期望的 `emitter.value = 2`。
- `App.vue` 引入 `Login` 但 `onUnmounted` 调用 `Logout` 未导入。
- `Login.vue` 中 “7天内保持登录” checkbox 使用 `@input="term = !$event.target.checked"`，语义可能反了。
- Go 版本写为 `1.26.2`，这不是当前常规稳定 Go 版本；构建环境必须确认可用。

## 17. 重构迁移清单

后端迁移必须覆盖：

- 所有公开路由：`/txt`、`/img`、`/rand`、`/web`、`/download`、`/session`、`/dash`、`/assets`。
- 每个 API 的 Query 参数、默认值、`format=json` 与 302 跳转行为。
- 任务记录字段和中文任务类型映射。
- SQLite 数据迁移或兼容策略，尤其是 `Setting.Value` 的 JSON 数组下标语义。
- `SettingMap`、`RepoMap`、`SessionMap` 的运行时缓存语义，或等价替代。
- Referer 白名单、频率限制、功能开关、API 白名单、Prompt 通配符限制、记录例外。
- 任务缓存、过期时间、并发限制、`Temp` 标记、图片有效性检查。
- 外部 API 调用行为与失败时 `Status/Return/ReturnError` 的结果语义。
- `assets/img` 临时图片目录和 `/download` 读取方式。
- 嵌入资源：字体、图标、Logo、空图片、默认设置。
- 后台 session token 的创建、过期、删除、Authorization 传递方式。
- 后台 `operation` 名称和请求/响应字段，避免破坏现有前端。

前端迁移必须覆盖：

- `/dash` 下所有页面和路由。
- 登录、Cookie token、`/session` 调用协议。
- 任务查看、筛选、分页、详情复制等管理能力。
- 接口设置、功能设置、安全设置各配置分组。
- 仓库新增、刷新、删除、读取。
- 工作台 URL 生成和结果预览。

部署迁移必须覆盖：

- 单体部署：Go 后端内嵌前端 dist。
- 前后端分离部署：Nginx 静态托管 + `/txt`、`/img`、`/rand`、`/web`、`/download`、`/session` 反代到后端。
- Docker 镜像端口：后端 `2233`，前端 `80`。
- 数据持久化目录：README 约定把宿主机目录挂载到 `/app/assets`。

## 18. 建议的重构边界

为降低迁移风险，建议按以下边界重构：

1. 先冻结接口协议：公开 API 参数、返回、后台 `/session` operation 不变。
2. 把 `Setting` 数组下标转为命名结构，但提供从旧数组迁移/读取的兼容层。
3. 把 `handler -> security -> processor -> database` 流程显式化，减少全局变量依赖。
4. 将任务缓存和并发限制抽成独立服务，保留原 `TaskQueueFilter` 语义。
5. 将外部 API 客户端独立封装，统一超时、错误和 JSON 编码。
6. 图片绘制可作为独立包迁移，优先保留输出效果和 `/download` 协议。
7. 前端可以后迁移，只要 `/session` 协议保持稳定。

## 19. 快速功能对照表

| 功能 | 后端入口 | Processor | 配置 | 数据记录 |
| --- | --- | --- | --- | --- |
| 文本生成 | `GET /txt` | `processor.TxtGen` | `txt`、`txtgenenabled`、`context`、`prompt`、各 API 配置 | `Task` |
| 图片生成 | `GET /img` | `processor.ImgGen` | `img`、`imgacceptprompt`、`openai`、`alibaba` | `Task` |
| 随机图片 | `GET /rand` | `processor.Rand` | `rand`、`RepoMap` | `Task` |
| 网页缩略图 | `GET /web` | `processor.WebImg` | `web`、`webimgallowed`、`txt/context` | `Task` |
| 下载 | `GET /download` | `processor.Download` | `web[4]` 失败图 | 默认不记录 |
| 后台会话 | `POST /session` | `processor.Session` | `dash` | `Session` |
| 后台配置 | `POST /session editSetting/fetchSetting` | `processor.Session` | `PartMap` | `Setting` |
| 后台任务 | `POST /session fetchTask` | `processor.Session` | 无 | `Task` |
| 仓库管理 | `POST /session newRepo/refreshRepo/delRepo/fetchRepo` | `processor.Session` | `rand` | `Repo` |
