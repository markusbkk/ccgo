# Copyright 2017 The CCGO Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

.PHONY:	all clean cover cpu editor internalError later mem nuke todo edit

grep=--include=*.go --include=*.l --include=*.y --include=*.yy
ngrep='TODOOK\|parser\.go\|scanner\.go\|.*_string\.go'

all: editor
	go vet 2>&1 | grep -v $(ngrep) || true
	golint 2>&1 | grep -v $(ngrep) || true
	make todo
	misspell *.go
	staticcheck || true
	maligned || true
	grep -n 'FAIL\|PASS' log
	git status
	go version
	date

clean:
	go clean
	rm -f *~ *.test *.out

cover:
	t=$(shell mktemp) ; go test -coverprofile $$t && go tool cover -html $$t && unlink $$t

cpu: clean
	# go test -run @ -bench . -cpuprofile cpu.out
	go test -run TCL -todo -cpuprofile cpu.out
	go tool pprof -lines *.test cpu.out

edit:
	gvim -p Makefile *.go

editor:
	cp log log0 || true
	go version | tee log
	./unconvert.sh
	gofmt -l -s -w *.go
	#TODO GOOS=linux GOARCH=arm go build
	#TODO GOOS=linux GOARCH=386 go build
	GOOS=linux GOARCH=amd64 go build
	#TODO GOOS=windows GOARCH=386 go build
	#TODO GOOS=windows GOARCH=amd64 go build
	go test -i
	date 2>&1 | tee -a log
	go test -timeout 24h -edit 2>&1 | tee -a log

internalError:
	egrep -ho '"internal error.*"' *.go | sort | cat -n

later:
	@grep -n $(grep) LATER * || true
	@grep -n $(grep) MAYBE * || true

mem: clean
	#go test -run @ -bench . -memprofile mem.out -memprofilerate 1 -timeout 24h
	go test -run TCL -todo -memprofile mem.out -memprofilerate 1 -timeout 24h
	go tool pprof -lines -web -alloc_space *.test mem.out

nuke: clean
	go clean -i

todo:
	@grep -nr $(grep) ^[[:space:]]*_[[:space:]]*=[[:space:]][[:alpha:]][[:alnum:]]* * | grep -v $(ngrep) || true
	@grep -nr $(grep) TODO * | grep -v $(ngrep) || true
	@grep -nr $(grep) BUG * | grep -v $(ngrep) || true
	@grep -nr $(grep) [^[:alpha:]]println * | grep -v $(ngrep) || true
