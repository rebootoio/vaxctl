package editors

import (
	"fmt"
	"vaxctl/tui/common"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	keyValueListMode = "list mode"
	keyValueItemMode = "item mode"

	keyValueItemCreation     = "item create"
	keyValueItemModification = "item modification"

	keyInput   = "Enter metadata key:"
	valueInput = "Enter metadata value:"
)

type KeyValueEditorModel struct {
	mode              string
	metadataList      list.Model
	textInput         textinput.Model
	selectedItemType  string
	selectedItemIndex int
	selectedItemMode  string
	help              help.Model
}

func NewKeyValueEditorModel(title string, keyValueMap map[string]string) KeyValueEditorModel {
	var listOfItems []list.Item
	for key, value := range keyValueMap {
		listOfItems = append(listOfItems, common.KeyValueItem{Key: key, Value: value})
	}

	metadataList := list.New(listOfItems, common.KeyValueItemDelegate{}, 1, 1)
	metadataList.Title = title
	metadataList.SetShowStatusBar(false)
	metadataList.SetFilteringEnabled(false)
	metadataList.SetShowHelp(false)
	metadataList.DisableQuitKeybindings()

	textInput := textinput.New()
	textInput.Focus()
	textInput.CharLimit = 64
	textInput.Width = 48

	help := common.GetHelpModel()
	help.ShowAll = true

	return KeyValueEditorModel{
		mode:         keyValueListMode,
		metadataList: metadataList,
		textInput:    textInput,
		help:         help,
	}
}

func (m *KeyValueEditorModel) SetValue(keyValueMap map[string]string) {
	var listOfItems []list.Item
	for key, value := range keyValueMap {
		listOfItems = append(listOfItems, common.KeyValueItem{Key: key, Value: value})
	}
	m.metadataList.SetItems(listOfItems)
}

func (m *KeyValueEditorModel) SetSize(width int, height int) {
	m.metadataList.SetWidth(width / 2)
	m.metadataList.SetHeight(height)
	m.help.Width = width / 2
}

func (m KeyValueEditorModel) Value() map[string]string {
	metadataMap := make(map[string]string)
	for _, item := range m.metadataList.Items() {
		itemData := item.(common.KeyValueItem)
		metadataMap[itemData.Key] = itemData.Value
	}
	return metadataMap
}

func (m KeyValueEditorModel) Update(msg tea.Msg) (KeyValueEditorModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch m.mode {
	case keyValueListMode:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, common.ConfirmKeys.ApplyData):
				cmd = common.ExitDynamicView(true)
				cmds = append(cmds, cmd)
			case key.Matches(msg, common.ConfirmKeys.ExitMode):
				cmd = common.ExitDynamicView(false)
				cmds = append(cmds, cmd)

			case key.Matches(msg, common.KeyValueListKeys.RemoveItem):
				_, ok := m.metadataList.SelectedItem().(common.KeyValueItem)
				if ok {
					m.metadataList.RemoveItem(m.metadataList.Index())
				}

			case key.Matches(msg, common.KeyValueListKeys.EditItemKey, common.KeyValueListKeys.EditItemValue):
				item, ok := m.metadataList.SelectedItem().(common.KeyValueItem)
				if ok {
					m.selectedItemIndex = m.metadataList.Index()
					m.selectedItemMode = keyValueItemModification
					m.mode = keyValueItemMode
					if key.Matches(msg, common.KeyValueListKeys.EditItemKey) {
						m.selectedItemType = keyInput
						m.textInput.SetValue(item.Key)
					} else {
						m.selectedItemType = valueInput
						m.textInput.SetValue(item.Value)
					}
					m.textInput.SetCursor(len(m.textInput.Value()))
				}

			case key.Matches(msg, common.KeyValueListKeys.AddItem):
				m.selectedItemMode = keyValueItemCreation
				m.selectedItemType = keyInput
				m.mode = keyValueItemMode
				m.textInput.SetValue("")
			}
		}
		m.metadataList, cmd = m.metadataList.Update(msg)
		cmds = append(cmds, cmd)

	case keyValueItemMode:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, common.ConfirmKeys.ApplyData):
				switch m.selectedItemMode {
				case keyValueItemCreation:
					if m.textInput.Value() != "" {
						switch m.selectedItemType {
						case keyInput:
							newIndex := len(m.metadataList.Items())
							cmd = m.metadataList.InsertItem(newIndex, common.KeyValueItem{Key: m.textInput.Value()})
							m.metadataList.Select(newIndex)
							m.selectedItemType = valueInput
							m.textInput.SetValue("")
							m.selectedItemIndex = newIndex
						case valueInput:
							item := m.metadataList.Items()[m.selectedItemIndex].(common.KeyValueItem)
							item.Value = m.textInput.Value()
							cmd = m.metadataList.SetItem(m.selectedItemIndex, item)
							m.mode = keyValueListMode
						}
					}

				case keyValueItemModification:
					newValue := m.textInput.Value()
					if newValue != "" {
						item := m.metadataList.Items()[m.selectedItemIndex].(common.KeyValueItem)
						switch m.selectedItemType {
						case keyInput:
							item.Key = newValue
						case valueInput:
							item.Value = newValue
						}
						cmd = m.metadataList.SetItem(m.selectedItemIndex, item)
						m.mode = keyValueListMode
					}
				}
				cmds = append(cmds, cmd)

			case key.Matches(msg, common.ConfirmKeys.ExitMode):
				m.mode = keyValueListMode
			}
		}
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m KeyValueEditorModel) View() string {
	var str string
	switch m.mode {
	case keyValueListMode:
		str = lipgloss.JoinHorizontal(
			lipgloss.Center,
			m.metadataList.View(),
			lipgloss.PlaceHorizontal(m.help.Width, lipgloss.Right, m.help.View(common.KeyValueListKeys)))

	case keyValueItemMode:
		str = fmt.Sprintf("%s\n%s", m.selectedItemType, m.textInput.View())
	}
	return str
}
