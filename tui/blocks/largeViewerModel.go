package blocks

import (
	"fmt"
	"strings"
	"vaxctl/tui/common"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LargeViewerModel struct {
	viewport viewport.Model
	width    int
	height   int
	percent  float64
	active   bool
}

func NewLargeViewerModel() LargeViewerModel {
	viewport := viewport.New(1, 1)
	return LargeViewerModel{
		viewport: viewport,
	}
}

func (m *LargeViewerModel) SetContent(content string) {
	m.viewport.SetContent(content)
}

func (m *LargeViewerModel) SetActive(active bool) {
	m.active = active
}

func (m *LargeViewerModel) SetSize(width int, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = m.width
	m.viewport.Height = m.height - lipgloss.Height(common.ActiveBorderColorStyle.Render("═")) - lipgloss.Height(common.BorderredStyle.Render("100%"))
}

func (m LargeViewerModel) Update(msg tea.Msg) (LargeViewerModel, tea.Cmd) {
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
	m.percent = m.viewport.ScrollPercent() * 100
	return m, cmd
}

func (m LargeViewerModel) View() string {
	var borderChar string
	var borderStyle lipgloss.Style
	if m.active {
		borderChar = "═"
		borderStyle = common.ActiveBorderColorStyle
	} else {
		borderChar = "─"
		borderStyle = common.BorderColorStyle
	}
	header := borderStyle.Render(strings.Repeat(borderChar, max(0, m.width)))
	percent := common.BorderredStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	bottomLine := borderStyle.Render(strings.Repeat(borderChar, max(0, m.viewport.Width-lipgloss.Width(percent))))
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Center, bottomLine, percent)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		m.viewport.View(),
		bottomRow,
	)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
