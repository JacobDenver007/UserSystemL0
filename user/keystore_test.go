package user

import (
	"testing"
)

func TestAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateAccount()
	}
}
