package k8s

import (
	"log"
)

type LabSession struct {
	userID string //id_user
	Namespace string
	Task string
}

type LabManager struct {
	active_session map[string]*LabSession
	Config *K8SClient
}

var session = make(map[string]*LabSession)

func CreateLabManager(user_id string , task_name string) (*LabManager,error) {

	client_set, err := InitClient()
	if err != nil {
		log.Fatal("Something is happed",err)
	}
	
	ns, err := client_set.CreateNamespace(user_id)
	if err != nil {
		log.Fatalf("❌ Error: %v\n", err)
	}
	log.Println("✅ Namespace successfull create ")
    
	newSession := &LabSession{
		userID: user_id,
		Namespace: ns,
		Task: task_name,
	}
	session[user_id]=newSession

	return &LabManager{active_session: session,Config: client_set},nil
}

func(l *LabManager) GetCurrentTask (user_id string) (string, error) {

	if _, exists := l.active_session[user_id]; !exists {
		log.Fatal("The task is not available for the user %v", user_id)

	}
	return l.active_session[user_id].Task, nil

}