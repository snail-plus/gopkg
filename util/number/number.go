package number

import (
	"fmt"
	"strconv"
)

type Float interface {
	float32 | float64
}

func FormatFloat[T Float](f T, scale int) float64 {
	formatStr := "%." + strconv.FormatInt(int64(scale), 10) + "f"
	fString := fmt.Sprintf(formatStr, f)
	result, _ := strconv.ParseFloat(fString, 64)
	return result
}
