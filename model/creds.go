package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"vaxctl/api"
	"vaxctl/helpers"
)

type CredsResponse struct {
	Creds []Cred `json:"creds"`
}

type Cred struct {
	Name        string `json:"name" yaml:"name" header:"Name"`
	Username    string `json:"username" yaml:"username" header:"Username"`
	Password    string `json:"password" yaml:"password" header:"Password"`
	IsDefault   bool   `json:"is_default,omitempty" yaml:"is_default,omitempty" header:"Default"`
	LastUpdated string `json:"last_updated,omitempty" yaml:"last_updated,omitempty"`
	CreatedAt   string `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

func PrintCreds(name string, output string) error {
	creds, err := GetCreds(name)
	if err != nil {
		return err
	}
	var reportObject interface{}
	if name != "" {
		if output != "json" && output != "yaml" {
			creds[0].Password = strings.Repeat("*", len(creds[0].Password))
		}
		reportObject = creds[0]
	} else {
		if output != "json" && output != "yaml" {
			for idx, cred := range creds {
				creds[idx].Password = strings.Repeat("*", len(cred.Password))
			}
		}
		reportObject = creds
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

func GetCreds(name string) ([]Cred, error) {
	var responseData []byte
	var err error
	if name != "" {
		responseData, err = api.GetResourceByName("creds", name)
	} else {
		responseData, err = api.GetResource("creds")
	}
	if err != nil {
		return nil, err
	}
	var responseObject CredsResponse
	json.Unmarshal(responseData, &responseObject)
	return responseObject.Creds, nil
}

func GetCredNames() ([]string, error) {
	var names []string
	creds, err := GetCreds("")
	if err != nil {
		return nil, err
	}
	for _, cred := range creds {
		names = append(names, cred.Name)
	}
	return names, nil
}

func SetCredsAsDefault(name string) error {
	_, err := api.SetCredsAsDefault(name)
	return err
}

func GenerateCred(filename string, mandatoryFlag bool, commentsFlag bool) error {
	props := []helpers.PropInfo{
		{
			Name:      "name",
			Type:      "string",
			Desc:      "logical name for the credentials",
			Mandatory: true,
		},
		{
			Name:      "username",
			Type:      "string",
			Desc:      "username for credentials",
			Mandatory: true,
		},
		{
			Name:      "password",
			Type:      "string",
			Desc:      "password for credentials",
			Mandatory: true,
		},
	}
	return GenerateResource(props, filename, mandatoryFlag, commentsFlag)
}
