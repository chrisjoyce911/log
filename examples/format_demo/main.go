package main

import (
	log "github.com/chrisjoyce911/log"
)

func main() {
	// Basic formatting with f-variants
	log.SetFlags(log.LstdFlags)
	log.Infof("service %s started on port %d", "catalog", 8080)
	log.Debugf("loaded %d features", 12)

	// Mix formatted and structured styles
	user := "alice"
	log.Warnf("quota near limit for %s", user)
	log.Info("quota", "user", user, "used", 92, "limit", 100)

	// Errors
	id := 42
	log.Errorf("failed to fetch record id=%d: %s", id, "not found")
}
