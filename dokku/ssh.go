package dokku

import (
	"golang.org/x/crypto/ssh"
	"os"
	"io/ioutil"
)

func PublicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), err
}

func RunCommand(cmd string) (string, error) {
	key, err := PublicKeyFile(os.Getenv("HOME")+"/.ssh/id_rsa")
	if err != nil {
		return "", err
	}
	
	sshConfig := &ssh.ClientConfig {
		User: os.Getenv("DOKKU_USER"),
		// we should really change this in the future
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			key,
		},
	}

	host := os.Getenv("DOKKU_HOST") + ":" + os.Getenv("DOKKU_SSH_PORT")

	connection, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return "", err
	}

	session, err := connection.NewSession()
	if err != nil {
		return "", err
	}

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return string(output), err
	}
	
	return string(output), nil

		
}
