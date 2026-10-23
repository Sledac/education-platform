package main

import (
	"k8s-connect/internal/config"
	"k8s-connect/internal/k8s"
	"log"
)



func main() {

	user:="user-test"

	lab_id:= "lab1-basics"

	client_set, err := k8s.InitClient()
	if err != nil {
		log.Fatalf("failed to initialize kubernetes client: %w", err)
	}
	
	cfg := config.Load()
	mgr := k8s.NewLabManager(client_set, cfg)

	
	session, err := mgr.CreateLabSession(user,lab_id)
	if err != nil {
		log.Fatal("failed to inizialization session for user '%s' labratory '%v'",user,lab_id)
	}

	for session.CurrentTaskIndex < len(mgr.Config.Labs["lab1-basics"].Tasks) {
		session.Task = &mgr.Config.Labs["lab1-basics"].Tasks[session.CurrentTaskIndex]

		err = mgr.ApplySetupForCurrentTask(session)
		if err != nil {
			log.Fatal("Cant apply Setup for task '%s': %w",session.Task.Title, err)
		}
		log.Println(session.Task.Title)
		log.Println(session.Task.Description)

		if k8s.CheckCurrentTask(session) {
			session.CurrentTaskIndex++
	   		session.CompletedTasks = append(session.CompletedTasks, session.Task.ID)
		}

	}
	mgr.DeleteNamespace(session)
	
}