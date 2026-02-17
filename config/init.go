package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	envConfigLocation = "CONFIG_LOCATION"
	envEnvironment    = "ENVIRONMENT"
	envDeployableName = "DEPLOYABLE_NAME"
)

// InitGlobalConfig loads the application.env file, reads the static YAML config,
// replaces ${ENV_VAR} placeholders, and unmarshals into the provided GlobalConf.
func InitGlobalConfig(conf GlobalConf) {
	initEnv()

	environment, deployableName := getEnvironmentAndDeployableName()

	configLocation := viper.GetString(envConfigLocation)
	if configLocation == "" {
		configLocation = "./configs/brainbash"
	}

	staticConfigFilePath := fmt.Sprintf("%s/application-%s.yml", configLocation, environment)

	if err := loadStaticConfig(staticConfigFilePath, conf.GetStaticConfig(), deployableName); err != nil {
		panic(fmt.Sprintf("Failed to load static config: %v", err))
	}
}

// initEnv loads key=value pairs from application.env into OS environment variables.
// Existing env vars are not overwritten so that real env takes precedence.
func initEnv() {
	data, err := os.ReadFile("application.env")
	if err != nil {
		log.Printf("Warning: could not read application.env: %v", err)
		return
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Only set if not already defined, so real env vars take precedence
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}

	// Bind all env vars into viper so they are accessible via viper.GetString()
	viper.AutomaticEnv()
}

func loadStaticConfig(filePath string, config interface{}, deployableName string) error {
	configBytes, err := readConfig(filePath, deployableName)
	if err != nil {
		return err
	}
	if len(configBytes) == 0 {
		return nil
	}

	// Replace ${ENV_VAR} placeholders with actual environment values
	expanded := os.ExpandEnv(string(configBytes))

	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewBuffer([]byte(expanded))); err != nil {
		return fmt.Errorf("failed to read config into viper: %w", err)
	}

	// Flatten nested keys (e.g. app.port -> APP_PORT) into viper
	setFlattenedKeys(v, v.AllSettings(), "")

	if err := v.Unmarshal(config); err != nil {
		return fmt.Errorf("failed to unmarshal config into struct: %w", err)
	}

	return nil
}

// readConfig reads a multi-document YAML file and merges the deployable-specific
// section into the default (first) section.
//
// Example YAML:
//
//	app:
//	  name: default-app
//	  port: 8080
//	---
//	deployable-name: brainbash
//	app:
//	  name: brainbash-app
//
// For deployableName="brainbash", the result merges the second section over the first.
func readConfig(filePath, deployableName string) ([]byte, error) {
	sections, err := getSections(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get sections from file: %w", err)
	}
	if len(sections) == 0 {
		return nil, nil
	}

	defaultSection := sections[0]
	deployableSpecificSection := make(map[string]interface{})

	for _, section := range sections {
		if section["deployable-name"] == deployableName {
			deployableSpecificSection = section
			break
		}
	}

	mergedSection := mergeMaps(deployableSpecificSection, defaultSection)

	mergedConfigData, err := yaml.Marshal(mergedSection)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal merged config: %w", err)
	}

	return mergedConfigData, nil
}

func getSections(filePath string) ([]map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var sections []map[string]interface{}
	decoder := yaml.NewDecoder(file)

	for {
		var section map[string]interface{}
		if err := decoder.Decode(&section); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("failed to decode YAML section: %w", err)
		}
		sections = append(sections, section)
	}
	return sections, nil
}

// mergeMaps merges source into target recursively. Source values overwrite target values.
// Nested maps are merged recursively rather than replaced wholesale.
func mergeMaps(source, target map[string]interface{}) map[string]interface{} {
	for key, sourceValue := range source {
		if targetValue, exists := target[key]; exists {
			if targetMap, isTargetMap := targetValue.(map[string]interface{}); isTargetMap {
				if sourceMap, isSourceMap := sourceValue.(map[string]interface{}); isSourceMap {
					target[key] = mergeMaps(sourceMap, targetMap)
					continue
				}
			}
		}
		target[key] = sourceValue
	}
	return target
}

func getEnvironmentAndDeployableName() (string, string) {
	environment := os.Getenv(envEnvironment)
	deployableName := os.Getenv(envDeployableName)

	if environment == "" || deployableName == "" {
		panic("ENVIRONMENT and DEPLOYABLE_NAME must be set")
	}

	return environment, deployableName
}

// setFlattenedKeys flattens nested viper keys into uppercase underscore format.
// e.g. app.port -> APP_PORT
func setFlattenedKeys(v *viper.Viper, configMap map[string]interface{}, prefix string) {
	for key, value := range configMap {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		if nestedMap, ok := value.(map[string]interface{}); ok {
			setFlattenedKeys(v, nestedMap, fullKey)
		} else {
			underscoreKey := strings.ToUpper(strings.ReplaceAll(fullKey, ".", "_"))
			v.Set(underscoreKey, value)
		}
	}
}
