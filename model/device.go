package model

import (
	"encoding/json"
	"fmt"
	"vaxctl/api"
	"vaxctl/helpers"
)

type DevicesResponse struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	UID                string `json:"uid" yaml:"uid" header:"UID"`
	IpmiIp             string `json:"ipmi_ip" yaml:"ipmi_ip" header:"IPMI IP"`
	CredsName          string `json:"creds_name,omitempty" yaml:"creds_name,omitempty" header:"Creds Name"`
	Model              string `json:"model" yaml:"model" header:"Model"`
	Zombie             bool   `json:"zombie" yaml:"zombie" header:"Zombie"`
	AgentVersion       string `json:"agent_version,omitempty" yaml:"agent_version,omitempty" header:"Agent Version"`
	HeartbeatTimestamp string `json:"heartbeat_timestamp,omitempty" yaml:"heartbeat_timestamp,omitempty" header:"Last Heartbeat"`
	LastUpdated        string `json:"last_updated,omitempty" yaml:"last_updated,omitempty"`
	CreatedAt          string `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

func PrintDevices(name string, output string) error {
	devices, err := GetDevices(name)
	if err != nil {
		return err
	}
	var reportObject interface{}
	if name != "" {
		reportObject = devices[0]
	} else {
		reportObject = devices
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

func GetDevices(name string) ([]Device, error) {
	var responseData []byte
	var err error
	if name != "" {
		responseData, err = api.GetResourceByUID("device", name)
	} else {
		responseData, err = api.GetResource("device")
	}
	if err != nil {
		return nil, err
	}
	var responseObject DevicesResponse
	json.Unmarshal(responseData, &responseObject)
	return responseObject.Devices, nil
}

func GetDeviceNames() ([]string, error) {
	var names []string
	responseData, err := api.GetResource("device")
	if err != nil {
		return names, err
	}
	var responseObject DevicesResponse
	json.Unmarshal(responseData, &responseObject)
	for _, device := range responseObject.Devices {
		names = append(names, device.UID)
	}
	return names, nil
}

func GenerateDevice(filename string, mandatoryFlag bool, commentsFlag bool) error {
	allCreds, err := GetCredNames()
	if err != nil {
		return err
	}
	credNameConstraint := make(map[string]interface{})
	credNameConstraint["enum"] = append([]string{"default"}, allCreds...)

	props := []helpers.PropInfo{
		{
			Name:      "uid",
			Type:      "string",
			Desc:      "UID of the device",
			Mandatory: true,
		},
		{
			Name:      "ipmi_ip",
			Type:      "string",
			Desc:      "IPMI IP of the device",
			Mandatory: true,
		},
		{
			Name:      "model",
			Type:      "string",
			Desc:      "device HW model",
			Mandatory: true,
		},
		{
			Name:        "creds_name",
			Type:        "string",
			Desc:        "Credential name to use (if not set on creation default will be used)",
			Mandatory:   false,
			Constraints: credNameConstraint,
		},
		{
			Name:         "zombie",
			Type:         "boolean",
			Desc:         "Whether the device will be a zombie",
			DefaultValue: false,
			Mandatory:    false,
		},
	}
	return GenerateResource(props, filename, mandatoryFlag, commentsFlag)
}
