package k8s

import (
	"fmt"
	"k8s-connect/internal/config"
	"sync"
)

type LabSession struct {
	labID            string
	userID           string //id_user
	Namespace        string
	TenantClient     *K8SClient
	Task             *config.Task
	CurrentTaskIndex int
	CompletedTasks   []string
}

type LabManager struct {
	Sessions   map[string]*LabSession
	HostClient *K8SClient
	Config     *config.LabRegistry
	mu         sync.RWMutex
}

func NewLabManager(client *K8SClient, cfg *config.LabRegistry) *LabManager {

	return &LabManager{
		HostClient: client,
		Config:     cfg,
		Sessions:   make(map[string]*LabSession),
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

		err := mng.DeployResourse(s)
		if err != nil {
			return fmt.Errorf("Error deploy Resourese: %w", err)
		}

	}

	return nil
}

func (mng *LabManager) CreateLabSession(user_id string, number_labs string) (*LabSession, error) {

	var newSession *LabSession

	ns, err := mng.HostClient.CreateNamespace(user_id)
	if err != nil {
		return &LabSession{}, fmt.Errorf("Can't apply Setup this task %w", err)
	}

	err = mng.DeployVclusterWithHelm(ns, user_id)
	if err != nil {
		return nil, fmt.Errorf("Error deploy vcluster %w", err)
	}

	tenantClient, err := mng.GetTenantClient(ns, user_id)
	if err != nil {
		return nil, fmt.Errorf("Error get kubeconfig from %s %v %w", ns, user_id, err)
	}

	newSession = &LabSession{
		labID:            number_labs,
		userID:           user_id,
		Namespace:        ns,
		TenantClient:     tenantClient,
		CurrentTaskIndex: 0,
	}

	mng.mu.Lock()
	mng.Sessions[user_id] = newSession
	mng.mu.Unlock()

	return newSession, nil
}
