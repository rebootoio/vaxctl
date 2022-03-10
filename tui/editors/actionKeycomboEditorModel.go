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
	listMode = "list mode"
	itemMode = "item mode"

	itemCreation     = "item create"
	itemModification = "item modification"

	specialType  = "special type"
	sequenceType = "sequence type"
)

type ActionKeycomboEditorModel struct {
	mode              string
	actionList        list.Model
	sequenceInput     textinput.Model
	specialInput      SpecialKeyEditorModel
	selectedItemType  string
	selectedItemIndex int
	selectedItemMode  string
	help              help.Model
}

func NewActionKeycomboEditorModel(title string, itemList []string, specialKeys []string) ActionKeycomboEditorModel {
	var listOfItems []list.Item
	for _, item := range itemList {
		var keyType string
		if strings.HasPrefix(item, "Keys.") {
			keyType = specialType
		} else {
			keyType = sequenceType
		}
		listOfItems = append(listOfItems, common.TypedItem{Value: item, Type: keyType})
	}

	actionList := list.New(listOfItems, common.TypedItemDelegate{}, 1, 1)
	actionList.Title = title
	actionList.SetShowStatusBar(false)
	actionList.SetFilteringEnabled(false)
	actionList.SetShowHelp(false)
	actionList.DisableQuitKeybindings()

	specialInput := NewASpecialKeyEditorModel("Select multiple keys to be pressed", specialKeys)

	sequenceInput := textinput.New()
	sequenceInput.Focus()
	sequenceInput.CharLimit = 64
	sequenceInput.Width = 48

	help := common.GetHelpModel()
	help.ShowAll = true

	return ActionKeycomboEditorModel{
		mode:          listMode,
		actionList:    actionList,
		sequenceInput: sequenceInput,
		specialInput:  specialInput,
		help:          help,
	}
}

func (m *ActionKeycomboEditorModel) SetValue(itemList []string) {
	var listOfItems []list.Item
	for _, item := range itemList {
		var keyType string
		if strings.HasPrefix(item, "Keys.") {
			keyType = specialType
		} else {
			keyType = sequenceType
		}
		listOfItems = append(listOfItems, common.TypedItem{Value: item, Type: keyType})
	}
	m.actionList.SetItems(listOfItems)
}

func (m *ActionKeycomboEditorModel) EdittingMode() bool {
	return m.mode == itemMode
}

func (m *ActionKeycomboEditorModel) SetSize(width int, height int) {
	m.actionList.SetWidth(width / 2)
	m.actionList.SetHeight(height)
	m.specialInput.SetSize(width, height)
	m.help.Width = width / 2
}

func (m ActionKeycomboEditorModel) Value() []string {
	var itemList []string
	for _, item := range m.actionList.Items() {
		itemList = append(itemList, item.(common.TypedItem).Value)
	}
	return itemList
}

func (m ActionKeycomboEditorModel) Update(msg tea.Msg) (ActionKeycomboEditorModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch m.mode {
	case listMode:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, common.KeycomboListKeys.MoveUp, common.KeycomboListKeys.MoveDown):
				item, ok := m.actionList.SelectedItem().(common.TypedItem)
				if ok {
					m.actionList.RemoveItem(m.actionList.Index())
					if key.Matches(msg, common.ListKeys.MoveUp) {
						cmd = m.actionList.InsertItem(m.actionList.Index()-1, item)
						m.actionList.CursorUp()
					} else {
						cmd = m.actionList.InsertItem(m.actionList.Index()+1, item)
						m.actionList.CursorDown()
					}
				}
				cmds = append(cmds, cmd)

			case key.Matches(msg, common.KeycomboListKeys.RemoveItem):
				_, ok := m.actionList.SelectedItem().(common.TypedItem)
				if ok {
					m.actionList.RemoveItem(m.actionList.Index())
				}

			case key.Matches(msg, common.KeycomboListKeys.EditItem):
				item, ok := m.actionList.SelectedItem().(common.TypedItem)
				if ok {
					m.selectedItemIndex = m.actionList.Index()
					m.selectedItemMode = itemModification
					m.mode = itemMode
					switch item.Type {
					case specialType:
						m.selectedItemType = specialType
						cmd = m.specialInput.SetValue(item.Value)
						cmds = append(cmds, cmd)
					case sequenceType:
						m.selectedItemType = sequenceType
						m.sequenceInput.SetValue(item.Value)
						m.sequenceInput.SetCursor(len(item.Value))
					}
				}

			case key.Matches(msg, common.KeycomboListKeys.AddSpecialItem, common.KeycomboListKeys.AddSequenceItem):
				m.selectedItemMode = itemCreation
				m.mode = itemMode
				if key.Matches(msg, common.KeycomboListKeys.AddSpecialItem) {
					m.selectedItemType = specialType
					m.specialInput.SetValue("")
				} else {
					m.selectedItemType = sequenceType
					m.sequenceInput.SetValue("")
				}
			}
		}
		m.actionList, cmd = m.actionList.Update(msg)
		cmds = append(cmds, cmd)

	case itemMode:
		if !(m.selectedItemType == specialType && m.specialInput.EdittingMode()) {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch {
				case key.Matches(msg, common.ConfirmKeys.ApplyData):
					m.mode = listMode
					var newValue string
					switch m.selectedItemType {
					case sequenceType:
						newValue = m.sequenceInput.Value()
					case specialType:
						newValue = m.specialInput.Value()
					}
					switch m.selectedItemMode {
					case itemCreation:
						newIndex := len(m.actionList.Items())
						cmd = m.actionList.InsertItem(newIndex, common.TypedItem{Value: newValue, Type: m.selectedItemType})
						m.actionList.Select(newIndex)
					case itemModification:
						cmd = m.actionList.SetItem(m.selectedItemIndex, common.TypedItem{Value: newValue, Type: m.selectedItemType})
					}
					cmds = append(cmds, cmd)
				case key.Matches(msg, common.ConfirmKeys.ExitMode):
					m.mode = listMode
				}
			}
		}
		switch m.selectedItemType {
		case sequenceType:
			m.sequenceInput, cmd = m.sequenceInput.Update(msg)
			cmds = append(cmds, cmd)
		case specialType:
			m.specialInput, cmd = m.specialInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m ActionKeycomboEditorModel) View() string {
	var str string
	switch m.mode {
	case listMode:
		str = lipgloss.JoinHorizontal(
			lipgloss.Center,
			m.actionList.View(),
			lipgloss.PlaceHorizontal(m.help.Width, lipgloss.Right, m.help.View(common.KeycomboListKeys)))

	case itemMode:
		switch m.selectedItemType {
		case sequenceType:
			str = fmt.Sprintf("Enter key sequence:\n%s", m.sequenceInput.View())
		case specialType:
			str = m.specialInput.View()
		}
	}
	return str
}
