package models

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"vaxctl/api"
	"vaxctl/helpers"
	"vaxctl/model"
	"vaxctl/tui/blocks"
	"vaxctl/tui/common"
	"vaxctl/tui/editors"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	credNameOption      = "Name"
	credUsernameOption  = "Username"
	credPasswordOption  = "Password"
	credIsDefaultOption = "Default"
)

type CredModel struct {
	CredName               string
	Username               string
	Password               string
	views                  []string
	CurrentView            string
	DynamicView            string
	height                 int
	width                  int
	help                   help.Model
	mainList               blocks.MainSelectionModel
	yamlNameInput          editors.TextInputEditorModel
	credNameInput          editors.TextInputEditorModel
	usernameInput          editors.TextInputEditorModel
	passwordInput          editors.TextInputEditorModel
	credResourceData       blocks.ResourceDataModel
	yamlViewer             blocks.SmallViewerModel
	tableModel             blocks.TableModel
	hidePasswords          bool
	credPasswords          map[string]string
	defaultCredName        string
	pendingDefaultCredName string
	StatusMessage          string
}

func InitialCredModel(interactiveData common.InteractiveData) CredModel {
	credName := interactiveData.CredName
	currentDir := interactiveData.CurrentDir
	var err error
	var statusMessage string

	var startCurrentView string
	if interactiveData.CurrentSubMenu == credSubMenu {
		startCurrentView = dataView
	} else {
		startCurrentView = viewerView
	}

	var username, password string
	if credName != "" {
		username, password, err = fetchCred(credName)
		if err != nil {
			statusMessage = getStatusMessage(err.Error(), true)
		}
	}

	var displayPassword string
	if password != "" {
		displayPassword = getDisplayPassword(password, true)
	}

	help := common.GetHelpModel()

	mainList := blocks.NewMainSelectionModel("Cred Actions", mainActions, showYamlAction)
	yamlNameInput := editors.NewTextInputEditor(fmt.Sprintf("Enter Yaml Name:\nAbsolute or relative to: '%s'", currentDir), "")
	credNameInput := editors.NewTextInputEditor("Enter Cred Name:", credName)
	usernameInput := editors.NewTextInputEditor("Enter Username:", displayPassword)
	passwordInput := editors.NewTextInputEditor("Enter Password:", password)
	passwordInput.SetEchoMode(textinput.EchoPassword)
	yamlViewer := blocks.NewSmallViewerModel()

	credResourceData := blocks.NewResourceDataModel("Credential", 0, 0, []blocks.ResourceItem{
		{
			Title:      credNameOption,
			Value:      credName,
			Selectable: true,
		},
		{
			Title:      credUsernameOption,
			Value:      username,
			Selectable: true,
		},
		{
			Title:      credPasswordOption,
			Value:      displayPassword,
			Selectable: true,
		},
	})

	tableColumns := []blocks.TableColumn{
		{Title: credNameOption, Size: 1},
		{Title: credUsernameOption, Size: 1},
		{Title: credPasswordOption, Size: 1},
		{Title: credIsDefaultOption, Size: 1},
	}
	tableRows, credPasswords, defaultCredName, err := generateCredTableRows(true)
	if err != nil {
		statusMessage = getDisplayPassword(err.Error(), true)
	}
	tableModel := blocks.NewTableModel(tableColumns, tableRows)
	tableModel.SetAdditionalKeys(common.CredTableKeys)

	return CredModel{
		CredName:         credName,
		Username:         username,
		Password:         password,
		views:            views,
		CurrentView:      startCurrentView,
		DynamicView:      "",
		credNameInput:    credNameInput,
		usernameInput:    usernameInput,
		passwordInput:    passwordInput,
		mainList:         mainList,
		credResourceData: credResourceData,
		yamlNameInput:    yamlNameInput,
		help:             help,
		tableModel:       tableModel,
		yamlViewer:       yamlViewer,
		hidePasswords:    true,
		credPasswords:    credPasswords,
		defaultCredName:  defaultCredName,
		StatusMessage:    statusMessage,
	}
}

