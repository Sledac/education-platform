package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type LabConfig struct {
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

var (
    GlobalConfig *LabConfig
    once         sync.Once
)


func Load() *LabConfig {
    once.Do(func() {
        data, err := os.ReadFile("config.json")
        if err != nil {
            log.Fatalf("Не удалось прочитать config.json: %v", err)
        }
        var cfg LabConfig
        if err := json.Unmarshal(data, &cfg); err != nil {
            log.Fatalf("Не удалось распарсить config.json: %v", err)
        }

        GlobalConfig = &cfg
    })

    return GlobalConfig
}