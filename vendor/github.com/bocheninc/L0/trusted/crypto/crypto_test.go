package main

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/bocheninc/L0/components/crypto"
)

const (
	testPrivateKey       = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	testRipemd160HashStr = "8eb208f7e05d987a9b044a8e98c6b087f15a0bfc"
	testSha256HashStr    = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
)

func BenchmarkSignAndRecover(b *testing.B) {
	priv, _ := crypto.HexToECDSA(testPrivateKey)
	for i := 0; i < b.N; i++ {
		msg := crypto.Sha256([]byte("hello"))
		sig, err := priv.Sign(msg[:])
		if err != nil {
			b.Errorf("Sign Error %s", err)
		}
		pub, err := sig.RecoverPublicKey(msg[:])
		if err != nil {
			b.Errorf("SigToPub Error %v - %v - %s", sig, pub, err)
		}
		pub2 := priv.Public()
		if !bytes.Equal(pub.Bytes(), pub2.Bytes()) {
			b.Errorf("public key not match! %0x - %0x ", pub.Bytes(), pub2.Bytes())
		}
	}
}
func BenchmarkSha256AndCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := crypto.Sha256([]byte("hello")).Bytes()
		if hex.EncodeToString(h[:]) != testSha256HashStr {
			b.Errorf("Sha256(%s) = %s, except %s !", "hello", testSha256HashStr, hex.EncodeToString(h[:]))
		}
	}
}
func BenchmarkRipemd160AndCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := crypto.Ripemd160([]byte("abc"))
		if hex.EncodeToString(h[:]) != testRipemd160HashStr {
			b.Errorf("Rimped160(%s) = %s, except %s !", "abc", testRipemd160HashStr, hex.EncodeToString(h))
		}
	}
}
