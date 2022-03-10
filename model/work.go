package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"vaxctl/api"
	"vaxctl/helpers"
)

type WorksResponse struct {
	Works []Work `json:"works"`
}

type Work struct {
	Id          int      `json:"work_id" yaml:"work_id" header:"Id"`
	StateId     int      `json:"state_id" yaml:"state_id" header:"State Id"`
	DeviceUID   string   `json:"device_uid" yaml:"device_uid" header:"Device"`
	Actions     []Action `json:"actions" yaml:"actions"`
	Trigger     string   `json:"trigger" yaml:"trigger" header:"Trigger"`
	Assigned    string   `json:"assigned" yaml:"assigned" header:"Assigned At"`
	Status      string   `json:"status" yaml:"status" header:"Status"`
	LastUpdated string   `json:"last_updated" yaml:"last_updated"`
	CreatedAt   string   `json:"created_at" yaml:"created_at"`
}

type WorkAssignment struct {
	DeviceUID string   `json:"device_uid" yaml:"device_uid"`
	Rule      string   `json:"rule,omitempty" yaml:"rule,omitempty"`
	Actions   []string `json:"actions,omitempty" yaml:"actions,omitempty"`
}

type WorkCompleted struct {
	WorkId int    `json:"work_id" yaml:"work_id"`
	Status string `json:"status" yaml:"status"`
}

func GetWorks(workId string, deviceUID string, showDetails bool, latest bool, output string) error {
	var responseData []byte
	var err error
	if workId != "" {
		responseData, err = api.GetWorkByID(workId)
	} else if deviceUID != "" {
		responseData, err = api.GetWorksByDevice(deviceUID)
	} else {
		responseData, err = api.GetResource("work")
	}
	if err != nil {
		return err
	}
	var responseObject WorksResponse
	json.Unmarshal(responseData, &responseObject)

	if showDetails {
		if latest {
			err = ShowExecutionsByWork(responseObject.Works[len(responseObject.Works)-1].Id, output)
		} else {
			for _, work := range responseObject.Works {
				err = ShowExecutionsByWork(work.Id, output)
			}
		}
		if err != nil {
			return err
		}
	} else {
		var reportObject interface{}
		if latest {
			reportObject = responseObject.Works[len(responseObject.Works)-1]
		} else {
			reportObject = responseObject.Works
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
	}
	return nil
}

func AssignWork(deviceUID string, ruleName string, actionsList []string, filename string) error {
	var err error
	if filename != "" {
		_, err = api.PostResourceFromFile("work", filename)
	} else {
		workAssignmentData, _ := json.Marshal(WorkAssignment{deviceUID, ruleName, actionsList})
		_, err = api.CreateWorkAssignment(workAssignmentData)
	}
	return err
}

func SetWork(deviceUID string, status string) error {
	var responseData []byte
	var err error
	responseData, err = api.GetWorksByDevice(deviceUID)
	if err != nil {
		return nil
	}
	var responseObject WorksResponse
	json.Unmarshal(responseData, &responseObject)

	latestWork := responseObject.Works[len(responseObject.Works)-1]

	if latestWork.Status != "PENDING" {
		return errors.New("No pending work found for device: " + deviceUID)
	}

	err = SetExecution(latestWork.Id, latestWork.Trigger, status)
	if err != nil {
		return err
	}

	workCompletedData, err := json.Marshal(WorkCompleted{latestWork.Id, status})
	if err != nil {
		return err
	}

	_, err = api.ReportWorkCompleted(workCompletedData)
	if err != nil {
		return err
	}

	return err
}

func GetWorkIds() ([]string, error) {
	var ids []string
	responseData, err := api.GetResource("work")
	if err != nil {
		return ids, err
	}
	var responseObject WorksResponse
	json.Unmarshal(responseData, &responseObject)
	for _, work := range responseObject.Works {
		ids = append(ids, strconv.Itoa(work.Id))
	}
	return ids, nil
}
