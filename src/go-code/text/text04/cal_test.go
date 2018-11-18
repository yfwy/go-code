package main

import "testing"


func TestAddup(t *testing.T)  {
	cal := addup(5)
	if cal != 55 {
		t.Fatalf("错误")
	}
	t.Logf("正确")
}