
dev:
	go run main/ript.go cheat gocli --projectname=ript --dest=../scratch/one


cleanscratch:
  test -d ../scratch/one/ && rm -rf ../scratch/one

