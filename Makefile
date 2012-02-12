include $(GOROOT)/src/Make.inc

all:
	gomake -C counter
	gomake -C model
	gomake -C app

install:
	gomake install -C counter 
	gomake install -C model install
	gomake install -C app 

clean:	
	gomake -C counter clean
	gomake -C model clean
	gomake -C app clean

gofmt:
	gofmt -w *.go
	gofmt -w */*.go
