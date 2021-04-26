#!/bin/bash
# This script deploys all of HaaS's config over SSH

HAAS_IP=45.55.45.5

# Deploy config in `dokku_deploy/``
echo "Deploying config..."
rsync -e "ssh -o StrictHostKeyChecking=no" -r dokku_deploy/* root@$HAAS_IP:/

# Restart nginx
echo "Restarting nginx..."
ssh root@$HAAS_IP service nginx restart

# Deploy Dokku plugin
echo "Deploying Dokku plugin..."
ssh root@$HAAS_IP /bin/bash << EOF
    dokku plugin:install git@github.com:hackclub/hack-as-a-service.git --name haas \
    || dokku plugin:update haas
EOF

# Deploy dokkud
echo "Deploying dokkud..."
ssh root@$HAAS_IP /bin/bash << EOF
    REPO_DIR=$(mktemp -d)

    git clone git@github.com:hackclub/hack-as-a-service.git \$REPO_DIR
    cd \$REPO_DIR/dokkud
    go build .

    service dokkud stop

    mkdir -p /opt/dokkud
    cp ./dokkud /opt/dokkud/dokkud

    service dokkud start

    rm -rf \$REPO_DIR
EOF