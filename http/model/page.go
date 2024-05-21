package model

type PageResult[T any] struct {
	TotalCount int64 `json:"totalCount" example:"10"` //总数
	Rows       []*T  `json:"rows"`
}

func (r *PageResult[T]) Append(item *T) {
	r.Rows = append(r.Rows, item)
}
