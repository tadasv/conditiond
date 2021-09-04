package main

import (
	"encoding/json"
	"github.com/tadasv/conditiond"
	"os"
)

type EvaluatorConfig struct {
	// FunctionWhitelist is a list of allowed functions from
	// ExpressionRegistry. If this slice is nil or empty, then all functions
	// from the registry will be allowed.
	FunctionWhitelist []string `json:"func_whitelist"`

	// FunctionMap is a mapping of expression from ExpressionRegistry to a
	// function name. Function mapping can be used to remap default function
	// names to other names.
	FunctionMap map[string]string `json:"func_map"`
}

type Config struct {
	EvaluatorConfig EvaluatorConfig `json:"evaluator"`
}

func (c Config) String() string {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (c Config) Validate() error {
	// TODO
	return nil
}

func getDefaultConfig() *Config {
	config := &Config{
		EvaluatorConfig: EvaluatorConfig{
			FunctionWhitelist: []string{},
			FunctionMap:       map[string]string{},
		},
	}

	for key := range condition.ExpressionRegistry {
		config.EvaluatorConfig.FunctionMap[key] = key
	}

	return config
}

func GetConfig(configPath string) (*Config, error) {
	config := getDefaultConfig()

	configData, err := os.ReadFile(configPath)
	if err != nil {
		return config, nil
	}

	if err := json.Unmarshal(configData, config); err != nil {
		return nil, err
	}

	return config, nil
}
