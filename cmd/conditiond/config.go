package main

import (
	"encoding/json"
	"fmt"
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
	for funcName, registeredFunc := range c.EvaluatorConfig.FunctionMap {
		if _, ok := condition.ExpressionRegistry[registeredFunc]; !ok {
			return fmt.Errorf("function map points to unavailable function: %q -> %q", funcName, registeredFunc)
		}
	}

	for _, allowedFunc := range c.EvaluatorConfig.FunctionWhitelist {
		if _, ok := condition.ExpressionRegistry[allowedFunc]; !ok {
			return fmt.Errorf("whitelisted function %q is not part of the registry. Remove or fix the whitelist value", allowedFunc)
		}
	}
	return nil
}

func getDefaultConfig() *Config {
	config := &Config{
		EvaluatorConfig: EvaluatorConfig{
			FunctionWhitelist: nil,
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

	defaultFuncMap := config.EvaluatorConfig.FunctionMap
	config.EvaluatorConfig.FunctionMap = map[string]string{}
	if err := json.Unmarshal(configData, config); err != nil {
		return nil, err
	}

	if len(config.EvaluatorConfig.FunctionMap) == 0 {
		// If function map is empty this means that no mapping was provided in the config.
		// Let's reset it back to the original one.
		config.EvaluatorConfig.FunctionMap = defaultFuncMap
	}

	return config, nil
}
