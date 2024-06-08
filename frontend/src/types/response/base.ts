// 定义基础类型接口
interface BaseResponse<T> {
    code: number;
    msg: string;
    data: T;
}