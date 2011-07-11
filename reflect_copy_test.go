package perse

import (
	"testing"
)


func TestBasicCopy(t *testing.T) {
	var v1, v2 int
	v1 = 1
	c := Copy(v1)
	v2, ok := c.(int)
	if !ok {
		t.Fail()
	}
	v2 = 2
	if v1 == v2 {
		t.Fail()
	}
}

type simple struct {
	A int
}

func TestStructExportedCopy(t *testing.T) {
	var v1, v2 simple
	v1 = simple{0}
	v1.A = 1
	c := Copy(v1)
	v2, ok := c.(simple)
	if !ok {
		t.Fail()
	}
	v2.A = 2
	if v1.A == v2.A {
		t.Fail()
	}
}

type simpleun struct {
	a int
}

func TestStructUnExportedCopy(t *testing.T) {
	var v1, v2 simpleun
	v1 = simpleun{0}
	v1.a = 1
	c := Copy(v1)
	v2, ok := c.(simpleun)
	if !ok {
		t.Fail()
	}
	v2.a = 2
	if v1.a == v2.a {
		t.Fail()
	}
}
func TestBasicPtrCopy(t *testing.T) {
	var v1, v2 *int
	v1 = new(int)
	*v1 = 1
	c := Copy(v1)
	v2, ok := c.(*int)
	if !ok {
		t.Error("Type missmatch")
	}
	*v2 = 2
	if *v1 == *v2 {
		t.Error("Its not a copy")
	}
}

func TestStructPtrCopy(t *testing.T) {
	var v1, v2 *simple
	v1 = &simple{0}
	v1.A = 1
	c := Copy(v1)
	v2, ok := c.(*simple)
	if !ok {
		t.Fail()
	}
	v2.A = 2
	if v1.A == v2.A {
		t.Fail()
	}
}

type ptrsimpleptr struct {
	B *simple
}

func TestPtrStructPtrCopy(t *testing.T) {
	var v1, v2 *ptrsimpleptr
	v1 = &ptrsimpleptr{&simple{0}}
	v1.B.A = 1
	c := Copy(v1)
	v2, ok := c.(*ptrsimpleptr)
	if !ok {
		t.Fail()
	}
	v2.B.A = 2
	if v1.B.A == v2.B.A {
		t.Fail()
	}
}
