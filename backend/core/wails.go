package core

import (
	c "HostManager/backend/config"
	"HostManager/backend/controller"
	"HostManager/backend/util"
	"context"
	"embed"
	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

// InitWails 初始化
func InitWails(assets embed.FS, port int, icon []byte, wIcon []byte) {
	// Create an instance of the app structure\
	system := controller.NewSystem()
	config := controller.NewConfig()
	config.SetAppIcon(icon)
	auth := controller.NewAuth()
	fs := controller.NewFileSystem()

	// 创建控制器实例
	//config := NewConfig()
	//server := NewServer()
	//wallpaper := NewWallpaper()
	//audio := NewAudio()
	//photo := NewPhoto()

	//应用实例
	app := application.NewWithOptions(&options.App{
		Title:     c.GetAppSettingAttr("appName").(string),
		Frameless: true,
		Width:     c.GetAppSettingAttr("appWidth").(int),
		Height:    c.GetAppSettingAttr("appHeight").(int),
		MinWidth:  c.GetAppSettingAttr("appMinWidth").(int),
		MinHeight: c.GetAppSettingAttr("appMinHeight").(int),
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		//BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		//启动时回调
		OnStartup: func(ctx context.Context) {
			//设置 context 对象
			system.SetCtx(ctx)
			config.SetCtx(ctx)
			auth.SetCtx(ctx)
			fs.SetCtx(ctx)
			//server.setCtx(ctx)
			//audio.setCtx(ctx)
			//photo.setCtx(ctx)
			//wallpaper.setCtx(ctx)
			//
			////启动自定义服务，初始化数据表
			//server.start(port).log("程序已启动").schema(&model.Config{})
		},
		OnShutdown: func(ctx context.Context) {
			auth.Logout() // 退出
		},
		Bind: []interface{}{
			system,
			config,
			auth,
			fs,
			//server,
			//wallpaper,
			//audio,
			//photo,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarDefault(),
			About: &mac.AboutInfo{
				Title:   c.GetAppSettingAttr("appName").(string),
				Message: c.GetAppSettingAttr("appDesc").(string),
				Icon:    icon,
			},
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false, // 透明
			WindowIsTranslucent:  false, // 半透明
			// Auto	让 Windows 决定使用哪个背景
			// None	不要使用半透明
			// Acrylic	使用 亚克力 效果
			// Mica	使用 Mica 效果
			// Tabbed	使用 Tabbed。 这是一个类似于 Mica 的背景
			//BackdropType:                        windows.Acrylic,
			BackdropType:                        windows.None,
			DisableWindowIcon:                   false,
			IsZoomControlEnabled:                false,
			ZoomFactor:                          0,
			DisablePinchZoom:                    false,
			DisableFramelessWindowDecorations:   false, // 无边框
			WebviewUserDataPath:                 "",
			WebviewBrowserPath:                  "",
			Theme:                               0,
			CustomTheme:                         nil,
			Messages:                            nil,
			ResizeDebounceMS:                    0,
			OnSuspend:                           nil,
			OnResume:                            nil,
			WebviewGpuIsDisabled:                false,
			WebviewDisableRendererCodeIntegrity: false,
			EnableSwipeGestures:                 false,
		},

		Linux: &linux.Options{
			ProgramName:         c.GetAppSettingAttr("appName").(string),
			Icon:                icon,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyOnDemand,
			WindowIsTranslucent: false,
		},
	})

	//设置托盘
	go func() {
		util.NewSystrayUtil(wIcon, system.GetCtx(), app).Run()
	}()

	// 启动wails服务
	err := app.Run()
	if err != nil {
		println("Error:", err.Error())
	}
}
