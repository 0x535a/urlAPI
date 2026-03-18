# 项目重写说明（功能、接口、鉴权、安全、数据库）

本文为重写版本的完整功能与数据设计说明，不依赖现有实现细节。

## 功能模块（完整描述）

- 文本生成：基于外部模型生成短文/句子/自定义内容
- 图片生成：基于提示词生成图像
- 随机图片：从指定仓库随机返回图片
- 网页缩略图：对指定页面生成缩略图/封面图
- 下载中转：对生成图片进行下载或回源
- 后台管理：管理访问、查询任务、系统配置
- 安全与防护：防盗链、频率限制、白名单、异常过滤

## 对外接口总览

- `GET /txt` 文本生成
- `GET /img` 图片生成
- `GET /rand` 随机图片
- `GET /web` 网页缩略图
- `GET /download` 下载图片（内部回源/中转）
- `POST /session` 后台会话与管理操作
- `/dash` 后台管理前端（静态页面）

## 接口参数与行为（重写版本）

### `GET /txt` 文本生成
- 参数：
  - `prompt` 必填，内容或预置关键字
  - `api` 选填，来源：`openai`/`alibaba`/`deepseek`/`otherapi`
  - `model` 选填，模型名称
  - `more` 选填，附加信息用于记录/筛选
  - `format` 选填，`json` 输出 JSON，否则 302 跳转结果地址
- 返回：
  - `format=json` 返回 JSON（`prompt/response/url` 等字段）
  - 否则 302 跳转到生成结果地址

### `GET /img` 图片生成
- 参数：
  - `prompt` 必填，提示词
  - `api` 选填，来源：`openai`/`alibaba`
  - `model` 选填，模型名称
  - `size` 选填，图像尺寸
  - `more` 选填，附加信息用于记录/筛选
  - `format` 选填，`json` 输出 JSON，否则 302 跳转结果地址

### `GET /rand` 随机图片
- 参数：
  - `api` 选填，来源：`github`/`gitee`
  - `user` 必填，仓库用户
  - `repo` 必填，仓库名
  - `more` 选填，附加信息
  - `format` 选填，`json` 输出 JSON，否则 302 跳转结果地址

### `GET /web` 网页缩略图
- 参数：
  - `img` 或 `url` 必填，目标页面 URL
  - `more` 选填，附加信息
  - `format` 选填，`json` 输出 JSON，否则 302 跳转结果地址
- 说明：
  - `img/url` 解析域名用于 API 路由与安全检查

### `GET /download` 图片下载/中转
- 参数：
  - `img` 必填，内部生成的图片标识
- 行为：
  - 校验通过后以 `image/png` 返回
  - 失败时 302 跳转至错误图片地址

### `POST /session` 后台会话与管理
- 请求体：`Session` 结构（JSON），`operation` 决定操作类型
- Header：`Authorization` 作为会话 token
- 支持操作：
  - `login` 登录并生成 token
  - `logout` 登出
  - `exit` 退出（短期 token 失效）
  - `fetchTask` 查询任务
  - `fetchSetting` 读取安全与功能配置（只读）
  - `editSetting` 修改安全与功能配置
- 返回：
  - 成功时返回包含会话与数据的 JSON
  - 失败时 `400` 并返回错误信息

## 鉴权与安全过滤（重写版本）

### 通用安全检查（所有公开接口）
统一在 `GeneralChecker()` 中执行：
- 频率限制：同 IP + 同类型 0.25 秒内超过 10 次视为频率异常
- Referer 检测：
  - 读取配置 `allowedref`
  - Referer 为空或域名不匹配则拒绝
- 例外过滤：
  - 命中 `taskexceptdomain` 或 `taskexceptinfo` 时标记 `SkipDB`，不写库

### 功能与 API 白名单
- 文本生成：
  - 需要 `txt` 功能启用
  - `txtgenenabled` 决定允许的 prompt 类型
  - `txtacceptprompt` 正则限制自由 prompt
  - API 必须在 `openai/alibaba/deepseek/otherapi` 列表
- 图片生成：
  - 需要 `img` 功能启用
  - `imgacceptprompt` 正则限制 prompt
  - API 必须在 `openai/alibaba` 列表
- 随机图片：
  - 需要 `rand` 功能启用
  - API 必须在 `github/gitee` 列表
- 网页缩略图：
  - 需要 `web` 功能启用
  - `webimgallowed` 限制域名/API
  - `ithome` 特例要求文本汇总能力启用

### 后台鉴权
- `login`：使用 `dash` 配置中的密码进行校验
- 成功后生成 `session_token`，有效期：
  - 勾选长期：7 天
  - 否则：1 天
- 后续操作通过 `Authorization` 传递 token

## 数据库设计（仅保留 task）

### `task` 表/集合（完整字段）
- `uuid` string 主键，任务 ID
- `time` datetime 任务创建时间
- `ip` string 发起 IP
- `type` string 任务类型（`txt.gen`/`img.gen`/`web.img`/`rand`/`download`）
- `status` string 任务状态（`pending`/`running`/`success`/`failed`）
- `target` string 任务目标（prompt、repo、url 等）
- `return` string 结果（JSON 或 URL）
- `region` string 地理区域
- `referer` string Referer
- `device` string 设备类型
- `more_info` string 附加信息
- `api` string 使用的外部 API 标识
- `model` string 使用的模型
- `temp` string 是否缓存复用（`Yes/No`）
- `size` string 图像尺寸

### 任务查询能力
- 支持按 `type`/`status`/`api`/`model`/`ip`/`referer`/`time` 查询
- 分页读取与排序（默认按时间倒序）
