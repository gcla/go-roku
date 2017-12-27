
all_source = $(shell git ls-files -x '.go' | grep -v '^vendor' | grep -v '_test.go')

.PHONY: install clean build test tests doc docs exe debug backup run xbox tv pause on off

build install: ; go install

run: install
	go-roku 

xbox: install
	go-roku --channel=Xbox

tv: install
	go-roku --channel=Cable

pause: install
	go-roku --pause

on: install
	go-roku --on

off: install
	go-roku --off

clean: ; -rm /tmp/abcdef* /tmp/xyzxyz* 

SRC=go-roku.go
EXE=./go-roku

exe: ; go build -gcflags='-N -l' ${SRC}

debug: exe
	dlv exec ${EXE}

debugtests debugtest: ; dlv test 

backup: ; git push google --all

doc docs: ; godoc -http=:6060 
