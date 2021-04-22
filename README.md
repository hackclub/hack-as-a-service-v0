# Hack as a Service (HaaS)

Run your own backend applications for just 5hn/app/month!

## Development

Contributions to the project are welcome! Make sure to discuss changes with other contributors in a GitHub issue or on Slack before beginning development on new features or bugfixes.

### Local Development

Make sure you have Docker and Docker Compose installed, then run `docker-compose up` to run a local instance of Dokku and the bot. Any edits you save should make the bot restart automatically.

To open a shell in a container, run `docker exec -it hack-as-a-service_bot_1 bash` or `docker exec -it hack-as-a-service_dokku_1 bash`

### File Guide

- `main.go` - main entrypoint for the built binary, contains web server
- `dokku/` - package for interacting with Dokku
- `dokku_data/`, `dokku_deploy/` - data folders used for running Dokku in development (ignore)
- `dokkud/` - daemon that connects to Dokku. this should remain relatively stagnant
- `frontend/` - the Next.js frontend for the app

### Deployment
Dokku recommends Debian for its installation, so we will assume a Debian-based system is being used. To start a fresh deployment, make sure to set up `dokkud` by building it (`cd` into `dokkud` and run `go build -o dokkud .`), moving the binary to `/opt/dokkud/`, placing the service file at `/etc/systemd/system` and enabling it (`systemctl enable --now dokkud`). Then follow the installation directions for Dokku, as the HaaS bot is managed by it, and create a new app which deploys from the Dockerfile at the root of this repo.

On our installation, the bot and frontend should auto-deploy once CI passes on master. Other files, including those in `dokku_deploy` and `dokkud` must be manually updated by someone with direct access to the server.