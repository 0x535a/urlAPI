# 新设置 API 文档

本文记录后端新增的结构化设置 API，供后续改造 Vue 前端时参考。旧 `fetchSetting/editSetting` 仍保留兼容，但新页面应使用本文中的 `fetchSettings/editSettings`。

## 1. 基本约定

请求入口仍然是当前后台会话接口：

```text
POST /session
```

认证方式不变：

```text
Authorization: <session_token>
Content-Type: application/json
```

新增两个 operation：

```text
fetchSettings
editSettings
```

请求体通用格式：

```json
{
  "operation": "fetchSettings",
  "setting_part": "feature.text"
}
```

编辑请求体通用格式：

```json
{
  "operation": "editSettings",
  "setting_part": "feature.text",
  "setting_body": {}
}
```

成功响应仍然返回完整 session 对象，其中新接口主要读取：

```json
{
  "setting_part": "feature.text",
  "setting_body": {}
}
```

失败响应沿用现有错误格式：

```json
{
  "error": "..."
}
```

## 2. 与旧 API 的区别

旧 API：

```text
operation=fetchSetting
operation=editSetting
```

旧 API 使用 `setting_data: [][]string` 和 `setting_edit: [][]string`，前端需要记数组下标。

新 API：

```text
operation=fetchSettings
operation=editSettings
```

新 API 使用 `setting_body` 具名对象，不再暴露数据库数组结构。

旧 API 暂时保留，现有前端可以继续运行。新前端页面应逐步迁移到新 API。

## 3. setting_part 列表

推荐使用新 part 名称。后端也兼容部分旧名称，方便渐进迁移。

| 新 part | 旧 part 兼容 | 用途 |
|---|---|---|
| `provider.openai` | `openai` | OpenAI 供应商设置 |
| `provider.deepseek` | `deepseek` | DeepSeek 供应商设置 |
| `provider.alibaba` | `alibaba` | Alibaba 供应商设置 |
| `provider.otherapi` | `otherapi` | 其他 OpenAI 兼容接口设置 |
| `feature.text` | `txt` | 文字功能设置 |
| `feature.image` | `img` | 图片功能设置 |
| `feature.web` | `web` | 网页图片功能设置 |
| `feature.random` | `rand` | 随机图功能设置 |
| `security.dashboard` | `security` | 后台安全设置 |
| `prompt` | `contxt` | Prompt 和上下文设置 |
| `security.task_behavior` | `taskBehavior` | 任务记录例外设置 |
| `security.text_prompt` | `txtSecurity` | 文字 Prompt 安全设置 |
| `security.image_prompt` | `imgSecurity` | 图片 Prompt 安全设置 |

## 4. 供应商设置

适用 part：

```text
provider.openai
provider.deepseek
provider.alibaba
provider.otherapi
```

### fetchSettings 响应

```json
{
  "api_key_set": true,
  "text_model": "gpt-4o",
  "summary_model": "gpt-4o-mini",
  "image_model": "dall-e-3",
  "image_size": "1024x1024",
  "endpoint": "https://api.openai.com/v1/chat/completions"
}
```

字段说明：

| 字段 | 类型 | 说明 |
|---|---|---|
| `api_key_set` | bool | 是否已配置 API Key |
| `api_key` | string | 只在编辑时提交；fetch 不返回明文 |
| `text_model` | string | 默认文字生成模型 |
| `summary_model` | string | 默认总结模型 |
| `image_model` | string | 默认图片模型，DeepSeek/OtherAPI 可为空 |
| `image_size` | string | 默认图片尺寸，DeepSeek/OtherAPI 可为空 |
| `endpoint` | string | API endpoint |

### editSettings 请求

```json
{
  "operation": "editSettings",
  "setting_part": "provider.openai",
  "setting_body": {
    "api_key": "sk-xxx",
    "text_model": "gpt-4o",
    "summary_model": "gpt-4o-mini",
    "image_model": "dall-e-3",
    "image_size": "1024x1024",
    "endpoint": "https://api.openai.com/v1/chat/completions"
  }
}
```

