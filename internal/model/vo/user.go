package vo

import "entrytask/internal/model"

type UserLoginVo struct {
	Id        uint64 `json:"id"`
	Aliasname string `json:"aliasname"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Headpic   string `json:"headpic"`
	Role      int    `json:"role"`
}

func FormatUserLoginVo(obj *model.User) *UserLoginVo {
	return &UserLoginVo{
		obj.Id,
		obj.Aliasname,
		obj.Username,
		obj.Email,
		obj.Headpic,
		obj.Role,
	}
}
