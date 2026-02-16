package k8s

import (
	"fmt"
	"context"
	"time"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
			log.Printf("✅ Pod '%s' запущен\n", podName)
			return nil
		case corev1.PodSucceeded:
			log.Printf("ℹ️  Pod '%s' завершился успешно\n", podName)
			return nil
		case corev1.PodFailed:
			return fmt.Errorf("pod '%s' завершился с ошибкой", podName)
		case corev1.PodPending:
			log.Printf("⏳ Pod '%s' в состоянии Pending...\n", podName)
		default:
			log.Printf("⏳ Состояние pod '%s': %s\n", podName, pod.Status.Phase)
		}

		// Ждем 1 секунду перед следующей проверкой
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("таймаут ожидания pod '%s'", podName)
}

func (c *K8SClient) GetPodsCount(namespace string) (int, error) {
	ctx := context.TODO()
	
	pods, err := c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return 0, fmt.Errorf("не удалось получить список подов: %w", err)
	}
	
	return len(pods.Items), nil
}