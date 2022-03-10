package editors

import (
	"fmt"
	"strings"
	"vaxctl/tui/common"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	keyListView  = "key list view"
	keyInputView = "key input view"

	charType       = "char type"
	specialKeyType = "special key type"
)

type SpecialKeyEditorModel struct {
	view           string
	keysList       list.Model
	charInput      textinput.Model
	specialInput   list.Model
	currentKeyType string
	help           help.Model
}

func NewASpecialKeyEditorModel(title string, specialKeys []string) SpecialKeyEditorModel {
	keysList := list.New([]list.Item{}, common.TypedItemDelegate{}, 1, 1)
	keysList.Title = title
	keysList.SetShowStatusBar(false)
	keysList.SetFilteringEnabled(false)
	keysList.SetShowHelp(false)
	keysList.DisableQuitKeybindings()

	charInput := textinput.New()
	charInput.Focus()
	charInput.CharLimit = 1
	charInput.Width = 16

	var listOfItems []list.Item
	for _, specialKey := range specialKeys {
		listOfItems = append(listOfItems, common.SimpleItem(specialKey))
	}
	specialInput := list.New(listOfItems, common.SimpleItemDelegate{}, 1, 1)
	specialInput.Title = "Choose special key"
	specialInput.SetShowStatusBar(false)
	specialInput.SetFilteringEnabled(true)
	specialInput.SetShowHelp(false)
	specialInput.DisableQuitKeybindings()

	help := common.GetHelpModel()
	help.ShowAll = true

	return SpecialKeyEditorModel{
		view:         keyListView,
		keysList:     keysList,
		charInput:    charInput,
		specialInput: specialInput,
		help:         help,
	}
}

func (m *SpecialKeyEditorModel) EdittingMode() bool {
	return m.view == keyInputView
}

func (m *SpecialKeyEditorModel) SetSize(width int, height int) {
	m.keysList.SetWidth(width / 2)
	m.keysList.SetHeight(height)
	m.specialInput.SetWidth(width / 2)
	m.specialInput.SetHeight(height)
	m.help.Width = width / 2
}

func (m *SpecialKeyEditorModel) SetValue(newValue string) tea.Cmd {
	var listOfItems []list.Item
	if newValue != "" {
		for _, keyStr := range strings.Split(newValue, "+") {
			var keyType string
			if strings.HasPrefix(keyStr, "Keys.") {
				keyType = specialType
			} else {
				keyType = sequenceType
			}
			listOfItems = append(listOfItems, common.TypedItem{Type: keyType, Value: keyStr})
		}
	}
	return m.keysList.SetItems(listOfItems)
}

func (m SpecialKeyEditorModel) Value() string {
	var itemList []string
	for _, item := range m.keysList.Items() {
		itemList = append(itemList, item.(common.TypedItem).Value)
	}
	return strings.Join(itemList, "+")
}

func (m SpecialKeyEditorModel) Update(msg tea.Msg) (SpecialKeyEditorModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch m.view {
	case keyListView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, common.SpecialKeyListKeys.RemoveItem):
				_, ok := m.keysList.SelectedItem().(common.TypedItem)
				if ok {
					m.keysList.RemoveItem(m.keysList.Index())
				}

			case key.Matches(msg, common.SpecialKeyListKeys.AddSpecialItem, common.SpecialKeyListKeys.AddSequenceItem):
				m.view = keyInputView
				if key.Matches(msg, common.SpecialKeyListKeys.AddSpecialItem) {
					m.currentKeyType = specialKeyType
				} else {
					m.currentKeyType = charType
					m.charInput.SetValue("")
				}
			}
		}
		m.keysList, cmd = m.keysList.Update(msg)
		cmds = append(cmds, cmd)

	case keyInputView:

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, common.ConfirmKeys.ApplyData):
				if !(m.currentKeyType == specialKeyType && m.specialInput.FilterState() == list.Filtering) {
					m.view = keyListView
					var newValue string
					switch m.currentKeyType {
					case charType:
						newValue = m.charInput.Value()
					case specialKeyType:
						newKey, ok := m.specialInput.SelectedItem().(common.SimpleItem)
						if ok {
							newValue = fmt.Sprintf("Keys.%s", string(newKey))
						}
					}
					cmd = m.keysList.InsertItem(len(m.keysList.Items()), common.TypedItem{Value: newValue, Type: m.currentKeyType})
					cmds = append(cmds, cmd)
				}
			case key.Matches(msg, common.ConfirmKeys.ExitMode):
				if !(m.currentKeyType == specialKeyType && m.specialInput.FilterState() != list.Unfiltered) {
					m.view = keyListView
				}
			}
		}
		switch m.currentKeyType {
		case charType:
			m.charInput, cmd = m.charInput.Update(msg)
			cmds = append(cmds, cmd)
		case specialKeyType:
			m.specialInput, cmd = m.specialInput.Update(msg)
			cmds = append(cmds, cmd)
		}

	}
	return m, tea.Batch(cmds...)
}

func (m SpecialKeyEditorModel) View() string {
	var str string
	switch m.view {
	case keyListView:
		str = lipgloss.JoinHorizontal(
			lipgloss.Center,
			m.keysList.View(),
			lipgloss.PlaceHorizontal(m.help.Width, lipgloss.Right, m.help.View(common.SpecialKeyListKeys)))

	case keyInputView:
		switch m.currentKeyType {
		case charType:
			str = fmt.Sprintf("Enter character:\n%s", m.charInput.View())
		case specialKeyType:
			str = lipgloss.JoinHorizontal(
				lipgloss.Bottom,
				m.specialInput.View(),
				lipgloss.PlaceHorizontal(m.help.Width, lipgloss.Right, m.help.View(m.specialInput)),
			)

		}
	}
	return str
}
