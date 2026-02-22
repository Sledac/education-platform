package k8s

import (
	"fmt"
	"k8s-connect/internal/config"
	"log"
	"sync"
)

type LabSession struct {
	labID string
	userID string //id_user
	Namespace string
	Task *config.Task
	CurrentTaskIndex int
	CompletedTasks []string
}

type LabManager struct {
	Sessions map[string]*LabSession
	Client *K8SClient
	Config *config.LabRegistry
	mu sync.RWMutex
}

func NewLabManager(client *K8SClient, cfg *config.LabRegistry) *LabManager {
    return &LabManager{
        Client:   client,
		Config: cfg,
        Sessions: make(map[string]*LabSession),
    }
}



func CheckCurrentTask(task *LabSession) bool {

	var check_answer string

	for check_answer != task.Task.Validation.Expected {
		fmt.Scan(&check_answer)
	}
	return true
	
}

func (mng *LabManager) ApplySetupForCurrentTask(s *LabSession) error {
	mng.mu.Lock()
    defer mng.mu.Unlock()

	if s.Task.Setup != nil {
		namespace, err := mng.Client.CreateNamespace(s.userID)
		if err != nil {
			return fmt.Errorf("Can't apply Setup this task %w", err)
		}
		s.Namespace = namespace
		log.Println("✅ Namespace successfull create ")
		
		err = mng.DeployResourse(s)

		if err != nil {
			return fmt.Errorf("Error deploy Resourese: %w", err)
		}

	}
	return nil
} 


func (mng *LabManager) CreateLabSession(user_id string, number_labs string) (*LabSession,error) {
	mng.mu.Lock()
    defer mng.mu.Unlock()

	var newSession *LabSession
	
	newSession = &LabSession{
		labID: number_labs,
		userID: user_id,
		Namespace: "default",
		CurrentTaskIndex: 0,
	}
	
	mng.Sessions[user_id]=newSession
	return newSession,nil
}