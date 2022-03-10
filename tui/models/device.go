package models

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"vaxctl/api"
	"vaxctl/helpers"
	"vaxctl/model"
	"vaxctl/tui/blocks"
	"vaxctl/tui/common"
	"vaxctl/tui/editors"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	deviceUIDOption    = "UID"
	deviceIPOption     = "IPMI IP"
	deviceCredsOption  = "Creds Name"
	deviceModelOption  = "Model"
	deviceZombieOption = "Zombie"

	useDefault = "default"
)

type DeviceModel struct {
	DeviceUID          string
	IpmiIP             string
	CredsName          string
	Model              string
	Zombie             bool
	views              []string
	CurrentView        string
	DynamicView        string
	height             int
	width              int
	help               help.Model
	mainList           blocks.MainSelectionModel
	yamlNameInput      editors.TextInputEditorModel
	deviceUIDInput     editors.TextInputEditorModel
	ipmiIPInput        editors.TextInputEditorModel
	credsSelection     editors.SimpleListEditorModel
	modelInput         editors.TextInputEditorModel
	deviceResourceData blocks.ResourceDataModel
	yamlViewer         blocks.SmallViewerModel
	tableModel         blocks.TableModel
	StatusMessage      string
}

func InitialDeviceModel(interactiveData common.InteractiveData) DeviceModel {
	deviceUID := interactiveData.DeviceUID
	credNames := interactiveData.CredNames
	currentDir := interactiveData.CurrentDir
	var err error
	var statusMessage string

	var startCurrentView string
	if interactiveData.CurrentSubMenu == deviceSubMenu {
		startCurrentView = dataView
	} else {
		startCurrentView = viewerView
	}

	var ipmiIP, credsName, model string
	var isZombie bool
	if deviceUID != "" {
		ipmiIP, credsName, model, isZombie, err = fetchDevice(deviceUID)
		if err != nil {
			statusMessage = getStatusMessage(err.Error(), true)
		}
	}
	if credsName == "" {
		credsName = useDefault
	}

	help := common.GetHelpModel()

	mainList := blocks.NewMainSelectionModel("Device Actions", mainActions, showYamlAction)
	yamlNameInput := editors.NewTextInputEditor(fmt.Sprintf("Enter Yaml Name:\nAbsolute or relative to: '%s'", currentDir), "")
	deviceUIDInput := editors.NewTextInputEditor("Enter Device UID:", deviceUID)
	ipmiIPInput := editors.NewTextInputEditor("Enter Device IPMI IP:", ipmiIP)
	ipmiIPInput.SetValidationFunc(func(input string) bool { return net.ParseIP(input) != nil })
	credsSelection := editors.NewSimpleListEditor("Choose Credentials Name:", append([]string{useDefault}, credNames...), credsName)
	modelInput := editors.NewTextInputEditor("Enter Device Model:", model)
	yamlViewer := blocks.NewSmallViewerModel()

	deviceResourceData := blocks.NewResourceDataModel("Device", 0, 0, []blocks.ResourceItem{
		{
			Title:      deviceUIDOption,
			Value:      deviceUID,
			Selectable: true,
		},
		{
			Title:      deviceIPOption,
			Value:      ipmiIP,
			Selectable: true,
		},
		{
			Title:      deviceCredsOption,
			Value:      credsName,
			Selectable: true,
		},
		{
			Title:      deviceModelOption,
			Value:      model,
			Selectable: true,
		},
		{
			Title:     deviceZombieOption,
			Value:     strconv.FormatBool(isZombie),
			Toggleble: true,
		},
	})

	tableColumns := []blocks.TableColumn{
		{Title: deviceUIDOption, Size: 2},
		{Title: deviceIPOption, Size: 2},
		{Title: deviceCredsOption, Size: 1},
		{Title: deviceModelOption, Size: 2},
		{Title: deviceZombieOption, Size: 1},
	}
	rows, err := generateDeviceTableRows()
	if err != nil {
		statusMessage = getStatusMessage(err.Error(), true)
	}
	tableModel := blocks.NewTableModel(tableColumns, rows)

	return DeviceModel{
		DeviceUID:          deviceUID,
		IpmiIP:             ipmiIP,
		CredsName:          credsName,
		Model:              model,
		Zombie:             isZombie,
		views:              views,
		CurrentView:        startCurrentView,
		DynamicView:        "",
		deviceUIDInput:     deviceUIDInput,
		ipmiIPInput:        ipmiIPInput,
		credsSelection:     credsSelection,
		deviceResourceData: deviceResourceData,
		modelInput:         modelInput,
		mainList:           mainList,
		yamlNameInput:      yamlNameInput,
		help:               help,
		tableModel:         tableModel,
		yamlViewer:         yamlViewer,
		StatusMessage:      statusMessage,
	}
}

