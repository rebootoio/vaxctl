package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"vaxctl/api"
	"vaxctl/helpers"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/muesli/reflow/wrap"
)

type ActionsResponse struct {
	Actions []Action `json:"actions"`
}

type Action struct {
	Name        string `json:"name" yaml:"name" header:"Name"`
	Type        string `json:"action_type" yaml:"action_type" header:"Type"`
	Data        string `json:"action_data" yaml:"action_data" header:"Data"`
	LastUpdated string `json:"last_updated" yaml:"last_updated,omitempty"`
	CreatedAt   string `json:"created_at" yaml:"created_at,omitempty"`
}

type ActionTypeResponse struct {
	ActionTypes []string `json:"action_types" yaml:"action_types"`
}

type PowerOptionsResponse struct {
	PowerOptions []string `json:"power_options" yaml:"power_options"`
}

type SpecialKeysResponse struct {
	SpecialKeys []string `json:"special_keys" yaml:"special_keys"`
}

type DetailedActionData struct {
	ActionType     string `header:"Action Type"`
	StringValue    string `header:"String Value"`
	AdditionalData string `header:"Additional Data"`
	Example        string `header:"Example"`
}

func GetActions(name string) ([]Action, error) {
	var responseData []byte
	var err error
	if name != "" {
		responseData, err = api.GetResourceByName("action", name)
	} else {
		responseData, err = api.GetResource("action")
	}
	if err != nil {
		return nil, err
	}
	var responseObject ActionsResponse
	json.Unmarshal(responseData, &responseObject)

	return responseObject.Actions, nil
}

func PrintActions(name string, output string) error {
	allActions, err := GetActions(name)
	if err != nil {
		return err
	}
	var reportObject interface{}
	if name != "" {
		reportObject = allActions[0]
	} else {
		reportObject = allActions
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

func GetActionNames() ([]string, error) {
	var names []string
	responseData, err := api.GetResource("action")
	if err != nil {
		return names, err
	}
	var responseObject ActionsResponse
	json.Unmarshal(responseData, &responseObject)
	for _, action := range responseObject.Actions {
		names = append(names, action.Name)
	}
	return names, nil
}

func GetActionTypes() ([]string, error) {
	responseData, err := api.GetActionTypes()
	if err != nil {
		return nil, err
	}
	var responseObject ActionTypeResponse
	json.Unmarshal(responseData, &responseObject)
	return responseObject.ActionTypes, nil
}

func GetPowerOptions() ([]string, error) {
	responseData, err := api.GetPowerOptions()
	if err != nil {
		return nil, err
	}
	var responseObject PowerOptionsResponse
	json.Unmarshal(responseData, &responseObject)
	return responseObject.PowerOptions, nil
}

func GetSpecialKeys() ([]string, error) {
	responseData, err := api.GetSpecialKeys()
	if err != nil {
		return nil, err
	}
	var responseObject SpecialKeysResponse
	json.Unmarshal(responseData, &responseObject)
	return responseObject.SpecialKeys, nil
}

func ListActionTypes() error {
	actionTypes, err := GetActionTypes()
	if err != nil {
		return nil
	}
	helpers.PrintTable(actionTypes)
	return nil
}

func ListDetailedActionTypes() error {
	powerOptions, err := GetPowerOptions()
	if err != nil {
		return err
	}
	specialKeys, err := GetSpecialKeys()
	if err != nil {
		return err
	}
	keyStrokeAdditionalData := `A list of key combos (either a string to send or multiple keys to be pressed at once):
 - string: the string as it should be entered (each key will be pressed in sequence)
 - multiple keys: seperated by '+' (special keys can be used by prepending 'Keys.' to the key name)
List of special keys:
[%s]`

	actionTypesData := []DetailedActionData{
		{
			ActionType:     "sleep",
			StringValue:    "number of seconds to sleep",
			AdditionalData: "only digits are allowed",
			Example:        "10",
		},
		{
			ActionType:     "power",
			StringValue:    "power action to send via ipmitool",
			AdditionalData: fmt.Sprintf("one of: [%s]", strings.Join(powerOptions, ", ")),
			Example:        "on",
		},
		{
			ActionType:     "ipmitool",
			StringValue:    "free text",
			AdditionalData: "will be appended to ipmitool command",
			Example:        "lan print",
		},
		{
			ActionType:     "keystroke",
			StringValue:    "a ';' seperated list of key combos",
			AdditionalData: wrap.String(fmt.Sprintf(keyStrokeAdditionalData, strings.Join(specialKeys, ", ")), 100),
			Example:        "Keys.Control+c;exit;Keys.Enter",
		},
		{
			ActionType:     "request",
			StringValue:    "URI for GET request",
			AdditionalData: "must start with protocol",
			Example:        "http://myservice:8080/resolve?uid={device::uid}",
		},
	}

	var rowsData [][]table.Row
	for _, data := range actionTypesData {
		rows := []table.Row{
			{"Action Type:", data.ActionType},
			{"String Value:", data.StringValue},
			{"Additional Data:", data.AdditionalData},
			{"Example:", data.Example},
		}
		rowsData = append(rowsData, rows)
	}

	helpers.PrintTableWithBorders(rowsData, "")
	attributesTableTitle := "You can use a device's attributes in the action data string.\nAvailable attributes:"
	attributesTableHeader := []table.Row{
		{"BASE KEY", "NESTED KEY", "DEFINED BY", "DESCRIPTION", "USAGE"},
	}
	attributesTableRows := []table.Row{
		{"device", "uid", "System", "the device's UID", "{device::uid}"},
		{"device", "ipmi_ip", "System", "the device's ipmi IP", "{device::ipmi_ip}"},
		{"device", "model", "System", "the device's model", "{device::model}"},
		{"cred", "username", "System", "the device's cred username", "{cred::username}"},
		{"cred", "password", "System", "the device's cred password", "{cred::password}"},
		{"metadata", "*", "User", "the value of the nested key from the device's metadata", "{metadata::ANY_KEY}"},
		{"cred_store", "CRED_NAME::username", "System", "a username from an existing cred", "{cred_store::CRED_NAME::username}"},
		{"cred_store", "CRED_NAME::password", "System", "a password from an existing cred", "{cred_store::CRED_NAME::password}"},
	}
	helpers.PrintTableWithBorders([][]table.Row{attributesTableHeader, attributesTableRows}, attributesTableTitle)
	return nil
}

func GenerateAction(filename string, mandatoryFlag bool, commentsFlag bool) error {
	allActionTypes, err := GetActionTypes()
	if err != nil {
		return err
	}
	actionTypeConstraint := make(map[string]interface{})
	actionTypeConstraint["enum"] = allActionTypes

	props := []helpers.PropInfo{
		{
			Name:      "name",
			Type:      "string",
			Desc:      "logical name of the action",
			Mandatory: true,
		},
		{
			Name:        "action_type",
			Type:        "string",
			Desc:        "action type",
			Mandatory:   true,
			Constraints: actionTypeConstraint,
		},
		{
			Name:      "action_data",
			Type:      "string",
			Desc:      "string with the action data (based on type), for more info use 'create action -I -v'",
			Mandatory: true,
		},
	}
	return GenerateResource(props, filename, mandatoryFlag, commentsFlag)
}
