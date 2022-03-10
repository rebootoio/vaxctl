package blocks

import (
	"fmt"
	"strconv"
	"strings"
	"vaxctl/tui/common"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ResourceDataModel struct {
	title               string
	items               []ResourceItem
	availableSelections []string
	currentSelected     string
	selectionIndex      int
	selectionStart      int
	width               int
	height              int
	help                help.Model
}

type ResourceItem struct {
	Title      string
	Value      string
	Selectable bool
	Toggleble  bool
}

func NewResourceDataModel(title string, initialSelection int, selectionStart int, items []ResourceItem) ResourceDataModel {
	help := common.GetHelpModel()
	help.ShowAll = true
	var availableSelections []string
	for _, item := range items {
		if item.Selectable || item.Toggleble {
			availableSelections = append(availableSelections, item.Title)
		}
	}
	var currentSelected string
	if len(availableSelections) == 0 {
		currentSelected = ""
	} else {
		currentSelected = availableSelections[initialSelection]
	}
	return ResourceDataModel{
		title:               title,
		items:               items,
		availableSelections: availableSelections,
		currentSelected:     currentSelected,
		selectionIndex:      initialSelection,
		selectionStart:      selectionStart,
		help:                help,
	}
}

func (m *ResourceDataModel) SetValue(title string, value string) {
	var idx int
	for i, item := range m.items {
		if item.Title == title {
			idx = i
			break
		}
	}
	m.items[idx].Value = value
}

func (m *ResourceDataModel) SetSize(width int, height int) {
	m.width = width
	m.height = height
}

func (m ResourceDataModel) Update(msg tea.Msg) (ResourceDataModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.ResourceDataKeys.Edit):
			currentItem := m.items[m.selectionStart+m.selectionIndex]
			if currentItem.Selectable {
				cmd = common.SetDynamicView(m.currentSelected)
			}

		case key.Matches(msg, common.ResourceDataKeys.Up, common.ResourceDataKeys.Down):
			if key.Matches(msg, common.ResourceDataKeys.Up) {
				m.selectionIndex--
			} else {
				m.selectionIndex++
			}
			if m.selectionIndex == len(m.availableSelections) {
				m.selectionIndex = 0
			} else if m.selectionIndex < 0 {
				m.selectionIndex = len(m.availableSelections) - 1
			}
			m.currentSelected = m.availableSelections[m.selectionIndex]

		case key.Matches(msg, common.ResourceDataKeys.Toggle):
			currentItem := m.items[m.selectionStart+m.selectionIndex]
			if currentItem.Toggleble {
				currentValue, err := strconv.ParseBool(currentItem.Value)
				if err == nil {
					m.SetValue(currentItem.Title, strconv.FormatBool(!currentValue))
					cmd = common.UpdateToggleValue(currentItem.Title, !currentValue)
				}
			}

		case key.Matches(msg, common.ResourceDataKeys.PageUp, common.ResourceDataKeys.PageDown):
			if key.Matches(msg, common.ResourceDataKeys.PageUp) {
				m.selectionIndex = 0
			} else {
				m.selectionIndex = len(m.availableSelections) - 1
			}
			m.currentSelected = m.availableSelections[m.selectionIndex]
		}
	}
	return m, cmd
}

func (m ResourceDataModel) View() string {
	lines := []string{}

	var spacingIndex int
	for _, item := range m.items {
		l := len(item.Title)
		if l > spacingIndex {
			spacingIndex = l
		}
	}

	for _, item := range m.items {
		spaces := strings.Repeat(" ", spacingIndex-len(item.Title))
		itemText := fmt.Sprintf("%s%s: ", item.Title, spaces)
		if item.Value == "" {
			itemText += common.NotSetString
		} else {
			itemText += item.Value
		}
		if len(itemText) > m.width*3/4 {
			itemText = itemText[:m.width*3/4] + "..."
		}
		if m.currentSelected == item.Title {
			if item.Selectable {
				lines = append(lines, common.ResourceSelectedItemStyle.Render("> "+itemText))
			} else {
				lines = append(lines, common.ResourceSelectedItemStyle.Render("- "+itemText))
			}
		} else {
			lines = append(lines, common.ResourceItemStyle.Render("  "+itemText))
		}
	}

	data := lipgloss.JoinVertical(
		lipgloss.Left,
		lines...,
	)

	m.help.Width = m.width

	return lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, common.ResourceTitleStyle.Render(m.title)),
		data,
		"\n",
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, "-------"),
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.help.View(common.ResourceDataKeys)),
	)
}
