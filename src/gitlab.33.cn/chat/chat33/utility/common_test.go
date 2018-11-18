package utility_test

import (
	"testing"

	"gitlab.33.cn/chat/chat33/utility"
)

func TestUUID(t *testing.T) {
	v := utility.RandomID()
	t.Log(v)
}
