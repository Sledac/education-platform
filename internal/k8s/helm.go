package k8s

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func (mng *LabManager) DeployVclusterWithHelm(namespace string, releaseName string) error {

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(mng.HostClient.Flags, namespace, "secret", log.Printf); err != nil {
		// Ошибка инициализации
		return fmt.Errorf("Сan't inizialization vcluster, %w", err)
	}

	// 3. Создаем клиента для установки (Install action)
	client := action.NewInstall(actionConfig)
	client.ReleaseName = releaseName
	client.Namespace = namespace
	client.Wait = true //Helm ждет поднятия всех подов!
	client.Timeout = 3 * time.Minute

	// 4. Загружаем наш скачанный локальный чарт
	chartPath := "./charts/vcluster-0.33.1.tgz"
	loadedChart, err := loader.Load(chartPath)
	if err != nil {
		return fmt.Errorf("Error downalod charsts %w", err)
	}

	// 5. Запускаем установку! (Здесь можно передать кастомные values вместо nil)
	release, err := client.Run(loadedChart, nil)
	if err != nil {
		return fmt.Errorf("Error install Vcluster", err)
	}

	log.Printf("✅ Vcluster задеплоен! Статус: %s", release.Info.Status)
	return nil
}

func (mng *LabManager) GetTenantClient(namespace string, releaseName string) (*K8SClient, error) {

	secretName := fmt.Sprintf("vc-%s", releaseName)
	secret, err := mng.HostClient.Clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("Error get Tenant-Secret %w", err)
	}

	kubeconfigBytes := secret.Data["config"]
	restConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeconfigBytes)
	if err != nil {
		return nil, fmt.Errorf("Error parse kubeconfig: %w", err)
	}

	// --- МАГИЯ СЕТЕВОГО МОСТА ДЛЯ ЛОКАЛКИ ---
	localPort := 8443
	// 1. Запускаем фоновый проброс порта
	err = startPortForward(namespace, releaseName, localPort)
	if err != nil {
		return nil, fmt.Errorf("Can't start port-forward: %w", err)
	}
	// 2. Явно переписываем хост в конфиге на наш туннель!
	restConfig.Host = fmt.Sprintf("https://localhost:%d", localPort)
	// ----------------------------------------

	// 🔥 ОТКЛЮЧАЕМ ПРОВЕРКУ СЕРТИФИКАТОВ ДЛЯ ЛОКАЛЬНОГО ТУННЕЛЯ
	restConfig.TLSClientConfig.Insecure = true
	// Обязательно нужно занулить старые данные CA, иначе Go все равно попытается их проверить
	restConfig.TLSClientConfig.CAData = nil
	restConfig.TLSClientConfig.CAFile = ""
	// ----------------------------------------

	tenantClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("Error create tenant client: %w", err)
	}

	wrapper := &K8SClient{
		Clientset: tenantClient,
		Flags:     nil,
	}

	return wrapper, nil
}

// для локальной отладки
func startPortForward(namespace string, releaseName string, localPort int) error {
	svcName := fmt.Sprintf("svc/%s", releaseName) // Имя совпадает с releaseName (user_id)

	// 🔥 ВОТ ЗДЕСЬ ИСПРАВЛЕНИЕ: пробрасываем локальный порт на 443 порт сервиса
	portMapping := "8443:443"

	cmd := exec.Command("kubectl", "port-forward", svcName, portMapping, "-n", namespace)

	// Оставляем логирование для контроля
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("ошибка запуска kubectl: %w", err)
	}

	log.Printf("🔌 Открыт туннель: localhost:%d -> %s:443", localPort, svcName)
	time.Sleep(3 * time.Second)

	return nil
}
