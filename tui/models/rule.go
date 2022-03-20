package models

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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
	ruleColumnKeyName       = "Name"
	ruleColumnKeyState      = "State ID"
	ruleColumnKeyRegex      = "Regex"
	ruleColumnKeyActions    = "Actions"
	ruleColumnKeyIgnoreCase = "Ignore Case"
	ruleColumnKeyEnabled    = "Enabled"
	ruleColumnKeyPosition   = "Position"
)

type RuleModel struct {
	StateId          string
	RegexString      string
	RuleName         string
	ChosenActionList []string
	IgnoreCase       bool
	Enabled          bool
	ocrText          string
	views            []string
	CurrentView      string
	DynamicView      string
	Editting         bool
	height           int
	width            int
	help             help.Model
	regexInput       editors.TextInputEditorModel
	ruleNameInput    editors.TextInputEditorModel
	yamlNameInput    editors.TextInputEditorModel
	actionListEditor editors.OrderedSelectionEditorModel
	ruleResourceData blocks.ResourceDataModel
	ocrViewer        blocks.LargeViewerModel
	mainList         blocks.MainSelectionModel
	yamlViewer       blocks.SmallViewerModel
	tableModel       blocks.TableModel
	reorderedRules   map[string]int
	StatusMessage    string
}

func InitialRuleModel(interactiveData common.InteractiveData) RuleModel {
	ruleName := interactiveData.RuleName
	stateId := interactiveData.StateId
	actionNameList := interactiveData.AllActionNames
	currentDir := interactiveData.CurrentDir
	var err error
	var statusMessage string

	var startCurrentView, startDynamicView string
	if interactiveData.CurrentSubMenu == ruleSubMenu {
		startCurrentView = dynamicView
		startDynamicView = ruleColumnKeyRegex
	} else {
		startCurrentView = viewerView
		startDynamicView = ""
	}

	var regexString, ocrText string
	var chosenActionList []string
	ignoreCase := true
	enabled := true
	if ruleName != "" {
		stateId, regexString, chosenActionList, ignoreCase, enabled, ocrText, err = fetchRule(ruleName)
	} else {
		if stateId != "" {
			ocrText, err = fetchOcrTextFromState(stateId)
		}
	}
	if err != nil {
		statusMessage = getStatusMessage(err.Error(), true)
	}

	help := common.GetHelpModel()

	mainList := blocks.NewMainSelectionModel("Rule Actions", mainActions, showYamlAction)
	ruleNameEditInput := editors.NewTextInputEditor("Enter Rule Name:", ruleName)
	regexTextInput := editors.NewTextInputEditor("Enter Regex:", regexString)
	yamlNameInput := editors.NewTextInputEditor(fmt.Sprintf("Enter Yaml Name:\nAbsolute or relative to: '%s'", currentDir), "")
	regexTextInput.Focus()
	actionListEditor := editors.NewOrderedSelectionEditor("Select actions for rule:", actionNameList, chosenActionList)
	ocrViewer := blocks.NewLargeViewerModel()
	yamlViewer := blocks.NewSmallViewerModel()

	ruleResourceData := blocks.NewResourceDataModel("Rule", 1, 0, []blocks.ResourceItem{
		{
			Title:      ruleColumnKeyName,
			Value:      ruleName,
			Selectable: true,
		},
		{
			Title:      ruleColumnKeyRegex,
			Value:      regexString,
			Selectable: true,
		},
		{
			Title:      ruleColumnKeyActions,
			Value:      strings.Join(chosenActionList, ", "),
			Selectable: true,
		},
		{
			Title:     ruleColumnKeyIgnoreCase,
			Value:     strconv.FormatBool(ignoreCase),
			Toggleble: true,
		},
		{
			Title:     ruleColumnKeyEnabled,
			Value:     strconv.FormatBool(enabled),
			Toggleble: true,
		},
	})

	ocrContent, err := createRegexAndColorOcrText(regexString, ignoreCase, ocrText)
	if err != nil {
		statusMessage = getStatusMessage(err.Error(), true)
	}
	ocrViewer.SetContent(ocrContent)

	views := []string{mainView, dataView, dynamicView, viewerView}

	tableColumns := []blocks.TableColumn{
		{Title: ruleColumnKeyName, Size: 2},
		{Title: ruleColumnKeyState, Size: 1},
		{Title: ruleColumnKeyRegex, Size: 2},
		{Title: ruleColumnKeyActions, Size: 4},
		{Title: ruleColumnKeyIgnoreCase, Size: 1},
		{Title: ruleColumnKeyEnabled, Size: 1},
		{Title: ruleColumnKeyPosition, Size: 1},
	}
	rows, err := generateRuleTableRows()
	if err != nil {
		statusMessage = getStatusMessage(err.Error(), true)
	}
	tableModel := blocks.NewTableModel(tableColumns, rows)
	tableModel.SetAdditionalKeys(common.RuleTableKeys)

	return RuleModel{
		RuleName:         ruleName,
		RegexString:      regexString,
		StateId:          stateId,
		ChosenActionList: chosenActionList,
		ocrText:          ocrText,
		views:            views,
		CurrentView:      startCurrentView,
		DynamicView:      startDynamicView,
		IgnoreCase:       ignoreCase,
		Enabled:          enabled,
		regexInput:       regexTextInput,
		ocrViewer:        ocrViewer,
		ruleNameInput:    ruleNameEditInput,
		actionListEditor: actionListEditor,
		ruleResourceData: ruleResourceData,
		mainList:         mainList,
		yamlNameInput:    yamlNameInput,
		help:             help,
		tableModel:       tableModel,
		yamlViewer:       yamlViewer,
		Editting:         ocrText != "",
		reorderedRules:   make(map[string]int),
		StatusMessage:    statusMessage,
	}
}

