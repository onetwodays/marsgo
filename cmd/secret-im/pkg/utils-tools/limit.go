package utils

// 最小整数
const IntMin = ^IntMax

// 最大整数
const IntMax = 2147483647

// 最小无符号整数
const UintMin uint = 0

// 最大无符号整数
const UintMax = ^uint(0)

// 最大Int64整数
const Int64Max = int64(^uint64(0) >> 1)
