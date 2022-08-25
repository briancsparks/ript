
dev:
	go run ript.go cheat gocli --projectname=ript --dest=../scratch/one


cleandev:
	rm -rf ../scratch/one/ && mkdir -p ../scratch/one/

devandtestrun: dev
	mkdir -p ../scratch/one/; cd ../scratch/one/; go mod tidy; go run ript.go

generate: ript/config_generated.go
	go generate ./...

ript.exe: generate
	go build -o ./ript.exe ript.go

buildrel: ript.exe

installlocal: generate
	go build -o $(HOME)/bin/ript ript.go