func (m RuleModel) Update(msg tea.Msg) (RuleModel, tea.Cmd) {

	// Handle global keys, size changes, custom messages
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		rowHeight, mainWidth, dataWidth, dynamicWidth, viewerHeight := getSizing(m.width, m.height)

		m.actionListEditor.SetSize(dynamicWidth, rowHeight)
		m.ruleResourceData.SetSize(dataWidth, rowHeight)
		m.mainList.SetSize(mainWidth, rowHeight)
		m.ocrViewer.SetSize(m.width, viewerHeight)
		m.tableModel.SetSize(m.width, viewerHeight)
		m.yamlViewer.SetSize(dynamicWidth, rowHeight)
		return m, nil

	case common.SetDynamicViewMsg:
		m.DynamicView = msg.View
		m.CurrentView = dynamicView

		switch m.DynamicView {
		case ruleColumnKeyName:
			m.ruleNameInput.Focus()
		case ruleColumnKeyRegex:
			m.regexInput.Focus()
		case yamlSaveOption:
			m.yamlNameInput.Focus()
		}

	case common.ExitDynamicViewMsg:
		if msg.Save {
			switch m.DynamicView {
			case ruleColumnKeyName:
				m.RuleName = m.ruleNameInput.Value()
				m.ruleResourceData.SetValue(ruleColumnKeyName, m.RuleName)

			case ruleColumnKeyRegex:
				err := m.updateRegex(m.regexInput.Value())
				if err != nil {
					return m, nil
				} else {
					m.RegexString = m.regexInput.Value()
					m.ruleResourceData.SetValue(ruleColumnKeyRegex, m.RegexString)
				}

			case ruleColumnKeyActions:
				m.ChosenActionList = m.actionListEditor.Value()
				m.ruleResourceData.SetValue(ruleColumnKeyActions, strings.Join(m.ChosenActionList, ", "))

			case yamlSaveOption:
				path := common.GetPath(m.yamlNameInput.Value())
				err := os.WriteFile(path, m.generateYaml(), 0644)
				if err != nil {
					m.StatusMessage = getStatusMessage(err.Error(), true)
				} else {
					m.StatusMessage = getStatusMessage(fmt.Sprintf("Saved to '%s'!", path), false)
				}
			}
		} else {
			if m.DynamicView == ruleColumnKeyRegex {
				m.updateRegex(m.RegexString)
			}
		}
		if m.DynamicView == yamlSaveOption {
			m.CurrentView = mainView
		} else {
			m.CurrentView = dataView
		}
		m.DynamicView = ""

	case common.UpdateToggleValueMsg:
		switch msg.Name {
		case ruleColumnKeyIgnoreCase:
			m.IgnoreCase = msg.Value
			m.updateRegex(m.RegexString)
		case ruleColumnKeyEnabled:
			m.Enabled = msg.Value
		}

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
			_, err := api.UpdateResourceFromBytes("rule", m.RuleName, jsonData)
			if err != nil {
				m.StatusMessage = getStatusMessage(err.Error(), true)
			} else {
				m.updateRules()
				m.StatusMessage = getStatusMessage("Saved to Server!", false)
			}
		case clearFieldsAction:
			m.RuleName = ""
			m.RegexString = ""
			m.ChosenActionList = []string{}
			m.IgnoreCase = true
			m.ruleResourceData.SetValue(ruleColumnKeyName, m.RuleName)
			m.ruleResourceData.SetValue(ruleColumnKeyRegex, m.RegexString)
			m.ruleResourceData.SetValue(ruleColumnKeyActions, strings.Join(m.ChosenActionList, ", "))
			m.ruleResourceData.SetValue(ruleColumnKeyIgnoreCase, strconv.FormatBool(m.IgnoreCase))
			m.ocrText = ""
			m.ocrViewer.SetContent(m.ocrText)
			m.actionListEditor.ClearSelected()
			m.DynamicView = ""
			m.CurrentView = viewerView
			m.Editting = false

		case backAction:
			m.CurrentView = viewerView
			m.StatusMessage = ""
			if m.Editting {
				m.Editting = false
			} else {

				return m, common.BackToNavigation()
			}
		}

	case common.EditItemMsg:
		currentTableData := m.tableModel.GetCurrentRow().Data
		m.RuleName = currentTableData[ruleColumnKeyName].(string)
		m.StateId = currentTableData[ruleColumnKeyState].(string)
		m.RegexString = currentTableData[ruleColumnKeyRegex].(string)
		chosenActionsStr := currentTableData[ruleColumnKeyActions].(string)
		m.ChosenActionList = strings.Split(chosenActionsStr, ", ")
		m.IgnoreCase, _ = currentTableData[ruleColumnKeyIgnoreCase].(bool)
		m.Enabled, _ = currentTableData[ruleColumnKeyEnabled].(bool)
		m.ruleResourceData.SetValue(ruleColumnKeyName, m.RuleName)
		m.ruleResourceData.SetValue(ruleColumnKeyRegex, m.RegexString)
		m.ruleResourceData.SetValue(ruleColumnKeyActions, chosenActionsStr)
		m.ruleResourceData.SetValue(ruleColumnKeyIgnoreCase, strconv.FormatBool(m.IgnoreCase))
		m.ruleResourceData.SetValue(ruleColumnKeyEnabled, strconv.FormatBool(m.Enabled))
		m.ruleNameInput.SetValue(m.RuleName)
		m.regexInput.SetValue(m.RegexString)
		m.actionListEditor.SetSelected(m.ChosenActionList)
		var err error
		m.ocrText, err = fetchOcrTextFromState(m.StateId)
		if err != nil {
			m.StatusMessage = getStatusMessage(err.Error(), true)
		}
		m.updateRegex(m.RegexString)
		m.CurrentView = dataView
		m.DynamicView = ""
		m.Editting = true

	case common.UpdateActionNamesMsg:
		actionNames, err := model.GetActionNames()
		if err != nil {
			m.StatusMessage = getStatusMessage(err.Error(), true)
		} else {
			m.actionListEditor.UpdateItemList(actionNames)
		}

	case common.RefreshDataMsg:
		m.updateRules()

	case common.ChangeOrderMsg:
		currentRow := m.tableModel.GetCurrentRow()
		currentRows := m.tableModel.GetRows()
		var newRows []table.Row
		var newIndex int

		ruleName := fmt.Sprintf("%v", currentRow.Data[ruleColumnKeyName])
		originalIndex, found := m.reorderedRules[ruleName]
		for idx, row := range currentRows {
			if row.Data[ruleColumnKeyName] == ruleName {
				if !found {
					m.reorderedRules[ruleName] = idx
				}
				if msg.Up {
					newIndex = idx - 1
				} else {
					newIndex = idx + 1
				}
			} else {
				newRows = append(newRows, row)
			}
		}
		if found {
			if newIndex == originalIndex {
				delete(m.reorderedRules, ruleName)
			}
		}
		if newIndex == -1 {
			newRows = append([]table.Row{currentRow}, newRows...)
		} else if newIndex >= len(newRows) {
			newRows = append(newRows, currentRow)
		} else {
			newRows = append(newRows[:newIndex+1], newRows[newIndex:]...)
			newRows[newIndex] = currentRow
		}
		m.tableModel.SetRows(newRows)
		m.tableModel.SetHighlightedRow(newIndex)
		m.tableModel.ManualChanges = len(m.reorderedRules) != 0

	case common.ApplyTableChangeMsg:
		if msg.Save {
			var afterRule string
			tableRows := m.tableModel.GetRows()
			for _, row := range tableRows {
				ruleName := fmt.Sprintf("%v", row.Data[ruleColumnKeyName])
				if _, ok := m.reorderedRules[ruleName]; ok {
					var err error
					if afterRule == "" {
						err = updateRuleOrderBefore(ruleName, fmt.Sprintf("%v", tableRows[1].Data[ruleColumnKeyName]))
					} else {
						err = updateRuleOrderAfter(ruleName, afterRule)
					}
					if err != nil {
						m.StatusMessage = getStatusMessage(err.Error(), true)
					} else {
						m.StatusMessage = getStatusMessage("Changes saved to server!", false)
					}
				}
				afterRule = ruleName
			}
		}

		m.reorderedRules = make(map[string]int)
		m.updateRules()
		m.tableModel.ManualChanges = false
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Send message to the appropriate component
	switch m.CurrentView {
	case mainView:
		m.mainList, cmd = m.mainList.Update(msg)
		cmds = append(cmds, cmd)
	case viewerView:
		if m.Editting {
			m.ocrViewer, cmd = m.ocrViewer.Update(msg)
		} else {
			m.tableModel, cmd = m.tableModel.Update(msg)
		}
		cmds = append(cmds, cmd)
	case dataView:
		m.ruleResourceData, cmd = m.ruleResourceData.Update(msg)
		cmds = append(cmds, cmd)
	case dynamicView:
		switch m.DynamicView {
		case ruleColumnKeyName:
			m.ruleNameInput, cmd = m.ruleNameInput.Update(msg)
			cmds = append(cmds, cmd)
		case ruleColumnKeyRegex:
			m.regexInput, cmd = m.regexInput.Update(msg)
			cmds = append(cmds, cmd)
			m.updateRegex(m.regexInput.Value())
		case ruleColumnKeyActions:
			m.actionListEditor, cmd = m.actionListEditor.Update(msg)
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

func (m RuleModel) View() string {

	mainStyle, dataStyle, dynamicStyle, viewerStyle := getWindowStyles(m.width, m.height)

	m.ocrViewer.SetActive(m.CurrentView == viewerView)

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
	case ruleColumnKeyName:
		dynamicText = m.ruleNameInput.View()
	case ruleColumnKeyRegex:
		dynamicText = m.regexInput.View()
	case ruleColumnKeyActions:
		dynamicText = m.actionListEditor.View()
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

	var viewerStr string
	if m.Editting {
		viewerStr = m.ocrViewer.View()
	} else {
		viewerStr = viewerStyle.Align(lipgloss.Center).Render(m.tableModel.View())
	}
	top := lipgloss.JoinHorizontal(
		lipgloss.Top,
		mainStyle.Render(m.mainList.View()),
		dataStyle.Render(lipgloss.PlaceHorizontal(lipgloss.Width(dataStyle.String()), lipgloss.Center, m.ruleResourceData.View())),
		dynamicStyle.Render(dynamicText))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, top),
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.StatusMessage),
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, viewerStr),
	)
}

