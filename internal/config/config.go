package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Setup struct {
	Pods []Pod `json:"pods"`
	// Можно позже добавить другие поля, если появятся:
	// Services  []Service  `json:"services"`
	// ConfigMaps []ConfigMap `json:"configmaps"`
	// Secrets    []Secret    `json:"secrets"`
	// etc.
}

type Pod struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	// Дополнительные поля, если понадобятся в будущем:
	// Replicas   *int              `json:"replicas"`
	// Namespace  string            `json:"namespace"`
	// Env        []EnvVar          `json:"env"`
	// Args       []string          `json:"args"`
	// Command    []string          `json:"command"`
}

type Lab struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
	EstimatedTime int  `json:"estimated_time"`
	Tasks       []Task `json:"tasks"`
}

type Task struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
    Setup       *Setup     `json:"setup,omitempty"`    // может отсутствовать
	Validation  Validation  `json:"validation"`
}

type Validation struct {
	Type        string `json:"type"`          // node-count, version-prefix, pod-count и т.д.
	Expected    string `json:"expected"`      // может отсутствовать
	ExpectedMin *int   `json:"expected_min"`  // может отсутствовать
	Namespace   string `json:"namespace"`     // опционально, для pod-count
	// можно добавить другие поля по мере появления новых типов валидации
}

type LabRegistry struct {
	Labs map[string]Lab `json:"labs"`
}

var (
    GlobalConfig *LabRegistry
    once         sync.Once
)


func Load() *LabRegistry {
    once.Do(func() {
        data, err := os.ReadFile("config1.json")
        if err != nil {
            log.Fatalf("Не удалось прочитать config.json: %v", err)
        }
        var cfg LabRegistry
        if err := json.Unmarshal(data, &cfg); err != nil {
            log.Fatalf("Не удалось распарсить config.json: %v", err)
        }

        GlobalConfig = &cfg
    })

    return GlobalConfig
}