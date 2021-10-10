package metadata

import "entrytask/internal/model"

type ListActReq struct {
	Page    int32   `form:"page" binding:"required"`
	Size    int32   `form:"size" binding:"required"`
	ActType *int    `form:"act_type"`
	BeginAt *uint64 `form:"begin_at"`
	EndAt   *uint64 `form:"end_at"`
}

type ListActUserReq struct {
	Page int32 `form:"page" binding:"required"`
	Size int32 `form:"size" binding:"required"`
}

func (l *ListActReq) ToWhere() map[string]interface{} {
	result := make(map[string]interface{})
	if l.ActType != nil {
		result["act_type"] = *l.ActType
	}
	if l.BeginAt != nil {
		result["begin_at"] = *l.BeginAt
	}
	if l.EndAt != nil {
		result["end_at"] = *l.EndAt
	}
	return result
}

type ActCreateReq struct {
	Title       string `json:"title" binding:"required"`
	BeginAt     uint64 `json:"begin_at" binding:"required"`
	EndAt       uint64 `json:"end_at" binding:"required"`
	Description string `json:"description" binding:"required"`
	ActType     uint64 `json:"act_type" binding:"required"`
	Address     string `json:"address" binding:"required"`
}

func (a *ActCreateReq) ToActivity(userId uint64) *model.Activity {
	return &model.Activity{
		Title:       a.Title,
		BeginAt:     a.BeginAt,
		EndAt:       a.EndAt,
		Description: a.Description,
		ActType:     a.ActType,
		Address:     a.Address,
		Creator:     userId,
	}
}

type ActUpdateReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        uint64 `json:"type"`
	BeginAt     uint64 `json:"begin_at"`
	EndAt       uint64 `json:"end_at"`
	Address     string `json:"address"`
}

func (a *ActUpdateReq) ToActivity() *model.Activity {
	return &model.Activity{
		Title:       a.Title,
		Description: a.Description,
		ActType:     a.Type,
		BeginAt:     a.BeginAt,
		EndAt:       a.EndAt,
		Address:     a.Address,
	}
}

type GetActResp struct {
	Id          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	BeginAt     uint64 `json:"begin_at"`
	EndAt       uint64 `json:"end_at"`
	Type        uint64 `json:"type"`
	Address     string `json:"address"`
}

type GetActWithJoinResp struct {
	Id          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	BeginAt     uint64 `json:"begin_at"`
	EndAt       uint64 `json:"end_at"`
	Type        uint64 `json:"type"`
	Address     string `json:"address"`
	IsJoin      bool   `json:"is_join"`
}

type GetActDetailResp struct {
	Id          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	BeginAt     uint64 `json:"begin_at"`
	EndAt       uint64 `json:"end_at"`
	TypeId      uint64 `json:"type_id"`
	TypeName    string `json:"type_name"`
	IsJoin      bool   `json:"is_join"`
	Address     string `json:"address"`
}

type ListActResp struct {
	ListBase
	Data []*GetActWithJoinResp `json:"data"`
}

func FormatActCreateResp(act *model.Activity) *GetActResp {
	return &GetActResp{
		Id:          act.Id,
		Title:       act.Title,
		BeginAt:     act.BeginAt,
		EndAt:       act.EndAt,
		Description: act.Description,
		Type:        act.ActType,
		Address:     act.Address,
	}
}

func FormatActUpdateResp(act *model.Activity) *GetActResp {
	return &GetActResp{
		Id:          act.Id,
		Title:       act.Title,
		BeginAt:     act.BeginAt,
		EndAt:       act.EndAt,
		Description: act.Description,
		Type:        act.ActType,
		Address:     act.Address,
	}
}

func FormatGetActDetailResp(act *model.ActivityDetail, isIn bool) *GetActDetailResp {
	return &GetActDetailResp{
		Id:          act.Id,
		Title:       act.Title,
		BeginAt:     act.BeginAt,
		EndAt:       act.EndAt,
		Description: act.Description,
		TypeName:    act.ActTypeName,
		TypeId:      act.ActType,
		Address:     act.Address,
		IsJoin:      isIn,
	}
}

func FormatListActResp(acts []*model.Activity, page int, size int, joinActs []uint64) *ListActResp {
	var data []*GetActWithJoinResp
	actIdMap := make(map[uint64]struct{})
	if len(joinActs) != 0 {
		for _, actId := range joinActs {
			actIdMap[actId] = struct{}{}
		}
	}
	for _, a := range acts {
		_, isJoin := actIdMap[a.Id]
		data = append(data, &GetActWithJoinResp{
			Id:          a.Id,
			Title:       a.Title,
			BeginAt:     a.BeginAt,
			EndAt:       a.EndAt,
			Description: a.Description,
			Type:        a.ActType,
			IsJoin:      isJoin,
			Address:     a.Address,
		})
	}
	return &ListActResp{
		ListBase{
			page,
			size,
		},
		data,
	}
}

func FormateGetActivityUser(users []*model.User, page int, size int) *ListUserResp {
	data := make([]*GetUserResp, 0)
	for _, u := range users {
		data = append(data, &GetUserResp{
			u.Id,
			u.Username,
			u.Email,
			u.Aliasname,
			u.Headpic,
		})
	}
	return &ListUserResp{
		ListBase{
			page,
			size,
		},
		data,
	}
}
