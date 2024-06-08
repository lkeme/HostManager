package controller

import (
	"HostManager/backend/internal/db"
	"HostManager/backend/util"
	"fmt"
	"os"
)

// Auth 认证API
type Auth struct {
	Base
	dbPath string
}

// NewAuth 创建Auth
func NewAuth() *Auth {
	return &Auth{
		dbPath: util.GetExecutablePath() + "data" + string(os.PathSeparator) + "data.db",
	}
}

// IsRegister 判断是否注册
func (a *Auth) IsRegister() Res {
	type IsRegisterResponse struct {
		IsRegister bool `json:"is_register"`
	}

	// 判断数据库文件是否存在
	if util.PathOrFileExist(a.dbPath) {
		return a.success(IsRegisterResponse{IsRegister: true})
	}
	return a.success(IsRegisterResponse{IsRegister: false})
}

// IsLogin 判断是否登录
func (a *Auth) IsLogin() Res {
	type IsLoginResponse struct {
		IsLogin bool   `json:"is_login"`
		Message string `json:"message"`
	}

	err := db.GetHealthStatus()
	if err != nil {
		return a.success(IsLoginResponse{IsLogin: false, Message: fmt.Sprintf("检查数据库连接失败, %v", err)})
	}
	return a.success(IsLoginResponse{IsLogin: true, Message: "检查数据库连接成功"})
}

// Register 注册
func (a *Auth) Register(password, confirmPassword string) Res {
	type RegisterResponse struct {
		Status    bool   `json:"status"`
		Message   string `json:"message"`
		Encrypted bool   `json:"encrypted"`
	}

	//
	if password != confirmPassword {
		return a.success(RegisterResponse{Status: false, Message: "注册失败, 两次密码不一致"})
	}
	//
	err := db.Connect(a.dbPath, password)
	if err != nil {
		return a.success(RegisterResponse{Status: false, Message: "注册失败, " + err.Error()})
	}
	//
	err = db.Migrate()
	if err != nil {
		return a.success(RegisterResponse{Status: false, Message: "注册失败, " + err.Error()})
	}
	//
	encrypted, err := db.IsEncrypted(a.dbPath)
	if err != nil {
		return a.success(RegisterResponse{Status: false, Message: "注册失败, " + err.Error()})
	}
	//
	return a.success(RegisterResponse{Status: true, Message: "注册成功", Encrypted: encrypted})
}

// Login 登录
func (a *Auth) Login(password string) Res {
	type LoginResponse struct {
		Status    bool   `json:"status"`
		Message   string `json:"message"`
		Encrypted bool   `json:"encrypted"`
	}
	//
	err := db.Connect(a.dbPath, password)
	if err != nil {
		return a.success(LoginResponse{Status: false, Message: "登录失败, " + err.Error()})
	}
	//
	err = db.Migrate()
	if err != nil {
		return a.success(LoginResponse{Status: false, Message: "登录失败, " + err.Error()})
	}
	//
	encrypted, err := db.IsEncrypted(a.dbPath)
	if err != nil {
		return a.success(LoginResponse{Status: false, Message: "登录失败, " + err.Error()})
	}
	//
	return a.success(LoginResponse{Status: true, Message: "登录成功", Encrypted: encrypted})
}

// IsEncrypted 判断是否加密
func (a *Auth) IsEncrypted() Res {
	type IsEncryptedResponse struct {
		Encrypted bool   `json:"encrypted"`
		Message   string `json:"message"`
	}
	//
	encrypted, err := db.IsEncrypted(a.dbPath)
	if err != nil {
		return a.success(IsEncryptedResponse{Encrypted: false, Message: "获取加密状态失败, " + err.Error()})
	}
	//
	return a.success(IsEncryptedResponse{Encrypted: encrypted, Message: "获取加密状态成功"})
}

// Logout 退出登录
func (a *Auth) Logout() Res {
	type LogoutResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}
	err := db.Close()
	if err != nil {
		return a.success(LogoutResponse{Status: false, Message: fmt.Sprintf("退出登录失败, %v", err)})
	}
	return a.success(LogoutResponse{Status: true, Message: "退出登录成功"})
}

// ChangePassword 修改密码
func (a *Auth) ChangePassword(oldPassword, newPassword, confirmPassword string) Res {
	type ChangePasswordResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}

	//
	if newPassword != confirmPassword {
		return a.success(ChangePasswordResponse{Status: false, Message: "修改密码失败, 两次密码不一致"})
	}
	//
	err := db.ChangePassword(a.dbPath, oldPassword, newPassword)
	if err != nil {
		return a.success(ChangePasswordResponse{Status: false, Message: "修改密码失败, " + err.Error()})
	}
	return a.success(ChangePasswordResponse{Status: true, Message: "修改密码成功"})
}
