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

func MergeUint16ToUint32(high, low uint16) uint32 {
	return (uint32(high) << 16) | uint32(low)
}

// SplitUint16 将 uint16 拆分成两个 uint8，高位在前，低位在后
func SplitUint16(value uint16) (uint8, uint8) {
	high := uint8(value >> 8) // 右移 8 位获取高 8 位
	low := uint8(value)       // 直接转换获取低 8 位
	return high, low
}

// SplitUint32 将 uint32 拆分成两个 uint16，高位在前，低位在后
func SplitUint32(value uint32) (uint16, uint16) {
	high := uint16(value >> 16) // 右移 16 位获取高 16 位
	low := uint16(value)        // 直接转换获取低 16 位
	return high, low
}

// ExtractUint16Bit 从 uint16 提取指定 bit 的值，返回一个 []uint8 数组。
// 每个元素是 0 或 1，对应传入的 bit 位置 (0 ~ 15)，例如: ExtractBitsToArray(0, 1, 2, 3)
func ExtractUint16Bit(value uint16, bit int) uint8 {
	return uint8((value >> bit) & 1)
}

// ExtractUint16BitsToArray 从 uint16 提取指定 bit 的值，返回一个 []uint8 数组。
// 每个元素是 0 或 1，对应传入的 bit 位置 (0 ~ 15)，例如: ExtractBitsToArray(0, 1, 2, 3)
func ExtractUint16BitsToArray(value uint16, bits ...int) []uint8 {
	var result []uint8
	for _, bit := range bits {
		if bit < 0 || bit > 15 {
			continue // 忽略非法 bit 位
		}
		// 提取对应 bit 的值，并追加到结果切片中
		result = append(result, uint8((value>>bit)&1))
	}
	return result
}

// ExtractUint16Bits 从 uint16 提取指定 bit 的值，组合成一个新的 uint16。
// 支持传入多个 bit 位置 (0 ~ 15)，例如: ExtractBits(0, 1, 2, 3)
func ExtractUint16Bits(value uint16, bits ...int) uint16 {
	var result uint16
	for _, bit := range bits {
		if bit < 0 || bit > 15 {
			continue // 忽略非法 bit 位
		}
		// 提取对应 bit 的值，并左移至结果中的正确位置
		result = (result << 1) | ((value >> bit) & 1)
	}
	return result
}
