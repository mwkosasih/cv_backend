package main

import "log"

func main() {
	cfg := LoadConfig()
	srv := NewServer(cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
