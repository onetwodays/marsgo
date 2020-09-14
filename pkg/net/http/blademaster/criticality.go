package blademaster

import (
	criticalityPkg "marsgo/pkg/net/criticality"
	"marsgo/pkg/net/metadata"

	"github.com/pkg/errors"
)

// Criticality is
func Criticality(pathCriticality criticalityPkg.Criticality) HandlerFunc {
	if !criticalityPkg.Exist(pathCriticality) {
		panic(errors.Errorf("This criticality is not exist: %s", pathCriticality))
	}
	return func(ctx *Context) {
		md, ok := metadata.FromContext(ctx)
		if ok { //只是判断一下是否有元数据
			md[metadata.Criticality] = string(pathCriticality)
		}
	}
}
