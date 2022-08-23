
dev:
	go run main/ript.go cheat gocli --projectname=ript --dest=../scratch/one


cleandev:
	rm -rf ../scratch/one/ && mkdir -p ../scratch/one/

devntrun: dev
	cd ../scratch/one/
	go mod tidy
	go run main/ript.go


