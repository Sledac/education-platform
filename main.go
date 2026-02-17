package main

import (
	"k8s-connect/internal/config"
	"k8s-connect/internal/k8s"
	"log"
)

func main() {

	client_set, err := k8s.InitClient()
	if err != nil {
		log.Fatalf("failed to initialize kubernetes client: %w", err)
	}
	cfg := config.Load()

	mgr := k8s.NewLabManager(client_set,cfg)

    CliLabSession, err := mgr.CreateLabSession("user-test","count-pods")
	if err != nil {
		log.Fatal("Can't create LabSession, '%v'", err)
	}

    log.Printf("Connect to cluser use this command '%s' \n",CliLabSession)
	
	mgr.Checker("user-test")
	mgr.DeleteNamespace("user-test")


	


}
