package psql

import (
	"testing"
	ps "github.com/mjarco/perse"
	//    "fmt"
	"strconv"
)

type Test struct {
	Someid   *ps.Serial
	SomeText *ps.Char
	I        *ps.Int
}

func NewTest(text string, i int) *Test {
	return &Test{ps.NewSerial(0), ps.NewChar(text), ps.NewInt(i)}
}


func TestInit(t *testing.T) {
	Conn.Connect(CD)
	Conn.Query(&ps.Stmt{make([]string, 0), "CREATE TABLE test(Someid Serial, SomeText Char(255), i int);"})

	for i := 0; i < 10; i++ {
		ps.Save(NewTest("test"+strconv.Itoa(i), i%3))
	}
}

func TestConds(t *testing.T) {
	c1 := ps.NewCond("i", ps.EQ, "1")
	c2 := ps.NewCond("i", ps.EQ, "0")
	res := ps.NewCollection(NewTest("", 0)).Filter(c1).Query()

	mm := res.AllValues(NewTest("", 0))

	//    println(mm[1],mm[2],mm[3])
	if len(mm) != 3 {
		t.Error("Expected 3 got", len(mm))
	}
	res2 := ps.NewCollection(NewTest("", 0)).Filter(c1.Or(c2)).Query()
	mm2 := res2.AllMaps()
	if len(mm2) != 7 {
		t.Error("Expected 7 got", len(mm2))
	}

	res3 := ps.NewCollection(NewTest("", 0)).Filter(c1.Or(c2)).Limit(2).Query()
	mm3 := res3.AllMaps()
	if len(mm3) != 2 {
		t.Error("Expected 2 got", len(mm3))
	}

}


var CD = &ps.ConnData{"host", "5432", "user", "pass", "dbname", ""}

func TestDefaultConnection(t *testing.T) {
	//TODO: Gdy nie ma polaczenia idzie panic... slaaabo!
	Conn.Connect(CD)
	//insert
	tt := NewTest("łopata", 10)
	v, id := ps.Save(tt)
	//update
	tt2, ok := v.(*Test)
	if ok {
		t.Log(tt2.Someid.String())
		tt2.SomeText.Value("update")
		_, id2 := ps.Save(tt2)
		if id != id2 {
			t.Error("Inserted twice " + id + " " + id2)
		}
	}

}


func TestCustomConnection(t *testing.T) {
	con := NewConnection()
	ok := con.Connect(CD)
	if !ok {
		t.Error("Postgresql connection cannot be established")
	}
}

func TestShutdown(t *testing.T) {
	Conn.Query(&ps.Stmt{make([]string, 0), "drop table test;"})
}
