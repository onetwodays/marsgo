package websocket

import (
	"context"
	"io"
	"net/http"
	"net/url"
	shared "secret-im/service/signalserver/cmd/api/shared"
	"secret-im/service/signalserver/cmd/api/textsecure"
	"secret-im/service/signalserver/cmd/api/websocket/factory"
	"strings"
)

// 去除空格
func strip(s string) string {
	for {
		if len(s) == 0 {
			break
		}
		if s[0] != ' ' && s[0] != '\t' {
			break
		}
		s = s[1:]
	}
	for {
		if len(s) == 0 {
			break
		}
		if s[len(s)-1] != ' ' && s[len(s)-1] != '\t' {
			break
		}
		s = s[:len(s)-1]
	}
	return s
}
// 消息体
type netHTTPBody struct {
	b []byte
}

func (r *netHTTPBody) Read(p []byte) (int, error) {
	if len(r.b) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.b)
	r.b = r.b[n:]
	return n, nil
}

func (r *netHTTPBody) Close() error {
	r.b = r.b[:0]
	return nil
}

// 响应编写器
type netHTTPResponseWriter struct {
	statusCode int
	h          http.Header
	body       []byte
}

func (w *netHTTPResponseWriter) StatusCode() int {
	if w.statusCode == 0 {
		return http.StatusOK
	}
	return w.statusCode
}

func (w *netHTTPResponseWriter) Header() http.Header {
	if w.h == nil {
		w.h = make(http.Header)
	}
	return w.h
}

func (w *netHTTPResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *netHTTPResponseWriter) Write(p []byte) (int, error) {
	w.body = append(w.body, p...)
	return len(p), nil
}



// 处理HTTP请求
func HandleHTTPRequest(ctx *SessionContext, h http.Handler,
	request *textsecure.WebSocketRequestMessage) *textsecure.WebSocketMessage {

	var r http.Request

	r.Method = request.GetVerb()
	r.Proto = "HTTP/1.1"
	r.ProtoMajor = 1
	r.ProtoMinor = 1
	r.RequestURI = request.GetPath()
	r.ContentLength = int64(len(request.GetBody()))
	r.Host = "websocket"
	r.RemoteAddr = ctx.Session.RemoteAddr()

	header := make(http.Header)
	for _, s := range request.GetHeaders() {
		slice := strings.SplitN(s, ":", 2)
		if len(slice) != 2 {
			continue
		}
		header.Set(strip(slice[0]), strip(slice[1]))
	}
	header.Set("Flag","ws") //标识请求来自ws



	r.Header = header
	r.Body = &netHTTPBody{request.Body}
	rURL, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		return factory.CreateResponse(request.GetId(), http.StatusNotFound, err.Error(), nil, nil)
	}
	r.URL = rURL

	w := netHTTPResponseWriter{statusCode: http.StatusOK}
	// 是来自ws的请求，这里上下文都设置几个值

	c:=context.WithValue(context.Background(), "ws",ctx) //没有被使用
	if ctx.Device!=nil{
		header.Set("Ws-Auth","ok") //标识请求来自ws
		c= context.WithValue(c, shared.CONTENTKEYUUID,ctx.Device.UUID)
		c= context.WithValue(c, shared.CONTENTKEYDEVICEID,ctx.Device.Device.ID)
		c=context.WithValue(c, shared.HttpReqContextAccountKey, ctx.account)
	}

	// //////////////////////////////////


	h.ServeHTTP(&w, r.WithContext(c))

	headers := make([]string, 0)
	for name, items := range w.Header() {
		for _, val := range items {
			headers = append(headers, name+":"+val)
		}
	}

	return factory.CreateResponse(request.GetId(), w.statusCode, "", headers, w.body)
}
