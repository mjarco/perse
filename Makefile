include $(GOROOT)/src/Make.inc

TARG=perse
GOFILES=\
	collection.go\
	fieldtypes.go\
	interfaces.go\
	utils.go\
	crud.go\


include $(GOROOT)/src/Make.pkg
