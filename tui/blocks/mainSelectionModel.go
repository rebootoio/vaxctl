package blocks

import (
	"vaxctl/tui/common"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type MainSelectionModel struct {
	li    list.Model
	width int
}

func NewMainSelectionModel(title string, itemList []string, selectedItem string) MainSelectionModel {
	var listOfItems []list.Item
	var selectedIndex int
	for idx, item := range itemList {
		listOfItems = append(listOfItems, common.SimpleItem(item))
		if item == selectedItem {
			selectedIndex = idx
		}
	}

	li := list.New(listOfItems, common.SimpleItemDelegate{}, 1, 1)
	li.Title = title
	li.SetShowStatusBar(false)
	li.SetFilteringEnabled(false)
	li.SetShowHelp(false)
	li.DisableQuitKeybindings()
	li.SetShowPagination(false)
	li.Select(selectedIndex)

	return MainSelectionModel{
		li: li,
	}
}

func (m *MainSelectionModel) SetSize(width int, height int) {
	m.width = width
	m.li.SetWidth(width)
	m.li.SetHeight(height)
}

func (m MainSelectionModel) Value() string {
	item, _ := m.li.SelectedItem().(common.SimpleItem)
	return string(item)
}

func (m MainSelectionModel) Update(msg tea.Msg) (MainSelectionModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.ConfirmKeys.ApplyData):
			selected, ok := m.li.SelectedItem().(common.SimpleItem)
			if ok {
				cmd = common.ApplyMainAction(string(selected))
				return m, cmd
			}
		case key.Matches(msg, common.FooterKeys.PageDown, common.FooterKeys.PageUp):
			if m.li.Paginator.TotalPages == 1 {
				if key.Matches(msg, common.FooterKeys.PageDown) {
					m.li.Select(len(m.li.Items()) - 1)
				} else {
					m.li.Select(0)
				}
			}
		}
	}
	m.li, cmd = m.li.Update(msg)
	return m, cmd
}

func (m MainSelectionModel) View() string {
	return m.li.View()
}
