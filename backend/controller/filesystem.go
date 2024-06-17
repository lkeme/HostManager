package controller

import (
	"HostManager/backend/model"
	"HostManager/backend/service"
)

// FileSystem 文件系统API
type FileSystem struct {
	Base
	ser *service.FileSystemSrv
}

// NewFileSystem 创建FileSystem
func NewFileSystem() *FileSystem {
	return &FileSystem{
		ser: service.GetFileSystemSrv(),
	}
}

// GetAllFolders 获取所有文件夹
func (f *FileSystem) GetAllFolders() Res {
	type GetAllFoldersResponse struct {
		Status          bool                    `json:"status"`
		Message         string                  `json:"message"`
		FolderHierarchy []model.FolderHierarchy `json:"folderHierarchy"`
	}
	//
	results := f.ser.GetAllFolders()
	return f.success(GetAllFoldersResponse{Status: true, Message: "获取文件夹成功", FolderHierarchy: results})
}

// CreateFolder 创建目录
func (f *FileSystem) CreateFolder(name string, parentID uint64, key string) Res {
	type CreateFolderResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}
	//
	if name == "" || parentID < 0 {
		return f.success(CreateFolderResponse{Status: false, Message: "目录参数错误"})
	}
	//
	err := f.ser.CreateFolder(name, parentID)
	if err != nil {
		return f.success(CreateFolderResponse{Status: false, Message: err.Error()})
	}
	return f.success(CreateFolderResponse{Status: true, Message: "创建目录成功"})
}
