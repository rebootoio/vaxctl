package models

import (
	"encoding/json"
	"fmt"
	"os"
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
	actionNameOption = "Name"
	actionTypeOption = "Type"
	actionDataOption = "Data"
)

type ActionModel struct {
	ActionName          string
	ActionType          string
	ActionData          string
	views               []string
	CurrentView         string
	DynamicView         string
	height              int
	width               int
	actionTypes         []string
	help                help.Model
	mainList            blocks.MainSelectionModel
	yamlNameInput       editors.TextInputEditorModel
	actionNameInput     editors.TextInputEditorModel
	actionTypeSelection editors.SimpleListEditorModel
	actionResourceData  blocks.ResourceDataModel
	actionDataEditor    editors.ActionDataEditorModel
	yamlViewer          blocks.SmallViewerModel
	tableModel          blocks.TableModel
	StatusMessage       string
}

func InitialActionModel(interactiveData common.InteractiveData) ActionModel {
	actionName := interactiveData.ActionName
	actionTypes := interactiveData.ActionTypes
	currentDir := interactiveData.CurrentDir
	var err error
	var statusMessage string

	var startCurrentView, startDynamicView string
	if interactiveData.CurrentSubMenu == actionSubMenu {
		startCurrentView = dynamicView
		startDynamicView = actionTypeOption
	} else {
		startCurrentView = viewerView
		startDynamicView = ""
	}

	var actionType, actionData string
	if actionName != "" {
		actionType, actionData, err = fetchAction(actionName)
		if err != nil {
			statusMessage = getStatusMessage(err.Error(), true)
		}
	}

	help := common.GetHelpModel()

	mainList := blocks.NewMainSelectionModel("Action Actions", mainActions, showYamlAction)
	yamlNameInput := editors.NewTextInputEditor(fmt.Sprintf("Enter Yaml Name:\nAbsolute or relative to: '%s'", currentDir), "")
	actionNameInput := editors.NewTextInputEditor("Enter Action Name:", actionName)
	actionTypeSelection := editors.NewSimpleListEditor("Choose Action Type:", actionTypes, actionTypes[0])
	actionDataEditor := editors.NewActionDataEditor(actionType, actionData, interactiveData.PowerOptions, interactiveData.SpecialKeys)
	yamlViewer := blocks.NewSmallViewerModel()

	actionResourceData := blocks.NewResourceDataModel("Action", 1, 0, []blocks.ResourceItem{
		{
			Title:      actionNameOption,
			Value:      actionName,
			Selectable: true,
		},
		{
			Title:      actionTypeOption,
			Value:      actionType,
			Selectable: true,
		},
		{
			Title:      actionDataOption,
			Value:      actionData,
			Selectable: true,
		},
	})

	views := []string{mainView, dataView, dynamicView, viewerView}

	tableColumns := []blocks.TableColumn{
		{Title: actionNameOption, Size: 1},
		{Title: actionTypeOption, Size: 1},
		{Title: actionDataOption, Size: 4},
	}
	tableRows, err := generateActionTableRows()
	if err != nil {
		statusMessage = getStatusMessage(err.Error(), true)
	}
	tableModel := blocks.NewTableModel(tableColumns, tableRows)

	return ActionModel{
		ActionName:          actionName,
		ActionType:          actionType,
		ActionData:          actionData,
		views:               views,
		CurrentView:         startCurrentView,
		DynamicView:         startDynamicView,
		actionNameInput:     actionNameInput,
		actionTypeSelection: actionTypeSelection,
		actionResourceData:  actionResourceData,
		actionDataEditor:    actionDataEditor,
		mainList:            mainList,
		yamlNameInput:       yamlNameInput,
		help:                help,
		tableModel:          tableModel,
		yamlViewer:          yamlViewer,
		StatusMessage:       statusMessage,
	}
}

