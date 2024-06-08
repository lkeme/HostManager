package main

import (
	"HostManager/backend/core"
	"embed"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

//go:embed build/windows/icon.ico
var wIcon []byte

func main() {
	// 初始化环境
	core.InitEnv()
	//使用 10086 作为ginServer的端口
	core.InitWails(assets, 10086, icon, wIcon)
}
