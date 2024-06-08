package core

import "HostManager/backend/util"

func InitEnv() {
	// 创建data文件夹
	util.CreateFolder(util.GetExecutablePath() + "/data")
	// 创建logs文件夹
	util.CreateFolder(util.GetExecutablePath() + "/logs")
}
