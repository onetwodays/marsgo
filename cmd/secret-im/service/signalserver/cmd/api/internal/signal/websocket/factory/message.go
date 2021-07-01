package factory

import (
	"net/http"
	"secret-im/service/signalserver/cmd/api/textsecure"
)

func CreateRequest ( requestID uint64,verb,path string,headers []string,body []byte) *textsecure.WebSocketMessage{

	requestMessage:=&textsecure.WebSocketRequestMessage{
		Id: requestID,
		Verb: verb,
		Path: path,
	}
	if len(body) >0 {
		requestMessage.Body = body
	}

	if len(headers) >0 {
		requestMessage.Headers=headers
	}
	message := &textsecure.WebSocketMessage{}
	message.Type = textsecure.WebSocketMessage_REQUEST
	message.Request = requestMessage
	return message
}


func CreateResponse(requstID uint64,status int,messageString string,headers []string,body []byte) *textsecure.WebSocketMessage{
	responseMessage :=&textsecure.WebSocketResponseMessage{
		Id:requstID,
		Status: uint32(status),
		Message: messageString,
	}
	if len(messageString) == 0{
		responseMessage.Message = http.StatusText(status)
	}
	if len(body) >0 {
		responseMessage.Body = body
	}
	if len(headers) >0 {
		responseMessage.Headers =headers
	}


	message := &textsecure.WebSocketMessage{}
	message.Type = textsecure.WebSocketMessage_RESPONSE
	message.Response = responseMessage
	return message
}
