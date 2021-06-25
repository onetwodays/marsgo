package auth

import (
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"secret-im/pkg/utils-tools"
	"secret-im/service/signalserver/cmd/api/internal/crypto"
	"secret-im/service/signalserver/cmd/api/textsecure"
)

// 证书生成器
type CertificateGenerator struct {
	PrivateKey        crypto.ECPrivateKeyable
	ExpiresDays       int
	ServerCertificate *textsecure.ServerCertificate
}

// 创建证书生成器
func NewCertificateGenerator(
	serverCertificate []byte, privateKey crypto.ECPrivateKeyable, expiresDays int) (*CertificateGenerator, error) {
	var message textsecure.ServerCertificate
	if err := proto.Unmarshal(serverCertificate, &message); err != nil {
		return nil, err
	}

	return &CertificateGenerator{
		PrivateKey:        privateKey,
		ExpiresDays:       expiresDays,
		ServerCertificate: &message,
	}, nil
}

// 创建证书
func (generator *CertificateGenerator) CreateFor(
	number, identityKey, uuid string, deviceID int64, includeUuid bool) ([]byte, error) {
	senderDevice := uint32(deviceID)
	identityKeyData, err := base64.StdEncoding.DecodeString(identityKey)
	if err != nil {
		return nil, err
	}
	expires := uint64(utils.CurrentTimeMillis() + utils.DaysToMillis(generator.ExpiresDays))

	builder := textsecure.SenderCertificateMessage{
		Sender:       number,
		SenderDevice: senderDevice,
		Expires:      expires,
		IdentityKey:  identityKeyData,
		Signer:       generator.ServerCertificate,
	}
	if includeUuid {
		builder.SenderUuid = uuid
	}

	certificate, err := proto.Marshal(&builder)
	if err != nil {
		return nil, err
	}
	signature := crypto.CalculateSignature(generator.PrivateKey, certificate)

	message := textsecure.ServerCertificate{
		Certificate: certificate,
		Signature:   signature[:],
	}
	return proto.Marshal(&message)
}