func (m ActionModel) Update(msg tea.Msg) (ActionModel, tea.Cmd) {

	// Handle global keys, size changes, custom messages
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		rowHeight, mainWidth, dataWidth, dynamicWidth, viewerHeight := getSizing(m.width, m.height)

		m.actionTypeSelection.SetSize(dynamicWidth, rowHeight)
		m.actionDataEditor.SetSize(dynamicWidth, rowHeight)
		m.actionResourceData.SetSize(dataWidth, rowHeight)
		m.mainList.SetSize(mainWidth, rowHeight)
		m.tableModel.SetSize(m.width, viewerHeight)
		m.yamlViewer.SetSize(dynamicWidth, rowHeight)
		return m, nil

	case common.SetDynamicViewMsg:
		m.DynamicView = msg.View
		m.CurrentView = dynamicView

		switch m.DynamicView {
		case actionNameOption:
			m.actionNameInput.Focus()
		case yamlSaveOption:
			m.yamlNameInput.Focus()
		}

	case common.ExitDynamicViewMsg:
		if msg.Save {
			switch m.DynamicView {
			case actionNameOption:
				m.ActionName = m.actionNameInput.Value()
				m.actionResourceData.SetValue(actionNameOption, m.ActionName)

			case actionTypeOption:
				newValue := m.actionTypeSelection.Value()
				if newValue != m.ActionType {
					m.actionResourceData.SetValue(actionDataOption, "")
					m.ActionData = ""
				}
				m.ActionType = newValue
				m.actionResourceData.SetValue(actionTypeOption, m.ActionType)
				m.actionDataEditor.SetActionType(m.ActionType)

			case actionDataOption:
				m.ActionData = m.actionDataEditor.Value()
				m.actionResourceData.SetValue(actionDataOption, m.ActionData)

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
			_, err := api.UpdateResourceFromBytes("action", m.ActionName, jsonData)
			if err != nil {
				m.StatusMessage = getStatusMessage(err.Error(), true)
			} else {
				m.StatusMessage = getStatusMessage("Saved to Server!", false)
				m.updateActions()
				return m, common.UpdateActionNames()
			}
		case clearFieldsAction:
			m.ActionName = ""
			m.ActionType = ""
			m.ActionData = ""
			m.actionResourceData.SetValue(actionNameOption, m.ActionName)
			m.actionResourceData.SetValue(actionTypeOption, m.ActionType)
			m.actionResourceData.SetValue(actionDataOption, m.ActionData)
			m.DynamicView = ""
			m.CurrentView = viewerView
		case backAction:
			m.CurrentView = viewerView
			m.StatusMessage = ""
			return m, common.BackToNavigation()
		}

	case common.EditItemMsg:
		currentTableData := m.tableModel.GetCurrentRow().Data
		m.ActionName = currentTableData[actionNameOption].(string)
		m.ActionType = currentTableData[actionTypeOption].(string)
		m.ActionData = currentTableData[actionDataOption].(string)
		m.actionResourceData.SetValue(actionNameOption, m.ActionName)
		m.actionResourceData.SetValue(actionTypeOption, m.ActionType)
		m.actionResourceData.SetValue(actionDataOption, m.ActionData)
		m.actionNameInput.SetValue(m.ActionName)
		m.actionTypeSelection.SetValue(m.ActionType)
		m.actionDataEditor.SetActionType(m.ActionType)
		m.actionDataEditor.SetValue(m.ActionData)
		m.CurrentView = dataView
		m.DynamicView = ""

	case common.RefreshDataMsg:
		m.updateActions()
		return m, common.UpdateActionNames()
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
		m.actionResourceData, cmd = m.actionResourceData.Update(msg)
		cmds = append(cmds, cmd)
	case dynamicView:
		switch m.DynamicView {
		case actionNameOption:
			m.actionNameInput, cmd = m.actionNameInput.Update(msg)
			cmds = append(cmds, cmd)
		case actionTypeOption:
			m.actionTypeSelection, cmd = m.actionTypeSelection.Update(msg)
			cmds = append(cmds, cmd)
		case actionDataOption:
			m.actionDataEditor, cmd = m.actionDataEditor.Update(msg)
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

func (m ActionModel) View() string {
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
	case actionNameOption:
		dynamicText = m.actionNameInput.View()
	case actionTypeOption:
		dynamicText = m.actionTypeSelection.View()
	case actionDataOption:
		dynamicText = m.actionDataEditor.View()
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
		dataStyle.Render(lipgloss.PlaceHorizontal(lipgloss.Width(dataStyle.String()), lipgloss.Center, m.actionResourceData.View())),
		dynamicStyle.Render(dynamicText))

	viewerStr := viewerStyle.Render(m.tableModel.View())
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, top),
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.StatusMessage),
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, viewerStr),
	)
}

func (m ActionModel) generateYaml() []byte {
	yamlString, _ := helpers.EncodeToYaml(m.generateResource())
	return yamlString
}

func (m ActionModel) generateJson() []byte {
	jsonString, _ := json.MarshalIndent(m.generateResource(), "", "  ")
	return jsonString
}

func (m ActionModel) generateResource() model.Action {

	return model.Action{
		Name: m.ActionName,
		Type: m.ActionType,
		Data: m.ActionData,
	}
}

func (m *ActionModel) updateActions() {
	rows, err := generateActionTableRows()
	if err != nil {
		m.StatusMessage = getStatusMessage(err.Error(), true)
	} else {
		m.tableModel.SetRows(rows)
	}
}
func generateActionTableRows() ([]table.Row, error) {
	allActions, err := model.GetActions("")
	if err != nil {
		return nil, err
	}

	var tableRows []table.Row
	for _, action := range allActions {
		tableRow := table.NewRow(table.RowData{
			actionNameOption: action.Name,
			actionTypeOption: action.Type,
			actionDataOption: action.Data,
		})
		tableRows = append(tableRows, tableRow)
	}

	return tableRows, nil
}

func fetchAction(actionName string) (string, string, error) {
	responseData, err := api.GetResourceByName("action", actionName)
	if err != nil {
		return "", "", err
	}
	var actionsResponse model.ActionsResponse
	json.Unmarshal(responseData, &actionsResponse)

	action := actionsResponse.Actions[0]
	return action.Type, action.Data, nil
}
