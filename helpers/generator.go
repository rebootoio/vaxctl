package helpers

import (
	"fmt"
	"sort"
	"strings"
)

type PropInfo struct {
	Name         string
	Type         string
	Desc         string
	IsArray      bool
	DefaultValue interface{}
	Constraints  map[string]interface{}
	Mandatory    bool
	ItemsData    map[string]interface{}
}

func GenerateProp(propInfo PropInfo, indent string, mandatoryFlag bool, commentsFlag bool) string {
	propString := ""

	if mandatoryFlag == true && propInfo.Mandatory == false {
		return ""
	}
	if commentsFlag == true {
		comment := generatePropComment(propInfo, indent)
		propString += comment
	}

	var valueStr string
	if propInfo.IsArray {
		valueStr = fmt.Sprintf("%v- ", indent)
	} else {
		valueStr = indent
	}
	if propInfo.Mandatory == false {
		propString += fmt.Sprintf("%v#%v: ", valueStr, propInfo.Name)
	} else {
		propString += fmt.Sprintf("%v%v: ", valueStr, propInfo.Name)
	}
	if propInfo.ItemsData != nil {
		var newIndent string
		if propInfo.Mandatory {
			newIndent = indent
		} else {
			newIndent = indent + "#"
		}
		propString += fmt.Sprintf("\n%v  -", newIndent)
	}

	return propString
}

func generatePropComment(data PropInfo, indent string) string {
	comment := fmt.Sprintf("%v### %v\n", indent, data.Name)
	comment += fmt.Sprintf("%v# %v\n", indent, data.Desc)
	comment += fmt.Sprintf("%v# Type - %v\n", indent, data.Type)
	comment += fmt.Sprintf("%v# Mandatory - %v\n", indent, data.Mandatory)
	if data.DefaultValue == nil {
		comment += fmt.Sprintf("%v# Default - None\n", indent)
	} else {
		comment += fmt.Sprintf("%v# Default - %v\n", indent, data.DefaultValue)
	}
	comment += fmt.Sprintf("%v# Constraints - ", indent)
	if len(data.Constraints) == 0 {
		comment += fmt.Sprint("None\n")
	} else {
		var constraints []string
		for k, v := range data.Constraints {
			var valueStr string
			switch v.(type) {
			case []string:
				valueStr = fmt.Sprintf("[%s]", strings.Join(v.([]string), ", "))
			default:
				valueStr = fmt.Sprintf("%v", v)
			}
			constraints = append(constraints, fmt.Sprintf("%v: %v", k, valueStr))
		}
		comment += fmt.Sprint(strings.Join(constraints, ", ") + "\n")
	}
	if data.ItemsData != nil {
		comment += fmt.Sprintf("%v# Array Items:\n", indent)
		for k, v := range data.ItemsData {
			var valueStr string
			switch v.(type) {
			case []string:
				valueStr = fmt.Sprintf("[%s]", strings.Join(v.([]string), ", "))
			default:
				valueStr = fmt.Sprintf("%v", v)
			}
			comment += fmt.Sprintf("%v  # %v - %s\n", indent, strings.Title(k), valueStr)
		}
	}
	return comment
}

func sortGenerateKeys(data map[string]map[string]interface{}) []string {
	finalKeys := make([]string, 0, len(data))
	mandatoryKeys := make([]string, 0)
	optionalKeys := make([]string, 0)
	for k := range data {
		if k == "name" {
			finalKeys = append(finalKeys, k)
		} else if data[k]["mandatory"] == true {
			mandatoryKeys = append(mandatoryKeys, k)
		} else {
			optionalKeys = append(optionalKeys, k)
		}
	}
	sort.Strings(mandatoryKeys)
	sort.Strings(optionalKeys)
	finalKeys = append(finalKeys, mandatoryKeys...)
	finalKeys = append(finalKeys, optionalKeys...)
	return finalKeys
}
