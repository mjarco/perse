package perse

import (
	"time"
)

type ConnData struct {
	Host     string
	Port     string
	User     string
	Pass     string
	Name     string
	Location string
}

/*func NewConnData(c map[string]string) *ConnData {

}*/

type Stmt struct {
	Fields []string
	Query  string
}

type CRUDStatementBuilder interface {
	// insert statement
	//v is ModelLike or ModelValue
	Create(v interface{}) *Stmt
	// select where PK = id
	//v is ModelLike or ModelValue
	Read(v interface{}) *Stmt
	// update
	//v is ModelLike or ModelValue
	Update(v interface{}) *Stmt
	//delete
	//v is ModelLike or ModelValue
	Delete(v interface{}) *Stmt
}

type Connection interface {
	Connect(...*ConnData) bool
	Query(*Stmt) Result
	Last() *time.Time
	LastId(Result) string
	CrudBuilder() CRUDStatementBuilder
}

type ConnPool struct {
	cons []Connection
}

type Result interface {
	//Returns map which keys are column names and values are stringified data from row, second value tells if operation was successful
	NextMap() (res map[string]string, ok bool)
	//Given ModelType value is filled up with data from row and returned
	NextValue(v interface{}) (res interface{}, ok bool)
	//Returns slice of maps with all fetched records
	AllMaps() (res []map[string]string)
	//returns slice of values filled up with data from all fetched rows
	AllValues(v interface{}) (res []interface{})
}

type Cond interface {
	IsValid() bool
	String() string
	And(c Cond) Cond
	Or(c Cond) Cond
	Not() Cond
}


type Collection interface {
	Filter(c Cond) Collection
	Sort(fields ...string) Collection
	Limit(limit int64) Collection
	Offset(offset int64) Collection
	Query(...Connection) Result
}


type Field interface {
	String() string
	Value(string) (ok bool)
	setDirty(bool)
	isDirty() bool
}

type escaper interface {
	Escape() string
	UnEscape(string)
}

type ModelLike interface {
	RelName() string
	Fields() []string
	SetValues(vals ...string)
	Map() map[string]string
}
