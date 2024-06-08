package config

var AppSetting map[string]interface{}

func GetAppSettingAttr(key string) interface{} {
	return AppSetting[key]
}

func init() {
	AppSetting = make(map[string]interface{})
	// AppVersion 应用版本
	AppSetting["appVersion"] = "0.1.0"
	// AppName 应用名称
	AppSetting["appName"] = "HostManager"
	// AppDescription 应用描述
	AppSetting["appDesc"] = "by Lkeme"
	// AppWidth 应用宽度
	AppSetting["appWidth"] = 1100
	// AppHeight 应用高度
	AppSetting["appHeight"] = 700
	// AppMinWidth 应用最小宽度
	AppSetting["appMinWidth"] = 1000
	// AppMinHeight 应用最小高度
	AppSetting["appMinHeight"] = 600
	// AppIcon 应用图标
	AppSetting["appIcon"] = ""
}
