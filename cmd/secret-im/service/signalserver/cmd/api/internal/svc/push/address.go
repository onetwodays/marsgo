package push

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// 设备地址
type Address struct {
	Number   string
	DeviceID int64
}

// 创建设备地址
func NewAddress(serialized string) (Address, error) {
	parts := strings.Split(serialized, ":")
	if len(parts) != 2 || len(parts[0]) < 2 {
		return Address{}, errors.New("bad address")
	}

	number := parts[0]
	if number[0] != '{' || number[len(number)-1] != '}' {
		return Address{}, errors.New("bad address")
	}

	addr := Address{
		Number: number[1 : len(number)-1],
	}
	deviceID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return Address{}, err
	}
	addr.DeviceID = deviceID
	return addr, nil
}

// 序列化地址
func (address Address) Serialize() string {
	return fmt.Sprintf("{%s}:%d", address.Number, address.DeviceID)
}

// 服务调配地址
type ProvisioningAddress struct {
	Address
}

// 生成新地址
func (ProvisioningAddress) Generate() ProvisioningAddress {
	var random [16]byte
	rand.Read(random[:])

	s := base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(random[:])
	s = strings.ReplaceAll(s, "+", "-")
	s = strings.ReplaceAll(s, "/", "_")
	return ProvisioningAddress{
		Address: Address{Number: s},
	}
}
