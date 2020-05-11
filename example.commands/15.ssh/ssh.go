package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"time"

	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
)

var (
	host = flag.String("Host", "", "Hostname to connect")
	port = flag.Int("Port", 22, "Port SSH daemon is listening")
	user = flag.String("User", "", "Remote user")
)

func init() {
	application.Configure()
}

func main() {
	cmd := "ps"

	// get host public key
	hostKey := getHostKey(*host)

	// ssh client config
	config := &ssh.ClientConfig{
		User: *user,
		Auth: []ssh.AuthMethod{
			//ssh.Password(pass),
			auth(),
		},
		// allow any host key to be used (non-prod)
		// HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		// verify host public key
		HostKeyCallback: ssh.FixedHostKey(hostKey),
		// optional host key algo list
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoRSA,
			ssh.KeyAlgoDSA,
			ssh.KeyAlgoECDSA256,
			ssh.KeyAlgoECDSA384,
			ssh.KeyAlgoECDSA521,
			ssh.KeyAlgoED25519,
		},
		// optional tcp connect timeout
		Timeout: 3 * time.Second,
	}

	// connect
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", *host, *port), config)
	logging.Panic(err)

	defer client.Close()

	// start session
	sess, err := client.NewSession()
	logging.Panic(err)

	defer sess.Close()

	// setup standard out and error
	// uses writer interface
	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	// run single command
	err = sess.Run(cmd)
	logging.Panic(err)
}

func auth() ssh.AuthMethod {
	key, err := ioutil.ReadFile(fmt.Sprintf("/home/%v/.ssh/test_ecdsa", *user))
	logging.Panic(err)

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	logging.Panic(err)

	return ssh.PublicKeys(signer)
}

func getHostKey(host string) ssh.PublicKey {
	// parse OpenSSH known_hosts file
	// ssh or use ssh-keyscan to get initial key
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	logging.Panic(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			logging.Panic(err)

			break
		}
	}

	logging.Panic(err)

	return hostKey
}
