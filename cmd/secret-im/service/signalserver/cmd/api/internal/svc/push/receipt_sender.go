package push

import (
	"secret-im/service/signalserver/cmd/api/internal/storage"
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

	destinationAccount, err := storage.AccountManager{}.GetByNumber(destination)
	if err != nil || destinationAccount==nil {
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
		sender.pushSender.SendMessage(destinationAccount.Number, &destinationDevice.Device, &message, false)
	}
	return nil
}


