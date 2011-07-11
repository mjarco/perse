package perse

import (
	"strconv"
	"time"
	"os"
	"strings"
	//    "fmt"
)

type dirtTeller struct {
	dirty bool
}

func (dt *dirtTeller) setDirty(a bool) {
	dt.dirty = a
}

func (dt *dirtTeller) isDirty() bool {
	return dt.dirty
}

/** 
  Char type is counterpart of VARCHAR or CHAR types in RDBMS
  It realizes perse.Field and perse.Escaper interfaces
  Warning: escaping is not completed for now Char and Text field escapes only ' and "
*/
type Char struct {
	value  string
	length int
	*dirtTeller
}

func (c *Char) String() string {
	return c.value
}

/*
   Warning: escaping is not completed for now Char and Text field escapes only ' and "
*/

func (c *Char) Escape() string {
	//TODO: \x00, \n, \r, \, ', " and \x1a
	return "'" + strings.Replace(c.String(), `'`, `\'`, -1) + "'"
}
/*
   Warning: escaping is not completed for now Char and Text field escapes only ' and "
*/
func (c *Char) UnEscape(s string) {
	c.Value(strings.Replace(s, `\'`, `'`, -1))
	c.setDirty(false)
}

func (c *Char) Value(v string) (ok bool) {
	ok = false
	if len(v) <= c.length {
		c.value = v
		ok = true
		c.setDirty(true)
	}
	return
}

func (c *Char) SetMaxLen(i int) (ok bool) {
	ok = false
	if i < 256 {
		ok = true
		c.length = i
	}
	return
}

func NewChar(s string, l ...int) *Char {
	//FIXME: consider this code once again....
	if len(l) > 0 {
		if len(s) > l[0] {
			return nil
		}
		return &Char{s, l[0], &dirtTeller{false}}
	}

	return &Char{s, len(s), &dirtTeller{false}}
}

type Text struct {
	value string
	*dirtTeller
}

func (t *Text) String() string {
	return t.value
}

func (t *Text) Value(v string) (ok bool) {
	t.value = v
	t.setDirty(true)
	return true
}
/*
   Warning: escaping is not completed for now Char and Text field escapes only ' and "
*/

func (t *Text) Escape() string {
	return "'" + strings.Replace(t.String(), `'`, `\'`, -1) + "'"
}
/*
   Warning: escaping is not completed for now Char and Text field escapes only ' and "
*/

func (t *Text) UnEscape(s string) {
	if strings.HasPrefix(s, `'`) && strings.HasSuffix(s, `'`) {
		s = string(s[1 : len(s)-1])
	}
	t.Value(strings.Replace(s, `\'`, `'`, -1))
	t.setDirty(false)
}


func NewText(s string) *Text {
	return &Text{s, &dirtTeller{false}}
}

type Int struct {
	value int
	*dirtTeller
}

func (i *Int) String() string {
	return strconv.Itoa(i.value)
}
//Set converted string as a value of Int
func (i *Int) Value(v string) (ok bool) {
	i.value, _ = strconv.Atoi(v)

	i.setDirty(true)
	return true
}
//Set int value
func (i *Int) ValueInt(v int) (ok bool) {
	i.value = v

	i.setDirty(true)
	return true
}
//Create new struct
func NewInt(i int) *Int {
	return &Int{i, &dirtTeller{false}}
}

type Serial struct {
	value uint
	*dirtTeller
}

func (s *Serial) String() string {
	return strconv.Uitoa(s.value)
}

func (s *Serial) Value(v string) (ok bool) {
	s.value, _ = strconv.Atoui(v)

	s.setDirty(true)
	return true
}

func (s *Serial) ValueInt(v uint) (ok bool) {
	s.value = v

	s.setDirty(true)
	return true
}
func NewSerial(s uint) *Serial {
	return &Serial{s, &dirtTeller{false}}
}

type Date struct {
	format string
	value  *time.Time
	*dirtTeller
}

func (d *Date) Valuef(v string, f ...string) {
	var t *time.Time
	var err os.Error
	if len(f) > 0 {
		t, err = time.Parse(f[0], v)
	} else {
		t, err = time.Parse(d.format, v)
	}
	if err == nil {
		d.value = t
	}
}

func (d *Date) Value(v string) (ok bool) {
	t, err := time.Parse(d.format, v)
	if err != nil {
		println(err.String())
	}
	d.value = t
	return err == nil
}

func (d *Date) Stringf(f ...string) string {
	if len(f) > 0 {
		return d.value.Format(f[0])
	}
	return d.value.Format(d.format)
}

func (d *Date) String() string {
	return d.value.Format(d.format)
}
func (d *Date) Escape() string {
	return `'` + d.value.Format("2006-01-02 15:04:05") + `'`
}

func (d *Date) UnEscape(s string) {
	s = strings.Trim(s, `'`)
	d.value, _ = time.Parse("2006-01-02 15:04:05", s)
}

// Yep! default formater is mandatory
// Just remember about std go formatting, which is... hmm... unique 
func NewDate(format string, value ...string) *Date {
	d := &Date{format, new(time.Time), &dirtTeller{false}}
	if len(value) > 0 {
		d.Value(value[0])
	} else {
		d.value = time.LocalTime()
	}
	return d
}
