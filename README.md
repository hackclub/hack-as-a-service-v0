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
