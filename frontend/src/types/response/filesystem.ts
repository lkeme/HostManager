//
export interface CreateFolderResponse extends BaseResponse<{
    status: boolean;
    message: string;
}> {
}

export interface Folder{
    ID: number;
    name: string;
    key: string;
    level: number;
    parentID: number;
}

export interface FolderHierarchy{
    folder: Folder;
    children: Array<FolderHierarchy>;
}

export interface GetAllFoldersResponse extends BaseResponse<{
    status: boolean;
    message: string;
    folderHierarchy: Array<FolderHierarchy>;
}> {
}