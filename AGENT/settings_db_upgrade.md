# 设置数据库改进与迁移方案

本文记录当前设置数据库的主要问题、目标设计、迁移步骤和验证清单。后续重构设置模块时，应优先参考本文，避免继续扩大 `settings.value` 字符串数组和魔法下标的使用范围。

## 1. 背景

当前设置数据由 `database.Setting` 表承载：

```go
type Setting struct {
    Name  string `json:"name" gorm:"primaryKey"`
    Value string `json:"value"`
}
```

数据库中每行代表一个配置组，`value` 是 JSON 字符串数组。例如：

```text
name  = "txt"
value = ["false", "alibaba", "alibaba", "60", "https://..."]
```

默认值来自 `file/setting.json` 的 `names` 与 `edits` 两个数组。运行时通过 `database.SettingMap map[string][]string` 缓存配置，业务代码和前端页面都直接按数组下标读取和写入配置。

该设计在早期实现简单，但当前设置项已经覆盖模型供应商、功能开关、缓存策略、安全策略、白名单、Prompt 模板、API Key、Token、fallback URL 等多个领域，继续使用无字段名数组会显著增加维护成本和运行风险。

## 2. 当前设计问题

### 2.1 数据库无法表达业务语义

`settings` 表只知道 `name` 和 `value`，不知道 `txt[0]` 是开关，`txt[3]` 是缓存时间，`web[5]` 是 GitHub/Gitee token。配置语义散落在默认 JSON、后端业务代码和 Vue 组件里，数据库层无法做类型约束、字段校验、迁移判断或审计。

### 2.2 魔法下标导致高耦合

典型读取方式：

```go
database.SettingMap["txt"][0] // enabled
database.SettingMap["txt"][1] // generation api
database.SettingMap["txt"][3] // cache minutes
database.SettingMap["web"][6] // YouTube token
```

只要默认配置、前端写入、后端读取中任一位置发生下标错位，就会产生难排查的问题。当前代码已经存在 `Deepseek.vue` 读取 `settings[0][3]` 但写入 `settings[0][5]` 的问题。

### 2.3 前端暴露并修改内部存储结构

`fetchSetting` 返回 `[][]string`，Vue 页面直接执行 `settings[0][3] = value`。这使前端协议等同于数据库内部格式，后端无法在不破坏页面的情况下重构存储结构。

### 2.4 所有字段都是字符串

布尔值、数字、URL、Token、模型名、IP、通配符列表都以字符串保存。后端需要到处做字符串比较和临时转换，且目前部分转换失败会被忽略，例如缓存过期时间 `strconv.Atoi(...)` 忽略错误。

### 2.5 缺少校验和边界保护

`editSetting` 基本直接保存前端提交的数组，没有验证 `SettingPart` 是否存在、配置组数量是否匹配、数组长度是否足够、字段类型是否合法、URL/IP/endpoint 是否有效、敏感字段是否允许返回。

### 2.6 默认配置迁移能力弱

`settingInit()` 只会在已有数组长度小于默认数组时追加尾部字段，无法处理中间插入字段、删除字段、字段改名、字段含义变化、配置拆分和配置合并。

### 2.7 全局缓存并发风险

`database.SettingMap` 是普通 map，业务请求会频繁读，后台保存设置时会写。没有锁或不可变快照机制，存在并发读写风险。

### 2.8 敏感配置无法字段级治理

API Key、后台密码 hash、第三方 token 和普通配置混在相同数组中。`fetchSetting` 会完整返回配置数组，未来做脱敏、权限拆分、审计和密钥轮换会很困难。

## 3. 改进目标

设置模块重构目标如下：

1. 消除业务代码和前端中的魔法下标。
2. 使用具名字段表达配置语义。
3. 为 bool、int、URL、token、列表等字段建立类型和校验规则。
4. 后端 API 返回稳定 DTO，不暴露数据库内部结构。
5. 支持从旧数组配置自动迁移到新结构。
6. 支持配置版本管理，后续字段变更可控迁移。
7. 对敏感字段建立脱敏和写入策略。
8. 使用线程安全的配置读取机制。
9. 保持现有行为兼容，迁移后原有功能不应改变。

## 4. 推荐目标设计

推荐采用“结构化配置 + 单行 JSON 文档 + 运行时快照”的方案作为第一阶段目标。该方案改动小于完全拆表，但能解决当前最主要的数组下标问题。

