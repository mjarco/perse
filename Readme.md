Getting started
===============

Simple CRUD:
------------


    package main

    import (ps "perse", psql "perse/psql")
    /*
        Connecting to database
    */
    var Conn ps.Connection
    func init() {
        cd = new(ps.ConnData)
        cd.User = "username"
        cd.Pass = "password"
        cd.Name = "dbname" 
        Conn = psql.NewConnection()
        Conn.Connect(cd)
    }

    /*
        Model declaration
        This is recommended model declaration.
    */

    type Test struct { //name of type is name of DB relation by default
        TextField *ps.Char //attributes have to be exported fields with one of perses Model Attributes Types (ps.MAT) pointer(!)
        Test_id   *ps.Serial //primary auto increement field
        smth   string // model types can have additional non ps.MAT fields, and any methods you want, perse doesn't bother.
    }

    // Before you run this example please make sure that you have relation test(textfield Char(255), test_id Serial) in your RDBMS
    func NewTest(text string) *Test {
        t = new(Test)
        t.TextField = ps.NewChar(text)
        return t
    }

    func main() {
        //Create
        v := NewTest("first")
        ps.Save(v, Conn) //Save takes value and Connection instance

        //Tada! Done. Record should be stored in database. So how about retriving it from db? 
        //Primary method is geting it by id
        v1 = &Test{ps.NewChar("", 255), ps.NewSerial(1)}
        ps.Get(v1, Conn)
        //For more info about "Read" queries go to "Collections" section

        //Update
        v.TextField.Value("updated")
        ps.Save(v, Conn) //yep! Save method knows if its update or create call
        //Delete
        ps.Delete(v, Conn)
    }


Custom queries
--------------


Collections
-----------

Basic example

    package main
    import (ps "perse", psql "perse/psql", "strconv", "fmt")
    type Test struct { //name of type is name of DB relation by default
        TextField *ps.Char 
        Test_id   *ps.Serial
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
        resLT5 = c.Filter(ps.NewCond("test_id", ps.LT, "5"))
        fmt.Println(resLT5.AllMaps()) 
    }

    func main() {
        create10records()
        collections()
    }

Collection may be created in two other ways. Statemets below are equivalent 

1)  c := NewCollection(NewTest(""))
2)  c := NewCollection("test", "test_id", "text_field")
3)  type testML struct{ fieldNames []string; fieldValues []string}
    func(m *testML) SetValues(vals ... string) { m.fieldValues = vals }
    func(m *testML) Fields() []string { return fieldNames }
    func(m *testML) RelName() string { return "test" }
    c := NewCollection(&testML{[]string{"test_id","textfield"}, make([]string, 0))
