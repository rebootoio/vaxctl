package editors

import (
	"fmt"
	"strings"
	"vaxctl/tui/common"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	keysrokeType = "keystroke"
	ipmitoolType = "ipmitool"
	powerType    = "power"
	sleepType    = "sleep"
)

type ActionDataEditorModel struct {
	actionType    string
	ipmitoolInput textinput.Model
	powerList     list.Model
	waitInput     textinput.Model
	keycomboInput ActionKeycomboEditorModel
	helpText      string
}

func NewActionDataEditor(actionType string, actionData string, powerOptions []string, specialKeys []string) ActionDataEditorModel {
	ipmitoolInput := textinput.New()
	ipmitoolInput.Focus()
	ipmitoolInput.CharLimit = 64
	ipmitoolInput.Width = 48

	var powerItems []list.Item
	for _, powerOption := range powerOptions {
		powerItems = append(powerItems, common.SimpleItem(powerOption))
	}

	powerList := list.New(powerItems, common.SimpleItemDelegate{}, 1, 1)
	powerList.Title = "Select Power Action:"
	powerList.SetShowStatusBar(false)
	powerList.SetFilteringEnabled(false)
	powerList.SetShowHelp(false)
	powerList.DisableQuitKeybindings()

	waitInput := textinput.New()
	waitInput.Focus()
	waitInput.CharLimit = 16
	waitInput.Width = 24

	help := common.GetHelpModel()

	var keycomboList []string
	switch actionType {
	case keysrokeType:
		keycomboList = strings.Split(actionData, ";")
	case ipmitoolType:
		ipmitoolInput.SetValue(actionData)
	case powerType:
		for idx, item := range powerItems {
			if string(item.(common.SimpleItem)) == actionData {
				powerList.Select(idx)
				break
			}
		}
	case sleepType:
		waitInput.SetValue(actionData)
	}

	keycomboInput := NewActionKeycomboEditorModel("Key Combo list", keycomboList, specialKeys)
	return ActionDataEditorModel{
		actionType:    actionType,
		ipmitoolInput: ipmitoolInput,
		powerList:     powerList,
		waitInput:     waitInput,
		keycomboInput: keycomboInput,
		helpText:      help.View(common.ConfirmKeys),
	}
}

func (m *ActionDataEditorModel) SetValue(value string) {
	switch m.actionType {
	case keysrokeType:
		m.keycomboInput.SetValue(strings.Split(value, ";"))
	case ipmitoolType:
		m.ipmitoolInput.SetValue(value)
	case powerType:
		for idx, item := range m.powerList.Items() {
			if string(item.(common.SimpleItem)) == value {
				m.powerList.Select(idx)
				break
			}
		}
	case sleepType:
		m.waitInput.SetValue(value)
	}
}

func (m *ActionDataEditorModel) SetSize(width int, height int) {
	m.powerList.SetWidth(width / 2)
	m.powerList.SetHeight(height - lipgloss.Height(m.helpText))
	m.keycomboInput.SetSize(width, height-lipgloss.Height(m.helpText))
}

func (m *ActionDataEditorModel) SetActionType(actionType string) {
	m.actionType = actionType
}

func (m ActionDataEditorModel) Value() string {
	var value string
	switch m.actionType {
	case keysrokeType:
		value = strings.Join(m.keycomboInput.Value(), ";")
	case ipmitoolType:
		value = m.ipmitoolInput.Value()
	case powerType:
		item, _ := m.powerList.SelectedItem().(common.SimpleItem)
		value = string(item)
	case sleepType:
		value = m.waitInput.Value()
	}
	return value
}

func (m ActionDataEditorModel) Update(msg tea.Msg) (ActionDataEditorModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if !(m.actionType == keysrokeType && m.keycomboInput.EdittingMode()) {
			switch {
			case key.Matches(msg, common.ConfirmKeys.ApplyData):
				cmd = common.ExitDynamicView(true)
				cmds = append(cmds, cmd)
			case key.Matches(msg, common.ConfirmKeys.ExitMode):
				cmd = common.ExitDynamicView(false)
				cmds = append(cmds, cmd)
			}
		}
	}
	switch m.actionType {
	case keysrokeType:
		m.keycomboInput, cmd = m.keycomboInput.Update(msg)
		cmds = append(cmds, cmd)
	case ipmitoolType:
		m.ipmitoolInput, cmd = m.ipmitoolInput.Update(msg)
		cmds = append(cmds, cmd)
	case powerType:
		m.powerList, cmd = m.powerList.Update(msg)
		cmds = append(cmds, cmd)
	case sleepType:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, common.NumberKeys.Number):
				m.waitInput, cmd = m.waitInput.Update(msg)
				cmds = append(cmds, cmd)
			default:
				switch msg.Type {
				case tea.KeyRunes:
				default:
					m.waitInput, cmd = m.waitInput.Update(msg)
					cmds = append(cmds, cmd)
				}
			}
		default:
			m.waitInput, cmd = m.waitInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m ActionDataEditorModel) View() string {
	var dynamicText string
	switch m.actionType {
	case keysrokeType:
		dynamicText = m.keycomboInput.View()
	case ipmitoolType:
		dynamicText = fmt.Sprintf("Enter ipmitool action data:\n%s", m.ipmitoolInput.View())
	case powerType:
		dynamicText = m.powerList.View()
	case sleepType:
		dynamicText = fmt.Sprintf("Enter wait action data (in seconds):\n%s\nOnly numbers are allowed", m.waitInput.View())
	default:
		return "Action Type must be set to edit the data"
	}
	return lipgloss.JoinVertical(
		lipgloss.Top,
		dynamicText,
		m.helpText)
}
