package middleware

import "time"

type General struct {
	Referer string    `json:"referer"` //Complete Referer
	IP      string    `json:"ip"`
	Type    string    `json:"type"` // 任务类型
	Target  string    `json:"target"`
	Time    time.Time `json:"time"`
	Unsafe  bool      `json:"unsafe"`
	SkipDB  bool      `json:"skip_db"`
	Info    string    `json:"info"`
}

type FrequencyFilter struct {
	Type string `json:"type"`
	IP   string `json:"ip"`
}

type FrequencyData struct {
	Counter int       `json:"counter"`
	Time    time.Time `json:"time"`
}
