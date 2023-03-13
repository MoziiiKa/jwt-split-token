package main

import "os"

// getting config from environment
func SetENVs() (connectionString string, port string) {
	connectionString = os.Getenv("CONNECTION_STRING")
	port = os.Getenv("PORT")
	return connectionString, port
}