func (m DeviceModel) Update(msg tea.Msg) (DeviceModel, tea.Cmd) {

	// Handle global keys, size changes, custom messages
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		rowHeight, mainWidth, dataWidth, dynamicWidth, viewerHeight := getSizing(m.width, m.height)

		m.credsSelection.SetSize(dynamicWidth, rowHeight)
		m.deviceResourceData.SetSize(dataWidth, rowHeight)
		m.mainList.SetSize(mainWidth, rowHeight)
		m.tableModel.SetSize(m.width, viewerHeight)
		m.yamlViewer.SetSize(dynamicWidth, rowHeight)
		return m, nil

	case common.SetDynamicViewMsg:
		m.DynamicView = msg.View
		m.CurrentView = dynamicView
		switch m.DynamicView {
		case deviceUIDOption:
			m.deviceUIDInput.Focus()
		case deviceIPOption:
			m.ipmiIPInput.Focus()
		case deviceModelOption:
			m.modelInput.Focus()
		case yamlSaveOption:
			m.yamlNameInput.Focus()
		}

	case common.ExitDynamicViewMsg:
		if msg.Save {
			switch m.DynamicView {
			case deviceUIDOption:
				m.DeviceUID = m.deviceUIDInput.Value()
				m.deviceResourceData.SetValue(deviceUIDOption, m.DeviceUID)

			case deviceIPOption:
				m.IpmiIP = m.ipmiIPInput.Value()
				m.deviceResourceData.SetValue(deviceIPOption, m.IpmiIP)

			case deviceCredsOption:
				m.CredsName = m.credsSelection.Value()
				m.deviceResourceData.SetValue(deviceCredsOption, m.CredsName)

			case deviceModelOption:
				m.Model = m.modelInput.Value()
				m.deviceResourceData.SetValue(deviceModelOption, m.Model)

			case yamlSaveOption:
				path := common.GetPath(m.yamlNameInput.Value())
				err := os.WriteFile(path, m.generateYaml(), 0644)
				if err != nil {
					m.StatusMessage = getStatusMessage(err.Error(), true)
				} else {
					m.StatusMessage = getStatusMessage(fmt.Sprintf("Saved to '%s'!", path), false)
				}
			}
		}
		if m.DynamicView == yamlSaveOption {
			m.CurrentView = mainView
		} else {
			m.CurrentView = dataView
		}
		m.DynamicView = ""

	case common.ApplyMainActionMsg:
		switch msg.Action {
		case showYamlAction:
			m.yamlViewer.SetContent(string(m.generateYaml()))
			m.DynamicView = yamlViewerOption
		case saveToFileAction:
			cmd := common.SetDynamicView(yamlSaveOption)
			return m, cmd
		case saveToServerAction:
			jsonData := m.generateJson()
			_, err := api.UpdateResourceFromBytes("device", m.DeviceUID, jsonData)
			if err != nil {
				m.StatusMessage = getStatusMessage(err.Error(), true)
			} else {
				m.updateDevices()
				m.StatusMessage = getStatusMessage("Saved to Server!", false)
			}
		case clearFieldsAction:
			m.DeviceUID = ""
			m.IpmiIP = ""
			m.CredsName = ""
			m.Model = ""
			m.Zombie = false
			m.deviceResourceData.SetValue(deviceUIDOption, m.DeviceUID)
			m.deviceResourceData.SetValue(deviceIPOption, m.IpmiIP)
			m.deviceResourceData.SetValue(deviceCredsOption, m.CredsName)
			m.deviceResourceData.SetValue(deviceModelOption, m.Model)
			m.deviceResourceData.SetValue(deviceZombieOption, strconv.FormatBool(m.Zombie))
			m.DynamicView = ""
			m.CurrentView = viewerView
		case backAction:
			m.CurrentView = viewerView
			m.StatusMessage = ""
			return m, common.BackToNavigation()
		}

	case common.EditItemMsg:
		currentTableData := m.tableModel.GetCurrentRow().Data
		m.DeviceUID = currentTableData[deviceUIDOption].(string)
		m.IpmiIP = currentTableData[deviceIPOption].(string)
		m.CredsName = currentTableData[deviceCredsOption].(string)
		m.Model = currentTableData[deviceModelOption].(string)
		m.Zombie = currentTableData[deviceZombieOption].(bool)
		m.deviceResourceData.SetValue(deviceUIDOption, m.DeviceUID)
		m.deviceResourceData.SetValue(deviceIPOption, m.IpmiIP)
		m.deviceResourceData.SetValue(deviceCredsOption, m.CredsName)
		m.deviceResourceData.SetValue(deviceModelOption, m.Model)
		m.deviceResourceData.SetValue(deviceZombieOption, strconv.FormatBool(m.Zombie))
		m.deviceUIDInput.SetValue(m.DeviceUID)
		m.ipmiIPInput.SetValue(m.IpmiIP)
		m.credsSelection.SetValue(m.CredsName)
		m.modelInput.SetValue(m.Model)
		m.CurrentView = dataView
		m.DynamicView = ""

	case common.UpdateToggleValueMsg:
		switch msg.Name {
		case deviceZombieOption:
			m.Zombie = msg.Value
		}

	case common.UpdateCredNamesMsg:
		credNames, err := model.GetCredNames()
		if err != nil {
			m.StatusMessage = getStatusMessage(err.Error(), true)
		} else {
			credNames = append([]string{useDefault}, credNames...)
			m.credsSelection.UpdateItemList(credNames)
		}

	case common.RefreshDataMsg:
		m.updateDevices()
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Send message to the appropriate component
	switch m.CurrentView {
	case mainView:
		m.mainList, cmd = m.mainList.Update(msg)
		cmds = append(cmds, cmd)
	case viewerView:
		m.tableModel, cmd = m.tableModel.Update(msg)
		cmds = append(cmds, cmd)
	case dataView:
		m.deviceResourceData, cmd = m.deviceResourceData.Update(msg)
		cmds = append(cmds, cmd)
	case dynamicView:
		switch m.DynamicView {
		case deviceUIDOption:
			m.deviceUIDInput, cmd = m.deviceUIDInput.Update(msg)
			cmds = append(cmds, cmd)
		case deviceIPOption:
			m.ipmiIPInput, cmd = m.ipmiIPInput.Update(msg)
			cmds = append(cmds, cmd)
		case deviceCredsOption:
			m.credsSelection, cmd = m.credsSelection.Update(msg)
			cmds = append(cmds, cmd)
		case deviceModelOption:
			m.modelInput, cmd = m.modelInput.Update(msg)
			cmds = append(cmds, cmd)
		case yamlSaveOption:
			m.yamlNameInput, cmd = m.yamlNameInput.Update(msg)
			cmds = append(cmds, cmd)
		case yamlViewerOption:
			m.yamlViewer, cmd = m.yamlViewer.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m DeviceModel) View() string {
	mainStyle, dataStyle, dynamicStyle, viewerStyle := getWindowStyles(m.width, m.height)

	switch m.CurrentView {
	case mainView:
		mainStyle = common.ApplyActiveStyle(mainStyle)
	case dataView:
		dataStyle = common.ApplyActiveStyle(dataStyle)
	case dynamicView:
		dynamicStyle = common.ApplyActiveStyle(dynamicStyle)
	case viewerView:
		viewerStyle = common.ApplyActiveStyle(viewerStyle)
	}
	var dynamicText string
	switch m.DynamicView {
	case deviceUIDOption:
		dynamicText = m.deviceUIDInput.View()
	case deviceIPOption:
		dynamicText = m.ipmiIPInput.View()
	case deviceCredsOption:
		dynamicText = m.credsSelection.View()
	case deviceModelOption:
		dynamicText = m.modelInput.View()
	case yamlViewerOption:
		dynamicText = m.yamlViewer.View()
	case yamlSaveOption:
		dynamicText = lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.yamlNameInput.View(),
		)
	default:
		dynamicText = m.DynamicView
	}

	top := lipgloss.JoinHorizontal(
		lipgloss.Top,
		mainStyle.Render(m.mainList.View()),
		dataStyle.Render(lipgloss.PlaceHorizontal(lipgloss.Width(dataStyle.String()), lipgloss.Center, m.deviceResourceData.View())),
		dynamicStyle.Render(dynamicText))

	viewerStr := viewerStyle.Render(m.tableModel.View())
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, top),
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.StatusMessage),
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, viewerStr),
	)
}

