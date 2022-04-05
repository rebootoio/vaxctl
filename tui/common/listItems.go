package common

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type SelectableItem struct {
	Name     string
	Selected bool
}

func (i SelectableItem) FilterValue() string { return i.Name }

type SelectableItemDelegate struct{}

func (d SelectableItemDelegate) Height() int                               { return 1 }
func (d SelectableItemDelegate) Spacing() int                              { return 0 }
func (d SelectableItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d SelectableItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(SelectableItem)
	if !ok {
		return
	}
	var str string
	if i.Selected {
		str = fmt.Sprintf("[x] %s", i.Name)
	} else {
		str = fmt.Sprintf("[ ] %s", i.Name)
	}

	fn := ItemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return SelectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

type SimpleItem string

func (i SimpleItem) FilterValue() string { return string(i) }

type SimpleItemDelegate struct{}

func (d SimpleItemDelegate) Height() int                               { return 1 }
func (d SimpleItemDelegate) Spacing() int                              { return 0 }
func (d SimpleItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d SimpleItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(SimpleItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i)

	fn := ItemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return SelectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

type TypedItem struct {
	Type  string
	Value string
}

func (i TypedItem) FilterValue() string { return i.Value }

type TypedItemDelegate struct{}

func (d TypedItemDelegate) Height() int                               { return 1 }
func (d TypedItemDelegate) Spacing() int                              { return 0 }
func (d TypedItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d TypedItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(TypedItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.Value)

	fn := ItemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return SelectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

type KeyValueItem struct {
	Key   string
	Value string
}

func (i KeyValueItem) FilterValue() string { return "" }

type KeyValueItemDelegate struct{}

func (d KeyValueItemDelegate) Height() int                               { return 1 }
func (d KeyValueItemDelegate) Spacing() int                              { return 0 }
func (d KeyValueItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d KeyValueItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(KeyValueItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s = %s", i.Key, i.Value)

	fn := ItemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return SelectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}