### 4.1 新数据库表

新增 `app_settings` 表：

```go
type AppSetting struct {
    Key       string    `json:"key" gorm:"primaryKey"`
    Version   int       `json:"version"`
    Value     string    `json:"value" gorm:"type:text"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

建议只保存一行：

```text
key     = "app"
version = 1
value   = JSON object of AppSettings
```

保留旧 `settings` 表一段时间用于迁移和回滚，不要在第一阶段直接删除。

### 4.2 新配置结构

建议在 `database` 或新包 `settings` 中定义：

```go
type AppSettings struct {
    Version  int              `json:"version"`
    Providers ProviderSettings `json:"providers"`
    Features  FeatureSettings  `json:"features"`
    Security  SecuritySettings `json:"security"`
    Text      TextSettings     `json:"text"`
    Image     ImageSettings    `json:"image"`
    Web       WebSettings      `json:"web"`
    Random    RandomSettings   `json:"random"`
    Prompts   PromptSettings   `json:"prompts"`
    Task      TaskSettings     `json:"task"`
}
```

供应商配置：

```go
type ProviderSettings struct {
    OpenAI   ProviderConfig `json:"openai"`
    DeepSeek ProviderConfig `json:"deepseek"`
    Alibaba  ProviderConfig `json:"alibaba"`
    OtherAPI ProviderConfig `json:"otherapi"`
}

type ProviderConfig struct {
    APIKey       string `json:"api_key" sensitive:"true"`
    TextModel    string `json:"text_model"`
    SummaryModel string `json:"summary_model"`
    ImageModel   string `json:"image_model,omitempty"`
    ImageSize    string `json:"image_size,omitempty"`
    Endpoint     string `json:"endpoint"`
}
```

功能配置：

```go
type FeatureSettings struct {
    TextEnabled   bool `json:"text_enabled"`
    ImageEnabled  bool `json:"image_enabled"`
    WebImgEnabled bool `json:"web_img_enabled"`
    RandomEnabled bool `json:"random_enabled"`
}
```

文字配置：

```go
type TextSettings struct {
    GenerationAPI      string   `json:"generation_api"`
    SummaryAPI         string   `json:"summary_api"`
    CacheMinutes       int      `json:"cache_minutes"`
    FallbackImageURL   string   `json:"fallback_image_url"`
    EnabledPromptKeys  []string `json:"enabled_prompt_keys"`
    AcceptedPromptGlob []string `json:"accepted_prompt_glob"`
}
```

图片配置：

```go
type ImageSettings struct {
    API                string   `json:"api"`
    CacheMinutes       int      `json:"cache_minutes"`
    FallbackImageURL   string   `json:"fallback_image_url"`
    AcceptedPromptGlob []string `json:"accepted_prompt_glob"`
}
```

网页图片配置：

```go
type WebSettings struct {
    SummaryAPI       string   `json:"summary_api"`
    CacheMinutes     int      `json:"cache_minutes"`
    FallbackImageURL string   `json:"fallback_image_url"`
    RepoToken        string   `json:"repo_token" sensitive:"true"`
    YouTubeToken     string   `json:"youtube_token" sensitive:"true"`
    AllowedHosts     []string `json:"allowed_hosts"`
}
```

随机图配置：

```go
type RandomSettings struct {
    SourceRewriteFrom string `json:"source_rewrite_from"`
    FallbackImageURL  string `json:"fallback_image_url"`
    DefaultAPI        string `json:"default_api"`
}
```

安全配置：

```go
type SecuritySettings struct {
    DashboardPasswordHash string   `json:"dashboard_password_hash" sensitive:"true"`
    DashboardAllowedIPs   []string `json:"dashboard_allowed_ips"`
    AllowedReferers       []string `json:"allowed_referers"`
}
```

Prompt 配置：

```go
type PromptSettings struct {
    GenerationContext string            `json:"generation_context"`
    SummaryContext    string            `json:"summary_context"`
    Templates         map[string]string `json:"templates"`
}
```

任务过滤配置：

```go
type TaskSettings struct {
    ExceptDomains []string `json:"except_domains"`
    ExceptInfos   []string `json:"except_infos"`
}
```

### 4.3 运行时缓存

不要继续暴露 `map[string][]string`。推荐使用不可变快照加锁：

