package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"strings"
)

func GenRsaKey(bits int) (priHex string, pubHex string, priPem string, pubPem string, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}

	publicKey := &privateKey.PublicKey
	derPubStream := x509.MarshalPKCS1PublicKey(publicKey)
	pubHex = hex.EncodeToString(derPubStream)
	block := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPubStream,
	}
	pubPemBytes := pem.EncodeToMemory(&block)
	pubPem = string(pubPemBytes)

	derPriStream := x509.MarshalPKCS1PrivateKey(privateKey)
	priHex = hex.EncodeToString(derPriStream)
	block = pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derPriStream,
	}
	priPemBytes := pem.EncodeToMemory(&block)
	priPem = string(priPemBytes)
	return
}

func RsaEncrypt(data []byte, pubKey string) ([]byte, error) {
	pubBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}

	pub, err := x509.ParsePKCS1PublicKey(pubBytes)
	if err != nil {
		return nil, err
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

func RsaDecrypt(ciphertext []byte, privateKey string) ([]byte, error) {
	privateBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}

	priv, err := x509.ParsePKCS1PrivateKey(privateBytes)
	if err != nil {
		return nil, err
	}
	key := priv

	data, err := rsa.DecryptPKCS1v15(rand.Reader, key, ciphertext)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func RemovePrefixAndBase64(pemStr string) string {
	str := strings.TrimPrefix(pemStr, "-----BEGIN RSA PUBLIC KEY-----")
	newstr := strings.TrimSuffix(str, "-----END RSA PUBLIC KEY-----\n")
	base64ed := base64.StdEncoding.EncodeToString([]byte(newstr))
	return base64ed
}
