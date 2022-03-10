package editors

import (
	"vaxctl/tui/common"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SimpleListEditorModel struct {
	li   list.Model
	help help.Model
}

func NewSimpleListEditor(title string, itemList []string, selectedItem string) SimpleListEditorModel {
	var listOfItems []list.Item
	for _, item := range itemList {
		listOfItems = append(listOfItems, common.SimpleItem(item))
	}

	li := list.New(listOfItems, common.SimpleItemDelegate{}, 1, 1)
	li.Title = title
	li.SetShowStatusBar(false)
	li.SetFilteringEnabled(false)
	li.SetShowHelp(false)
	li.DisableQuitKeybindings()

	help := common.GetHelpModel()

	return SimpleListEditorModel{
		li:   li,
		help: help,
	}
}

func (m *SimpleListEditorModel) SetValue(value string) {
	for idx, item := range m.li.Items() {
		if string(item.(common.SimpleItem)) == value {
			m.li.Select(idx)
			break
		}
	}
}

func (m *SimpleListEditorModel) UpdateItemList(itemList []string) {
	var listOfItems []list.Item
	for _, item := range itemList {
		listOfItems = append(listOfItems, common.SimpleItem(item))
	}
	m.li.SetItems(listOfItems)
}

func (m SimpleListEditorModel) Value() string {
	item, _ := m.li.SelectedItem().(common.SimpleItem)
	return string(item)
}

func (m *SimpleListEditorModel) SetSize(width int, height int) {
	m.li.SetSize(width/2, height)
}

func (m SimpleListEditorModel) Update(msg tea.Msg) (SimpleListEditorModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.ConfirmKeys.ApplyData):
			cmd = common.ExitDynamicView(true)
			cmds = append(cmds, cmd)
		case key.Matches(msg, common.ConfirmKeys.ExitMode):
			cmd = common.ExitDynamicView(false)
			cmds = append(cmds, cmd)
		case key.Matches(msg, common.FooterKeys.PageDown, common.FooterKeys.PageUp):
			if m.li.Paginator.TotalPages == 1 {
				if key.Matches(msg, common.FooterKeys.PageDown) {
					m.li.Select(len(m.li.Items()))
				} else {
					m.li.Select(0)
				}
			}
		}
	}
	m.li, cmd = m.li.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m SimpleListEditorModel) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		m.li.View(),
		m.help.View(common.ConfirmKeys))
}
