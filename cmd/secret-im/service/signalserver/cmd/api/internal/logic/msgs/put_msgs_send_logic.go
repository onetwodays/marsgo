package logic

import (
	"context"
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"secret-im/service/signalserver/cmd/api/textsecure"
	"secret-im/service/signalserver/cmd/api/util"
	"secret-im/service/signalserver/cmd/model"
	"time"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type PutMsgsSendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPutMsgsSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) PutMsgsSendLogic {
	return PutMsgsSendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PutMsgsSendLogic) PutMsgsSend(req types.PutMsgsSendReq,number string) (*types.PutMsgSendRes, error) {
	// todo: add your logic here and delete this line
	_,account,err:=l.svcCtx.GetOneAccountByNumber(number)
	if err!=nil{
		return nil, err
	}
	isSyncMessage:=false //是否是一个帐号下面同步消息
	if number==req.Destination {
		isSyncMessage=true //
	}
	if req.Relay==""{
		_,dstAccount,err:=l.svcCtx.GetOneAccountByNumber(req.Destination)
		if err!=nil{
			return nil, err
		}

		// 发消息
		ok:=l.sendLocalMessage(account,dstAccount,&req,isSyncMessage)
		if !ok{
			//保存到数据库或者redis
			_,err=l.svcCtx.MsgsModel.Insert(model.TMessages{
				Content: req.MsgList[0].Content,
				Tm: req.Timestamp,
				Type: int64(req.MsgList[0].Type),
				Source: account.Number,
				SourceDevice: account.Devices[0].Id,
				Destination: dstAccount.Number,
				DestinationDevice: dstAccount.Devices[0].Id,
			})
			if err!=nil{
				return nil, err
			}
		}
	}
	return &types.PutMsgSendRes{}, nil
}

func (l *PutMsgsSendLogic) sendLocalMessage(src,dst *types.Account,req *types.PutMsgsSendReq,isSyncMsgList bool) (ok bool) {
	if req.Timestamp == 0 {
		req.Timestamp = time.Now().UnixNano() / 1e6
	}
	for i := range req.MsgList {
		msg := &req.MsgList[i]
		outmsq := &textsecure.OutMessage{}
		body, _ := base64.StdEncoding.DecodeString(msg.Body)
		outmsq.Body = body
		content, _ := base64.StdEncoding.DecodeString(msg.Content)
		outmsq.Content = content
		outmsq.Timestamp = *(proto.Int64(req.Timestamp))
		outmsq.SourceDevice = *(proto.Int64(src.Devices[0].Id))
		outmsq.Source = src.Number
		outmsq.Type = *(textsecure.OutMessage_Type(int32(msg.Type)).Enum())
		out, err := proto.Marshal(outmsq)
		if err != nil {
			l.Logger.Error(err)
		}
		sm := new(types.SingleMessage)
		sm.Id = dst.Number
		sm.Message, _ = util.AESCBCEncrypt(out, dst.Devices[0].AccountAttributes.SignalingKey)
		l.Logger.Infof("%v", sm.Message)
		/*
		if clinet, ok := l.svcCtx.Hub.ClientMap[sm.Id]; ok {
			clinet.Hub.SendSingle <- sm
			break
		}*/
	}
	return
}
