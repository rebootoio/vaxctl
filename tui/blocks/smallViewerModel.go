package blocks

import (
	"vaxctl/tui/common"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wrap"
)

type SmallViewerModel struct {
	viewport viewport.Model
	width    int
	height   int
}

func NewSmallViewerModel() SmallViewerModel {
	viewport := viewport.New(1, 1)
	return SmallViewerModel{
		viewport: viewport,
	}
}

func (m *SmallViewerModel) SetContent(content string) {
	m.viewport.SetContent(wrap.String(content, m.width-4))
}

func (m *SmallViewerModel) SetSize(width int, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = m.width
	m.viewport.Height = m.height
}

func (m SmallViewerModel) Update(msg tea.Msg) (SmallViewerModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(
			msg,
			common.FooterKeys.PageUp,
			common.FooterKeys.PageDown,
			common.ResourceDataKeys.Up,
			common.ResourceDataKeys.Down,
		):
			m.viewport, cmd = m.viewport.Update(msg)
		}
	}
	return m, cmd
}

func (m SmallViewerModel) View() string {
	return m.viewport.View()
}
