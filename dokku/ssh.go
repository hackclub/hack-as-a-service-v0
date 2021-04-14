package dokku

import (
	"golang.org/x/crypto/ssh"
	"os"
	"fmt"
	"io/ioutil"
)

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		fmt.Printf("%s", err)
		return nil
	}
	return ssh.PublicKeys(key)
}

func RunCommand(cmd string) string {

	sshConfig := &ssh.ClientConfig {
		User: os.Getenv("DOKKU_USER"),
		// we should really change this in the future
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			PublicKeyFile(os.Getenv("HOME")+"/.ssh/id_rsa"),
		},
	}

	 host := os.Getenv("DOKKU_HOST") + ":" + os.Getenv("DOKKU_SSH_PORT")

	connection, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		fmt.Printf("Failed to dial: %s", err)
		return ""
	}

	session, err := connection.NewSession()
	if err != nil {
		fmt.Printf("Failed to create session: %s", err)
		return ""
	}

	output, err := session.CombinedOutput(cmd)
	
	return string(output)

		
}
