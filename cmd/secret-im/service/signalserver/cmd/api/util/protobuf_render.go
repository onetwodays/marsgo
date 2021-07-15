package util
import (
	"github.com/tal-tech/go-zero/core/logx"
	"net/http"

	"github.com/golang/protobuf/proto"
)


var protobufContentType = []string{"application/x-protobuf"}

// Render (ProtoBuf) marshals the given interface object and writes data with custom ContentType.
func Render(w http.ResponseWriter,Data interface{}) error {
	logx.Infof("protobuf body resp {%s}",Data.(proto.Message).String())

	bytes, err := proto.Marshal(Data.(proto.Message))
	if err != nil {
		return err
	}
	writeContentType(w, protobufContentType)
	_, err = w.Write(bytes)
	return err
}


func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

