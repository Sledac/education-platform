package k8s

import (
	"fmt"
	"k8s-connect/internal/config"
	"log"
	"sync"
)

type LabSession struct {
	userID string //id_user
	Namespace string
	Task string
}

type LabManager struct {
	Sessions map[string]*LabSession
	Client *K8SClient
	Config *config.LabConfig
	mu sync.RWMutex
}

func NewLabManager(client *K8SClient, cfg *config.LabConfig) *LabManager {
    return &LabManager{
        Client:   client,
		Config: cfg,
        Sessions: make(map[string]*LabSession),
    }
}

func (mng *LabManager) CreateLabSession(user_id string , task_name string) (string,error) {
	mng.mu.Lock()
    defer mng.mu.Unlock()

	ns, err := mng.Client.CreateNamespace(user_id)
	if err != nil {
		return "",fmt.Errorf("❌ Error: %v\n", err)
	}
	log.Println("✅ Namespace successfull create ")
    
	newSession := &LabSession{
		userID: user_id,
		Namespace: ns,
		Task: task_name,
	}
	mng.Sessions[user_id]=newSession

	err = mng.DeployResourse(user_id)

	if err != nil {
		return "",fmt.Errorf("Error deloy Resourese: %w", err)
	}

	return fmt.Sprintf("kubectl config set-context --current --namespace=%s", ns),nil

}