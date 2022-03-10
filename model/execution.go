package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"vaxctl/api"
	"vaxctl/helpers"
)

type ExecutionsResponse struct {
	Executions []Execution `json:"executions"`
}

type Execution struct {
	Id          int         `json:"execution_id" yaml:"execution_id"`
	WorkId      int         `json:"work_id" yaml:"work_id" header:"Work Id"`
	StateId     int         `json:"state_id" yaml:"state_id" header:"State Id"`
	ActionNmae  string      `json:"action_name" yaml:"action_name" header:"Action"`
	Trigger     string      `json:"trigger" yaml:"trigger" header:"Trigger"`
	Status      string      `json:"status" yaml:"status" header:"Status"`
	ElapsedTime float32     `json:"elapsed_time" yaml:"elapsed_time"`
	LastUpdated string      `json:"last_updated" yaml:"last_updated" header:"Completed At"`
	RunData     interface{} `json:"run_data" yaml:"run_data" header:"Run Data"`
	CreatedAt   string      `json:"created_at" yaml:"created_at"`
}

func ShowExecutionsByWork(workId int, output string) error {
	responseData, err := api.GetExecutionsByWork(strconv.Itoa(workId))
	if err != nil {
		return err
	}
	var responseObject ExecutionsResponse
	json.Unmarshal(responseData, &responseObject)

	var reportObject interface{}
	reportObject = responseObject.Executions

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

func SetExecution(workId int, trigger string, status string) error {
	var err error

	executionData := Execution{WorkId: workId, ActionNmae: "Manual report", Trigger: trigger, Status: status, ElapsedTime: 0.0}

	executionCompletedData, err := json.Marshal(executionData)
	if err != nil {
		return err
	}
	_, err = api.ReportExecutionCompleted(executionCompletedData)

	if err != nil {
		return err
	}

	return err
}
