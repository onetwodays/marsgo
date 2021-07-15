package util
import (
	"github.com/tal-tech/go-zero/core/logx"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/proto"
)




func PFParse(req *http.Request, obj interface{}) error {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	return pfParseBody(buf, obj)
}

func pfParseBody(body []byte, obj interface{}) error {
	if err := proto.Unmarshal(body, obj.(proto.Message)); err != nil {
		return err
	}
	logx.Infof("protobuf body req :{%s}",obj.(proto.Message).String())
	// Here it's same to return validate(obj), but util now we can't add
	// `binding:""` to the struct which automatically generate by gen-proto
	return nil
	// return validate(obj)
}