package model

import (
	"gorm.io/gorm"
)

type Folder struct {
	gorm.Model        // ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string `gorm:"non null" json:"name"`     // 文件夹名
	Key        string `gorm:"not null" json:"key"`      // 文件夹唯一标识 UUID_V4{9,23}
	Level      uint64 `gorm:"not null" json:"level"`    // 文件夹层级`0`为根目录
	ParentID   uint   `gorm:"not null" json:"parentID"` // 父目录ID `0`为根目录
}

type FolderTree struct {
	Root    *Folder       `json:"root"`
	Files   []*File       `json:"files"`
	Folders []*FolderTree `json:"folders"`
}

type FolderContent struct {
	Files  []*File   `gorm:"foreignKey:ID" json:"files"`
	Folder []*Folder `gorm:"foreignKey:ID" json:"folders"`
}

type FolderHierarchy struct {
	Folder   Folder            `json:"folder"`
	Children []FolderHierarchy `json:"children"`
}
