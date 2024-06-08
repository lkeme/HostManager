package model

import "gorm.io/gorm"

type Server struct {
	gorm.Model        // ID, CreatedAt, UpdatedAt, DeletedAt
	ParentId   uint   `gorm:"not null" json:"parent_id"` // 父目录ID
	Key        string `gorm:"not null" json:"key"`       // 服务器唯一标识 UUID_V4{9,23}
	Name       string `gorm:"not null" json:"name"`      // 服务器名
	Host       string `gorm:"not null" json:"host"`      // 服务器地址
	Port       uint   `gorm:"not null" json:"port"`      // 服务器端口 0-65535
	Username   string `gorm:"not null" json:"username"`  // 用户名
	Password   string `gorm:"not null" json:"password"`  // 密码
	PrivateKey string `gorm:"null" json:"private_key"`   // 私钥 可选
}
