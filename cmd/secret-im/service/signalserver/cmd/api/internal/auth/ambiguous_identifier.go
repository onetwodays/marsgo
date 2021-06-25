package auth

import "strings"

// 模糊身份(使用number或者uuid)
type AmbiguousIdentifier struct {
	UUID   string
	Number string
}

// 创建模糊身份
func NewAmbiguousIdentifier(target string) AmbiguousIdentifier {
	if strings.Contains(target,`-`) {
		return AmbiguousIdentifier{
			UUID: target,
		}
	}
	return AmbiguousIdentifier{
		Number: target,
	}
}
