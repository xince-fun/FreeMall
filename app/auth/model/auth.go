package model

type AuthAccount struct {
	Uid      int64  // 全平台用户唯一id
	UserName string // 用户名
	Password string // 密码
	CreateIp string // 注册ip
	Status   int8   // 状态 1:启用 0:禁用 -1:删除
	SysType  int8   // 系统类型 0:普通用户系统 1:商家系统
	UserId   int64  // 用户id
	TenantId int64  // 所属租户id
	IsAdmin  bool   // 是否是管理员
}
