package main

import (
	"encoding/hex"
	"fmt"

	"github.com/bocheninc/L0/components/crypto"
)

func main() {
	fmt.Println("secp256k1")
	privateKey, _ := crypto.GenerateKey()
	fmt.Println("privateKey hex:\t\t", hex.EncodeToString(privateKey.SecretBytes()))
	fmt.Println("publicKey hex:\t\t", hex.EncodeToString(privateKey.Public().Bytes()))

	content := "hello world"
	fmt.Println("sign content:\t\t", content)

	hash := crypto.Sha256([]byte(content[:]))
	fmt.Println("sha256 hash:\t\t", hex.EncodeToString(hash[:]))

	hash160 := crypto.Ripemd160([]byte(content[:]))
	fmt.Println("ripemd160 hash:\t\t", hex.EncodeToString(hash160[:]))

	signature, _ := privateKey.Sign(hash[:])
	fmt.Println("signature(sha256):\t\t", hex.EncodeToString(signature.Bytes()))

	pub, _ := signature.RecoverPublicKey(hash[:])
	fmt.Println("recover publicKey hex:\t\t", hex.EncodeToString(pub.Bytes()))
}
