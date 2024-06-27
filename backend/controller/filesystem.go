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

func (f *FileSystem) GetStorage() Res {
	type GetStorageResponse struct {
		Status  bool           `json:"status"`
		Message string         `json:"message"`
		Storage []model.Folder `json:"storage"`
	}
	//
	storage, err := f.ser.GetStorage()
	if err != nil && err.Error() != "record not found" {
		return f.success(GetStorageResponse{Status: false, Message: err.Error()})
	}
	return f.success(GetStorageResponse{Status: true, Message: "获取文件系统成功", Storage: storage})
}

// CreateFolder 创建目录
func (f *FileSystem) CreateFolder(name string, parentID uint64) Res {
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

// RenameFolder 重命名目录
func (f *FileSystem) RenameFolder(name string, id uint64) Res {
	type RenameFolderResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}
	//
	if id < 0 || name == "" {
		return f.success(RenameFolderResponse{Status: false, Message: "目录参数错误"})
	}
	//
	err := f.ser.RenameFolder(name, id)
	if err != nil {
		return f.success(RenameFolderResponse{Status: false, Message: err.Error()})
	}
	return f.success(RenameFolderResponse{Status: true, Message: "重命名目录成功"})
}

// DeleteFolder 删除文件夹
func (f *FileSystem) DeleteFolder(id uint64) Res {
	type DeleteFolderResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}
	//
	if id < 0 {
		return f.success(DeleteFolderResponse{Status: false, Message: "目录参数错误"})
	}
	//
	err := f.ser.DeleteFolder(id)
	if err != nil {
		return f.success(DeleteFolderResponse{Status: false, Message: err.Error()})
	}
	return f.success(DeleteFolderResponse{Status: true, Message: "删除目录成功"})
}
