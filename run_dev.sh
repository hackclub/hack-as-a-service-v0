ssh-keygen -f ~/.ssh/id_rsa -N ""
cp ~/.ssh/id_rsa.pub ./dokku_deploy/misc/id_rsa.pub
cp -r ./dokku_deploy/* /dokku
reflex --start-service -r '\.go$' -- go run /code/main.go
