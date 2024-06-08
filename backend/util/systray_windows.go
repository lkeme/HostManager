//go:build windows

package util

import (
	"context"
	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	systray.Run(func() {
		//监听托盘事件
		go func() {
			systray.SetIcon(u.icon)
			systray.SetTitle("GoDesk")
			systray.SetTooltip("基于wails2实现的快速开发模板")
			//设置托盘选项
			openWindow := systray.AddMenuItem("显示主界面", "显示主界面")
			hideWindow := systray.AddMenuItem("隐藏到托盘", "隐藏到托盘")
			systray.AddSeparator()
			topWindow := systray.AddMenuItem("置顶", "置顶")
			bottomWindow := systray.AddMenuItem("取消置顶", "取消置顶")
			systray.AddSeparator()
			quitWindow := systray.AddMenuItem("退出程序", "Quit the whole app")
			for {
				select {
				case <-openWindow.ClickedCh:
					go func() {
						runtime.WindowShow(*u.ctx)
					}()
				case <-hideWindow.ClickedCh:
					go func() {
						runtime.WindowHide(*u.ctx)
					}()
				case <-topWindow.ClickedCh:
					go func() {
						runtime.WindowSetAlwaysOnTop(*u.ctx, true)
					}()
				case <-bottomWindow.ClickedCh:
					go func() {
						runtime.WindowSetAlwaysOnTop(*u.ctx, false)
					}()
				case <-quitWindow.ClickedCh:
					systray.Quit() //退出托盘
					return
				}
			}
		}()
	}, func() {
		//退出程序
		u.app.Quit()
	})
}
