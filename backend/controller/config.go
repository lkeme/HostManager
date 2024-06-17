package controller

import (
	"HostManager/backend/config"
	"encoding/base64"
)

type Config struct {
	Base
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) SetAppSettingAttr(key string, value interface{}) {
	config.AppSetting[key] = value
}

func (c *Config) GetAppSettingAttr(key string) Res {
	type GetAppSettingAttrResponse struct {
		Value interface{} `json:"value"`
	}

	return c.success(GetAppSettingAttrResponse{Value: config.AppSetting[key]})
}

func (c *Config) SetAppIcon(i []byte) {
	c.SetAppSettingAttr("appIcon", "data:image/png;base64,"+base64.StdEncoding.EncodeToString(i))
}
