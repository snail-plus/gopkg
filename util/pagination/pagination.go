// Copyright 2024 eve.  All rights reserved.

package pagination

func GetPageOffset(pageNum, pageSize int64) int64 {
	return (pageNum - 1) * pageSize
}