注意：`api_key` 为空或不传时，后端保留旧 API Key，不会清空。

## 5. 文字功能设置

part：

```text
feature.text
```

### setting_body

```json
{
  "enabled": false,
  "generation_api": "alibaba",
  "summary_api": "alibaba",
  "cache_minutes": 60,
  "fallback_image_url": "https://raw.githubusercontent.com/stephen-zeng/img/master/fallback.png",
  "enabled_prompt_keys": ["laugh", "poem", "sentence", "other"],
  "accepted_prompt_glob": ["*"]
}
```

字段说明：

| 字段 | 类型 | 说明 |
|---|---|---|
| `enabled` | bool | 文字功能总开关 |
| `generation_api` | string | 随机生成使用的供应商：`openai`、`deepseek`、`alibaba`、`otherapi` |
| `summary_api` | string | 总结使用的供应商 |
| `cache_minutes` | int | 缓存过期时间，单位分钟，必须大于等于 0 |
| `fallback_image_url` | string | 失败时返回图片 URL |
| `enabled_prompt_keys` | string[] | 启用的内置 Prompt key |
| `accepted_prompt_glob` | string[] | 允许的自定义 Prompt 通配符 |

## 6. 图片功能设置

part：

```text
feature.image
```

### setting_body

```json
{
  "enabled": false,
  "api": "alibaba",
  "cache_minutes": 60,
  "fallback_image_url": "https://raw.githubusercontent.com/stephen-zeng/urlAPI/img/master/fallback.png",
  "accepted_prompt_glob": ["*"]
}
```

字段说明：

| 字段 | 类型 | 说明 |
|---|---|---|
| `enabled` | bool | 图片生成功能总开关 |
| `api` | string | 图片生成供应商，目前支持 `openai`、`alibaba` |
| `cache_minutes` | int | 缓存过期时间，单位分钟 |
| `fallback_image_url` | string | 失败时返回图片 URL |
| `accepted_prompt_glob` | string[] | 允许的 Prompt 通配符 |

## 7. 网页图片功能设置

part：

```text
feature.web
```

### fetchSettings 响应

```json
{
  "enabled": false,
  "summary_api": "alibaba",
  "cache_minutes": 10,
  "fallback_image_url": "https://raw.githubusercontent.com/stephen-zeng/urlAPI/img/master/fallback.png",
  "repo_token_set": false,
  "youtube_token_set": false,
  "allowed_hosts": ["github.com", "gitee.com", "www.youtube.com"]
}
```

### editSettings 请求

```json
{
  "operation": "editSettings",
  "setting_part": "feature.web",
  "setting_body": {
    "enabled": true,
    "summary_api": "alibaba",
    "cache_minutes": 10,
    "fallback_image_url": "https://raw.githubusercontent.com/stephen-zeng/urlAPI/img/master/fallback.png",
    "repo_token": "github-or-gitee-token",
    "youtube_token": "youtube-token",
    "allowed_hosts": ["github.com", "gitee.com", "www.youtube.com"]
  }
}
```

注意：`repo_token` 和 `youtube_token` 为空或不传时，后端保留旧值。

## 8. 随机图功能设置

part：

```text
feature.random
```

### setting_body

```json
{
  "enabled": false,
  "source_rewrite_from": "https://raw.githubusercontent.com",
  "fallback_image_url": "https://raw.githubusercontent.com/stephen-zeng/urlAPI/master/fallback.png",
  "default_api": "github"
}
```

字段说明：

| 字段 | 类型 | 说明 |
|---|---|---|
| `enabled` | bool | 随机图功能总开关 |
| `source_rewrite_from` | string | GitHub raw URL 替换目标 |
| `fallback_image_url` | string | 失败时 fallback URL |
| `default_api` | string | 默认仓库 API：`github` 或 `gitee` |

## 9. 后台安全设置

