package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	_ "modernc.org/sqlite"
)

func main() {
	defer application.OnPanic()

	config := &Config{}
	application.Configure(config)

	if len(config.Proxy) == 0 {
		logging.Error(errors.New("SOCKS proxy address is required"))
	}

	dbPath, err := resolveDBPath(config.DB)
	logging.Error(err)

	db, err := sql.Open(DBDriverName, dbPath)
	logging.Error(err)
	defer close(db)

	logging.Error(createSessionTable(db))

	client, err := setupProxyClient(config.Proxy, config.HTTPTimeoutSeconds)
	logging.Error(err)

	startTime := time.Now()

	proxyPublicIP, err := getPublicIP(client, config.Service)
	if err != nil {
		log.Warn().Msgf("Could not determine public IP through proxy: %v", err)
		proxyPublicIP = "unknown"
	}

	devicePublicIP, err := getDevicePublicIP(config.Service, config.HTTPTimeoutSeconds)
	if err != nil {
		log.Warn().Msgf("Could not determine device public IP: %v", err)
		devicePublicIP = "unknown"
	}

	fmt.Printf("Connected via SOCKS proxy %s\n", config.Proxy)
	fmt.Printf("Public IP through proxy: %s\n", proxyPublicIP)
	fmt.Printf("Device public IP: %s\n", devicePublicIP)
	fmt.Printf("Session started at: %s\n", startTime.Format(time.RFC3339))

	application.Wait()

	endTime := time.Now()
	fmt.Printf("\nSession ended at: %s\n", endTime.Format(time.RFC3339))

	session := Session{
		ProxyAddr: config.Proxy,
		PublicIP:  proxyPublicIP, // Store the IP seen through the proxy
		StartTime: startTime,
		EndTime:   endTime,
	}

	err = saveSession(db, session)
	if err != nil {
		log.Warn().Msgf("Failed to save session to database: %v", err)
	} else {
		fmt.Printf("Session saved to database: %s\n", dbPath)
	}
}

func close(db *sql.DB) {
	logging.Warn(db.Close(), "Failed to close database")
}
