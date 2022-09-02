
dev:
	go run ript.go cheat gocli --projectname=ript --dest=../scratch/one


cleandev:
	rm -rf ../scratch/one/ && mkdir -p ../scratch/one/

devandtestrun: dev
	mkdir -p ../scratch/one/; cd ../scratch/one/; go mod tidy; make dev

generate:
	go generate ./...

ript/templates/gocli.tar: ript/gentemplatetars.go
	go generate ript/gentemplatetars.go

ript/config_generated.go: ript/genconfig.go
	go generate ript/genconfig.go

ript.exe: ript/templates/gocli.tar ript/config_generated.go
	go build -o ./ript.exe ript.go

ript-release.exe: generate
	go build -tags release -o ./ript-release.exe ript.go

builddbg: ript.exe

buildrel: ript-release.exe

installlocal: buildrel
	cp ript-release.exe $(HOME)/bin/ript