func (m RuleModel) generateYaml() []byte {
	yamlString, _ := helpers.EncodeToYaml(m.generateResource())
	return yamlString
}

func (m RuleModel) generateJson() []byte {
	jsonString, _ := json.MarshalIndent(m.generateResource(), "", "  ")
	return jsonString
}

func (m RuleModel) generateResource() model.Rule {
	stateId, _ := strconv.Atoi(m.StateId)

	return model.Rule{
		Name:       m.RuleName,
		StateId:    stateId,
		Regex:      m.RegexString,
		Actions:    m.ChosenActionList,
		IgnoreCase: m.IgnoreCase,
		Enabled:    m.Enabled,
	}
}

func (m *RuleModel) UpdateStateData(stateId string) {
	m.RuleName = ""
	m.RegexString = ""
	m.ChosenActionList = []string{}
	m.IgnoreCase = true
	m.StateId = stateId
	m.ruleResourceData.SetValue(ruleColumnKeyName, m.RuleName)
	m.ruleResourceData.SetValue(ruleColumnKeyRegex, m.RegexString)
	m.ruleResourceData.SetValue(ruleColumnKeyActions, strings.Join(m.ChosenActionList, ", "))
	m.ruleResourceData.SetValue(ruleColumnKeyIgnoreCase, strconv.FormatBool(m.IgnoreCase))
	var err error
	m.ocrText, err = fetchOcrTextFromState(stateId)
	if err != nil {
		m.StatusMessage = getStatusMessage(err.Error(), true)
	}
	m.updateRegex(m.RegexString)
	m.CurrentView = dynamicView
	m.DynamicView = ruleColumnKeyRegex
	m.Editting = true
}

