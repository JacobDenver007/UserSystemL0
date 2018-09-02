package user

import (
	"encoding/hex"
	"log"

	"github.com/bocheninc/L0/components/crypto"
	"github.com/bocheninc/L0/core/accounts"
)

func CreateAccount() string {
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
	DBClient.InsertAccount(addressStr, hexString)
	return addressStr
}
