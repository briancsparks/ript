
# To use git tags, see: https://www.forkingbytes.com/blog/dynamic-versioning-your-go-application/
# I.e.:   go build -ldflags "-X main.version=`git tag --sort=-version:refname | head -n 1`
VERSION="v0.9.9"

LDFLAGS_DBG=-ldflags "-X github.com/briancsparks/ript/ript.IsActiveDevelopment=true"
LDFLAGS_REL=-ldflags "-X github.com/briancsparks/ript/ript.Version=${VERSION} -X github.com/briancsparks/ript/ript.IsActiveDevelopment=false"

outdir:
	mkdir -p ../scratch/one

dev:
	go run ${LDFLAGS_DBG} cmd/main/ript.go cheat gocli --projectname=ript --dest=../scratch/one

showversiond:
	go run ${LDFLAGS_DBG} cmd/main/ript.go version

showversionr:
	go run ${LDFLAGS_REL} cmd/main/ript.go version


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

riptd.exe: ript/templates/gocli.tar ript/config_generated.go
	go build ${LDFLAGS_DBG} -o ./riptd.exe cmd/main/ript.go

riptr.exe: generate
	go build ${LDFLAGS_REL} -tags release -o ./riptr.exe cmd/main/ript.go

builddbg: riptd.exe

buildrel: riptr.exe

installlocal: buildrel
	cp riptr.exe $(HOME)/bin/ript
	rm riptr.exe

