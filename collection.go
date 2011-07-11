package perse

import (
	"strconv"
	"container/list"
	"strings"
)

type Operator int

const (
	AND Operator = iota
	OR
	NOT
)

func (o *Operator) String() string {
	switch *o {
	case AND:
		return " AND "
	case OR:
		return " OR "
	case NOT:
		return " NOT "
	}
	return ""
}


type Relation int

const (
	EQ Relation = iota
	NEQ
	IN
	LT
	GT
	BETWEEN
	LIKE
)

type condatom struct {
	fname    string
	fvalues  []string
	relation Relation
}

func (c *condatom) IsValid() (ok bool) {
	if len(c.fname) < 1 {
		return false
	}
	vallen := len(c.fvalues)
	if vallen < 1 {
		return false
	}
	switch c.relation {
	case EQ, NEQ, LT, GT:
		return vallen == 1
	case IN:
		return vallen >= 1
	case BETWEEN:
		return vallen == 2
	}
	return true
}

// returns sql WHERE clause
func (c *condatom) String() string {
	switch c.relation {
	case EQ:
		return c.fname + " = " + c.fvalues[0]
	case NEQ:
		return c.fname + " <> " + c.fvalues[0]
	case LT:
		return c.fname + " < " + c.fvalues[0]
	case GT:
		return c.fname + " > " + c.fvalues[0]
	case IN:
		return c.fname + " IN " + c.fvalues[0]
	case BETWEEN:
		return c.fname + " BETWEEEN " + c.fvalues[0] + " AND " + c.fvalues[1]
	case LIKE:
		return c.fname + " LIKE " + c.fvalues[0]
	}
	return ""

}

func check(c *cond, nc Cond) bool {
	if v, ok := nc.(*cond); ok {
		//FIXME: panic?
		return v != c
	}
	return true
}

type cond struct {
	atos *list.List
}

func (c *cond) And(ncond Cond) Cond {
	if check(c, ncond) {
		c.atos.PushBack(AND)
		c.atos.PushBack(ncond)
	}
	return c
}
func (c *cond) Or(ncond Cond) Cond {
	if check(c, ncond) {
		c.atos.PushBack(OR)
		c.atos.PushBack(ncond)
	}
	return c
}

func (c *cond) Not() Cond {
	c.atos.PushBack(NOT)

	return c
}

func (c *cond) IsValid() bool {
	//FIXME:
	return true
}

func (c *cond) String() string {
	result := ""
	stckcnt := 0
	for e := c.atos.Front(); e != nil; e = e.Next() {
		switch v := e.Value.(type) {
		case Operator:
			if v == NOT {
				result = v.String() + "(" + result + ")"
			} else {
				result += v.String() + "("
				stckcnt++
			}

		case *condatom:
			result += v.String()
			if stckcnt > 0 {
				result += ")"
				stckcnt--
			}

		case Cond:
			result += v.String()
			if stckcnt > 0 {
				result += ")"
				stckcnt--
			}

		}
	}
	for ; stckcnt > 0; stckcnt-- {
		result += " ) "
	}
	return result
}

//Conditions compares only field to concrete values
//There is no possibility to compare fields with each other, by now...
func NewCond(field string, r Relation, vals ...string) Cond {
	l := list.New()
	l.PushBack(&condatom{field, vals, r})

	return &cond{l}
}

type CustomML struct {
	name   string
	fields []string
	values []string
}

func (cml *CustomML) RelName() string {
	return cml.name
}

func (cml *CustomML) Fields() []string {
	return cml.fields
}

func (cml *CustomML) SetValues(val ...string) {
	cml.values = val
}

func (cml *CustomML) Map() map[string]string {
	///
	ret := make(map[string]string)
	for i, v := range cml.fields {
		ret[v] = cml.values[i]
	}
	return ret
}


type collection struct {
	model  ModelLike
	conds  Cond
	order  []string
	limit  int64
	offset int64
}

func (c *collection) Filter(conds Cond) Collection {
	c.conds = conds
	return c
}

func (c *collection) Sort(fields ...string) Collection {
	c.order = fields
	return c
}

func (c *collection) Limit(limit int64) Collection {
	c.limit = limit
	return c
}

func (c *collection) Offset(offset int64) Collection {
	c.offset = offset
	return c

}

func (c *collection) Query(co ...Connection) Result {
	conn := getConnection(co)
	q := "SELECT " + strings.Join(c.model.Fields(), ",") + " FROM " + c.model.RelName()
	cs := c.conds.String()

	if len(cs) > 0 {
		q += " WHERE " + c.conds.String()
	}

	if len(c.order) > 0 {
		q += "ORDER BY " + strings.Join(c.order, ",")
	}
	if c.limit > 0 {
		q += "LIMIT " + strconv.Itoa64(c.limit)
	}
	return conn.Query(&Stmt{c.model.Fields(), q})
}

func NewCollection(v interface{}, fields ...string) Collection {
	var ml ModelLike
	switch nv := v.(type) {
	case ModelLike:
		ml = nv
	case string:
		ml = &CustomML{nv, fields, nil}
	default:
		if IsModelValue(nv) {
			name, _, mapa, _, _ := CrudInfo(v)
			fields, _ := GetKVSlices(mapa, "", false)
			ml = &CustomML{name, fields, nil}
		} else {
			return nil
		}

	}
	c := &collection{ml, nil, nil, 0, 0}
	return c
}
