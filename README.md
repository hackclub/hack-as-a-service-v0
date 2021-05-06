# Hack as a Service (HaaS)

Run your own backend applications for just 5hn/app/month!

<font size="5"> \>> [Check out our progress reports here!](progress/README.md) << </font>

## Development

Contributions to the project are welcome! The backend is mostly Go, and on the frontend we use Next.js with Chakra UI + Hack Club Theme for styling. Make sure to discuss changes with other contributors in a GitHub issue or on Slack before beginning development on new features or bugfixes. Please assign yourself to issues you're working on, to help everyone know what still needs to be done, and keep everyone on the same page.

### Local Development

Before starting local development make sure you have docker and docker-compose installed as well as an empty `.env` file created at the top of the project. You can then run `docker-compose up` to run a local instance of Dokku and the bot. Any edits you save should make the bot restart automatically. The file watcher appears to be broken on Windows for currently unknown reasons - we recommend the use of Windows Subsystem for Linux (WSL) for development on Windows systems.

To open a shell in a container, run `docker exec -it hack-as-a-service_bot_1 bash` or `docker exec -it hack-as-a-service_dokku_1 bash`

### File Guide

- `main.go` - main entrypoint for the built binary, contains web server
- `assets/` - contains static assets including images which are served by the backend
- `dokkud/` - daemon that connects to Dokku. this should remain relatively stagnant
- `frontend/` - the Next.js frontend for the app.
- `pkg/` - contains most of the backend code, including API, routing for frontend, and DB interactions
  - `api` - this package contains API routes and most of the application's business logic
  - `biller` - this package manages real-time billing based on resource usage
  - `db` - this package manages all interactions with the database
  - `dokku` - this package connects to Dokku via our Dokku daemon
- `dokku_data/` - data folder used for running Dokku in development (ignore)
- `dokku_deploy/` - holds files which are automatically deployed to the server, including Nginx config and custom Dokku error pages
- `dokku_plugin/` - Dokku plugin that provides a variety of HaaS-specific commands.
- `docker/` - Docker images for HaaS services

### Deployment

Dokku recommends Debian for its installation, so we will assume a Debian-based system is being used. To start a fresh deployment, make sure to set up `dokkud` by building it (`cd` into `dokkud` and run `go build -o dokkud .`), moving the binary to `/opt/dokkud/`, placing the service file at `/etc/systemd/system` and enabling it (`systemctl enable --now dokkud`). Then follow the installation directions for Dokku, as the HaaS bot is managed by it, and create a new app which deploys from the Dockerfile at the root of this repo.

On our installation, the API, frontend, and server configuration should auto-deploy once CI passes on master.
