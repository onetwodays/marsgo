package entities

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/service/signalserver/cmd/api/internal/cipher"
	"secret-im/service/signalserver/cmd/api/textsecure"
)

const (
	MacSize       = 10
	MacKeySize    = 20
	CipherKeySize = 32
)

// 加密传出消息
type EncryptedOutgoingMessage struct {
	Serialized []byte
}

// 创建加密传出消息
func NewEncryptedOutgoingMessage(outgoingMessage *textsecure.Envelope, signalingKey string) (EncryptedOutgoingMessage, error) {
	var message EncryptedOutgoingMessage
	plaintext, err := proto.Marshal(outgoingMessage)
	if err != nil {
		return message, err
	}

	macKey, err := message.getMacKey(signalingKey)
	if err != nil {
		return message, err
	}
	cipherKey, err := message.getCipherKey(signalingKey)
	if err != nil {
		return message, err
	}

	message.Serialized, err = message.getCipherText(plaintext, cipherKey, macKey)
	return message, err
}

// 获取哈希密钥
func (EncryptedOutgoingMessage) getMacKey(signalingKey string) ([]byte, error) {
	var macKey [MacKeySize]byte
	signalingKeyBytes, err := base64.StdEncoding.DecodeString(signalingKey)
	if err != nil {
		return nil, err
	}

	if len(signalingKeyBytes) < CipherKeySize+MacKeySize {
		return nil, errors.New("signaling key too short")
	}

	copy(macKey[:], signalingKeyBytes)
	return macKey[:], nil
}

// 获取加密密钥
func (EncryptedOutgoingMessage) getCipherKey(signalingKey string) ([]byte, error) {
	var cipherKey [CipherKeySize]byte
	signalingKeyBytes, err := base64.StdEncoding.DecodeString(signalingKey)
	if err != nil {
		return nil, err
	}

	if len(signalingKeyBytes) < CipherKeySize {
		return nil, errors.New("signaling key too short")
	}

	copy(cipherKey[:], signalingKeyBytes)
	return cipherKey[:], nil
}

// 获取消息密文
func (EncryptedOutgoingMessage) getCipherText(plaintext, cipherKey, macKey []byte) ([]byte, error) {
	// 加密明文
	var iv [16]byte
	_, err := rand.Read(iv[:])
	if err != nil {
		return nil, err
	}

	cipherText, err := cipher.Encrypt(iv[:], cipherKey, plaintext)
	if err != nil {
		return nil, err
	}

	// 计算hash值
	version := []byte{0x01}
	var truncatedMac [MacSize]byte
	hmac := hmac.New(sha256.New, macKey)
	hmac.Write(version)
	hmac.Write(iv[:])
	mac := hmac.Sum(cipherText)
	copy(truncatedMac[:], mac)

	combine := bytes.Join([][]byte{version, iv[:], cipherText, truncatedMac[:]}, nil)

	logx.Info("[Test] new encrypted outgoing message,",
		"version:", version,
		"iv:", iv,
		"cipher_text:", cipherText,
		"mac:", truncatedMac)

	return combine, nil
}
