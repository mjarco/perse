Perse - concepts
================

This document describes basic therms and features of perse package.


perse is not an ORM. Its more like key-value storage helper where key is uint64 value. Perse will not support cross table queries and "foreign key"-like fetures, at least not now and it''s not in  plans.

ModelType
---------

There is no such thing as ModelType per se, but as a user you can define struct types that will map database tables structure into your program. This struct types are called ModelType. 
Example of ModelType:
    
    import "perse"
    type struct Tutorial {
        Title    *perse.Char
        Body     *perse.Text
        Date     *perse.Date
        Id       *perse.Serial
    }

Good practice in Go is to write initializer functions:

    func NewTutorial() *Tutorial {
        return &Tutorial{perse.NewChar("", 255), perse.NewText(""), perse.NewDate("2006-01-02 15:04:05"), perse.Serial(0)}
    }

Notice that field types satisfies perse.Field interface.

perse.Field
------------
All exported fields in ModelType struct that have corresponding column in table must implement perse.Field interface. parse.Field consists four methods:
    String() string
    Value(string) bool
    isDirty() bool
    setDirty(bool)

isDirty method is called when ModelType value is passed to perse.Save and it determines which columns have to be updated
setDirty should be called if field value was changed, as default dirty flag is set true on each Value() call, only exception to this rule is setting values from db, we explicitly set it to false.


If fields in ModelType struct should be serialized in speciffic way they should implement perse.Escaper interface. 

At the moment we have serveral types that satisfies perse.Field interface:
*  Char
*  Text
*  Int
*  Serial
*  Date

First three of them are obvious what they do and how they behave (see description in api documentation)... So lets focus on Serial and Date.

perse.Serial & primary keys
---------------------------
There should be only one *Serial type field in ModelType struct. This field is treated as mysql''s "Primary key autoincrement" or postgresql Serial column. So if you want to fetch a speciffic row from "tutorial" table you should pass *Tutorial value with proper Id field value to perse.Get func.

perse.Date
----------

perse.Date type internally stores *time.Time value so you have to pass date format in initialization phase. If you''re going to implement perse db driver you should be 


ModelType values convertion
---------------------------
It's one of most important features in perse package. With functions listed below it's possible to convert ModelType values into map[string]string values and other way round.
*  map[string]string => value:
    *  MapValue
    *  CrudValue
*  value => map[string]string:
    *  ModelMap
    *  EscapedMap

Example:
    mm := map[string]string{"title":"Tutorial title", "body":"Sample tutorial body", "date":"2011-02-01 08:14:00", "id":"0"}
    v := NewTutorial()
    perse.MapValue(v, mm)// Title, Body and Date fields have now real values
    //now we''ll transform *Tutorial into map
    mm2, _ := perse.ModelMap(v)
    //in this place in code mm and mm2 should have same keys and values

In above example we used only MapValue and ModelMap funcs, the others are simmilar only difference they use Escape/Unescape in favor of String/Value methods


Connection
----------

Each driver package have NewConnection() method which returns perse.Connection interface value. To establish real connection with database use Connect(conndata) method. It''s up to driver implementation which fields in ConnData struct have to be filled.

Once you created connection you can register it in perse package by assigning Conn value to perse.Conn package variable, later on (in runtime) it will be default connection for all perse CRUD and collection queries.
    
    import (
        ps "perse"
        psql "perse/psql" 
        "model"
    )

    func main() {
        conn := psql.NewConnection()
        ok := ps.Conn.Connect(ps.ConnData{"localhost","3456", "tutorial","tutorial","tutorial",""})
        if (!ok) {
            println("connection not established")
            exit(1)
        }
        ps.Conn = conn //setting conn as a default connection in perse package
        //from now on we don''t have to pass connection parameter to perse funcs
        //so assumming that in model package Tutorial type is defined you can do:
        v := model.NewTutorial()
        v.Section.Value("first record!")
        v.Body.Value("hello world!")
        perse.Save(v)
        //or alternatively
        mm := map[string]string{"title":"Tutorial title", "body":"Sample tutorial body", "date":"2011-02-01 08:14:00", "id":"0"}
        v2 := NewTutorial()
        perse.MapValue(v2, mm)// Title, Body and Date fields have now real values
        perse.Save(v2)

    }

CRUD:  perse.Save [CU], perse.Get [R], perse.Delete [D]
-------------------------------------------------------
All CRUD functions takes a ModelType value and connection in parameters, as mentioned above connection is optional if default connection (perse.Conn variable) was set.

Save, Get and Delete functions realize CRUD. 
Save make inserts and updates depending if there is non-zero value in Serial field otherwise inserts new record (maybe it's not the best solution but good enough for start).

Get fetches row with id given in Serial field of value and fill this value with proper data. 
Delete removes record with id given in Serial field of value from table

Collections & filters
---------------------

As far we were taking about ModelTypes, but perse has some functionalities that customize database querying. Collections are one of them. If you need something more than one record from table you can construct collection.

    import (
        ps "perse", 
        psql "perse/psql", 
        "strconv", 
        "fmt")
    /**
     *   Connecting to database
     */
    var cd = ps.ConnData
    func init() {
        cd = new(ps.ConnData)
        cd.User = "username"
        cd.Pass = "password"
        cd.Name = "dbname" 
        ps.Conn = psql.NewConnection()//default connection for all method of perse package
    }

    type Test struct { //name of type is name of DB relation by default
        TextField *ps.Char 
        TestId   *ps.Serial

    }

    func NewTest(text string) *Test {
        return &Test{ps.NewText(text), ps.NewSerial(0)}
    }

    func create10records() {
        for i:=0; i<10; i++ {
            ps.Save(NewTest("record"+strconv.Itoa(i));
        }
    }

    func collections() {
        c := NewCollection(NewTest(""))
        res = c.Query() //will fetch all test records
        fmt.Println(res.AllMaps()) 
        fmt.Println("-------------------------------------") 
        resLT5 = c.Filter(ps.NewCond("test_id", ps.LT, "5")).Query()
        fmt.Println(resLT5.AllMaps()) 
    }

    func main() {
        ok := ps.Conn.Connect(cd)
        if (!ok) {
            println("connection not established")
            exit(1)
        }
        create10records()
        collections()
    }

Collection may be created in two other ways. Statemets below are equivalent

1.  c := NewCollection(NewTest(""))
2.  c := NewCollection("test", "testid", "textfield")
3.  
    
    type testML struct{ fieldNames []string; fieldValues []string}
    func(m *testML) SetValues(vals ... string) { m.fieldValues = vals }
    func(m *testML) Fields() []string { return fieldNames }
    func(m *testML) RelName() string { return "test" }
    c := NewCollection(&testML{[]string{"test_id","textfield"}, make([]string, 0))

Custom Queries and Results
--------------------------

Everything above should work in all RDBMSes, however sometimes you need do something that is not included in this set of functionalities. Then you may create custom Query. As a returning value you will get value that implements Result interface. 


