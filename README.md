# Hello Fooprojectname


```sh
(cd ../../__go-project-template/one/ && git ls-files) | xargs dirname | sort | uniq | egrep -v '[.]' | xargs mkdir -p
tree
(cd ../../__go-project-template/one/ && git ls-files)
(cd ../../__go-project-template/one/ && git ls-files) | xargs -I SUB echo "../../__go-project-template/one/SUB" "SUB"
(cd ../../__go-project-template/one/ && git ls-files) | xargs -I SUB cp "../../__go-project-template/one/SUB" "SUB"
tree
perl -pi -e "s/RIPTENV_IMFO_USER_IRL_NAME/${IMFO_USER_IRL_NAME}/g" options.go
perl -pi -e "s/fooprojectname/ript/g" options.go
echo "${IMFO_GITHUB_USEREMAIL}" | sed 's/@/\\@/g'
perl -pi -e "s/RIPTENV_IMFO_GITHUB_USEREMAIL/${IMFO_GITHUB_USEREMAIL}/g" options.go
```

