package user

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func CreateAccount() string {
	k, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	address := crypto.PubkeyToAddress(k.PublicKey)
	addressStr := "0x" + hex.EncodeToString(address[:])

	content, err := json.Marshal(k)
	fmt.Println(string(content))
	DBClient.InsertAccount(addressStr, string(content))
	return addressStr
}
