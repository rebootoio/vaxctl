package model

import (
	"io/ioutil"
	"vaxctl/api"
)

func GetScreenshot(id string, deviceUID string, filename string) error {
	var responseData []byte
	var err error

	if id != "" {
		responseData, err = api.GetScreenshotByStateId(id)
	} else {
		responseData, err = api.GetScreenshotByDeviceUID(deviceUID)
	}
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, responseData, 0644)
	if err != nil {
		return err
	}
	return nil
}
