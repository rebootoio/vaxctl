package editors

import (
	"vaxctl/tui/common"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type OrderedSelectionEditorModel struct {
	li             list.Model
	selectedCount  int
	help           help.Model
	chosenItemList []string
}

func NewOrderedSelectionEditor(title string, itemList []string, chosenItemList []string) OrderedSelectionEditorModel {
	var listOfItems []list.Item
	for i := range chosenItemList {
		listOfItems = append(listOfItems, common.SelectableItem{Name: chosenItemList[i], Selected: true})
	}
	for _, item := range itemList {
		var found bool
		for i := range chosenItemList {
			if chosenItemList[i] == item {
				found = true
			}
		}
		if !found {
			listOfItems = append(listOfItems, common.SelectableItem{Name: item, Selected: false})
		}
	}

	li := list.New(listOfItems, common.SelectableItemDelegate{}, 1, 1)
	li.Title = title
	li.SetShowStatusBar(false)
	li.SetFilteringEnabled(true)
	li.SetShowHelp(false)
	li.DisableQuitKeybindings()
	li.Help = common.GetHelpModel()

	help := common.GetHelpModel()

	return OrderedSelectionEditorModel{
		li:             li,
		selectedCount:  len(chosenItemList),
		help:           help,
		chosenItemList: chosenItemList,
	}
}

func (m *OrderedSelectionEditorModel) ClearSelected() {
	for idx, item := range m.li.Items() {
		selectableItem := item.(common.SelectableItem)
		if selectableItem.Selected {
			selectableItem.Selected = false
			m.li.SetItem(idx, selectableItem)
		} else {
			break
		}
	}
	m.chosenItemList = []string{}
}

func (m *OrderedSelectionEditorModel) SetSelected(selectedItems []string) {
	var listOfItems []list.Item
	for i := range selectedItems {
		listOfItems = append(listOfItems, common.SelectableItem{Name: selectedItems[i], Selected: true})
	}
	for _, item := range m.li.Items() {
		var found bool
		itemName := item.(common.SelectableItem).Name
		for i := range selectedItems {
			if selectedItems[i] == itemName {
				found = true
				break
			}
		}
		if !found {
			listOfItems = append(listOfItems, common.SelectableItem{Name: itemName, Selected: false})
		}
	}
	m.li.SetItems(listOfItems)
	m.chosenItemList = selectedItems
	m.selectedCount = len(selectedItems)
}

func (m *OrderedSelectionEditorModel) UpdateItemList(itemList []string) {
	var listOfItems []list.Item
	for i := range m.chosenItemList {
		listOfItems = append(listOfItems, common.SelectableItem{Name: m.chosenItemList[i], Selected: true})
	}
	for _, item := range itemList {
		var found bool
		for _, chosenItem := range m.chosenItemList {
			if chosenItem == item {
				found = true
				break
			}
		}
		if !found {
			listOfItems = append(listOfItems, common.SelectableItem{Name: item, Selected: false})
		}
	}
	m.li.SetItems(listOfItems)
}

func (m *OrderedSelectionEditorModel) SetSize(width int, height int) {
	m.li.SetWidth(width / 2)
	m.li.SetHeight(height)
}

func (m *OrderedSelectionEditorModel) FilterState() list.FilterState {
	return m.li.FilterState()
}

func (m OrderedSelectionEditorModel) Value() []string {
	var chosenList []string
	allItems := m.li.Items()
	for i := range m.li.Items() {
		item := allItems[i].(common.SelectableItem)
		if item.Selected {
			chosenList = append(chosenList, item.Name)
		} else {
			break
		}
	}
	return chosenList
}

func (m OrderedSelectionEditorModel) Update(msg tea.Msg) (OrderedSelectionEditorModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.ConfirmKeys.ApplyData):
			if m.li.FilterState() == list.Unfiltered {
				cmd = common.ExitDynamicView(true)
				cmds = append(cmds, cmd)
			}

		case key.Matches(msg, common.ConfirmKeys.ExitMode):
			if m.li.FilterState() == list.Unfiltered {
				cmd = common.ExitDynamicView(false)
				cmds = append(cmds, cmd)
			}

		case key.Matches(msg, common.ListKeys.MoveUp, common.ListKeys.MoveDown):
			item, ok := m.li.SelectedItem().(common.SelectableItem)
			if ok && m.li.FilterState() == list.Unfiltered {
				m.li.RemoveItem(m.li.Index())
				if key.Matches(msg, common.ListKeys.MoveUp) {
					cmd = m.li.InsertItem(m.li.Index()-1, item)
					m.li.CursorUp()
				} else {
					cmd = m.li.InsertItem(m.li.Index()+1, item)
					m.li.CursorDown()
				}
			}
			cmds = append(cmds, cmd)

		case key.Matches(msg, common.ListKeys.SelectItem):
			item, ok := m.li.SelectedItem().(common.SelectableItem)
			if ok && m.li.FilterState() != list.Filtering {
				var indexToRemove int
				if m.li.FilterState() == list.Unfiltered {
					indexToRemove = m.li.Index()
				} else {
					items := m.li.Items()
					for i := range items {
						if items[i] == item {
							indexToRemove = i
							break
						}
					}
				}
				item.Selected = !item.Selected
				if item.Selected {
					m.li.RemoveItem(indexToRemove)
					cmd = m.li.InsertItem(m.selectedCount, item)
					if m.li.FilterState() == list.Unfiltered {
						m.li.Select(m.selectedCount)
					}
					m.selectedCount++
					m.chosenItemList = append(m.chosenItemList, item.Name)
				} else {
					m.li.RemoveItem(indexToRemove)
					m.selectedCount--
					cmd = m.li.InsertItem(m.selectedCount, item)
					var chosenItemList []string
					for _, chosenItem := range m.chosenItemList {
						if chosenItem != item.Name {
							chosenItemList = append(chosenItemList, chosenItem)
						}
					}
					m.chosenItemList = chosenItemList
				}
			}
			cmds = append(cmds, cmd)
		}
	}
	m.li, cmd = m.li.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m OrderedSelectionEditorModel) View() string {
	listHelp := m.li.Help.View(m.li)

	var joinedHelp string
	switch m.li.FilterState() {
	case list.Unfiltered:
		joinedHelp = lipgloss.JoinVertical(
			lipgloss.Top,
			m.help.View(common.ListKeys),
			listHelp,
			m.help.View(common.ConfirmKeys))

	case list.Filtering:
		joinedHelp = listHelp

	case list.FilterApplied:
		joinedHelp = lipgloss.JoinVertical(
			lipgloss.Top,
			m.help.View(common.FilteredListKeys),
			listHelp)

	}

	return lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		m.li.View(),
		joinedHelp,
	)
}
