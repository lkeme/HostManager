//go:build linux || darwin

package util

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/application"
)

// SystrayUtil 托盘工具
type SystrayUtil struct {
	icon []byte
	ctx  *context.Context
	app  *application.Application
}

func NewSystrayUtil(
	icon []byte,
	ctx *context.Context,
	app *application.Application) *SystrayUtil {
	return &SystrayUtil{
		icon: icon,
		ctx:  ctx,
		app:  app,
	}
}

// Run 运行托盘
func (u *SystrayUtil) Run() {
	//linux 和 mac 不需要托盘，所以关闭窗口不需要提醒是否最小化到托盘区
}
