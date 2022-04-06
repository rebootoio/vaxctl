package model

import (
	"encoding/json"
	"fmt"
	"vaxctl/api"
	"vaxctl/helpers"
)

type RulesResponse struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Name        string   `json:"name" yaml:"name" header:"Name"`
	Regex       string   `json:"regex,omitempty" yaml:"regex" header:"Regex"`
	Actions     []string `json:"actions,omitempty" yaml:"actions" header:"Actions"`
	IgnoreCase  bool     `json:"ignore_case" yaml:"ignore_case" header:"Ignore Case"`
	Enabled     bool     `json:"enabled" yaml:"enabled" header:"Enabled"`
	Position    int      `json:"position,omitempty" yaml:"position,omitempty" header:"Position"`
	AfterRule   string   `json:"after_rule,omitempty" yaml:"after_rule,omitempty"`
	BeforeRule  string   `json:"before_rule,omitempty" yaml:"before_rule,omitempty"`
	Screenshot  string   `json:"screenshot,omitempty" yaml:"screenshot,omitempty"`
	OcrText     string   `json:"ocr_text,omitempty" yaml:"ocr_text,omitempty"`
	LastUpdated string   `json:"last_updated,omitempty" yaml:"last_updated,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

func PrintRules(name string, verbose bool, output string) error {
	allRules, err := GetRules(name)
	if err != nil {
		return err
	}
	var reportObject interface{}
	if name != "" {
		if verbose {
			reportObject = allRules[0]
		} else {
			rule := allRules[0]
			rule.Screenshot = ""
			rule.OcrText = ""
			reportObject = rule
		}
	} else {
		if verbose {
			reportObject = allRules
		} else {
			var rules []Rule
			for _, rule := range allRules {
				rule.Screenshot = ""
				rule.OcrText = ""
				rules = append(rules, rule)
			}
			reportObject = rules
		}
	}

	switch output {
	case "json":
		returnObject, _ := json.MarshalIndent(reportObject, "", "  ")
		fmt.Println(string(returnObject))
	case "yaml":
		returnObject, _ := helpers.EncodeToYaml(reportObject)
		fmt.Println(string(returnObject))

	default:
		helpers.PrintTable(reportObject)
	}
	return nil
}

func GetRules(name string) ([]Rule, error) {
	var responseData []byte
	var err error
	if name != "" {
		responseData, err = api.GetResourceByName("rule", name)
	} else {
		responseData, err = api.GetOrderedRules()
	}
	if err != nil {
		return nil, err
	}
	var responseObject RulesResponse
	json.Unmarshal(responseData, &responseObject)

	return responseObject.Rules, nil
}

func GetRuleNames() ([]string, error) {
	var names []string
	responseData, err := api.GetResource("rule")
	if err != nil {
		return names, err
	}
	var responseObject RulesResponse
	json.Unmarshal(responseData, &responseObject)
	for _, rule := range responseObject.Rules {
		names = append(names, rule.Name)
	}
	return names, nil
}

func GenerateRule(filename string, mandatoryFlag bool, commentsFlag bool) error {
	allActions, err := GetActionNames()
	if err != nil {
		return err
	}
	actionsItemsData := make(map[string]interface{})
	actionsItemsData["type"] = "string"
	actionsItemsData["enum"] = allActions

	allRules, err := GetRuleNames()
	if err != nil {
		return err
	}
	afterRuleConstraint := make(map[string]interface{})
	afterRuleConstraint["enum"] = allRules
	afterRuleConstraint["unique"] = "only one of [before_rule, after_rule] can be set"

	stateOrScreenshotConstraint := make(map[string]interface{})
	stateOrScreenshotConstraint["unique"] = "one of [state_id, screenshot] MUST be set (screenshot takes precedence)"

	props := []helpers.PropInfo{
		{
			Name:      "name",
			Type:      "string",
			Desc:      "logical name of the rule",
			Mandatory: true,
		},
		{
			Name:        "state_id",
			Type:        "integer",
			Desc:        "ID of the state from which to take the screenshot from",
			Mandatory:   false,
			Constraints: stateOrScreenshotConstraint,
		},
		{
			Name:        "Screenshot",
			Type:        "string",
			Desc:        "base64 string of the screenshot image",
			Mandatory:   false,
			Constraints: stateOrScreenshotConstraint,
		},
		{
			Name:      "regex",
			Type:      "string",
			Desc:      "regex to use for matching states",
			Mandatory: true,
		},
		{
			Name:      "actions",
			Type:      "array",
			Desc:      "list of actions to use",
			Mandatory: true,
			ItemsData: actionsItemsData,
		},
		{
			Name:         "ignore_case",
			Type:         "boolean",
			Desc:         "whether the regex should ignore case (not case-sensitive)",
			DefaultValue: true,
			Mandatory:    false,
		},
		{
			Name:         "enabled",
			Type:         "boolean",
			Desc:         "should the rule be enabled",
			DefaultValue: true,
			Mandatory:    false,
		},
		{
			Name:        "after_rule",
			Type:        "string",
			Desc:        "after which rule name should it be placed (if not set new rules will be added last)",
			Mandatory:   false,
			Constraints: afterRuleConstraint,
		},
		{
			Name:        "before_rule",
			Type:        "string",
			Desc:        "before which rule name should it be placed (if not set new rules will be added last)",
			Mandatory:   false,
			Constraints: afterRuleConstraint,
		},
	}
	return GenerateResource(props, filename, mandatoryFlag, commentsFlag)
}