```go
type Store struct {
    mu       sync.RWMutex
    settings AppSettings
}

func (s *Store) Get() AppSettings {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.settings
}

func (s *Store) Replace(next AppSettings) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.settings = next
}
```

读取方拿到的是配置快照，不直接读写全局 map。保存配置时先校验，再写数据库，最后替换内存快照。

### 4.4 后端业务读取方式

迁移前：

```go
info.API = database.SettingMap["txt"][1]
expiredTime, _ := strconv.Atoi(database.SettingMap["txt"][3])
```

迁移后：

```go
cfg := settings.Store.Get()
info.API = cfg.Text.GenerationAPI
expiredTime := cfg.Text.CacheMinutes
```

Provider 读取示例：

```go
provider := cfg.Providers.ByName(info.API)
token := provider.APIKey
model := provider.TextModel
endpoint := provider.Endpoint
```

建议实现 `ByName(name string) (ProviderConfig, bool)`，避免在业务代码里重复 switch。

### 4.5 后台 API DTO

后台设置接口不应直接返回 `AppSettings` 全量对象，尤其不能默认返回敏感字段明文。建议按页面定义 DTO：

```go
type TextSettingsDTO struct {
    Enabled            bool     `json:"enabled"`
    GenerationAPI      string   `json:"generation_api"`
    SummaryAPI         string   `json:"summary_api"`
    CacheMinutes       int      `json:"cache_minutes"`
    FallbackImageURL   string   `json:"fallback_image_url"`
    EnabledPromptKeys  []string `json:"enabled_prompt_keys"`
    AcceptedPromptGlob []string `json:"accepted_prompt_glob"`
}
```

供应商 DTO 对 API Key 做特殊处理：

```go
type ProviderSettingsDTO struct {
    APIKeySet    bool   `json:"api_key_set"`
    APIKey       string `json:"api_key,omitempty"`
    TextModel    string `json:"text_model"`
    SummaryModel string `json:"summary_model"`
    ImageModel   string `json:"image_model,omitempty"`
    ImageSize    string `json:"image_size,omitempty"`
    Endpoint     string `json:"endpoint"`
}
```

更新时如果 `api_key` 为空，保留旧 key；如果显式传入新 key，则覆盖。这样可以避免每次打开页面都返回密钥明文。

## 5. 旧配置到新配置映射

迁移器需要从旧 `settings` 表或旧 `SettingMap` 读取数组，并构造 `AppSettings`。

### 5.1 Provider 映射

| 旧 key | 旧下标 | 新字段 |
|---|---:|---|
| `openai` | 0 | `Providers.OpenAI.APIKey` |
| `openai` | 1 | `Providers.OpenAI.TextModel` |
| `openai` | 2 | `Providers.OpenAI.SummaryModel` |
| `openai` | 3 | `Providers.OpenAI.ImageModel` |
| `openai` | 4 | `Providers.OpenAI.ImageSize` |
| `openai` | 5 | `Providers.OpenAI.Endpoint` |
| `deepseek` | 0 | `Providers.DeepSeek.APIKey` |
| `deepseek` | 1 | `Providers.DeepSeek.TextModel` |
| `deepseek` | 2 | `Providers.DeepSeek.SummaryModel` |
| `deepseek` | 3 | `Providers.DeepSeek.Endpoint` |
| `alibaba` | 0 | `Providers.Alibaba.APIKey` |
| `alibaba` | 1 | `Providers.Alibaba.TextModel` |
| `alibaba` | 2 | `Providers.Alibaba.SummaryModel` |
| `alibaba` | 3 | `Providers.Alibaba.ImageModel` |
| `alibaba` | 4 | `Providers.Alibaba.ImageSize` |
| `alibaba` | 5 | `Providers.Alibaba.Endpoint` |
| `otherapi` | 0 | `Providers.OtherAPI.APIKey` |
| `otherapi` | 1 | `Providers.OtherAPI.TextModel` |
| `otherapi` | 2 | `Providers.OtherAPI.SummaryModel` |
| `otherapi` | 3 | `Providers.OtherAPI.Endpoint` |

### 5.2 功能和业务映射

