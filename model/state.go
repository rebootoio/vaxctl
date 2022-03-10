package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"vaxctl/api"
	"vaxctl/helpers"
)

type StatesResponse struct {
	States []State `json:"states"`
}

type State struct {
	StateId     int    `json:"state_id" yaml:"state_id" header:"ID"`
	Screenshot  string `json:"screenshot" yaml:"screenshot"`
	OcrText     string `json:"ocr_text" yaml:"ocr_text" header:"OCR Text"`
	DeviceUID   string `json:"device_uid" yaml:"device_uid" header:"Device"`
	Resolved    bool   `json:"resolved" yaml:"resolved" header:"Resolved"`
	MatchedRule string `json:"matched_rule" yaml:"matched_rule" header:"Matched Rule"`
	LastUpdated string `json:"last_updated" yaml:"last_updated" header:"Last Modified"`
	CreatedAt   string `json:"created_at" yaml:"created_at" header:"Created At"`
}

type updateResolved struct {
	StateId  int  `json:"state_id"`
	Resolved bool `json:"resolved"`
}

func PrintStates(id string, stateType string, deviceUid string, regex string, verbose bool, output string) error {
	states, err := GetStates(id, stateType, deviceUid, regex)
	if err != nil {
		return err
	}

	var reportObject interface{}
	if id != "" {
		if output != "json" && output != "yaml" && !verbose {
			if len(states[0].OcrText) > 100 {
				states[0].OcrText = states[0].OcrText[:100] + "..."
			}
		}
		reportObject = states[0]
	} else {
		if output != "json" && output != "yaml" && !verbose {
			for idx, state := range states {
				if len(states[0].OcrText) > 100 {
					states[idx].OcrText = state.OcrText[:100] + "..."
				}
			}
		}
		reportObject = states
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

func GetStates(id string, stateType string, deviceUid string, regex string) ([]State, error) {
	var responseData []byte
	var err error
	if id != "" {
		responseData, err = api.GetResourceByID("state", id)
	} else {
		responseData, err = api.GetStatesByTypeAndDevice(stateType, deviceUid, regex)
	}
	if err != nil {
		return nil, err
	}
	var responseObject StatesResponse
	json.Unmarshal(responseData, &responseObject)

	return responseObject.States, nil
}

func GetStateIds() ([]string, error) {
	var ids []string
	responseData, err := api.GetResource("state")
	if err != nil {
		return ids, err
	}
	var responseObject StatesResponse
	json.Unmarshal(responseData, &responseObject)
	for _, state := range responseObject.States {
		ids = append(ids, strconv.Itoa(state.StateId))
	}
	return ids, nil
}

func CreateOrUpdateState(filename string) error {
	_, err := api.PutResourceFromFile("state", filename)
	return err
}

func SetStateAsResolved(deviceUID string) error {
	_, err := api.ReportStateAsResolved(deviceUID)
	return err
}

func UpdateResolvedState(stateId int, resolved bool) error {
	data, err := json.Marshal(updateResolved{StateId: stateId, Resolved: resolved})
	if err != nil {
		return err
	}
	_, err = api.UpdateResolvedState(data)
	return err
}

func GenerateState(filename string, mandatoryFlag bool, commentsFlag bool) error {
	props := []helpers.PropInfo{
		{
			Name:      "device_uid",
			Type:      "string",
			Desc:      "UID of the device the screenshot is from",
			Mandatory: true,
		},
		{
			Name:      "screnshot",
			Type:      "string",
			Desc:      "base64 string of the screenshot image",
			Mandatory: true,
		},
	}
	return GenerateResource(props, filename, mandatoryFlag, commentsFlag)
}
