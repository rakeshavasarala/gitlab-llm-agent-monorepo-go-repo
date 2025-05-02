package shared

import (
	"os"

	"gopkg.in/yaml.v3"
)

type RuleConfig struct {
	ID          string                 `yaml:"id"`
	Enabled     bool                   `yaml:"enabled"`
	Description string                 `yaml:"description"`
	Severity    string                 `yaml:"severity"`
	Params      map[string]interface{} `yaml:"params,omitempty"`
}

type Ruleset struct {
	Rules []RuleConfig `yaml:"rules"`
}

type RuleResult struct {
	RuleID     string `json:"rule_id"`
	Applicable bool   `json:"applicable"`
	Passed     bool   `json:"passed,omitempty"`
	Message    string `json:"message,omitempty"`
}

type RuleResults struct {
	Results []RuleResult
}

func LoadRulesConfig(path string) (*Ruleset, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var ruleset Ruleset
	err = yaml.Unmarshal(data, &ruleset)
	if err != nil {
		return nil, err
	}

	return &ruleset, nil
}