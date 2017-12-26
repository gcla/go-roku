
all_source = $(shell git ls-files -x '.go' | grep -v '^vendor' | grep -v '_test.go')

.PHONY: install clean build test tests doc trace exe debug backup $(all_examples)

build install: ; go install

run: install
	go-roku

clean: ; -rm /tmp/abcdef* /tmp/xyzxyz* 

SRC=go-roku.go
EXE=./go-roku

exe: ; go build -gcflags='-N -l' ${SRC}

debug: exe
	dlv exec ${EXE}

debugtests debugtest: ; dlv test 

backup: ; git push google --all

dot: ; dep status -dot | dot -T png | display

doc docs: ; godoc -http=:6060 
