namespace go auth

struct AuthAccount {
    1: i64 uid; // 全平台用户唯一id
    2: string username; // 用户名
    3: string password; // 密码
    4: string createIp; // 注册ip
    5: i8 status; // 状态 1:启用 0:禁用 -1:删除
    6: i8 sysType; // 系统类型 0:普通用户系统 1:商家系统
    7: i64 userId; // 用户id
    8: i64 tenantId; // 所属租户id
    9: bool isAdmin; // 是否是管理员
}

struct GetByUserIdAndTypeReq {
    1: i64 userId; // 用户id
    2: i8 sysType; // 系统类型 0:普通用户系统 1:商家系统
}

struct GetByUserIdAndTypeResp {
    1: AuthAccount authAccount;
}

struct GetByUidReq {
    1: i64 uid; // 全平台用户唯一id
}

struct GetByUidResp {
    1: AuthAccount authAccount;
}

struct UpdatePasswordReq {
    1: i64 uid; // 全平台用户唯一id
    2: i8 sysType; // 系统类型 0:普通用户系统 1:商家系统
    3: string password; // 密码
}

struct UpdatePasswordResp {
    1: bool success;
}

struct GetAccountByUserNameReq {
    1: string username; // 用户名
}

struct GetAccountByUserNameResp {
    1: AuthAccount authAccount;
}

service AuthService {
    GetByUserIdAndTypeResp getByUserIdAndType(1: GetByUserIdAndTypeReq req);
    GetByUidResp getByUid(1: GetByUidReq req);
    UpdatePasswordResp updatePassword(1: UpdatePasswordReq req);
    GetAccountByUserNameResp getAccountByUserName(1: GetAccountByUserNameReq req);
}