package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"reflect"
	"regexp"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

// 获取纤程ID
func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	re, err := regexp.Compile("goroutine (\\d+)")
	if err != nil {
		return 0
	}
	results := re.FindStringSubmatch(string(b))
	if len(results) <= 1 {
		return 0
	}
	n, _ := strconv.ParseUint(results[1], 10, 64)
	return n
}

// 安全函数调用
func NoPanic(f func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()

	f()
	return nil
}

// 是否有效UUID
func IsValidUUID(id string) bool {
	_, err := uuid.FromString(id)
	return err == nil
}

// 是否空字符串
func EmptyString(s *string) bool {
	if s == nil {
		return true
	}
	return len(*s) == 0
}

// 获取对象类型
func GetObjectType(v interface{}) string {
	reflectValue := reflect.ValueOf(v)
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue.Type().String()
}

// 字节数格式化
func BytesPretty(bytes int64) string {
	const B int64 = 1
	const KB = B * 1024
	const MB = KB * 1024
	const GB = MB * 1024
	const TB = GB * 1024
	const PB = TB * 1024
	units := []int64{PB, TB, GB, MB, KB, B}
	symbols := []string{"PB", "TB", "GB", "MB", "KB", "B"}
	var pos int
	for idx, unit := range units {
		pos = idx
		if bytes >= unit {
			break
		}
	}
	symbol := symbols[pos]
	numerator := decimal.New(bytes, 0)
	denominator := decimal.New(units[pos], 0)
	value := numerator.DivRound(denominator, 6).Truncate(2)
	return value.String() + " " + symbol
}

// 生成自签名证书
func GenerateCert(host string) (cert, key *bytes.Buffer, err error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 24 * 365)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, err
	}

	key = bytes.NewBuffer(nil)
	cert = bytes.NewBuffer(nil)
	if err := pem.Encode(cert, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return nil, nil, err
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, nil, err
	}
	if err := pem.Encode(key, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return nil, nil, err
	}
	return cert, key, nil
}
