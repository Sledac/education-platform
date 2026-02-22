package k8s

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func (c *K8SClient) CreateNamespace(user_id string) (string, error) {

	name := fmt.Sprintf("%s-%d", user_id, time.Now().Unix())

	exists, err := namespaceExists(c.clientset, name)
	if err != nil {
		return name, fmt.Errorf("Verification error namespace: %w", err)
	}

	if exists {
		return name, fmt.Errorf("namespace '%s' already exists", name)
	}

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: name},
	}
	_, err = c.clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		return name, fmt.Errorf("Couldn't create namespace '%s': %w", name, err)
	}

	return name, nil
}

func namespaceExists(clientset *kubernetes.Clientset, name string) (bool, error) {
	_, err := clientset.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		// Если namespace не найден
		return false, nil
	}
	// Namespace существует
	return true, nil
}

func (mng *LabManager) DeleteNamespace(s *LabSession) error {

	err := mng.Client.clientset.CoreV1().Namespaces().Delete(context.TODO(), s.Namespace, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("Couldn't delete namespace '%s': %w", s.Namespace, err)
	}

    fmt.Printf("✅ Namespace deleted: %s\n", s.Namespace)

	mng.mu.Lock()
	s.Namespace = "default"
	mng.mu.Unlock()

	
	return nil

}
