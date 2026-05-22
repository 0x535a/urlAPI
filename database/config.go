package database

import (
	"gorm.io/gorm"
	"time"
)

var (
	dbPath    = "assets/database.db"
	db        *gorm.DB
	PromptMap = map[string]int{
		"laugh":    0,
		"poem":     1,
		"sentence": 2,
	}
	RepoMap    = make(map[string][]string)
	SessionMap = make(map[string]Session)
)

type Repo struct {
	UUID    string `json:"uuid" gorm:"primaryKey"`
	API     string `json:"api"`
	Info    string `json:"info"`
	Content string `json:"content"`
}

type Session struct {
	Token  string    `json:"token" gorm:"primaryKey"`
	Expire time.Time `json:"expire"`
	Term   bool      `json:"term"`
}

type Task struct {
	// all
	UUID     string    `json:"uuid" gorm:"primaryKey"`
	Time     time.Time `json:"time"`
	IP       string    `json:"ip"`
	Type     string    `json:"type"`
	Status   string    `json:"status"`
	Target   string    `json:"target"`
	Return   string    `json:"return"`
	Region   string    `json:"region"`
	Referer  string    `json:"referer"`
	Device   string    `json:"device"`
	MoreInfo string    `json:"more_info" gorm:"more_info"`
	//txt, img, web, rand
	API string `json:"api"`

	// txt, img
	Model string `json:"model"`
	Temp  string `json:"temp"`

	// img
	Size string `json:"size"`
}

type AppSetting struct {
	Key   string `json:"key" gorm:"primaryKey"`
	Value string `json:"value"`
	// Keep timestamps for existing databases created by the previous design.
	CreatedAt time.Time `json:"created_at" gorm:"-"`
	UpdatedAt time.Time `json:"updated_at" gorm:"-"`
	Version   int       `json:"version" gorm:"-"`
}

type Provider struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"unique;not null"`
	APIKeyEnc    string    `json:"api_key_enc"`
	TextModel    string    `json:"text_model"`
	SummaryModel string    `json:"summary_model"`
	ImageModel   *string   `json:"image_model"`
	ImageSize    *string   `json:"image_size"`
	Endpoint     string    `json:"endpoint"`
	Enabled      bool      `json:"enabled" gorm:"default:1"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ServiceConfig struct {
	Service          string `json:"service" gorm:"primaryKey"`
	CacheMinutes     int    `json:"cache_minutes"`
	FallbackImageURL string `json:"fallback_image_url"`
	Settings         string `json:"settings" gorm:"type:json"`
}

type Prompt struct {
	Key      string `json:"key" gorm:"primaryKey"`
	Template string `json:"template" gorm:"not null"`
}

func (Prompt) TableName() string {
	return "prompts"
}

type ConfigListItem struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Scope     string `json:"scope" gorm:"not null;uniqueIndex:idx_config_list_scope_value"`
	Value     string `json:"value" gorm:"not null;uniqueIndex:idx_config_list_scope_value"`
	SortOrder int    `json:"sort_order" gorm:"default:0"`
}

type DBList struct {
	RepoList       []Repo
	TaskList       []Task
	SessionList    []Session
	AppSettingList []AppSetting
}
