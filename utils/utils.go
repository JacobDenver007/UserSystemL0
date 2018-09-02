package utils

import (
	"encoding/hex"
)

func HexToBytes(s string) []byte {
	h, _ := hex.DecodeString(s)
	return h
}

func HexToChainCoordinate(hex string) []byte {
	return NewChainCoordinate(HexToBytes(hex))
}

func NewChainCoordinate(c []byte) []byte {
	var cc = make([]byte, len(c))
	copy(cc, c)
	return cc
}