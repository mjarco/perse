package mysql

import (
	mysql "github.com/Philio/GoMySQL"
	ps "perse"
	"time"
	"fmt"
	"strings"
)

/* specific psql structs */

//type that satisfies Result interface
type myres struct {
	rs *mysql.Result
	c  *mysql.Client
}

func (pr *myres) NextMap() (res map[string]string, ok bool) {
	fields := pr.rs.FetchFields()
	mm := pr.rs.FetchRow()
	if mm != nil {
		l := len(mm)
		res = make(map[string]string, l)
		for i, field := range fields {
			value := fmt.Sprint(mm[i])
			res[field.Name] = value
		}
		return res, true
	}
	return nil, false
}

func (pr *myres) NextValue(v interface{}) (res interface{}, ok bool) {
	m, ok := pr.NextMap()
	if ok {
		return ps.CrudValue(v, m)
	}
	return

}

func (pr *myres) AllMaps() (res []map[string]string) {
	res = make([]map[string]string, 0)
	v, ok := pr.NextMap()
	for ok {
		res = append(res, v)
		v, ok = pr.NextMap()
	}
	pr.rs.Free()
	return res
}

func (pr *myres) AllValues(v interface{}) (res []interface{}) {
	res = make([]interface{}, 0)
	vn, ok := pr.NextValue(ps.Copy(v))
	for ok {
		res = append(res, vn)
		vn, ok = pr.NextValue(v)
	}
	return res
}
func (pr *myres) Close() {
	pr.rs.Free()
}

type mycon struct {
	conndata  *ps.ConnData
	conn      *mysql.Client
	lastrs    *mysql.Result
	time      int64
	connected bool
}

func (pc *mycon) Connect(cd ...*ps.ConnData) bool {
	if len(cd) > 0 {
		pc.conndata = cd[0]
	}

	conn, err := mysql.DialTCP(pc.conndata.Host, pc.conndata.User, pc.conndata.Pass, pc.conndata.Name)

	if err != nil {
		return false
	}
	pc.conn = conn
	pc.connected = true
	pc.time = time.Seconds()
	return true
}

func (pc *mycon) Query(stm *ps.Stmt) ps.Result {
	if pc.lastrs != nil {
		_ = pc.conn.FreeResult()
	}
	if !pc.connected {
		panic("Query on not established connection")
	}

	err := pc.conn.Query(stm.Query)
	if err != nil {
		panic(err)
	}
	rsp, _ := pc.conn.UseResult()
	pc.lastrs = rsp
	return &myres{rsp, pc.conn}

}
func (pc *mycon) LastId(r ps.Result) string {
	return fmt.Sprint(pc.conn.LastInsertId)
}
func (pc *mycon) Last() *time.Time {
	return time.SecondsToLocalTime(pc.time)
}

type MyCrudBuilder struct {
	ps.Cb
}

func insertStmt(name string, mapp map[string]string, pk string) *ps.Stmt {
	keys, values := ps.GetKVSlices(mapp, pk, false)
	stmt := "INSERT INTO " + name + "(" + strings.Join(keys, ",") + ") VALUES (" +
		strings.Join(values, ",") + ");"
	return &ps.Stmt{keys, stmt}
}

func (cb *MyCrudBuilder) Create(v interface{}) *ps.Stmt {
	name, pk, mapp, _, _ := ps.CrudInfo(v)
	return insertStmt(name, mapp, pk)
}

func (pc *mycon) CrudBuilder() ps.CRUDStatementBuilder {
	return &MyCrudBuilder{}
}

func NewConnection() ps.Connection {
	return new(mycon)
}


var Conn ps.Connection = NewConnection()

func init() {
	ps.Conn = Conn
}
