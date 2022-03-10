package model

import (
	"fmt"
	"io/ioutil"
	"vaxctl/api"
	"vaxctl/helpers"
)

func ApplyResource(resource string, filename string) error {
	_, err := api.UpdateResourceFromFile(resource, filename)
	return err
}

func CreateResource(resource string, filename string) error {
	_, err := api.PostResourceFromFile(resource, filename)
	return err
}

func DeleteResource(resource string, filename string, name string) error {
	var err error

	if name != "" {
		_, err = api.DeleteResource(resource, name)

	} else {
		_, err = api.DeleteResourceFromFile(resource, filename)

	}
	return err
}

func GenerateResource(props []helpers.PropInfo, filename string, mandatoryFlag bool, commentsFlag bool) error {
	generatedStr := ""

	for _, prop := range props {
		propString := helpers.GenerateProp(prop, "", mandatoryFlag, commentsFlag)
		if propString != "" {
			generatedStr += propString
			if commentsFlag {
				generatedStr += fmt.Sprint("\n\n")
			} else {
				generatedStr += fmt.Sprint("\n")
			}
		}
	}
	if filename == "" {
		fmt.Println(generatedStr)
	} else {
		err := ioutil.WriteFile(filename, []byte(generatedStr), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