func (m *RuleModel) updateRules() {
	rows, err := generateRuleTableRows()
	if err != nil {
		m.StatusMessage = getStatusMessage(err.Error(), true)
	} else {
		m.tableModel.SetRows(rows)
	}
}

func (m *RuleModel) updateRegex(regexString string) error {
	ocrContent, err := createRegexAndColorOcrText(regexString, m.IgnoreCase, m.ocrText)
	if err != nil {
		m.StatusMessage = getStatusMessage(fmt.Sprintf("Regex error: %s", err.Error()), true)
	} else {
		m.StatusMessage = ""
	}
	m.ocrViewer.SetContent(ocrContent)
	return err
}

func createRegexAndColorOcrText(regexString string, ignoreCase bool, ocrText string) (string, error) {
	var regex *regexp.Regexp
	var err error

	newLineCount := strings.Count(regexString, "\\n")

	// if we have newline in regex we separate them to match groups
	if newLineCount > 0 {
		regexString = fmt.Sprintf("(%s)", strings.ReplaceAll(regexString, "\\n", ")\\n("))
	}

	if ignoreCase {
		regex, err = regexp.Compile("(?i)" + regexString)
	} else {
		regex, err = regexp.Compile(regexString)
	}
	if err != nil {
		return ocrText, err
	}

	var parsedOcrText string
	if newLineCount > 0 {
		replString := common.OcrMatchedTextStyle.Render("${1}")
		// for every newline we add a \n + coloring the match group
		for i := 2; i < newLineCount+2; i++ {
			replString += "\n" + common.OcrMatchedTextStyle.Render(fmt.Sprintf("${%d}", i))
		}
		parsedOcrText = regex.ReplaceAllString(ocrText, replString)

	} else {
		parsedOcrText = regex.ReplaceAllString(ocrText, common.OcrMatchedTextStyle.Render("${0}"))
	}

	return parsedOcrText, nil
}

