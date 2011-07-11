package perse

import (
	"reflect"
	"strings"
)
//Returns copy of value
func Copy(src interface{}) (dst interface{}) {
	fmval := reflect.ValueOf(src)
	if ptr := fmval; ptr.Kind() == reflect.Ptr {
		/* value copy */
		interf := ptr.Elem().Interface()
		copied := Copy(interf)
		/* pointer creation */
		newptr := reflect.New(reflect.ValueOf(copied).Type()) //pointer to a value
		/* value assigment */
		newptr.Elem().Set(reflect.ValueOf(copied)) //set pointer value, to copy value
		dst = newptr.Interface()
	} else if sct := fmval; sct.Kind() == reflect.Struct {
		newsct := reflect.New(sct.Type())
		for i := 0; i < sct.NumField(); i++ {

			if newsct.Elem().Field(i).CanSet() {
				newsct.Elem().Field(i).Set(reflect.ValueOf(Copy(sct.Field(i).Interface())))
			}
		}
		dst = newsct.Elem().Interface()

	} else {
		dst = reflect.Zero(fmval.Type()).Interface()
	}
	return
}


func GetKVSlices(mapp map[string]string, pk string, withpk bool) (keys []string, values []string) {
	keys = make([]string, 0)
	values = make([]string, 0)
	for key, val := range mapp {
		if (withpk == false) && (strings.ToLower(key) == strings.ToLower(pk)) {
			continue
		}
		if (len(key) == 0) || (len(val) == 0) {
			continue
		}
		keys = append(keys, key)
		values = append(values, val)
	}
	return
}

type structMapFun func(reflect.StructField, interface{}) bool


func fieldIter(v interface{}, funs ...structMapFun) (typename string, ok bool) {
	value := reflect.ValueOf(v)
	if ptr := value; ptr.Kind() == reflect.Ptr {
		//well argument was a pointer, but we need a value
		v := ptr.Elem().Interface()
		//lets try again
		return fieldIter(v, funs...)
	}

	sct := value
	if !(sct.Kind() == reflect.Struct) { //maybe return sth more verbose
	    panic("v is not a struct nor pointer to struct")	
	}
	typ := sct.Type()
	typename = strings.ToLower(typ.Name())
	l := sct.NumField()

	for i := 0; i < l; i++ {
		ftype := typ.Field(i)
		fval := sct.Field(i)
		if fval.IsNil() {
			break
		}
		fvalintf := fval.Interface()

		ok = true
		for _, fun := range funs {
			ok = ok && fun(ftype, fvalintf)
		}
	}
	return
}

func produceMapAll(mapa *map[string]string) structMapFun {
	//println("return closure mapall")
	return func(field reflect.StructField, value interface{}) bool {
		val, ok := value.(Field)
		if ok {
			(*mapa)[strings.ToLower(field.Name)] = val.String()
		}
		return true
	}
}

func produceMapDirty(mapa *map[string]string) structMapFun {
	return func(field reflect.StructField, value interface{}) bool {
		if val, ok := value.(Field); ok {
			if val.isDirty() {
				(*mapa)[strings.ToLower(field.Name)] = val.String()
			}
		}
		return true
	}
}

func producePK(f *string) structMapFun {
	return func(field reflect.StructField, value interface{}) bool {
		if _, ok := value.(*Serial); ok {
			*f = strings.ToLower(field.Name)
		}
		return true
	}
}

type Escape func(v Field) string
type UnEscape func(s string, v Field) bool
/** 
Returns all info nedded for CRUD funcs that is:
   * name of type (name of db relation)
   * primary key field name
   * map field name => escaped field value 
   * map of dirty field name => escaped field value
   panics if value is not ModelValue
*/
func CrudInfo(value interface{}, es ...Escape) (typename, pk string, mapa, mapd map[string]string, ok bool) {
	mapa = make(map[string]string)
	mapd = make(map[string]string)
	var esc Escape
	if len(es) > 0 {
		esc = es[0]
	}
	typename, ok = fieldIter(value, producePK(&pk), produceMapEsc(&mapa, esc), produceMapDirtyEsc(&mapd, esc))
	return
}

func DefaultEsc(value Field) string {
	if valesc, ok := value.(escaper); ok {
		return valesc.Escape()
	}
	return value.String()
}

func produceMapDirtyEsc(mapa *map[string]string, esc Escape) structMapFun {
	return func(field reflect.StructField, value interface{}) bool {
		if val, ok := value.(Field); ok {
			if val.isDirty() {
				if esc != nil {
					(*mapa)[strings.ToLower(field.Name)] = esc(val)
				} else {
					(*mapa)[strings.ToLower(field.Name)] = DefaultEsc(val)
				}
			}
		}
		return true
	}
}

