package utils

import (
	"errors"
	"reflect"
)

// 反转切片
func ReverseSlice(slice interface{}) error {
	value := reflect.ValueOf(slice)
	if value.Kind() != reflect.Slice {
		return errors.New("non slice type")
	}

	size := value.Len()
	mid := size / 2
	if value.Len()%2 != 0 {
		mid++
	}

	for i := size - 1; i >= mid; i-- {
		head := size - i - 1
		temp := reflect.New(value.Index(i).Type()).Elem()
		temp.Set(value.Index(i))

		value.Index(i).Set(value.Index(head))
		value.Index(head).Set(temp)
	}
	return nil
}

// 字符串切片
type StringSlice []string

// 元素去重
func (StringSlice) Distinct(slice []string) []string {
	cp := make([]string, 0, len(slice))
	mapper := make(map[string]struct{})
	for i := 0; i < len(slice); i++ {
		s := (slice)[i]
		if _, ok := mapper[s]; !ok {
			cp = append(cp, s)
			mapper[s] = struct{}{}
		}
	}
	return cp
}

// 是否包含
func (StringSlice) Contains(slice []string, value string) bool {
	for _, s := range slice {
		if s == value {
			return true
		}
	}
	return false
}