func (m CredModel) Update(msg tea.Msg) (CredModel, tea.Cmd) {

	// Handle global keys, size changes, custom messages
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.CredFooterKeys.ShowPasswords):
			m.hidePasswords = !m.hidePasswords
			m.credResourceData.SetValue(credPasswordOption, getDisplayPassword(m.Password, m.hidePasswords))
			m.updateCreds()
			if m.hidePasswords {
				m.passwordInput.SetEchoMode(textinput.EchoPassword)
			} else {
				m.passwordInput.SetEchoMode(textinput.EchoNormal)
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		rowHeight, mainWidth, dataWidth, dynamicWidth, viewerHeight := getSizing(m.width, m.height)
		m.credResourceData.SetSize(dataWidth, rowHeight)
		m.mainList.SetSize(mainWidth, rowHeight)
		m.tableModel.SetSize(m.width, viewerHeight)
		m.yamlViewer.SetSize(dynamicWidth, rowHeight)
		return m, nil

	case common.SetDynamicViewMsg:
		m.DynamicView = msg.View
		m.CurrentView = dynamicView
		switch m.DynamicView {
		case credNameOption:
			m.credNameInput.Focus()
		case credUsernameOption:
			m.usernameInput.Focus()
		case credPasswordOption:
			m.passwordInput.Focus()
		case yamlSaveOption:
			m.yamlNameInput.Focus()
		}

	case common.ExitDynamicViewMsg:
		if msg.Save {
			switch m.DynamicView {
			case credNameOption:
				m.CredName = m.credNameInput.Value()
				m.credResourceData.SetValue(credNameOption, m.CredName)

			case credUsernameOption:
				m.Username = m.usernameInput.Value()
				m.credResourceData.SetValue(credUsernameOption, m.Username)

			case credPasswordOption:
				m.Password = m.passwordInput.Value()
				m.credResourceData.SetValue(credPasswordOption, getDisplayPassword(m.Password, m.hidePasswords))

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
			_, err := api.UpdateResourceFromBytes("creds", m.CredName, jsonData)
			if err != nil {
				m.StatusMessage = getStatusMessage(err.Error(), true)
			} else {
				m.updateCreds()
				m.StatusMessage = getStatusMessage("Saved to Server!", false)
				return m, common.UpdateCredNames()
			}
		case clearFieldsAction:
			m.CredName = ""
			m.Username = ""
			m.Password = ""
			m.credResourceData.SetValue(credNameOption, m.CredName)
			m.credResourceData.SetValue(credUsernameOption, m.Username)
			m.credResourceData.SetValue(credPasswordOption, m.Password)
			m.DynamicView = ""
			m.CurrentView = viewerView
		case backAction:
			m.CurrentView = viewerView
			m.StatusMessage = ""
			return m, common.BackToNavigation()
		}

	case common.EditItemMsg:
		currentTableData := m.tableModel.GetCurrentRow().Data
		m.CredName = currentTableData[credNameOption].(string)
		m.Username = currentTableData[credUsernameOption].(string)
		m.Password = m.credPasswords[m.CredName]
		displayPassword := getDisplayPassword(m.Password, m.hidePasswords)
		m.credResourceData.SetValue(credNameOption, m.CredName)
		m.credResourceData.SetValue(credUsernameOption, m.Username)
		m.credResourceData.SetValue(credPasswordOption, displayPassword)
		m.credNameInput.SetValue(m.CredName)
		m.usernameInput.SetValue(m.Username)
		m.passwordInput.SetValue(m.Password)
		m.CurrentView = dataView
		m.DynamicView = ""

	case common.SetAsDefaultMsg:
		currentTableData := m.tableModel.GetCurrentRow().Data
		credName := currentTableData[credNameOption].(string)
		if credName != m.defaultCredName {
			m.pendingDefaultCredName = credName
			var newTableRows []table.Row
			for _, row := range m.tableModel.GetRows() {
				if row.Data[credNameOption] == credName {
					row.Data[credIsDefaultOption] = true
				} else {
					row.Data[credIsDefaultOption] = false
				}
				newTableRows = append(newTableRows, row)
			}
			m.tableModel.SetRows(newTableRows)
			m.tableModel.ManualChanges = true
		} else {
			m.tableModel.ManualChanges = false
			m.pendingDefaultCredName = ""
			m.updateCreds()
		}

	case common.ApplyTableChangeMsg:
		var err error
		if msg.Save {
			err = model.SetCredsAsDefault(m.pendingDefaultCredName)
		}
		if err != nil {
			m.StatusMessage = getStatusMessage(err.Error(), true)
		} else {
			m.tableModel.ManualChanges = false
			m.pendingDefaultCredName = ""
			m.updateCreds()
		}

	case common.RefreshDataMsg:
		m.updateCreds()
		return m, common.UpdateCredNames()
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
		m.credResourceData, cmd = m.credResourceData.Update(msg)
		cmds = append(cmds, cmd)
	case dynamicView:
		switch m.DynamicView {
		case credNameOption:
			m.credNameInput, cmd = m.credNameInput.Update(msg)
			cmds = append(cmds, cmd)
		case credUsernameOption:
			m.usernameInput, cmd = m.usernameInput.Update(msg)
			cmds = append(cmds, cmd)
		case credPasswordOption:
			m.passwordInput, cmd = m.passwordInput.Update(msg)
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

func (m CredModel) View() string {
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
	case credNameOption:
		dynamicText = m.credNameInput.View()
	case credUsernameOption:
		dynamicText = m.usernameInput.View()
	case credPasswordOption:
		dynamicText = m.passwordInput.View()
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
		dataStyle.Render(lipgloss.PlaceHorizontal(lipgloss.Width(dataStyle.String()), lipgloss.Center, m.credResourceData.View())),
		dynamicStyle.Render(dynamicText))

	viewerStr := viewerStyle.Align(lipgloss.Center).Render(m.tableModel.View())
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, top),
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.StatusMessage),
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, viewerStr),
	)
}