func produceMapEsc(mapa *map[string]string, esc Escape) structMapFun {
	return func(field reflect.StructField, value interface{}) bool {
		if val, ok := value.(Field); ok {
			if esc != nil {
				(*mapa)[strings.ToLower(field.Name)] = esc(val)
			} else {
				(*mapa)[strings.ToLower(field.Name)] = DefaultEsc(val)
			}
		}
		return true
	}
}
//Collects escaped strings from ModelType value fields and returns map with lowercase keys
func EscapedMap(value interface{}, es ...Escape) (mapa map[string]string, ok bool) {
	mapa = make(map[string]string)
	var esc Escape
	if len(es) > 0 {
		esc = es[0]
	}
	_, _ = fieldIter(value, produceMapEsc(&mapa, esc))
	return
}
//Collects strings from ModelType value fields and returns map with lowercase keys
func ModelMap(value interface{}) (mapa map[string]string, ok bool) {
	mapa = make(map[string]string)
	//print ("beeeen here \n")
	_, ok = fieldIter(value, produceMapAll(&mapa))
	return
}

//
type fillFunc func(string, Field) bool

//Accepts pointer to model struct and fills its fields with values
//TODO: ?rewrite into fieldIter? have to decide
func pointerIter(zeropointer interface{}, data *map[string]string, funs ...fillFunc) (interface{}, bool) {
	rval := reflect.ValueOf(zeropointer)
	ptr := rval
	isPointer := ptr.Kind() == reflect.Ptr
	var sct reflect.Value
	var ok bool = false
	if isPointer {
		sct = reflect.Indirect(ptr)
		ok = sct.Kind() == reflect.Struct
		if !ok { //if pointer does not point to structvalue it's a failure
			//println("returning nil zeropointer was a pointer")
            panic("zeropointer does not points to struct")
		}
	} else {
		sct = rval
	}

	typ := sct.Type()
	if !(sct.Kind() == reflect.Struct) {
        panic("zeropointer is not a struct")
	}
	l := typ.NumField()
	var field reflect.StructField
	for i := 0; i < l; i++ {
		field = typ.Field(i)
		tval := sct.Field(i).Interface()
		if rat, ok := tval.(Field); ok {
			for _, fun := range funs {
				ok = ok && fun((*data)[strings.ToLower(field.Name)], rat)
				//    ok = rat.Value()
			}
		}
	}
	if isPointer {
		return ptr.Interface(), ok
	}
	return sct.Interface(), ok
}

func DefaultUnEsc(stringValue string, val Field) bool {
	if esc, ok := val.(escaper); ok {
		esc.UnEscape(stringValue)
		return true //FixMe: -> change all escape methods in Field types to return bool
	}
	return valueFill(stringValue, val)
}

func unescapeFill(unesc UnEscape) fillFunc {
	return func(stringValue string, val Field) bool {
		if unesc != nil {
			return unesc(stringValue, val)
		}
		return DefaultUnEsc(stringValue, val)
	}
}

func valueFill(stringValue string, val Field) bool {
	return val.Value(stringValue)
}

//takes ModelType value to fill it up with unescaped data
//panics if value is not a ModelType or *ModelType
func MapValue(value interface{}, mapp map[string]string, copy ...bool) (interface{}, bool) {
	if len(copy) > 0 && copy[0] {
		value = Copy(value)
	}
	_, ok := pointerIter(value, &mapp, valueFill)
	return value, ok
}

//takes ModelType value to fill up with escaped data
//panics if value is not a ModelType or *ModelType
func CrudValue(value interface{}, mapp map[string]string, ues ...UnEscape) (v interface{}, ok bool) {
	var unesc UnEscape
	if len(ues) > 0 {
		unesc = ues[0]
	}
	_, ok = pointerIter(value, &mapp, unescapeFill(unesc))
	return value, ok
}

//Checks if value (or pointer to value) is able to be passed to "read" functions
// if second argument is true it also checks against funs that modyfies data
func IsModelValue(val interface{}, saveable ...bool) (ok bool) {

	if _, ok := val.(ModelLike); ok {
		if len(saveable) > 0 && saveable[0] {
			return false
		}
		return true
	}
	return isModelV(val) //, saveable...)
}
func produceIsModelFun(ret *bool) structMapFun {
	*ret = false
	/*
		if checksave {
			_, ok = rat.(*Serial)
			if ok {
				//may be should be if !ok {continue}, at least it's less code
				ok = true
				break
			}
			ok = false
		} else {
			ok = true
			break
		}
	*/
	return func(field reflect.StructField, value interface{}) bool {
		if _, ok := value.(Field); ok {
			*ret = true
		}
		return true
	}
}
/* Checks if argument is Model-Like value, so it could be used as ModelValue */
func isModelV(value interface{}) bool {
	var ret bool
	fieldIter(value, produceIsModelFun(&ret))
	return ret
}

var Conn Connection = nil

func getConnection(conn []Connection) Connection {
	var c Connection
	if len(conn) > 0 {
		c = conn[0]
	} else {
		c = Conn
	}
	return c
}
