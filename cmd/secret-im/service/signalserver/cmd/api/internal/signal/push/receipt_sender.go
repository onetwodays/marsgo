package push

import (
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/textsecure"
)

// 回执发送器
type ReceiptSender struct {
	pushSender *Sender
}

// 创建回执发送器
func NewReceiptSender(sender *Sender) *ReceiptSender {
	return &ReceiptSender{pushSender: sender}
}

// 发送回执
func (sender *ReceiptSender) SendReceipt(number, uuid string, deviceID int64,
	destination string, timestamp int64) error {

	if number == destination {
		return nil
	}

	destinationAccount, err := sender.getDestinationAccount(destination)
	if err != nil {
		return err
	}

	var message textsecure.Envelope
	temp := uint64(timestamp)
	typ := textsecure.Envelope_RECEIPT
	sourceDeviceID := uint32(deviceID)
	message.Source = number
	message.SourceUuid = uuid
	message.SourceDevice = sourceDeviceID
	message.Timestamp = temp
	message.Type = typ

	destinationDevices := destinationAccount.Devices
	for _, destinationDevice := range destinationDevices {
		sender.pushSender.SendMessage(destinationAccount.Number, &destinationDevice.Device, message, false)
	}
	return nil
}

// 获取目标帐号
func (sender *ReceiptSender) getDestinationAccount(destination string) (*entities.Account, error) {
	
	/*
	account, err := storage.AccountsManager{}.GetByNumber(destination)
	if err != nil {
		return nil, err
	}
	return account, nil
	 */
	return nil, nil
}

