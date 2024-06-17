package test

//const TableNameFolder = "folder"
//
//func (*Folder) TableName() string {
//	return TableNameFolder
//}

//const (
//	FolderKeyPrefix = "-"
//)
//var FolderAlreadyExisted = errors.New("该目录已经存在")
//type FolderStore interface {
//	// 创建一个目录
//	CreateFolder(folder *Folder) (err error)
//	// 目录是否存在
//	ExistFolder(userId, parentId int64, folderName string) (isExist bool)
//	// 当 id != 0 则表示加载指定目录, 当 id == 0 则表示加载根目录
//	LoadFolder(id, userId int64, isLoadRelated bool) (folder *Folder, err error)
//	// 只加载目录和下面的文件
//	LoadSimpleFolder(id, userId int64) (folder *SimpleFolder, err error)
//	// 删除指定目录
//	DeleteFolder(ids []int64, userId int64) (allowDelFileHashList []string, err error)
//	// 移动目录
//	MoveFolder(to *Folder, ids []int64) (err error)
//	// 复制目录
//	CopyFolder(to *Folder, foders []*Folder) (totalSize uint64, err error)
//	// 重命名目录
//	RenameFolder(id, currentFolderId int64, newName string) (err error)
//	// 目录列表
//	ListFolder(folderIds []int64, userId int64) (folder []*Folder, err error)
//}
//
//type FolderService interface {
//	FolderStore
//}
