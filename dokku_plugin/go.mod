module github.com/hackclub/hack-as-a-service/dokku_json

go 1.16

replace github.com/hackclub/hack-as-a-service => ../

require (
	github.com/codeskyblue/go-sh v0.0.0-20200712050446-30169cf553fe // indirect
	github.com/dokku/dokku/plugins/common v0.0.0-20210501220036-f4b4752e20dc
	github.com/hackclub/hack-as-a-service v0.0.0
	github.com/ryanuber/columnize v2.1.2+incompatible // indirect
	go.lsp.dev/jsonrpc2 v0.9.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
)
