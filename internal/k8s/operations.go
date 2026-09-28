package k8s

import (
	"context"
	"fmt"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (mng *LabManager) DeployResourse(s *LabSession) error {
	
	//session := mng.Sessions[s.userID]


	config_setup := s.Task.Setup

	switch {
	case len(s.Task.Setup.Pods) > 0:
		pod := config_setup.Pods[0]
		err := mng.Client.CreatePod(s.Namespace, pod.Name, pod.Image)
		if err != nil {
			return fmt.Errorf("Failed to create pod %s: %w", pod.Name, err)
		} else {
			log.Printf("\n. Waiting Pod Start...\n")
			err = mng.Client.WaitForPodRunning(s.Namespace, pod.Name, 60)
			if err != nil {
				return fmt.Errorf("⚠️Error waiting Pod: %v", err)
			}
		}
	}

	return nil
}

func (mng *LabManager) Checker(user_id string) error {
	var check_pods int
	mng.mu.RLock()
	session := mng.Sessions[user_id]
	mng.mu.RUnlock()

	config_setup := mng.Config.Labs["lasbs"].Tasks[0].Setup

	if len(config_setup.Pods) > 0 {
		number_pods, err := mng.Client.GetPodsCount(session.Namespace)
		if err != nil {
			log.Fatal("Error when getting pods")
		}
		log.Println("How much pod is running ?")
		for check_pods != number_pods {
			fmt.Scan(&check_pods)
			fmt.Println("Wrong, try again ")
		}
		fmt.Printf("Exactly %v pods is running \n", check_pods)
	}
	return nil
}

func (c *K8SClient) CreatePod(namespace, podName, image string) error {
	// Создаем объект Pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels: map[string]string{
				"app":     podName,
				"created": "go-client",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  podName,
					Image: image,
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 80,
							Name:          "http",
						},
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyAlways,
		},
	}

	ctx := context.TODO()
	_, err := c.clientset.CoreV1().Pods(namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("не удалось создать pod '%s': %w", podName, err)
	}

	fmt.Printf("✅ Pod '%s' создан в namespace '%s' (image: %s)\n",
		podName, namespace, image)
	return nil
}

func (c *K8SClient) WaitForPodRunning(namespace, podName string, timeoutSeconds int) error {
	ctx := context.TODO()

	for i := 0; i < timeoutSeconds; i++ {
		pod, err := c.clientset.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("ошибка получения pod: %w", err)
		}

		switch pod.Status.Phase {
		case corev1.PodRunning:
			log.Printf("✅ Pod '%s' running in namespace '%v' \n", podName, namespace)
			return nil
		case corev1.PodSucceeded:
			log.Printf("ℹ️ Pod '%s' completed successfully\n", podName)
			return nil
		case corev1.PodFailed:
			return fmt.Errorf("Pod '%s' comleted with error", podName)
		case corev1.PodPending:
			log.Printf("⏳ Pod '%s' in Pending state...\n", podName)
		default:
			log.Printf("⏳ State Pod '%s': %s\n", podName, pod.Status.Phase)
		}

		// Ждем 1 секунду перед следующей проверкой
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("Timeout Pod '%s'", podName)
}

func (c *K8SClient) GetPodsCount(namespace string) (int, error) {
	ctx := context.TODO()

	pods, err := c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return 0, fmt.Errorf("Couldn't get a list of pods: %w", err)
	}

	return len(pods.Items), nil
}
