package target

import (
	"bytes"
	"testing"
)

func TestMethod(t *testing.T) {
	at := &MethodActivateTarget{TargetID: "abc"}

	if at.Domain() != "Target" {
		t.FailNow()
	}
	if at.Name() != ActivateTarget {
		t.FailNow()
	}
	data, err := at.Dump()
	if err != nil {
		t.FailNow()
	}
	if bytes.Compare(data, []byte(`{"targetId":"abc"}`)) != 0 {
		t.FailNow()
	}
	ret, err := at.Load([]byte("{}"))
	if err != nil {
		t.FailNow()
	}
	if _, ok := ret.(ActivateTargetReturns); !ok {
		t.FailNow()
	}
}

func TestBadRetType(t *testing.T) {
	at := &MethodCreateTarget{URL: "www.baidu.com"}
	val, err := at.Load([]byte(`{"targetId":"abc"}`))
	if err != nil {
		t.FailNow()
	}
	ret, ok := val.(*CreateTargetReturns)
	if !ok {
		t.FailNow()
	}
	if ret.TargetID != "abc" {
		t.FailNow()
	}
}
