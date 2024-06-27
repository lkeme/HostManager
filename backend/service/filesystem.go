package service

import (
	"HostManager/backend/global"
	"HostManager/backend/model"
	"HostManager/backend/util"
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
)

var FileSystemSrvIns *FileSystemSrv
var FileSystemSrvOnce sync.Once

type FileSystemSrv struct{}

func GetFileSystemSrv() *FileSystemSrv {
	FileSystemSrvOnce.Do(func() {
		FileSystemSrvIns = &FileSystemSrv{}
	})
	return FileSystemSrvIns
}

// GetParent 获取父文件夹
func (s *FileSystemSrv) GetParent(parentID uint64) (*model.Folder, error) {
	return s.GetFolderByID(parentID)
}

func preloadSubfolders(db *gorm.DB) *gorm.DB {
	//return d.Preload(clause.Associations, preload)
	return db.Preload("Subfolders", preloadSubfolders).Preload("Files")
}

// GetStorage 获取文件系统
func (s *FileSystemSrv) GetStorage() (rootFolders []model.Folder, err error) {
	err = global.DB().
		Where("parent_id = ?", 0).
		Preload("Subfolders", preloadSubfolders).
		Preload("Files").
		Find(&rootFolders).Error
	return
}

// RenameFolder 重命名文件夹
func (s *FileSystemSrv) RenameFolder(name string, id uint64) error {
	folder, err := s.GetFolderByID(id)
	if err != nil {
		return err
	}
	folder.Name = name
	return global.DB().Save(&folder).Error
}

// DeleteFolder 删除文件夹
func (s *FileSystemSrv) DeleteFolder(id uint64) error {
	folder, err := s.GetFolderByID(id)
	if err != nil {
		return err
	}
	return global.DB().Delete(&folder).Error
}

// CreateFolder 创建文件夹
func (s *FileSystemSrv) CreateFolder(name string, parentID uint64) error {
	// 获取父文件夹
	parent, err := s.GetParent(parentID)
	if err != nil {
		return err
	}
	// 获取父文件夹中的文件夹
	folders, err := s.GetFoldersInFolder(parent)
	if err != nil {
		return err
	}
	// 检查是否已经有同名文件夹
	for _, folder := range folders {
		if folder.Name == name {
			return errors.New("同级文件夹下已存在同名文件夹")
		}
	}
	// 构建路径并注册到数据库
	f := &model.Folder{
		Type:     "folder",
		Name:     name,
		Key:      util.GenerateKey(),
		Level:    parent.Level + 1,
		ParentID: parent.ID,
	}
	//
	log.Println(44444)

	return global.DB().Create(f).Error
}

// GetFolderByID 通过 ID 获取文件夹
func (s *FileSystemSrv) GetFolderByID(id uint64) (*model.Folder, error) {
	folder := model.Folder{}
	err := global.DB().Take(&folder, id).Error
	if err != nil {
		if err.Error() == "record not found" {
			folder.Name = "root"
			folder.Key = "root"
			folder.Level = 0
			folder.ParentID = 0
			log.Println("文件夹不存在，默认根目录")
			return &folder, nil
		}
	}
	return &folder, err
}

// GetFoldersInFolder 获取文件夹中的文件夹
func (s *FileSystemSrv) GetFoldersInFolder(parent *model.Folder) ([]*model.Folder, error) {
	var folders []*model.Folder
	return folders,
		global.DB().
			Where("parent_id = ?", parent.ID).
			Find(&folders).
			Error
}

//---------------------------------------------------

// setSubfolders 递归设置子文件夹属性
//func (s *FileSystemSrv) setSubfolders(folder *model.Folder, folderMap map[uint][]model.Folder) {
//	subfolders := folderMap[folder.ID]
//	folder.Children = subfolders
//
//	for i := range subfolders {
//		s.setSubfolders(&subfolders[i], folderMap)
//	}
//}

// GetFolders 获取文件夹
//func (s *FileSystemSrv) GetFolders() ([]model.Folder, error) {
//	var folders []model.Folder
//	err := global.DB().Find(&folders).Error
//	if err != nil {
//		return nil, err
//	}
//	// 使用 map 来存储父文件夹 ID 和子文件夹的映射关系
//	folderMap := make(map[uint][]model.Folder)
//	// 遍历所有文件夹，将它们按照父文件夹 ID 进行分组
//	for _, folder := range folders {
//		folderMap[folder.ParentID] = append(folderMap[folder.ParentID], folder)
//	}
//	// 获取根文件夹，即父文件夹 ID 为 0 的文件夹
//	rootFolders := folderMap[0]
//
//	// 递归设置子文件夹属性
//	for _, rootFolder := range rootFolders {
//		s.setSubfolders(&rootFolder, folderMap)
//	}
//
//	return rootFolders, nil
//}