| 旧 key | 旧下标 | 新字段 |
|---|---:|---|
| `txt` | 0 | `Features.TextEnabled` |
| `txt` | 1 | `Text.GenerationAPI` |
| `txt` | 2 | `Text.SummaryAPI` |
| `txt` | 3 | `Text.CacheMinutes` |
| `txt` | 4 | `Text.FallbackImageURL` |
| `txtgenenabled` | 全部 | `Text.EnabledPromptKeys` |
| `txtacceptprompt` | 全部 | `Text.AcceptedPromptGlob` |
| `img` | 0 | `Features.ImageEnabled` |
| `img` | 1 | `Image.API` |
| `img` | 2 | `Image.CacheMinutes` |
| `img` | 3 | `Image.FallbackImageURL` |
| `imgacceptprompt` | 全部 | `Image.AcceptedPromptGlob` |
| `web` | 1 | `Features.WebImgEnabled` |
| `web` | 2 | `Web.SummaryAPI` |
| `web` | 3 | `Web.CacheMinutes` |
| `web` | 4 | `Web.FallbackImageURL` |
| `web` | 5 | `Web.RepoToken` |
| `web` | 6 | `Web.YouTubeToken` |
| `webimgallowed` | 全部 | `Web.AllowedHosts` |
| `rand` | 0 | `Features.RandomEnabled` |
| `rand` | 1 | `Random.SourceRewriteFrom` |
| `rand` | 2 | `Random.FallbackImageURL` |
| `rand` | 3 | `Random.DefaultAPI` |

注意：旧 `web[0]` 当前默认值为空，现有业务未发现稳定使用，应迁移器保留到兼容日志中，但不映射到新字段。若后续确认含义，再补字段。

### 5.3 安全和 Prompt 映射

| 旧 key | 旧下标 | 新字段 |
|---|---:|---|
| `dash` | 0 | `Security.DashboardPasswordHash` |
| `dashallowedip` | 全部 | `Security.DashboardAllowedIPs` |
| `allowedref` | 全部 | `Security.AllowedReferers` |
| `context` | 0 | `Prompts.GenerationContext` |
| `context` | 1 | `Prompts.SummaryContext` |
| `prompt` | 0 | `Prompts.Templates["laugh"]` |
| `prompt` | 1 | `Prompts.Templates["poem"]` |
| `prompt` | 2 | `Prompts.Templates["sentence"]` |
| `taskexceptdomain` | 全部 | `Task.ExceptDomains` |
| `taskexceptinfo` | 全部 | `Task.ExceptInfos` |

## 6. 迁移策略

### 6.1 总体原则

1. 第一阶段只新增新表和新读写路径，不删除旧表。
2. 迁移必须幂等，多次启动不能重复污染数据。
3. 迁移失败时应保留旧设置表，服务可以回退到旧配置读取。
4. 迁移完成后写入 `app_settings`，业务读取新快照。
5. 旧 `settings` 表只读保留一个版本周期，确认稳定后再考虑删除。

### 6.2 启动流程调整

建议将当前 `database.init()` 中设置初始化流程改为：

1. 连接数据库。
2. `AutoMigrate(&AppSetting{})`，同时保留旧 `Setting{}` 迁移。
3. 尝试读取 `app_settings` 中 `key = "app"` 的新配置。
4. 如果新配置存在，校验并加载到运行时 Store。
5. 如果新配置不存在，读取旧 `settings` 表，执行旧数组到 `AppSettings` 的迁移。
6. 迁移结果校验通过后写入 `app_settings`。
7. 加载新配置到运行时 Store。
8. 如果迁移失败，记录明确错误并决定是否中止启动。建议中止启动，避免带着半损坏配置运行。

### 6.3 迁移器伪代码

```go
func LoadOrMigrateSettings(db *gorm.DB) (AppSettings, error) {
    next, ok, err := readAppSettings(db)
    if err != nil {
        return AppSettings{}, err
    }
    if ok {
        return validateAndNormalize(next)
    }

    legacy, err := readLegacySettings(db)
    if err != nil {
        return AppSettings{}, err
    }

    migrated, err := migrateLegacySettings(legacy)
    if err != nil {
        return AppSettings{}, err
    }

    migrated, err = validateAndNormalize(migrated)
    if err != nil {
        return AppSettings{}, err
    }

    if err := saveAppSettings(db, migrated); err != nil {
        return AppSettings{}, err
    }

    return migrated, nil
}
```

### 6.4 安全读取旧数组

迁移器不得直接 `legacy["txt"][3]`。应使用安全 helper：

```go
func legacyString(m map[string][]string, key string, index int, fallback string) string {
    values, ok := m[key]
    if !ok || index < 0 || index >= len(values) {
        return fallback
    }
    return values[index]
}
```

