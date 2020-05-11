package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"

	//"bytes"
	"flag"
	"github.com/mxk/go-imap/imap"
	"time"
)

var (
	serverAddressFlag = flag.String("ServerAddress", "imap.gmail.com:993", "Server Address")
	//logoutTimeoutFlag = flag.Int("LogoutTimeout", 30, "Logout Timeout in seconds")
	usernameFlag = flag.String("Username", "", "Username")
	passwordFlag = flag.String("Password", "", "Password")
	folderFlag   = flag.String("Folder", "", "Folder Name")
)

func init() {
	application.Configure()
}

// TODO:
// record approval / denial (to file, to ES)
// dynamically serve requests (set timeout for approval, assume denied if not approved by specified time, allow approval to be denied within a given time frame)
// record (client IP address, browser, other headers)
// DONE
// 1. inject request #
// 2. inject request description
func main() {
	i := connect()
	defer cleanup(i)

	processFolder(i)
}

func connect() *imap.Client {
	i, err := imap.DialTLS(*serverAddressFlag, nil)
	logging.Panic(err)

	_, err = i.Login(*usernameFlag, *passwordFlag)
	logging.Panic(err)

	return i
}

func processFolder(i *imap.Client) {
	_, err := i.Select(*folderFlag, false)
	logging.Panic(err)

	set, err := imap.NewSeqSet("1:*")
	logging.Panic(err)

	log.Info().Msgf("Set: %v", set)

	//cmd, err := i.Fetch(set, "FLAGS", "INTERNALDATE", "RFC822.SIZE", "ENVELOPE", "BODYSTRUCTURE", "SUBJECT")
	cmd, err := i.Fetch(set, "RFC822.HEADER", "RFC822.TEXT")
	logging.Panic(err)

	for cmd.InProgress() {
		logging.Panic(i.Recv(-1))

		for _, response := range cmd.Data {
			processMessage(i, response)
		}
	}
}

func processMessage(i *imap.Client, response *imap.Response) {
	/*
		message := response.MessageInfo()
		log.Info().Msgf("= %d\n", message.Seq)
		for _, k := range imap.AsList(message.Attrs["BODYSTRUCTURE"]) {
			log.Info().Msgf("== %#v\n", k)
		}
	*/
	header := imap.AsBytes(response.MessageInfo().Attrs["RFC822.HEADER"])
	uid := imap.AsNumber((response.MessageInfo().Attrs["UID"]))
	body := imap.AsBytes(response.MessageInfo().Attrs["RFC822.TEXT"])

	/*
		if msg, _ := mail.ReadMessage(bytes.NewReader(header)); msg != nil {
			log.Info().Msgf("Subject: %v", msg.Header.Get("Subject"))
			log.Info().Msgf("UID: %v", uid)

			log.Info().Msgf(string(body))
		}*/

	log.Info().Msgf(string(body))
	log.Info().Msgf(string(header))
	log.Info().Msgf(string(uid))
}

func cleanup(i *imap.Client) {
	//_, err := i.Logout(*logoutTimeoutFlag * time.Second)
	_, err := i.Logout(30 * time.Second)
	logging.Panic(err)
}