func (m DeviceModel) generateYaml() []byte {
	yamlString, _ := helpers.EncodeToYaml(m.generateResource())
	return yamlString
}

func (m DeviceModel) generateJson() []byte {
	jsonString, _ := json.MarshalIndent(m.generateResource(), "", "  ")
	return jsonString
}

func (m DeviceModel) generateResource() model.Device {
	return model.Device{
		UID:       m.DeviceUID,
		IpmiIp:    m.IpmiIP,
		CredsName: m.CredsName,
		Model:     m.Model,
		Zombie:    m.Zombie,
	}
}

func (m *DeviceModel) updateDevices() {
	rows, err := generateDeviceTableRows()
	if err != nil {
		m.StatusMessage = getStatusMessage(err.Error(), true)
	} else {
		m.tableModel.SetRows(rows)
	}
}

func generateDeviceTableRows() ([]table.Row, error) {
	allDevices, err := model.GetDevices("")
	if err != nil {
		return nil, err
	}

	var tableRows []table.Row
	for _, device := range allDevices {
		tableRow := table.NewRow(table.RowData{
			deviceUIDOption:    device.UID,
			deviceIPOption:     device.IpmiIp,
			deviceCredsOption:  device.CredsName,
			deviceModelOption:  device.Model,
			deviceZombieOption: device.Zombie,
		})
		tableRows = append(tableRows, tableRow)
	}

	return tableRows, nil
}

func fetchDevice(deviceUID string) (string, string, string, bool, error) {
	responseData, err := api.GetResourceByUID("device", deviceUID)
	if err != nil {
		return "", "", "", false, err
	}
	var devicesResponse model.DevicesResponse
	json.Unmarshal(responseData, &devicesResponse)

	device := devicesResponse.Devices[0]
	return device.IpmiIp, device.CredsName, device.Model, device.Zombie, nil
}
