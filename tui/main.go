package tui

import (
	"os"
	"vaxctl/model"
	"vaxctl/tui/common"
	"vaxctl/tui/models"

	tea "github.com/charmbracelet/bubbletea"
)

func CreateCred() error {
	return InteractiveCredParser("")
}

func EditCred(credName string) error {
	return InteractiveCredParser(credName)
}

func CreateDevice() error {
	return InteractiveDeviceParser("")
}

func EditDevice(deviceUID string) error {
	return InteractiveDeviceParser(deviceUID)
}

func CreateAction() error {
	return InteractiveActionParser("")
}

func EditAction(actionName string) error {
	return InteractiveActionParser(actionName)
}

func CreateRule(stateId string) error {
	return InteractiveRuleParser("", stateId)
}

func EditRule(ruleName string) error {
	return InteractiveRuleParser(ruleName, "")
}

func InteractiveCredParser(credName string) error {
	return InteractiveNavigationMenu(common.InteractiveData{CredName: credName, CurrentSubMenu: "Creds"})
}

func InteractiveDeviceParser(deviceUID string) error {
	return InteractiveNavigationMenu(common.InteractiveData{DeviceUID: deviceUID, CurrentSubMenu: "Devices"})
}

func InteractiveActionParser(actionName string) error {
	return InteractiveNavigationMenu(common.InteractiveData{ActionName: actionName, CurrentSubMenu: "Actions"})
}

func InteractiveRuleParser(ruleName string, stateId string) error {
	return InteractiveNavigationMenu(common.InteractiveData{RuleName: ruleName, CurrentSubMenu: "Rules", StateId: stateId})
}

func StartInteractiveMode() error {
	return InteractiveNavigationMenu(common.InteractiveData{})
}

func InteractiveNavigationMenu(data common.InteractiveData) error {
	var err error
	data.ActionTypes, err = model.GetActionTypes()
	if err != nil {
		return err
	}
	data.AllActionNames, err = model.GetActionNames()
	if err != nil {
		return err
	}
	data.CredNames, err = model.GetCredNames()
	if err != nil {
		return err
	}
	data.PowerOptions, err = model.GetPowerOptions()
	if err != nil {
		return err
	}
	data.SpecialKeys, err = model.GetSpecialKeys()
	if err != nil {
		return err
	}
	data.CurrentDir, err = os.Getwd()
	if err != nil {
		return err
	}
	p := tea.NewProgram(models.InitialNavigationModel(data))
	err = p.Start()
	return err
}