布尔和整数转换也应提供 helper：

```go
func legacyBool(m map[string][]string, key string, index int, fallback bool) bool {
    raw := legacyString(m, key, index, "")
    switch raw {
    case "true":
        return true
    case "false":
        return false
    default:
        return fallback
    }
}

func legacyInt(m map[string][]string, key string, index int, fallback int) int {
    raw := legacyString(m, key, index, "")
    n, err := strconv.Atoi(raw)
    if err != nil {
        return fallback
    }
    return n
}
```

### 6.5 默认值策略

建议将默认值从 `file/setting.json` 的数组形式迁移为具名 JSON 文件，例如 `file/settings.default.json`：

```json
{
  "version": 1,
  "features": {
    "text_enabled": false,
    "image_enabled": false,
    "web_img_enabled": false,
    "random_enabled": false
  }
}
```

短期可以保留旧 `setting.json`，但新增配置不应继续添加到 `names/edits` 数组中。迁移完成后，默认值应由 `DefaultAppSettings()` 函数或新 JSON 文件提供。

### 6.6 配置版本升级

`AppSettings.Version` 和 `AppSetting.Version` 应保持一致。后续字段变化时使用版本迁移函数：

```go
func UpgradeSettings(in AppSettings) (AppSettings, error) {
    for in.Version < CurrentSettingsVersion {
        switch in.Version {
        case 1:
            in = upgradeV1ToV2(in)
        default:
            return AppSettings{}, fmt.Errorf("unsupported settings version %d", in.Version)
        }
    }
    return validateAndNormalize(in)
}
```

规则：

1. 新增字段必须有默认值。
2. 字段改名必须保留迁移函数。
3. 字段删除前至少保留一个版本的兼容读取。
4. 禁止在没有版本迁移的情况下改变字段含义。

## 7. API 改造方案

### 7.1 保留旧 operation 名称的兼容方案

为了减少前端一次性重写成本，可以暂时保留：

```text
POST /session operation=fetchSetting
POST /session operation=editSetting
```

但 `setting_part` 返回的数据结构应逐步从 `[][]string` 切换为具名对象。为了兼容，可以新增字段：

```go
type Session struct {
    SettingPart string          `json:"setting_part"`
    SettingData [][]string      `json:"setting_data,omitempty"`      // legacy
    SettingBody json.RawMessage `json:"setting_body,omitempty"`      // new
}
```

前端迁移完成后删除 `SettingData`。

### 7.2 推荐新接口

更清晰的方案是新增独立设置接口：

```text
GET  /session/settings/:part
PUT  /session/settings/:part
```

如果不想增加路由，也可以继续用 `/session`，但内部应按 `part` 调用 typed handler：

```go
type SettingsHandler interface {
    Fetch(AppSettings) (any, error)
    Apply(AppSettings, json.RawMessage) (AppSettings, error)
}
```

### 7.3 页面 part 命名

建议统一 part 名称：

| 当前 part | 建议 part |
|---|---|
| `openai` | `provider.openai` |
| `deepseek` | `provider.deepseek` |
| `alibaba` | `provider.alibaba` |
| `otherapi` | `provider.otherapi` |
| `txt` | `feature.text` |
| `img` | `feature.image` |
| `web` | `feature.web` |
| `rand` | `feature.random` |
| `security` | `security.dashboard` |
| `contxt` | `prompt` |
| `taskBehavior` | `security.task_behavior` |
| `txtSecurity` | `security.text_prompt` |
| `imgSecurity` | `security.image_prompt` |

短期应保留旧 part 到新 part 的映射，尤其是拼写错误的 `contxt`，避免前端同步改造前接口失效。

## 8. 校验规则

保存配置必须先校验，校验成功后才能写数据库和替换内存快照。

### 8.1 通用规则

1. `CacheMinutes` 必须大于等于 0，建议最大值不超过 10080 分钟。
2. URL 字段允许为空；非空时必须是 `http` 或 `https` URL。
3. Endpoint 字段必须是 `http` 或 `https` URL。
4. Provider 名称必须在 `openai`、`deepseek`、`alibaba`、`otherapi` 范围内。
5. Image size 应符合对应供应商支持的枚举值。
6. 列表字段应去重、去空格、删除空字符串。
7. 白名单如果为空，应明确表示拒绝全部还是允许默认值，不能依赖 nil 的偶然行为。
8. 密码 hash 应为 64 位 hex SHA-256，除非后续升级密码算法。

