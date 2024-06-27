//
export interface CreateFolderResponse extends BaseResponse<{
    status: boolean;
    message: string;
}> {
}

export interface RenameFolderResponse extends BaseResponse<{
    status: boolean;
    message: string;
}> {
}

export interface DeleteFolderResponse extends BaseResponse<{
    status: boolean;
    message: string;
}> {
}

// // File 文件结构体
// type File struct {
//     gorm.Model        // ID, CreatedAt, UpdatedAt, DeletedAt
//     Type       string `gorm:"not null" json:"type"`  // 类型
//     Name       string `gorm:"not null" json:"name"`  // 文件名称
//     Key        string `gorm:"not null" json:"key"`   // 文件唯一标识 UUID_V4{9,23}
//     Level      uint   `gorm:"not null" json:"level"` // 文件层级
//     FolderID   uint
//     Folder     Folder `gorm:"foreignkey:FolderID"` // 所属文件夹对象
// }
//
// // Folder 文件夹结构体
// type Folder struct {
//     gorm.Model          // ID, CreatedAt, UpdatedAt, DeletedAt
//     Type       string   `gorm:"not null" json:"type"`     // 类型
//     Name       string   `gorm:"non null" json:"name"`     // 文件夹名
//     Key        string   `gorm:"not null" json:"key"`      // 文件夹唯一标识 UUID_V4{9,23}
//     Level      uint     `gorm:"not null" json:"level"`    // 文件夹层级
//     ParentID   uint     `gorm:"not null" json:"parentID"` // // 父文件夹ID，根目录的父文件夹ID设为0
//     Parent     *Folder  `gorm:"foreignkey:ParentID"`      // 父文件夹对象
//     Subfolders []Folder `gorm:"foreignkey:ParentID"`      // 子文件夹列表
//     Files      []File   // 文件列表
// }
//
// // Storage 文件系统根节点结构体
// type Storage struct {
//     Root Folder `json:"root"`
// }

export interface File {
    ID: number;
    type: string;
    name: string;
    key: string;
    level: number;
    folderID: number;
}

export interface Folder {
    ID: number;
    type: string;
    name: string;
    key: string;
    level: number;
    parentID: number;
    parent: Folder | null;
    subfolders: Array<Folder> | null;
    files: Array<File> | null;
}

export interface GetStorageResponse extends BaseResponse<{
    status: boolean;
    message: string;
    storage: Array<Folder>;
}> {
}