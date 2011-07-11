package perse

import (
	"strings"
)

func insertStmt(name string, mapp map[string]string, pk string) *Stmt {
	keys, values := GetKVSlices(mapp, pk, false)
	stmt := "INSERT INTO " + name + "(" + strings.Join(keys, ",") + ") VALUES (" +
		strings.Join(values, ",") + ") returning " + pk + " as id;"
	return &Stmt{keys, stmt}
}

func selectStmt(name, pk string, mapp map[string]string) *Stmt {
	keys, _ := GetKVSlices(mapp, pk, true)
	stmt := "SELECT " + strings.Join(keys, ",") +
		" FROM " + name + " WHERE " + pk + "=" + mapp[pk] + ";"
	return &Stmt{keys, stmt}

}

func updateStmt(name, pk string, mapd, mapp map[string]string) *Stmt {
	var changes []string
	keys, _ := GetKVSlices(mapd, pk, true)
	for _, val := range keys {
		changes = append(changes, val+"="+mapp[val])
	}
	stmt := "UPDATE " + name + " set " + strings.Join(changes, ",") + " WHERE " +
		pk + "=" + mapp[pk]
	return &Stmt{keys, stmt}

}

func deleteStmt(name, pk string, mapp map[string]string) *Stmt {
	keys, _ := GetKVSlices(mapp, pk, true)
	stmt := "DELETE FROM " + name + " WHERE " + pk + "=" + mapp[pk] + ";"
	return &Stmt{keys, stmt}
}

type Cb struct{}

func (c *Cb) Create(v interface{}) *Stmt {
	name, pk, mapp, _, _ := CrudInfo(v)
	return insertStmt(name, mapp, pk)
}

func (c *Cb) Read(v interface{}) *Stmt {
	name, pk, mapp, _, _ := CrudInfo(v)
	return selectStmt(name, pk, mapp)
}

func (c *Cb) Update(v interface{}) *Stmt {
	name, pk, mapp, mapd, _ := CrudInfo(v)
	return updateStmt(name, pk, mapd, mapp)

}

func (c *Cb) Delete(v interface{}) *Stmt {
	name, pk, mapp, _, _ := CrudInfo(v)
	return deleteStmt(name, pk, mapp)
}

var DefaultCrudBuilder CRUDStatementBuilder

func init() {
	DefaultCrudBuilder = &Cb{}
}
// Save ModelValue v in database, conn parameter is optional if not given uses default Connection instance
func Save(v interface{}, conn ...Connection) (val interface{}, id string) {
	c := getConnection(conn)
	//Fixme: add panic or smth
	//Error handling
	cr := c.CrudBuilder()
	_, pk, mapp, _, _ := CrudInfo(v)
	if len(pk) > 0 {
		if pp, ok := mapp[pk]; ok && (pp != "0") {
			c.Query(cr.Update(v))
			return v, mapp[pk]
		} else {
			stmt := cr.Create(v)
			res := c.Query(stmt)
			id = c.LastId(res)
			mapp[strings.ToLower(pk)] = id
			v, _ = CrudValue(v, mapp)
			return v, id
		}
	} else {

	}
	return
}

func Get(v interface{}, conn ...Connection) interface{} {
	c := getConnection(conn)
	cr := c.CrudBuilder()
	rs := c.Query(cr.Read(v))
	mmap, _ := rs.NextMap()
	ret, _ := CrudValue(v, mmap)
	return ret
}

func Delete(v interface{}, conn ...Connection) bool {
	c := getConnection(conn)
	cr := c.CrudBuilder()
	rs := c.Query(cr.Delete(v))
	_, _ = rs.NextMap()
	return true
}

func Query(query string, colnames []string, conn ...Connection) Result {
	c := getConnection(conn)
	return c.Query(&Stmt{colnames, query})
}
