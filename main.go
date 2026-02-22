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
	mgr := k8s.NewLabManager(client_set, cfg)

	session, err := mgr.CreateLabSession("user-test","lab1-basics")

	for session.CurrentTaskIndex < len(mgr.Config.Labs["lab1-basics"].Tasks) {
		session.Task = &mgr.Config.Labs["lab1-basics"].Tasks[session.CurrentTaskIndex]

		err = mgr.ApplySetupForCurrentTask(session)
		if err != nil {
			log.Fatal("Cant apply Setup for task '%s': %w",session.Task.Title, err)
		}
		log.Printf("Connect to cluster use this command: kubectl config set-context --current --namespace=%s",
    	session.Namespace)
		log.Println(session.Task.Title)
		log.Println(session.Task.Description)

		if k8s.CheckCurrentTask(session) {
			session.CurrentTaskIndex++
	   		session.CompletedTasks = append(session.CompletedTasks, session.Task.ID)
		}
		if session.Namespace != "default"{
			mgr.DeleteNamespace(session)
		}		

	}
	
}