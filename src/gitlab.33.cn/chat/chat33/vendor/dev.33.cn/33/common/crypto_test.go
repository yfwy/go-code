package common

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDesCBC(t *testing.T) {
	key := "01234567" //8
	plain := "this"

	must := require.New(t)
	encBytes, err := DesCBCEncrypt([]byte(plain), []byte(key))
	must.Nil(err)

	decBytes, err := DesCBCDecrypt(encBytes, []byte(key))
	must.Nil(err)
	must.Equal(plain, string(decBytes))
}

func TestDesECB(t *testing.T) {
	key := "01234567" //8
	plain := "this"

	must := require.New(t)
	encBytes, err := DesECBEncrypt([]byte(plain), []byte(key))
	must.Nil(err)

	decBytes, err := DesECBDecrypt(encBytes, []byte(key))
	must.Nil(err)
	must.Equal(plain, string(decBytes))
}

func TestDes3CBC(t *testing.T) {
	key := "6b65805371548541wert1234" //24
	plain := "this"

	must := require.New(t)
	encBytes, err := Des3CBCEncrypt([]byte(plain), []byte(key))
	must.Nil(err)

	decBytes, err := Des3CBCDecrypt(encBytes, []byte(key))
	must.Nil(err)
	must.Equal(plain, string(decBytes))
}

func TestDes3ECB(t *testing.T) {
	key := "6b65805371548541wert1234"
	plain := "this"

	must := require.New(t)
	encBytes, err := Des3ECBEncrypt([]byte(plain), []byte(key))
	must.Nil(err)

	decBytes, err := Des3ECBDecrypt(encBytes, []byte(key))
	must.Nil(err)
	must.Equal(plain, string(decBytes))
}

func TestVerify(t *testing.T) {
	key := "6b65805371548541wert1234"
	plain := "Scp5vvJOmPLC3DjeIbkV/ltWDK41NsL5fLrj5Q5rWPJ+aBgOaZIbrgSCq9vVQieC"

	must := require.New(t)
	encBytes, err := base64.StdEncoding.DecodeString(plain)
	must.Nil(err)

	decBytes, err := Des3ECBDecrypt(encBytes, []byte(key))
	must.Nil(err)
	t.Log(string(decBytes))
}
