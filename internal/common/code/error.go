package code

type ErrCode int32

const (
	E0000 ErrCode = 00000

	// about user
	E1001 ErrCode = 10001
	E1002 ErrCode = 10002
	E1003 ErrCode = 10003
	E1004 ErrCode = 10004
	E1005 ErrCode = 10005

	// about act
	E2001 ErrCode = 20001

	// about actType
	E3001 ErrCode = 30001
	E3002 ErrCode = 30002
	E3003 ErrCode = 30003

	// about comment
	E4001 ErrCode = 40001
	E4002 ErrCode = 40002

	// about auth
	E5001 ErrCode = 50001
	E5002 ErrCode = 50002

	// page not found
	E6001 ErrCode = 60001
	E6002 ErrCode = 60002

	E9999 ErrCode = 99999
)

var ErrCodeMap = map[ErrCode]string{
	E0000: "成功",

	E2001: "找不到对应的活动类型",

	E1001: "用户名冲突",
	E1002: "找不到对应用户",
	E1003: "用户密码错误",
	E1004: "用户已加入指定活动",
	E1005: "用户未加入指定活动",

	E3001: "找不到对应活动类型",
	E3002: "找不到对应的父类型",
	E3003: "该类型下存在活动，不允许删除",

	E4001: "找不到对应的活动实例",
	E4002: "找不到对应的用户实例",

	E5001: "未授权",
	E5002: "权限不足",

	E6001: "页面不存在",
	E6002: "参数格式化失败",

	E9999: "未知异常",
}

func GetErrMsg(code ErrCode) string {
	msg, ok := ErrCodeMap[code]
	if ok {
		return msg
	}
	return ErrCodeMap[E0000]
}
