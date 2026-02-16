package k8s

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)


type Config struct {
	Task        string     `json:"task"`
	Description string     `json:"description"`
	Setup       Setup      `json:"setup"`
	Validation  Validation `json:"validation"`
}

type Setup struct {
	Pods []Pod `json:"pods"`
}

type Pod struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Validation struct {
	Value int `json:"value"`
}

var session_connect string
var check_pods int


func ReadConfig() (*Config,error) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	var config *Config
	err = json.Unmarshal(data, &config)

	if err != nil {
		panic(err)
	}
	return config, nil
}

func (lab_mng *LabManager) CreateLabSession(user_id string) (string, error) {
	conf,err := ReadConfig()

	if err != nil {
		return "", fmt.Errorf("cant read configuration: %w", err)
	}
	config_setup := conf.Setup

	switch {
	case len(config_setup.Pods) > 0:
		pod := config_setup.Pods[0]
		err = lab_mng.Config.CreatePod(lab_mng.active_session[user_id].Namespace, pod.Name, pod.Image)
		if err != nil {
			return "", fmt.Errorf("failed to create pod %s: %w", pod.Name, err)
		} else { 
		  fmt.Printf("\n. Ожидаем запуск пода...\n")
		  err = lab_mng.Config.WaitForPodRunning(lab_mng.active_session[user_id].Namespace, pod.Name, 60)
		  if err != nil {
			log.Fatal("⚠️Ошибка ожидания пода: %v", err)
		  }
	    }
	}
	
	session_connect = fmt.Sprintf("kubectl config set-context --current --namespace=%s",
								  lab_mng.active_session[user_id].Namespace)

	return session_connect, nil

}

func(lab_mng *LabManager) Checker (user_id string) error {
	conf,err := ReadConfig()

	if err != nil {
		return fmt.Errorf("cant read configuration: %w", err)
	}
	log.Print(conf.Description)
	
	if len(conf.Setup.Pods) > 0 {
		number_pods, err := lab_mng.Config.GetPodsCount(lab_mng.active_session[user_id].Namespace)
		if err != nil {
			log.Fatal("Error when getting pods")
		}
		log.Println("How much pod is running ?")
		for check_pods != number_pods  {
		  	fmt.Scan(&check_pods)
		  	fmt.Println("Wrong, try again ")	
	    }
		fmt.Printf("Exactly %v pods is running \n", check_pods)
	}	
	return nil
} 

