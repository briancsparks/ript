
dev:
	go run main/ript.go cheat gocli --projectname=ript --dest=../scratch/one


cleandev:
	rm -rf ../scratch/one/ && mkdir -p ../scratch/one/

devandtestrun: dev
	cd ../scratch/one/
	go mod tidy
	go run main/ript.go

generate: ript/config_generated.go ript/templates/gocli.tar
	go generate ./...

ript.exe: generate
	go build -o ./ript.exe main/ript.go

buildrel: ript.exe


