package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"vaxctl/helpers"

	"github.com/spf13/viper"
)

type ErrorResponse struct {
	Errors  map[string]string `json:"errors"`
	Message string            `json:"message"`
}

type NameRequest struct {
	Name string `json:"name"`
}

type UIDRequest struct {
	UID string `json:"uid"`
}

type CredsIdRequest struct {
	Id int `json:"creds_id"`
}

type HttpError struct {
	Status  int
	Message string
	Errors  map[string]string
}

func (httpError *HttpError) Error() string {
	message := fmt.Sprintf("Request returned %v status.", httpError.Status)
	if len(httpError.Errors) > 0 {
		for k, v := range httpError.Errors {
			message += fmt.Sprintf("\nError: %v - %v", v, k)
		}
	} else if httpError.Message != "" {
		message += fmt.Sprintf("\nError: %v", httpError.Message)
	}
	return message
}

func GetStatesByTypeAndDevice(stateType string, deviceUid string, regex string) ([]byte, error) {
	paramValues := url.Values{}
	if stateType != "" {
		paramValues.Set("type", stateType)
	}
	if deviceUid != "" {
		paramValues.Set("uid", deviceUid)
	}
	if regex != "" {
		paramValues.Set("regex", regex)
	}
	return runQuery("state/all", "GET", nil, paramValues)
}

func GetScreenshotByStateId(stateId string) ([]byte, error) {
	paramValues := url.Values{}
	paramValues.Set("id", stateId)
	return runQuery("state-screenshot/by-id", "GET", nil, paramValues)
}

func GetScreenshotByDeviceUID(deviceUid string) ([]byte, error) {
	paramValues := url.Values{}
	paramValues.Set("uid", deviceUid)
	return runQuery("state-screenshot/by-device", "GET", nil, paramValues)
}

func GetWorksByDevice(deviceUid string) ([]byte, error) {
	paramValues := url.Values{}
	paramValues.Set("uid", deviceUid)
	return runQuery("work/all/by-device", "GET", nil, paramValues)
}

func GetExecutionsByWork(workId string) ([]byte, error) {
	paramValues := url.Values{}
	paramValues.Set("id", workId)
	return runQuery("execution/all/by-work-id", "GET", nil, paramValues)
}

func CreateWorkAssignment(workAssignmentData []byte) ([]byte, error) {
	return runQuery("work/", "POST", workAssignmentData, url.Values{})
}

func ReportWorkCompleted(workCompletedData []byte) ([]byte, error) {
	return runQuery("work/by-id", "POST", workCompletedData, url.Values{})
}

func ReportExecutionCompleted(executionCompletedData []byte) ([]byte, error) {
	return runQuery("execution/", "POST", executionCompletedData, url.Values{})
}

func ReportStateAsResolved(deviceUID string) ([]byte, error) {
	paramValues := url.Values{}
	paramValues.Set("uid", deviceUID)
	return runQuery("state/resolve", "POST", nil, paramValues)
}

func UpdateResolvedState(UpdateResolvedStateData []byte) ([]byte, error) {
	return runQuery("state/update-resolve", "POST", UpdateResolvedStateData, url.Values{})
}

func SetCredsAsDefault(name string) ([]byte, error) {
	paramValues := url.Values{}
	paramValues.Set("name", name)
	return runQuery("creds/default", "PUT", nil, paramValues)
}

func GetWorkByID(workId string) ([]byte, error) {
	paramValues := url.Values{}
	paramValues.Set("id", workId)
	return runQuery("work/by-id", "GET", nil, paramValues)
}

func GetActionTypes() ([]byte, error) {
	return runQuery("action/list-types", "GET", nil, url.Values{})
}

func GetPowerOptions() ([]byte, error) {
	return runQuery("action/list-power-options", "GET", nil, url.Values{})
}

func GetSpecialKeys() ([]byte, error) {
	return runQuery("action/list-special-keys", "GET", nil, url.Values{})
}

func GetOrderedRules() ([]byte, error) {
	return runQuery("rule/ordered", "GET", nil, url.Values{})
}

func GetResource(resource string) ([]byte, error) {
	return runQuery(resource+"/all", "GET", nil, url.Values{})
}

func PostResourceFromFile(resource string, filename string) ([]byte, error) {
	return ResourceFromFile("POST", resource+"/", filename)
}

