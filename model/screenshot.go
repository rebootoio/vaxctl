package model

import (
	"io/ioutil"
	"os/exec"
	"runtime"
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

	if filename != "" {
		err = ioutil.WriteFile(filename, responseData, 0644)
		if err != nil {
			return err
		}
	} else {
		file, err := ioutil.TempFile("", "vaxctl-screenshot*.png")
		if err != nil {
			return err
		}
		_, err = file.Write(responseData)
		if err != nil {
			return err
		}
		var openCommand string
		if runtime.GOOS == "darwin" {
			openCommand = "open"
		} else if runtime.GOOS == "linux" {
			openCommand = "display"
		}
		cmd := exec.Command(openCommand, file.Name())
		err = cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
