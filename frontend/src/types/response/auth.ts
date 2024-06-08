export interface IsRegisterResponse extends BaseResponse<{
    is_register: boolean;
}> {
}

export interface IsLoginResponse extends BaseResponse<{
    is_login: boolean;
    message: string;
}> {
}

export interface RegisterResponse extends BaseResponse<{
    status: boolean;
    message: string;
    encrypted: boolean;
}> {
}

export interface IsEncryptedResponse extends BaseResponse<{
    encrypted: boolean;
    message: string;
}> {
}

export interface LoginResponse extends BaseResponse<{
    status: boolean;
    message: string;
    encrypted: boolean;
}> {
}

export interface LogoutResponse extends BaseResponse<{
    status: boolean;
    message: string;
}> {
}