package main

import "github.com/getlantern/systray"

type Device struct {
	Sub         bool              `json:"sub,omitempty"`
	CreateTime  int64             `json:"create_time,omitempty"`
	LocalKey    string            `json:"local_key,omitempty"`
	OwnerID     string            `json:"owner_id,omitempty"`
	BizType     int64             `json:"biz_type,omitempty"`
	IP          string            `json:"ip,omitempty"`
	Icon        string            `json:"icon,omitempty"`
	TimeZone    string            `json:"time_zone,omitempty"`
	ProductName string            `json:"product_name,omitempty"`
	UUID        string            `json:"uuid,omitempty"`
	ActiveTime  int64             `json:"active_time,omitempty"`
	UID         string            `json:"uid,omitempty"`
	UpdateTime  int64             `json:"update_time,omitempty"`
	ProductID   string            `json:"product_id,omitempty"`
	Name        string            `json:"name,omitempty"`
	Online      bool              `json:"online,omitempty"`
	Model       string            `json:"model,omitempty"`
	ID          string            `json:"id,omitempty"`
	Category    string            `json:"category,omitempty"`
	Status      []DeviceStatus    `json:"status,omitempty"`
	State       bool              `json:"state,omitempty"`
	MenuItem    *systray.MenuItem `json:"menu_item,omitempty"`
}

type DeviceStatus struct {
	Code  string `json:"code"`
	Value string `json:"value"`
}

type DevicesResponse struct {
	Devices []Device `json:"devices"`
	Total   int64    `json:"total"`
	LastID  string   `json:"last_id"`
}
