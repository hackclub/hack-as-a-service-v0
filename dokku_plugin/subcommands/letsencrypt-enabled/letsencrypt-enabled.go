package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/dokku/dokku/plugins/common"
	"github.com/dokku/dokku/plugins/config"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: dokku haas:letsencrypt-enabled <app name>")
		os.Exit(1)
	}

	app_name := os.Args[2]

	if err := common.VerifyAppName(app_name); err != nil {
		fmt.Println("App not found")
		os.Exit(1)
	}

	active := LetsEncryptIsActive(app_name)
	if active {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
}

// https://github.com/dokku/dokku-letsencrypt/blob/1f0948aafaa187caa6d0e09346087b9b85dab7b3/functions#L94
func LetsEncryptIsActive(app string) bool {
	if !IsSslEnabled(app) {
		return false
	}
	domains := GetAppDomains(app)

	certSha1 := "not_found"
	serverLetsEncryptCrt := path.Join(os.Getenv("DOKKU_ROOT"), app, "tls", "server.letsencrypt.crt")
	serverCrt := path.Join(os.Getenv("DOKKU_ROOT"), app, "tls", "server.crt")
	if fileExists(serverLetsEncryptCrt) {
		file, err := os.Open(serverLetsEncryptCrt)
		if err == nil {
			contents, err := ioutil.ReadAll(file)
			if err == nil {
				hasher := sha1.New()
				hasher.Write(contents)
				certSha1 = hex.EncodeToString(hasher.Sum(nil))
			}
		}
	} else {
		file, err := os.Open(serverCrt)
		if err == nil {
			contents, err := ioutil.ReadAll(file)
			if err == nil {
				hasher := sha1.New()
				hasher.Write(contents)
				certSha1 = hex.EncodeToString(hasher.Sum(nil))
			}
		}
	}

	leSha1 := "not_found"
	domainFound := false
	if len(domains) > 0 {
		domain := domains[0]
		domainLetsEncryptCrt := path.Join(os.Getenv("DOKKU_ROOT"), app, "letsencrypt", "certs", "current", "certificates", domain+".pem")
		if fileExists(domainLetsEncryptCrt) {
			domainFound = true
			file, err := os.Open(domainLetsEncryptCrt)
			if err == nil {
				contents, err := ioutil.ReadAll(file)
				if err == nil {
					hasher := sha1.New()
					hasher.Write(contents)
					certSha1 = hex.EncodeToString(hasher.Sum(nil))
				}
			}
		}
	}
	if !domainFound {
		fullchainLetsEncryptCrt := path.Join(os.Getenv("DOKKU_ROOT"), app, "letsencrypt", "certs", "current", "fullchain.pem")
		if fileExists(fullchainLetsEncryptCrt) {
			domainFound = true
			file, err := os.Open(fullchainLetsEncryptCrt)
			if err == nil {
				contents, err := ioutil.ReadAll(file)
				if err == nil {
					hasher := sha1.New()
					hasher.Write(contents)
					certSha1 = hex.EncodeToString(hasher.Sum(nil))
				}
			}
		}
	}
	return certSha1 == leSha1
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

// https://github.com/dokku/dokku/blob/master/plugins/certs/functions
func IsSslEnabled(app string) bool {
	appSslPath := path.Join(os.Getenv("DOKKU_ROOT"), app, "tls")
	return fileExists(path.Join(appSslPath, "server.crt")) && fileExists(path.Join(appSslPath, "server.key"))
}

// https://github.com/dokku/dokku/blob/f4b4752e20dc87da31f6b656202ccb298d234a0d/plugins/domains/functions#L205
func GetAppDomains(app string) []string {
	appVhostPath := path.Join(os.Getenv("DOKKU_ROOT"), app, "VHOST")
	globalHostnamePath := path.Join(os.Getenv("DOKKU_ROOT"), "HOSTNAME")
	if IsAppVhostEnabled(app) {
		if file, err := os.Open(appVhostPath); err == nil {
			contents, err := ioutil.ReadAll(file)
			if err == nil {
				return strings.Split(string(contents), "\n")
			}
		} else if file, err := os.Open(globalHostnamePath); err == nil {
			contents, err := ioutil.ReadAll(file)
			if err == nil {
				return strings.Split(string(contents), "\n")
			}
		}
	}
	return []string{}
}

// https://github.com/dokku/dokku/blob/f4b4752e20dc87da31f6b656202ccb298d234a0d/plugins/domains/functions#L253
func IsAppVhostEnabled(app string) bool {
	noVhost, ok := config.Get(app, "NO_VHOST")
	if ok && noVhost == "1" {
		return false
	}
	return true
}