func PutResourceFromFile(resource string, filename string) ([]byte, error) {
	return ResourceFromFile("PUT", resource+"/", filename)
}

func ResourceFromFile(method string, resource string, filename string) ([]byte, error) {
	createData, err := helpers.ReadFileToJSON(filename)
	if err != nil {
		return nil, err
	}
	return runQuery(resource, method, createData, url.Values{})
}

func PostResourceFromBytes(resource string, data []byte) ([]byte, error) {
	return ResourceFromBytes("POST", resource+"/", data)
}

func PutResourceFromBytes(resource string, data []byte) ([]byte, error) {
	return ResourceFromBytes("PUT", resource+"/", data)
}

func ResourceFromBytes(method string, resource string, data []byte) ([]byte, error) {
	return runQuery(resource, method, data, url.Values{})
}

func DeleteResourceFromFile(resource string, filename string) ([]byte, error) {
	deleteData, err := helpers.ReadFileToJSON(filename)
	if err != nil {
		return nil, err
	}
	var deleteName string
	if resource == "device" {
		var deleteObject UIDRequest
		json.Unmarshal(deleteData, &deleteObject)
		deleteName = deleteObject.UID
	} else {
		var deleteObject NameRequest
		json.Unmarshal(deleteData, &deleteObject)
		deleteName = deleteObject.Name
	}
	response, err := DeleteResource(resource, deleteName)
	return response, err
}

func DeleteResource(resource string, name string) ([]byte, error) {
	paramValues := url.Values{}
	if resource == "device" {
		paramValues.Set("uid", name)
	} else {
		paramValues.Set("name", name)
	}
	return runQuery(resource+"/", "DELETE", nil, paramValues)
}

func GetResourceByUID(resource string, uid string) ([]byte, error) {
	paramValues := url.Values{}
	paramValues.Set("uid", uid)
	return runQuery(resource+"/", "GET", nil, paramValues)
}

func GetResourceByID(resource string, id string) ([]byte, error) {
	paramValues := url.Values{}
	paramValues.Set("id", id)
	return runQuery(resource+"/", "GET", nil, paramValues)
}

func GetResourceByName(resource string, name string) ([]byte, error) {
	paramValues := url.Values{}
	paramValues.Set("name", name)
	return runQuery(resource+"/", "GET", nil, paramValues)
}

func UpdateResourceFromBytes(resource string, name string, data []byte) ([]byte, error) {
	var err error
	if resource == "device" {
		_, err = GetResourceByUID(resource, name)
	} else {
		_, err = GetResourceByName(resource, name)
	}
	if err != nil {
		if errVal := err.(*HttpError); errVal.Status == http.StatusNotFound {
			return PostResourceFromBytes(resource, data)
		} else {
			return nil, err
		}
	} else {
		return PutResourceFromBytes(resource, data)
	}
}

func UpdateResourceFromFile(resource string, filename string) ([]byte, error) {
	updateData, err := helpers.ReadFileToJSON(filename)
	if err != nil {
		return nil, err
	}
	if resource == "device" {
		var getObject UIDRequest
		json.Unmarshal(updateData, &getObject)
		_, err = GetResourceByUID(resource, getObject.UID)
	} else {
		var getObject NameRequest
		json.Unmarshal(updateData, &getObject)
		_, err = GetResourceByName(resource, getObject.Name)
	}
	if err != nil {
		if errVal := err.(*HttpError); errVal.Status == http.StatusNotFound {
			return PostResourceFromFile(resource, filename)
		} else {
			return nil, err
		}
	} else {
		return PutResourceFromFile(resource, filename)
	}
}

func runQuery(urlPath string, method string, body []byte, params url.Values) ([]byte, error) {
	baseUrl := viper.GetString("url")
	if baseUrl == "" {
		baseUrl = "http://localhost:5000"
	}
	client := &http.Client{}
	reqUrl := baseUrl + "/api/v1/" + urlPath + "?" + params.Encode()
	requestBody := bytes.NewBuffer(body)

	request, err := http.NewRequest(method, reqUrl, requestBody)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		var errorObject ErrorResponse
		json.Unmarshal(responseData, &errorObject)

		return nil, &HttpError{response.StatusCode, errorObject.Message, errorObject.Errors}
	}
	return responseData, nil
}
