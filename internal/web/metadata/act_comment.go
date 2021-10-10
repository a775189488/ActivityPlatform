package metadata

import "entrytask/internal/model"

type GetActCommentResp struct {
	Id       uint64 `json:"id"`
	Message  string `json:"message"`
	Activity uint64 `json:"activity"`
	User     uint64 `json:"user"`
}

type ActCommentCreateReq struct {
	ActivityId int    `json:"activity_id" binding:"required"`
	Message    string `json:"message" binding:"required"`
}

type ActCommentUpdateReq struct {
	Message string `json:"message" binding:"required"`
}

func (c *ActCommentUpdateReq) ToActComment() *model.ActivityComment {
	return &model.ActivityComment{Message: c.Message}
}

func (c *ActCommentCreateReq) ToActComment(userId uint64) *model.ActivityComment {
	return &model.ActivityComment{ActId: uint64(c.ActivityId), UserId: userId, Message: c.Message}
}

type ListActCommentReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type ListActCommentResp struct {
	ListBase
	Data []*GetActCommentResp `json:"data"`
}

func FormatActCommentCreateResp(ac *model.ActivityComment) *GetActCommentResp {
	return &GetActCommentResp{
		ac.Id,
		ac.Message,
		ac.ActId,
		ac.UserId,
	}
}

func FormatActCommentUpdateResp(ac *model.ActivityComment) *GetActCommentResp {
	return &GetActCommentResp{
		ac.Id,
		ac.Message,
		ac.ActId,
		ac.UserId,
	}
}

func FormatActCommentListResp(acs []*model.ActivityComment, total int, page int, size int) *ListActCommentResp {
	var data []*GetActCommentResp
	for _, a := range acs {
		data = append(data, &GetActCommentResp{a.Id, a.Message, a.ActId, a.UserId})
	}
	return &ListActCommentResp{
		ListBase{
			page,
			size,
		},
		data,
	}
}
