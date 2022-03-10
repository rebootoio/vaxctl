package models

import (
	"encoding/json"
	"strconv"
	"vaxctl/api"
	"vaxctl/model"
	"vaxctl/tui/blocks"
	"vaxctl/tui/common"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	stateColumnKeyId          = "Id"
	stateColumnKeyOcr         = "Ocr Text"
	stateColumnKeyDevice      = "Device"
	stateColumnKeyResolved    = "Resolved"
	stateColumnKeyMatchedRule = "Matched Rule"
)

type StateModel struct {
	StateId           int
	Resolved          bool
	views             []string
	CurrentView       string
	DynamicView       string
	height            int
	width             int
	help              help.Model
	mainList          blocks.MainSelectionModel
	stateResourceData blocks.ResourceDataModel
	tableModel        blocks.TableModel
	StatusMessage     string
}

func InitialStateModel(interactiveData common.InteractiveData) StateModel {
	stateId := interactiveData.StateId
	var statusMessage string
	var err error

	var device string
	var resolved bool
	if stateId != "" {
		device, resolved, err = fetchState(stateId)
		if err != nil {
			statusMessage = getStatusMessage(err.Error(), true)
		}
	}

	help := common.GetHelpModel()

	mainList := blocks.NewMainSelectionModel("State Actions", stateMainActions, saveToServerAction)

	stateResourceData := blocks.NewResourceDataModel("State", 0, 2, []blocks.ResourceItem{
		{
			Title: stateColumnKeyId,
			Value: stateId,
		},
		{
			Title: stateColumnKeyDevice,
			Value: device,
		},
		{
			Title:     stateColumnKeyResolved,
			Value:     strconv.FormatBool(resolved),
			Toggleble: true,
		},
	})

	views := []string{mainView, dataView, dynamicView, viewerView}

	tableColumns := []blocks.TableColumn{
		{Title: stateColumnKeyId, Size: 1},
		{Title: stateColumnKeyOcr, Size: 4},
		{Title: stateColumnKeyDevice, Size: 2},
		{Title: stateColumnKeyResolved, Size: 1},
		{Title: stateColumnKeyMatchedRule, Size: 1},
	}
	rows, err := generateStateTableRows()
	if err != nil {
		statusMessage = getStatusMessage(err.Error(), true)
	}
	tableModel := blocks.NewTableModel(tableColumns, rows)

	parsedStateId, _ := strconv.Atoi(stateId)
	return StateModel{
		StateId:           parsedStateId,
		Resolved:          resolved,
		views:             views,
		CurrentView:       viewerView,
		stateResourceData: stateResourceData,
		mainList:          mainList,
		help:              help,
		tableModel:        tableModel,
		StatusMessage:     statusMessage,
	}
}

func (m StateModel) Update(msg tea.Msg) (StateModel, tea.Cmd) {

	// Handle global keys, size changes, custom messages
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		rowHeight, mainWidth, dataWidth, _, viewerHeight := getSizing(m.width, m.height)

		m.stateResourceData.SetSize(dataWidth, rowHeight)
		m.mainList.SetSize(mainWidth, rowHeight)
		m.tableModel.SetSize(m.width, viewerHeight)
		return m, nil

	case common.ApplyMainActionMsg:
		switch msg.Action {
		case saveToServerAction:
			err := model.UpdateResolvedState(m.StateId, m.Resolved)
			if err != nil {
				m.StatusMessage = getStatusMessage(err.Error(), true)
			} else {
				m.updateStates()
				m.StatusMessage = getStatusMessage("Saved to Server!", false)
			}
		case createRuleAction:
			if m.StateId != 0 {
				return m, common.CreateRule(strconv.Itoa(m.StateId))
			} else {
				m.StatusMessage = getStatusMessage("No state is selected", true)
			}
		case clearFieldsAction:
			m.StateId = 0
			m.Resolved = false
			m.stateResourceData.SetValue(stateColumnKeyId, "")
			m.stateResourceData.SetValue(stateColumnKeyDevice, "")
			m.stateResourceData.SetValue(stateColumnKeyResolved, strconv.FormatBool(m.Resolved))
			m.CurrentView = viewerView
		case backAction:
			m.CurrentView = viewerView
			m.StatusMessage = ""
			return m, common.BackToNavigation()
		}

	case common.EditItemMsg:
		currentTableData := m.tableModel.GetCurrentRow().Data
		stateId := currentTableData[stateColumnKeyId].(int)
		device := currentTableData[stateColumnKeyDevice].(string)
		resolved := currentTableData[stateColumnKeyResolved].(bool)
		m.stateResourceData.SetValue(stateColumnKeyId, strconv.Itoa(stateId))
		m.stateResourceData.SetValue(stateColumnKeyDevice, device)
		m.stateResourceData.SetValue(stateColumnKeyResolved, strconv.FormatBool(resolved))
		m.StateId = stateId
		m.Resolved = resolved
		m.CurrentView = dataView
		m.DynamicView = ""

	case common.UpdateToggleValueMsg:
		switch msg.Name {
		case stateColumnKeyResolved:
			m.Resolved = msg.Value
		}

	case common.RefreshDataMsg:
		m.updateStates()
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Send message to the appropriate component
	switch m.CurrentView {
	case mainView:
		m.mainList, cmd = m.mainList.Update(msg)
		cmds = append(cmds, cmd)
	case dataView:
		m.stateResourceData, cmd = m.stateResourceData.Update(msg)
		cmds = append(cmds, cmd)
	case viewerView:
		m.tableModel, cmd = m.tableModel.Update(msg)
		cmds = append(cmds, cmd)
	case dynamicView:
	}

	return m, tea.Batch(cmds...)
}

func (m StateModel) View() string {
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
	dynamicText := m.DynamicView

	top := lipgloss.JoinHorizontal(
		lipgloss.Top,
		mainStyle.Render(m.mainList.View()),
		dataStyle.Render(lipgloss.PlaceHorizontal(lipgloss.Width(dataStyle.String()), lipgloss.Center, m.stateResourceData.View())),
		dynamicStyle.Render(dynamicText))

	viewerStr := viewerStyle.Render(m.tableModel.View())
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, top),
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.StatusMessage),
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, viewerStr),
	)
}

func (m *StateModel) updateStates() {
	rows, err := generateStateTableRows()
	if err != nil {
		m.StatusMessage = getStatusMessage(err.Error(), true)
	} else {
		m.tableModel.SetRows(rows)
	}
}

func generateStateTableRows() ([]table.Row, error) {
	allStates, err := model.GetStates("", "", "", "")
	if err != nil {
		return nil, err
	}

	var tableRows []table.Row
	for _, state := range allStates {
		tableRow := table.NewRow(table.RowData{
			stateColumnKeyId:          state.StateId,
			stateColumnKeyOcr:         state.OcrText,
			stateColumnKeyDevice:      state.DeviceUID,
			stateColumnKeyResolved:    state.Resolved,
			stateColumnKeyMatchedRule: state.MatchedRule,
		})
		tableRows = append(tableRows, tableRow)
	}

	return tableRows, nil
}

func fetchState(stateId string) (string, bool, error) {
	responseData, err := api.GetResourceByID("state", stateId)
	if err != nil {
		return "", false, err
	}
	var statesResponse model.StatesResponse
	json.Unmarshal(responseData, &statesResponse)

	state := statesResponse.States[0]
	return state.DeviceUID, state.Resolved, nil
}
