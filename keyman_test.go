package keyman

import (
	"fmt"
	"testing"
)

func TestTokenInfoMarshal(t *testing.T) {
	tokeninfo := new(TokenInfo)
	tokeninfo.Key = "aaa"
	tokeninfo.Route = "/a/b"
	b, err := tokeninfo.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
	tokeninfo2 := new(TokenInfo)
	err = tokeninfo2.Unmarshal(b)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(tokeninfo2)
}

func TestTokenSign(t *testing.T) {
	key := "48409852818866747867224752556126404236692416387864301776209804402161055141729"
	priv := StrToPriv(key)
	if priv == nil {
		t.Error("key error")
		return
	}
	addrStr := KeyToAddrStr(key)
	t.Log(addrStr)
	token := MakeToken(priv)
	addrStr, err := TokenToAddrStr(token)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(addrStr)
	token = MakeToken(priv)
	addrStr, err = TokenToAddrStr(token)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(addrStr)
}
