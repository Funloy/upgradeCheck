package token

import (
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {

	tokenSecretKey = "malacode.com.token"
	tokenIssuer = "malacode.com"

	token := &Token{
		UserID:     "123456789",
		UserPhone:  "18566208215",
		ExpireTime: time.Now().Unix(),
	}

	sign, _ := token.CreateToken()
	t.Logf("signature: %s", sign)

	k, err := token.ValidateToken(sign)
	if err != nil {
		t.Errorf("token err: %v", err)
		return
	}
	if k.ExpireTime < time.Now().Unix() {
		t.Errorf("ExpireTime: %v", k.ExpireTime)
	}

}
