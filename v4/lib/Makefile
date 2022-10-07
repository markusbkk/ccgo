# Copyright 2022 The CC Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

.PHONY:	build_all_targets clean edit editor tmp x csmith fftest shorttest sqltest

build_all_targets:
	GOOS=darwin GOARCH=amd64 go build -v ./...
	GOOS=darwin GOARCH=amd64 go test -o /dev/null -c
	GOOS=darwin GOARCH=arm64 go build -v ./...
	GOOS=darwin GOARCH=arm64 go test -o /dev/null -c
	GOOS=freebsd GOARCH=386 go build -v ./...
	GOOS=freebsd GOARCH=386 go test -o /dev/null -c
	GOOS=freebsd GOARCH=amd64 go build -v ./...
	GOOS=freebsd GOARCH=arm go test -o /dev/null -c
	GOOS=freebsd GOARCH=arm64 go build -v ./...
	GOOS=linux GOARCH=386 go build -v ./...
	GOOS=linux GOARCH=386 go test -o /dev/null -c
	GOOS=linux GOARCH=amd64 go build -v ./...
	GOOS=linux GOARCH=amd64 go test -o /dev/null -c
	GOOS=linux GOARCH=arm go build -v ./...
	GOOS=linux GOARCH=arm go test -o /dev/null -c
	GOOS=linux GOARCH=arm64 go build -v ./...
	GOOS=linux GOARCH=arm64 go test -o /dev/null -c
	GOOS=linux GOARCH=ppc64le go build -v ./...
	GOOS=linux GOARCH=ppc64le go test -o /dev/null -c
	GOOS=linux GOARCH=riscv64 go build -v ./...
	GOOS=linux GOARCH=riscv64 go test -o /dev/null -c
	GOOS=linux GOARCH=s390x go build -v ./...
	GOOS=linux GOARCH=s390x go test -o /dev/null -c
	GOOS=netbsd GOARCH=amd64 go build -v ./...
	GOOS=netbsd GOARCH=arm go test -o /dev/null -c
	GOOS=openbsd GOARCH=amd64 go build -v ./...
	GOOS=openbsd GOARCH=amd64 go test -o /dev/null -c
	GOOS=openbsd GOARCH=arm64 go build -v ./...
	GOOS=openbsd GOARCH=arm64 go test -o /dev/null -c
	GOOS=windows GOARCH=386 go build -v ./...
	GOOS=windows GOARCH=386 go test -o /dev/null -c
	GOOS=windows GOARCH=amd64 go build -v ./...
	GOOS=windows GOARCH=amd64 go test -o /dev/null -c
	GOOS=windows GOARCH=arm64 go build -v ./...
	GOOS=windows GOARCH=arm64 go test -o /dev/null -c

clean:
	rm -f log-* cpu.test mem.test *.out
	go clean

edit:
	@touch log
	@if [ -f "Session.vim" ]; then gvim -S & else gvim -p Makefile *.go & fi

editor:
	stringer -output stringer.go -type mode,name
	gofmt -l -s -w *.go
	go test -o /dev/null -c
	go install -v -tags=ccgo.assert modernc.org/ccgo/v4 2>&1 | tee log-editor

tmp:
	go test -tags=ccgo.assert -v -run TestExec -trc -trctodo -shelltimeout 5s 2>&1 | tee log
	grep -a '^TODO' log | sort | uniq -c | sort -n | tee log-2
	grep -a 'TOTAL\|PASS' log || true

x:
	rm -f /tmp/ccgo.log
	go test -v -timeout 24h -tags=ccgo.assert,ccgo.dmesg -short -timeout 24h -run Exec -err1 -trc -shelltimeout 10s 2>&1 | tee log
	grep TOTAL log
	git status

csmith:
	go test -v -timeout 25h -tags=ccgo.assert -run CSmith -blackbox $(shell pwd)/log-csmith -csmith 24h -csmithc 10s -trc 2>&1 | tee log

fftest:
	go test -v -timeout 24h -short -failfast -tags=ccgo.assert -trc 2>&1 | tee log
	grep TOTAL log
	git status testdata
	git diff testdata/

shorttest:
	go test -v -timeout 24h -short -tags=ccgo.assert -trc 2>&1 | tee log
	grep TOTAL log
	git status testdata
	git diff testdata/

sqltest:
	go test -v -timeout 24h -run SQL -keep 2>&1 | tee log