### 8.2 敏感字段规则

1. Fetch 时默认不返回 API Key、Token、密码 hash 明文。
2. DTO 使用 `api_key_set`、`repo_token_set`、`youtube_token_set` 表示是否已配置。
3. Update 时空字符串表示保留旧值，特殊字段如 `clear_api_key=true` 才清空密钥。
4. 日志中不得打印完整设置 JSON。
5. 迁移日志不得打印敏感字段内容。

### 8.3 Prompt 和通配符规则

1. Prompt key 只能包含字母、数字、下划线、短横线。
2. 内置 key `laugh`、`poem`、`sentence` 迁移后必须存在。
3. 通配符列表需要限制最大条数和单项最大长度，避免恶意提交超大配置。

## 9. 前端迁移

### 9.1 当前问题

当前 Vue 页面大量使用：

```js
settings[0][0]
settings[0][3]
settings[1].push(value)
```

这类代码必须逐步替换为具名字段：

```js
settings.enabled
settings.cache_minutes
settings.accepted_prompt_glob.push(value)
```

### 9.2 迁移步骤

1. 后端先提供 `setting_body` 新格式，同时保留旧 `setting_data`。
2. 修改 `static/src/js/util.js` 的 `Setting()`，让新页面读取具名对象。
3. 逐个页面替换数组下标。
4. 每迁移一个页面，就在后端对应 part 上关闭旧数组返回或只保留兼容。
5. 所有页面迁移完成后删除 `[][]string` 协议。

### 9.3 页面对应 DTO

| 页面 | DTO |
|---|---|
| `Backend/OpenAI.vue` | `ProviderSettingsDTO` |
| `Backend/Deepseek.vue` | `ProviderSettingsDTO` |
| `Backend/Alibaba.vue` | `ProviderSettingsDTO` |
| `Backend/Other.vue` | `ProviderSettingsDTO` |
| `Backend/Context.vue` | `PromptSettingsDTO` |
| `Tool/Text.vue` | `TextSettingsDTO` |
| `Tool/Image.vue` | `ImageSettingsDTO` |
| `Tool/Web.vue` | `WebSettingsDTO` |
| `Tool/Rand.vue` | `RandomSettingsDTO` |
| `Security/DashSecurity.vue` | `DashboardSecurityDTO` |
| `Security/TxtSecurity.vue` | `TextPromptSecurityDTO` |
| `Security/ImgSecurity.vue` | `ImagePromptSecurityDTO` |
| `Security/TaskBehavior.vue` | `TaskBehaviorDTO` |

## 10. 分阶段实施计划

### 阶段 1：止血和兼容层

目标是不改变数据库结构也先降低风险。

1. 修复 `Deepseek.vue` endpoint 写入下标错误。
2. 为旧数组下标定义常量或 getter，禁止新增裸下标。
3. `editSetting` 增加 `SettingPart` 存在性检查和长度检查。
4. `SettingMap` 增加读写锁，或封装为只读快照。
5. 对缓存分钟数、URL、bool 做基础校验。
6. 添加设置相关单元测试，覆盖默认配置和关键下标。

### 阶段 2：新增结构化配置和迁移器

1. 新增 `AppSetting` 表。
2. 定义 `AppSettings` 结构和默认值。
3. 实现旧配置读取和迁移器。
4. 实现 `validateAndNormalize()`。
5. 启动时优先加载新配置，不存在时从旧配置迁移。
6. 保留旧 `settings` 表，不删除。

### 阶段 3：业务读取切换

1. `processor/txt.go` 切换到 `AppSettings.Text` 和 `Providers`。
2. `processor/imgGen.go` 切换到 `AppSettings.Image` 和 `Providers`。
3. `processor/webImg.go` 切换到 `AppSettings.Web`。
4. `processor/rand.go` 切换到 `AppSettings.Random`。
5. `security/function.go` 和 `security/general.go` 切换到 `Features`、`Security`、`Task`。
6. `handler/common.go` 的缓存时间读取切换到具名字段。
7. 删除业务代码中直接读取 `SettingMap` 的逻辑。

### 阶段 4：后台 API 和前端切换

1. 后端为每个 setting part 提供 DTO。
2. 前端逐页改为具名字段。
3. 供应商页面改为 API Key 不回显，只显示是否已配置。
4. 删除前端所有 `settings[x][y]` 写法。
5. 删除 `PartMap` 对数组配置组顺序的依赖。

