#!/bin/bash
# This script deploys all of HaaS's config over SSH

HAAS_IP=45.55.45.5

# Deploy config in `dokku_deploy/``
echo "Deploying config..."
rsync -r dokku_deploy/* root@$HAAS_IP:/

# Restart nginx
echo "Restarting nginx..."
ssh root@$HAAS_IP service nginx restart

# Deploy Dokku plugin
echo "Deploying Dokku plugin..."
ssh root@$HAAS_IP \
    dokku plugin:install https://github.com/hackclub/hack-as-a-service --committish config-deployment --name haas