func (m CredModel) generateYaml() []byte {
	yamlString, _ := helpers.EncodeToYaml(m.generateResource())
	return yamlString
}

func (m CredModel) generateJson() []byte {
	jsonString, _ := json.MarshalIndent(m.generateResource(), "", "  ")
	return jsonString
}

func (m CredModel) generateResource() model.Cred {
	return model.Cred{
		Name:     m.CredName,
		Username: m.Username,
		Password: m.Password,
	}
}

func (m *CredModel) updateCreds() {
	tableRows, credPasswords, defaultCredName, err := generateCredTableRows(m.hidePasswords)
	if err != nil {
		m.StatusMessage = getStatusMessage(err.Error(), true)
	} else {
		m.credPasswords = credPasswords
		m.defaultCredName = defaultCredName
		m.tableModel.SetRows(tableRows)
	}
}

func generateCredTableRows(hidePasswords bool) ([]table.Row, map[string]string, string, error) {
	allCreds, err := model.GetCreds("")
	credPasswords := make(map[string]string)
	var defaultCredName string
	if err != nil {
		return nil, credPasswords, defaultCredName, err
	}

	var tableRows []table.Row
	for _, cred := range allCreds {
		displayPassword := getDisplayPassword(cred.Password, hidePasswords)
		tableRow := table.NewRow(table.RowData{
			credNameOption:      cred.Name,
			credUsernameOption:  cred.Username,
			credPasswordOption:  displayPassword,
			credIsDefaultOption: cred.IsDefault,
		})
		tableRows = append(tableRows, tableRow)
		credPasswords[cred.Name] = cred.Password
		if cred.IsDefault {
			defaultCredName = cred.Name
		}
	}

	return tableRows, credPasswords, defaultCredName, nil
}

func fetchCred(credName string) (string, string, error) {
	responseData, err := api.GetResourceByName("creds", credName)
	if err != nil {
		return "", "", err
	}
	var credsResponse model.CredsResponse
	json.Unmarshal(responseData, &credsResponse)

	cred := credsResponse.Creds[0]
	return cred.Username, cred.Password, nil
}

func getDisplayPassword(password string, hidePassword bool) string {
	if hidePassword {
		return strings.Repeat("*", len(password))
	} else {
		return password
	}
}
