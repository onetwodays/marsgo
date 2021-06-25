package crypto

import (
	"crypto/rand"
	"errors"
	"golang.org/x/crypto/curve25519"
	"io"
)

// DjbType is the Diffie-Hellman curve type (curve25519) created by D. J. Bernstein.
const DjbType = 0x05

// DecodePoint will take the given bytes and offset and return an ECPublicKeyable object.
// This is used to check the byte at the given offset in the byte array for a special
// "type" byte that will determine the key type. Currently only DJB EC keys are supported.
func DecodePoint(bytes []byte, offset int) (ECPublicKeyable, error) {
	keyType := bytes[offset] & 0xFF

	switch keyType {
	case DjbType:
		keyBytes := [32]byte{}
		copy(keyBytes[:], bytes[offset+1:])
		return NewDjbECPublicKey(keyBytes), nil
	default:
		return nil, errors.New("Bad key type: " + string(keyType))
	}
}

// GenerateKeyPair returns an EC Key Pair.
func GenerateKeyPair() (*ECKeyPair, error) {
	random := rand.Reader
	var private, public [32]byte

	_, err := io.ReadFull(random, private[:])
	if err != nil {
		return nil, err
	}

	private[0] &= 248
	private[31] &= 127
	private[31] |= 64

	curve25519.ScalarBaseMult(&public, &private)

	djbECPub := NewDjbECPublicKey(public)
	djbECPriv := NewDjbECPrivateKey(private)
	keypair := NewECKeyPair(djbECPub, djbECPriv)
	return keypair, nil
}

// VerifySignature verifies that the message was signed with the given key.
func VerifySignature(signingKey ECPublicKeyable, message []byte, signature [64]byte) bool {
	publicKey := signingKey.PublicKey()
	valid := verify(publicKey, message, &signature)
	return valid
}

// CalculateSignature signs a message with the given private key.
func CalculateSignature(signingKey ECPrivateKeyable, message []byte) [64]byte {
	var random [64]byte
	r := rand.Reader
	io.ReadFull(r, random[:])
	privateKey := signingKey.Serialize()
	signature := sign(&privateKey, message, random)
	return *signature
}
