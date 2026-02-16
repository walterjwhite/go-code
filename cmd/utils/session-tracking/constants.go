package main

const (
	DefaultDBPath     = "~/.data/go/session-tracking.db"
	DefaultServiceURL = "https://httpbin.org/ip"
	
	HTTPTimeout = 30
	
	DBDriverName = "sqlite"
	
	SOCKSProtocol = "tcp"
	
	ResponseBufferSize = 1024
	
	HTTPStatusOK = 200
	
	OriginField = "origin"
	
	SessionTableSchema = `
	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		proxy_addr TEXT NOT NULL,
		public_ip TEXT NOT NULL,
		start_time DATETIME NOT NULL,
		end_time DATETIME NOT NULL
	);
	`
	
	InsertSessionQuery = `
	INSERT INTO sessions (proxy_addr, public_ip, start_time, end_time)
	VALUES (?, ?, ?, ?)
	`
)
