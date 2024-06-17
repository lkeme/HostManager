package service

import (
	"HostManager/backend/global"
	"HostManager/backend/model"
	"errors"
	"github.com/google/uuid"
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

// CreateFolder 创建文件夹
func (s *FileSystemSrv) CreateFolder(name string, parentID uint64) error {
	// 获取父文件夹
	parent, err := s.GetParent(parentID)
	log.Println(parent)
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
		Name:     name,
		ParentID: parent.ID,
		Level:    parent.Level + 1,
		Key:      uuid.New().String()[9:23],
	}
	//
	return global.DB().Create(f).Error
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

func GetFoldersInFolder(parent *model.Folder) ([]*model.Folder, error) {
	var folders []*model.Folder
	return folders,
		global.DB().
			Where("parent_id = ?", parent.ID).
			Find(&folders).
			Error
}

func GetFilesInFolder(parent *model.Folder) ([]*model.File, error) {
	var files []*model.File
	return files,
		global.DB().
			Where("parent_id = ?", parent.ID).
			Find(&files).
			Error

}

func GetChildrenInFolder(parent *model.Folder) (*model.FolderContent, error) {
	files, err := GetFilesInFolder(parent)
	if err != nil {
		return nil, err
	}
	folders, err := GetFoldersInFolder(parent)
	if err != nil {
		return nil, err
	}

	return &model.FolderContent{
		Files:  files,
		Folder: folders,
	}, nil
}

func GetChildrenRecursive(parent *model.Folder) (*model.FolderTree, error) {
	result := &model.FolderTree{Root: parent}
	content, err := GetChildrenInFolder(parent)
	if err != nil {
		return nil, err
	}

	result.Files = content.Files
	result.Folders = []*model.FolderTree{}

	for _, subDir := range content.Folder {
		t, err := GetChildrenRecursive(subDir)
		if err != nil {
			return nil, err
		}
		result.Folders = append(result.Folders, t)
	}

	return result, nil
}

func (s *FileSystemSrv) GetAllFolders() (folderHierarchy []model.FolderHierarchy) {
	// Query all folders
	var rootFolders []model.Folder
	global.DB().Where("parent_id = ?", 0).Find(&rootFolders)

	// Organize folders into a hierarchy
	for _, rootFolder := range rootFolders {
		hierarchy := s.buildFolderHierarchy(global.DB(), rootFolder)
		folderHierarchy = append(folderHierarchy, hierarchy)
	}
	return
}

// Recursive function to build folder hierarchy
func (s *FileSystemSrv) buildFolderHierarchy(db *gorm.DB, folder model.Folder) model.FolderHierarchy {
	var childrenFolders []model.Folder
	db.Where("parent_id = ?", folder.ID).Find(&childrenFolders)

	hierarchy := model.FolderHierarchy{Folder: folder}
	for _, childFolder := range childrenFolders {
		childHierarchy := s.buildFolderHierarchy(db, childFolder)
		hierarchy.Children = append(hierarchy.Children, childHierarchy)
	}
	return hierarchy
}