### 阶段 5：清理旧结构

1. 确认一个版本周期内没有回滚需求。
2. 删除 `file/setting.json` 的数组默认配置，改用新默认配置。
3. 删除 `SettingMap map[string][]string`。
4. 删除旧 `fetchSetting/editSetting` 的 `[][]string` 协议。
5. 可选：保留旧 `settings` 表备份，或提供一次性清理命令。

## 11. 回滚策略

第一阶段不要删除旧 `settings` 表，因此可以提供回滚开关：

```text
URLAPI_SETTINGS_LEGACY=1
```

开启后服务跳过 `app_settings`，重新使用旧 `settings` 表和旧 `SettingMap`。该开关只用于迁移期排障，不应长期保留。

迁移前建议备份 SQLite 文件：

```text
assets/database.db -> assets/database.db.bak.<timestamp>
```

备份由运维或启动迁移流程完成，不建议每次启动都自动创建无限备份。可以只在首次从旧设置迁移到新设置前创建一次备份，并记录迁移状态。

## 12. 测试计划

### 12.1 单元测试

1. 旧默认 `file/setting.json` 可以完整迁移为 `AppSettings`。
2. 缺失配置组时使用合理默认值。
3. 配置数组长度不足时不会 panic。
4. 非法 bool 使用默认值或返回校验错误。
5. 非法 cache minutes 返回校验错误。
6. 非法 endpoint 返回校验错误。
7. 敏感字段 DTO fetch 不返回明文。
8. 空 API Key 更新不会覆盖旧值。
9. `contxt` 旧 part 可以映射到新 `prompt`。

### 12.2 集成测试

1. 使用旧数据库启动，自动生成 `app_settings`。
2. 使用已存在新数据库启动，不重复迁移。
3. 修改文字设置后，`/txt` 使用新配置。
4. 修改图片设置后，`/img` 使用新配置。
5. 修改网页设置后，`/web` 使用新配置。
6. 修改安全白名单后，安全检查使用新配置。
7. 修改 Prompt 模板后，文本生成使用新模板。

### 12.3 前端手工验证

1. 打开接口设置页，API Key 不明文回显。
2. 修改 DeepSeek endpoint 后能生效。
3. 修改文字功能开关后 `/txt` 行为正确变化。
4. 修改图片功能开关后 `/img` 行为正确变化。
5. 修改网页允许站点后对应 host 行为正确变化。
6. 修改后台密码后旧密码失效，新密码可登录。
7. 刷新页面后所有设置保持一致。

## 13. 风险点

1. 旧数组中可能已经存在脏数据，迁移器必须安全读取并报告问题。
2. API Key 不回显会改变前端交互，需要明确“留空表示不修改”。
3. 如果同时保留新旧接口，可能出现双写不一致。建议迁移期只允许新接口写新表，旧表只读。
4. 配置快照替换必须在数据库写入成功后执行，避免内存和数据库不一致。
5. 如果迁移后仍有业务代码读取旧 `SettingMap`，会导致配置修改不生效或行为不一致。
6. SQLite 备份和迁移需要考虑容器文件系统权限。

## 14. 最小可落地版本

如果希望用最少改动先完成一次可靠升级，推荐范围如下：

1. 新增 `AppSetting` 表和 `AppSettings` 结构。
2. 启动时从旧 `settings` 表迁移到 `app_settings`。
3. 业务读取切换到 `AppSettings` 快照。
4. 前端暂时仍可使用旧页面，但后端 `fetchSetting/editSetting` 内部通过 DTO 和兼容转换读写 `AppSettings`。
5. 不立即删除旧表、不立即大改路由。

该版本已经可以消除业务层魔法下标和并发 map 风险，同时把前端迁移拆到后续阶段。

## 15. 完成标准

设置数据库升级完成应满足：

1. 后端业务代码不再直接访问 `database.SettingMap[...][...]`。
2. 前端设置页面不再使用 `settings[x][y]`。
3. 数据库中存在版本化结构化配置。
4. 从旧数据库启动可以自动迁移且不丢失用户设置。
5. 敏感字段不会在普通 fetch 接口中明文返回。
6. 设置保存有字段级校验。
7. 并发请求读取配置没有 data race 风险。
8. 配置字段变更有版本迁移函数。
