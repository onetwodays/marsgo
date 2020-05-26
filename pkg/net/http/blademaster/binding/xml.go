package binding

import (
	"encoding/xml"
	"net/http"

	"github.com/pkg/errors"
)

type xmlBinding struct{}

func (xmlBinding) Name() string {
	return "xml"
}

//把收到的数据解码成xml,并验证xml是佛正确
func (xmlBinding) Bind(req *http.Request, obj interface{}) error {
	decoder := xml.NewDecoder(req.Body)
	if err := decoder.Decode(obj); err != nil {
		return errors.WithStack(err)
	}
	return validate(obj)
}
