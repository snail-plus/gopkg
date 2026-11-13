package model

type PageRequest struct {
	PageSize int `form:"pageSize" example:"50" binding:"required"` //pageSize
	PageNum  int `form:"pageNum" example:"1" binding:"required"`   //pageNum
	Params   map[string]interface{}
}

func (r *PageRequest) AddParam(key string, value interface{}) {
	if r.Params == nil {
		r.Params = make(map[string]interface{})
	}
	r.Params[key] = value
}

func (r *PageRequest) Offset() int {
	// 防御非法分页参数，避免生成负数 offset 导致 SQL 报错
	pageNum := r.PageNum
	if pageNum < 1 {
		pageNum = 1
	}
	return (pageNum - 1) * r.PageSize
}
