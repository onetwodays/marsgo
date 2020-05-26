package binding

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "json"
}

//req 是src obj是dst,返回值
func (jsonBinding) Bind(req *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(obj); err != nil {
		return errors.WithStack(err)
	}
	return validate(obj)
}
