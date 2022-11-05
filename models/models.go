package models

import (
	"gorm.io/datatypes"
	"net"
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
	OutputData datatypes.JSON `json:"attributes" gorm:"type:json"`
	URL        string
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
