package main

import (
    "impulse/server"
    "log"
)


func main() {
    server, err := server.NewServer()
    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }
    
    err = server.Serve()
    if err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
