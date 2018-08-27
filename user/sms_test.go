package user

import (
	"fmt"
	"testing"
)

func TestSMS(t *testing.T) {
	//SendCode("18610019263", "1234")

	for i := 0; i < 0; i++ {
		code := MakeCode()
		fmt.Println(code, VailCode(code) == nil)
		SendCode("13466487843", code)
	}
}
