ssh-keygen -f ~/.ssh/id_rsa -N "" -C "admin@dokku"
DELETE_PLUGIN_LIST=0
if [ -f /dokku/plugin-list ]; then
    DELETE_PLUGIN_LIST=1
fi
cp -r ./dokku_deploy/* /dokku
[ "$DELETE_PLUGIN_LIST" -eq "1" ] && rm /dokku/plugin-list
mkdir -p /dokku/home/dokku/.ssh
cat ~/.ssh/id_rsa.pub >> /dokku/home/dokku/.ssh/authorized_keys
# dokku:dokku
chown -R 200:200 /dokku
rm -f /var/run/dokku-daemon/dokkud.sock
while [ ! -S /var/run/dokku-daemon/dokkud.sock ]; do
    sleep 1
done
reflex --start-service -r '\.go$' -- go run /code/main.go
