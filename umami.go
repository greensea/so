package main

import (
	"encoding/json"
	"fmt"

	"github.com/carlmjohnson/versioninfo"
	"github.com/go-resty/resty/v2"
	"github.com/greensea/so/common"
)

type Event struct {
	Hostname string `json:"hostname,omitempty"`
	Language string `json:"language,omitempty"`
	Referer  string `json:"referrer,omitempty"`
	Screen   string `json:"screen,omitempty"`
	Title    string `json:"title,omitempty"`
	Url      string `json:"url,omitempty"`
	Website  string `json:"website,omitempty"`
	Name     string `json:"name,omitempty"`
}

func Umami(cmd string) error {
	url := "https://umi.pingflash.com/api/send"
	website := "d53396b9-59bc-4e9a-9e8a-329ceba6ec65"

	body := map[string]any{
		"type": "event",
		"payload": Event{
			Website:  website,
			Hostname: "so.cli",
			Language: common.Lang(),
			Title:    cmd,
			Url:      "/",
		},
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", fmt.Sprintf("Mozilla/5.0 (so/%s)", versioninfo.Version)).
		SetBody(bodyBytes).
		Post(url)

	return nil
}
