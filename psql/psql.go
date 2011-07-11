package psql

import (
	psql "github.com/lxn/go-pgsql/src/pkg/pgsql"
	ps "github.com/mjarco/perse"
	"time"
	"fmt"
)

/* specific psql structs */

//type that satisfies Result interface
type pgres struct {
	rs *psql.ResultSet
}

func (pr *pgres) NextMap() (res map[string]string, ok bool) {
	ok, _ = pr.rs.FetchNext()
	if ok {
		l := pr.rs.FieldCount()
		res = make(map[string]string, l)
		for i := 0; i < l; i++ {
			name, _ := pr.rs.Name(i)
			value, _, _ := pr.rs.String(i)
			res[name] = value
		}
	} else { //if its insert 
		ok, _ = pr.rs.NextResult()
		if ok {
			return pr.NextMap()
		} else {
			pr.rs.Close()
		}
	}
	return
}

func (pr *pgres) NextValue(v interface{}) (res interface{}, ok bool) {
	m, ok := pr.NextMap()
	if ok {
		return ps.CrudValue(v, m)
	}
	return

}

func (pr *pgres) AllMaps() (res []map[string]string) {
	res = make([]map[string]string, 0)
	v, ok := pr.NextMap()
	for ok {
		res = append(res, v)
		v, ok = pr.NextMap()
	}
	return res
}

func (pr *pgres) AllValues(v interface{}) (res []interface{}) {
	res = make([]interface{}, 0)
	vn, ok := pr.NextValue(ps.Copy(v))
	for ok {
		res = append(res, vn)
		vn, ok = pr.NextValue(v)
	}
	return res
}
func (pr *pgres) Close() {
	pr.rs.Close()
}

type pgcon struct {
	conndata  *ps.ConnData
	conn      *psql.Conn
	lastrs    *psql.ResultSet
	time      int64
	connected bool
}

func connstring(cd *ps.ConnData) string {
	//FIXME!!
	return fmt.Sprintf("dbname=%s user=%s password=%s", cd.Name, cd.User, cd.Pass)
}

func (pc *pgcon) Connect(cd ...*ps.ConnData) bool {
	if len(cd) > 0 {
		pc.conndata = cd[0]
	}
	conn, err := psql.Connect(connstring(pc.conndata), psql.LogError)
	if err != nil {
		return false
	}
	pc.conn = conn
	pc.connected = true
	pc.time = time.Seconds()
	return true
}

func (pc *pgcon) Query(stm *ps.Stmt) ps.Result {
	if pc.conn.Status() != psql.StatusReady {
		pc.lastrs.Close()
	}
	if !pc.connected {
		panic("Query on not established connection")
	}

	rsp, err := pc.conn.Query(stm.Query)
	if err != nil {
		return nil
	}
	pc.lastrs = rsp
	return &pgres{rsp}
}
func (pc *pgcon) LastId(r ps.Result) string {
	res, _ := r.NextMap()
	return res["id"]
}
func (pc *pgcon) Last() *time.Time {
	return time.SecondsToLocalTime(pc.time)
}

func (pc *pgcon) CrudBuilder() ps.CRUDStatementBuilder {
	return ps.DefaultCrudBuilder
}

func NewConnection() ps.Connection {
	return new(pgcon)
}


var Conn ps.Connection = NewConnection()

func init() {
	ps.Conn = Conn
}
