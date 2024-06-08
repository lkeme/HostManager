package model

import (
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Filename string `gorm:"type:char(255); NOT NULL" json:"filename"`    // 文件名称
	Hash     string `gorm:"type:varchar(32);INDEX;NOT NULL" json:"hash"` // 文件Hash
	Format   string `gorm:"type:varchar(255);NOT NULL" json:"format"`    // 文件MimeType 例如: video/mp4 -> .mp4
	Extra    string `gorm:"NOT NULL;type:TEXT" json:"extra"`             // extra
	Size     int64  `gorm:"type:BIGINT" json:"size"`                     // 文件大小
}

type FileStore interface {
	// 保存文件到指定目录
	SaveFileToFolder(file *File, folderId int64) (err error)
	// 删除文件和目录之间的关联 返回允许删除的文件 Hash 列表
	DeleteFile(ids []int64, folderId int64) (allowDelFileHashList []string, err error)
	// 移动文件
	MoveFile(fromId, toId int64, fileIds []int64) (err error)
	// 复制文件
	CopyFile(fromId, toId int64, fileIds []int64) (totalSize uint64, err error)
	// 重命名文件
	RenameFile(folderId, fileId int64, newName string) (err error)
	// 加载文件
	LoadFile(folderId, fileId, userId int64) (file *File, err error)
	//GetHashByFileIds()
}

type FileService interface {
	FileStore
}

const TableNameServer = "server"

func (*File) TableName() string {
	return TableNameServer
}
