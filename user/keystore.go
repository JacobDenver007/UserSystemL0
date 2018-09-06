package user

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/bocheninc/L0/components/crypto"
	"github.com/bocheninc/L0/core/accounts"
)

func GernerateAccount() (string, string) {
	k, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	address := accounts.PublicKeyToAddress(*k.Public())
	addressStr := "0x" + hex.EncodeToString(address[:])

	secretByte := k.SecretBytes()
	hexString := hex.EncodeToString(secretByte)

	k1, _ := crypto.HexToECDSA(hexString)

	address = accounts.PublicKeyToAddress(*k1.Public())
	addressStr = "0x" + hex.EncodeToString(address[:])
	return addressStr, hexString
}

func GenerateAccountByPrivateKey(hexString string) (string, error) {
	k1, err := crypto.HexToECDSA(hexString)
	if err != nil {
		return "", fmt.Errorf("私钥不合法")
	}

	address := accounts.PublicKeyToAddress(*k1.Public())
	addressStr := "0x" + hex.EncodeToString(address[:])

	return addressStr, nil
}
