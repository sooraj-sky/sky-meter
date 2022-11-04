package models

type JsonInput []struct {
		URL     string `json:"url",omitempty`
		Timeout int    `json:"timeout",omitempty`
		SkipSsl bool   `json:"skip_ssl",omitempty`
		Frequency uint64    `json:"frequency",omitempty`
		Group     string `json:"group",omitempty`
}

type AllEndpoints struct {
  ID        uint  `gorm:"primaryKey"`
  CreatedAt int64 `gorm:"autoUpdateTime"`
  UpdatedAt int64 `gorm:"autoCreateTime"`
  URL       string
  Timeout   int
  SkipSsl   bool
  Frequency uint64
  Group     string
}