func fetchRule(ruleName string) (string, string, []string, bool, bool, string, error) {
	responseData, err := api.GetResourceByName("rule", ruleName)
	if err != nil {
		return "", "", nil, false, false, "", err
	}
	var rulesResponse model.RulesResponse
	json.Unmarshal(responseData, &rulesResponse)

	rule := rulesResponse.Rules[0]
	stateId := strconv.Itoa(rule.StateId)
	ocrText, err := fetchOcrTextFromState(stateId)

	return stateId, rule.Regex, rule.Actions, rule.IgnoreCase, rule.Enabled, ocrText, err
}

func fetchOcrTextFromState(stateId string) (string, error) {
	responseData, err := api.GetResourceByID("state", stateId)
	if err != nil {
		return "", err
	}
	var responseObject model.StatesResponse
	json.Unmarshal(responseData, &responseObject)
	state := responseObject.States[0]
	return state.OcrText, nil
}

func generateRuleTableRows() ([]table.Row, error) {
	allRules, err := model.GetRules("")
	if err != nil {
		return nil, err
	}

	var tableRows []table.Row
	for _, rule := range allRules {
		tableRow := table.NewRow(table.RowData{
			ruleColumnKeyName:       rule.Name,
			ruleColumnKeyState:      strconv.Itoa(rule.StateId),
			ruleColumnKeyRegex:      rule.Regex,
			ruleColumnKeyActions:    strings.Join(rule.Actions, ", "),
			ruleColumnKeyIgnoreCase: rule.IgnoreCase,
			ruleColumnKeyEnabled:    rule.Enabled,
			ruleColumnKeyPosition:   rule.Position,
		})
		tableRows = append(tableRows, tableRow)
	}
	return tableRows, nil
}

func updateRuleOrderAfter(ruleName string, afterRule string) error {
	jsonData, _ := json.Marshal(model.Rule{Name: ruleName, AfterRule: afterRule})
	_, err := api.UpdateResourceFromBytes("rule", ruleName, jsonData)
	return err
}

func updateRuleOrderBefore(ruleName string, beforeRule string) error {
	jsonData, _ := json.Marshal(model.Rule{Name: ruleName, BeforeRule: beforeRule})
	_, err := api.UpdateResourceFromBytes("rule", ruleName, jsonData)
	return err
}