part：

```text
security.dashboard
```

### fetchSettings 响应

```json
{
  "dashboard_allowed_ips": ["*"],
  "allowed_referers": ["localhost"]
}
```

### editSettings 请求

```json
{
  "operation": "editSettings",
  "setting_part": "security.dashboard",
  "setting_body": {
    "password_hash": "<sha256 hex>",
    "dashboard_allowed_ips": ["*"],
    "allowed_referers": ["localhost"]
  }
}
```

注意：`password_hash` 为空或不传时，后端保留旧密码 hash。

## 10. Prompt 设置

part：

```text
prompt
```

### setting_body

```json
{
  "generation_context": "你是一个助手...",
  "summary_context": "你是一个助手...",
  "templates": {
    "laugh": "讲一个笑话...",
    "poem": "做几句诗歌...",
    "sentence": "写几句心灵鸡汤..."
  }
}
```

字段说明：

| 字段 | 类型 | 说明 |
|---|---|---|
| `generation_context` | string | 文字生成系统上下文 |
| `summary_context` | string | 总结系统上下文 |
| `templates` | object | 内置 Prompt 模板，key 到模板内容 |

## 11. 任务记录例外设置

part：

```text
security.task_behavior
```

### setting_body

```json
{
  "except_domains": ["localhost"],
  "except_infos": []
}
```

字段说明：

| 字段 | 类型 | 说明 |
|---|---|---|
| `except_domains` | string[] | 命中这些 Referer domain 时跳过任务记录 |
| `except_infos` | string[] | 命中这些错误信息时跳过任务记录 |

## 12. 文字 Prompt 安全设置

part：

```text
security.text_prompt
```

### setting_body

```json
{
  "accepted_prompt_glob": ["*"]
}
```

该接口只修改 `Text.AcceptedPromptGlob`。如果前端已经使用 `feature.text` 统一保存文字设置，可以不单独调用此接口。

## 13. 图片 Prompt 安全设置

part：

```text
security.image_prompt
```

### setting_body

```json
{
  "accepted_prompt_glob": ["*"]
}
```

该接口只修改 `Image.AcceptedPromptGlob`。如果前端已经使用 `feature.image` 统一保存图片设置，可以不单独调用此接口。

## 14. 前端调用建议

新的前端封装可以类似：

```js
export async function Settings(operation, settingPart = "", settingBody = null) {
  const session = await Post({
    Token: Cookies.get("token"),
    Send: {
      operation,
      setting_part: settingPart,
      setting_body: settingBody,
    },
  })
  if (session.error) {
    Notification(session.error)
    return null
  }
  if (operation === "fetchSettings") {
    return session.setting_body
  }
  Notification("Saved")
  return true
}
```

示例：

```js
const textSettings = await Settings("fetchSettings", "feature.text")
textSettings.enabled = true
await Settings("editSettings", "feature.text", textSettings)
```

## 15. 校验规则

当前后端已做基础校验：

1. `cache_minutes` 必须大于等于 0。
2. URL 字段为空时允许；非空时必须是 `http` 或 `https` URL。
3. 文字/总结 API 支持 `openai`、`deepseek`、`alibaba`、`otherapi`。
4. 图片 API 支持 `openai`、`alibaba`。
5. 未知 `setting_part` 返回错误。
6. 敏感字段 fetch 不返回明文，edit 为空时保留旧值。

## 16. 迁移顺序建议

建议前端按以下顺序迁移：

1. 新增 `Settings()` 封装，不删除旧 `Setting()`。
2. 先迁移供应商页面，因为 DTO 简单且能验证敏感字段不回显。
3. 迁移 `Tool/Text.vue`、`Tool/Image.vue`、`Tool/Web.vue`、`Tool/Rand.vue`。
4. 迁移 `Backend/Context.vue`。
5. 迁移 `Security/*` 页面。
6. 所有页面迁移完成后，后端再考虑删除旧 `fetchSetting/editSetting`。
