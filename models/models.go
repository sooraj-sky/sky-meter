package models

import (
	"net"
	"time"

	"gorm.io/datatypes"
)

type JsonInput []struct {
	URL       string `json:"url",omitempty`
	Timeout   int    `json:"timeout",omitempty`
	SkipSsl   bool   `json:"skip_ssl",omitempty`
	Frequency uint64 `json:"frequency",omitempty`
	Group     string `json:"group",omitempty`
}

type HttpOutput struct {
	ID         uint           `gorm:"primaryKey"`
	CreatedAt  int64          `gorm:"autoUpdateTime"`
	OutputData datatypes.JSON `json:"attributes" gorm:"type:json"`
	URL        string
	StatusCode int
	Timeout    bool
	Error      string
}

type OpsgenieAlertData struct {
	ID        uint  `gorm:"primaryKey"`
	CreatedAt int64 `gorm:"autoUpdateTime"`
	URL       string
	RequestId string
	Error     string
	Active    bool
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
	NextRun   int
	Active    bool
}

type Debug struct {
	DNS struct {
		Start   string       `json:"start"`
		End     string       `json:"end"`
		Host    string       `json:"host"`
		Address []net.IPAddr `json:"address"`
		Error   error        `json:"error"`
	} `json:"dns"`
	Dial struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"dial"`
	Connection struct {
		Time string `json:"time"`
	} `json:"connection"`
	WroteAllRequestHeaders struct {
		Time string `json:"time"`
	} `json:"wrote_all_request_header"`
	WroteAllRequest struct {
		Time string `json:"time"`
	} `json:"wrote_all_request"`
	FirstReceivedResponseByte struct {
		Time string `json:"time"`
	} `json:"first_received_response_byte"`
}

type OpsGenieAlertStatus struct {
	Data struct {
		Success       bool      `json:"success"`
		Action        string    `json:"action"`
		ProcessedAt   time.Time `json:"processedAt"`
		IntegrationID string    `json:"integrationId"`
		IsSuccess     bool      `json:"isSuccess"`
		Status        string    `json:"status"`
		AlertID       string    `json:"alertId"`
		Alias         string    `json:"alias"`
	} `json:"data"`
	Took      float64 `json:"took"`
	RequestID string  `json:"requestId"`
}

type SmtpErr struct {
	URL      string
	Subject  string
	Downtime time.Time
	Reason   string
	Mailto   []string
}

type UserInput struct {
	Opegenie struct {
		Enabled bool
	}
	Email struct {
		Enabled bool
	}

	Groups []struct {
		Name   string
		Emails []string
	}
	Domains []struct {
		Name      string
		Enabled   bool
		Timeout   int
		SkipSsl   bool
		Frequency uint64
		Group     string
	}
}

type AlertGroups struct {
	ID        uint  `gorm:"primaryKey"`
	CreatedAt int64 `gorm:"autoUpdateTime"`
	Name      string
	Email     string
}

type AllEnvs struct {
	DnsServer      string
	EmailPass      string
	EmailFrom      string
	EmailPort      string
	EmailServer    string
	OpsgenieSecret string
	SentryDsn      string
	Mode           string
	DbUrl          string
}
