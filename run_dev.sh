ssh-keygen -f ~/.ssh/id_rsa -N "" -C "admin@dokku"
cp -r ./dokku_deploy/* /dokku
mkdir -p /dokku/home/dokku/.ssh
cat ~/.ssh/id_rsa.pub >> /dokku/home/dokku/.ssh/authorized_keys
# dokku:dokku
chown -R 200:200 /dokku
reflex --start-service -r '\.go$' -- go run /code/main.go
