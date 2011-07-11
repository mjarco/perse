package perse

import (
	"testing"
	"strings"
	"strconv"
	//    "fmt"
)

type TStruct struct {
	A *Text
	B *Int
}
/** field types tests **/
func TestChar(t *testing.T) {
	ch := NewChar("ttest")
	if ch.length != 5 {
		t.Fail()
	}
	ch = NewChar("ttest", 25)
	if ch.length != 25 {
		t.Fail()
	}

}
/** struct2map & map2struct tests**/
//ModelMap

type TestStruct struct {
	Date *Date
	Char *Char
	Int  *Int
	Text *Text
}

func NewTestStruct() *TestStruct {
	return &TestStruct{NewDate("2006 01 02 15:04"), NewChar(`non-significant' "text"`), NewInt(1), NewText(`testtest '"\" " testt`)}
}
func TestMaps(t *testing.T) {
	d1 := &TStruct{NewText(""), NewInt(1)} // &TStruct{&Text{"dupa"}, &Int{1}}
	d1.A.Value("valuea")
	d1.B.Value("2")
	mm, ok := ModelMap(d1)
	if !ok {
		t.Fail()
	}
	if v, ok := mm["a"]; ok {
		ok := strings.Contains(v, d1.A.String())
		if !ok {
			t.Fail()
		}
	} else {
		t.Error("Map doesn't contain key")
	}
}
/** escaped map **/
func TestEscaping(t *testing.T) {

	ts := NewTestStruct()
	mm, _ := EscapedMap(ts)
	v, _ := CrudValue(NewTestStruct(), mm)
	mm2, _ := EscapedMap(v)
	mv, _ := ModelMap(ts)
	mv2, _ := ModelMap(v)
	for k, v := range mm {
		v2 := mm2[k]
		if v != v2 {
			t.Error("Escape/Unescape fail for type  " + k +
				"\nexpected:  " + v +
				"\ngiven:     " + v2 +
				"\noriginal:  " + mv[k] +
				"\nunescread: " + mv2[k])
		}
	}

}


func TestFromMapCreation(t *testing.T) {
	d1 := &TStruct{NewText(""), NewInt(0)}
	svalue := "test"
	ivalue := "20"
	e, ok := MapValue(d1, map[string]string{"a": svalue, "b": "20"})
	t.Log(d1)
	if !ok {
		t.Error("NewComposite returned false")
	}

	if ele, ok := (e).(*TStruct); ok {
		ok := strings.Contains(ele.A.String(), svalue)
		if !ok {
			t.Error("String assigment fail")
		}
		v, _ := strconv.Atoi(ivalue)

		if ele.B.value != v {
			t.Error("Int assigment fail")
		}

	} else {
		t.Error("Type assertion failed", e)
	}
	t.Log("done")
}


/** ModelLike  tests section */
type Test struct {
	Someid   *Serial
	SomeText *Char
	I        *Int
}
//returns pointer to ModelValue
func NewTest(text string, i int) *Test {
	return &Test{NewSerial(0), NewChar(text), NewInt(i)}
}
// type e implemets ModelLike
type e struct{}

func (r *e) RelName() string {
	return "OK"
}

func (r *e) Fields() []string {
	return []string{"OKOK"}
}

func (r *e) SetValues(vals ...string) {}

func (r *e) Map() map[string]string {
	return make(map[string]string)
}
/* type f does not implements ModelLike */
type f struct{}

func (r *f) RelName() string {
	return "OK"
}
func (r *f) Fields() []string {
	return []string{"OKOK"}
}
func (r *f) SetValues(vals ...string) {}
func TestInterfaceCompatibility(t *testing.T) {
	var v interface{}
	v = NewSerial(0)
	if _, ok := v.(Field); !ok {
		t.Error("Serial doesn't satissfy Field interface")
	}
	v = NewChar("")
	if _, ok := v.(Field); !ok {
		t.Error("Char doesn't satissfy Field interface")
	}
	v = NewInt(0)
	if _, ok := v.(Field); !ok {
		t.Error("Int doesn't satissfy Field interface")
	}

}
func TestModelValueCheck(t *testing.T) {
	if true != IsModelValue(NewTest("", 0)) {
		t.Error("Model value check failure 1")
	}
	if true != IsModelValue(&e{}) {
		t.Error("Model value check failure1")
	}
	if true == IsModelValue(&f{}) {
		t.Error("Value that is not a ModelStruct and does not satissfy ModelLike interface passes test")
	}
}
/**/
