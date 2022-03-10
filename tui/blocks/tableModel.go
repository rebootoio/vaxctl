package blocks

import (
	"fmt"
	"vaxctl/tui/common"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

type TableColumn struct {
	Title string
	Size  int
}

type TableModel struct {
	table          table.Model
	columns        []TableColumn
	rows           []table.Row
	editRuleName   string
	totalUnits     int
	width          int
	help           help.Model
	additionalKeys help.KeyMap
	ManualChanges  bool
}

func NewTableModel(columns []TableColumn, tableRows []table.Row) TableModel {
	var totalUnits int
	for _, column := range columns {
		totalUnits += column.Size
	}

	help := common.GetHelpModel()

	return TableModel{
		columns:    columns,
		rows:       tableRows,
		totalUnits: totalUnits,
		help:       help,
	}
}

func (m *TableModel) SetAdditionalKeys(additionalKeys help.KeyMap) {
	m.additionalKeys = additionalKeys
}

func (m *TableModel) SetHighlightedRow(index int) {
	m.table = m.table.WithHighlightedRow(index)
}

func (m *TableModel) SetRows(rows []table.Row) {
	m.rows = rows
	m.table = m.table.WithRows(m.rows)
}

func (m TableModel) GetRows() []table.Row {
	return m.rows
}
func (m TableModel) GetCurrentRow() table.Row {
	return m.table.HighlightedRow()
}

func (m *TableModel) updateFooter() {
	var manualChangesStr string
	if m.ManualChanges {
		manualChangesStr = "The table currently shows unsaved changes (Enter to apply, Esc to cancel)"
	}
	footerText := fmt.Sprintf(
		"%s        Page %d/%d",
		common.OcrMatchedTextStyle.Render(manualChangesStr),
		m.table.CurrentPage(),
		m.table.MaxPages(),
	)
	m.table = m.table.WithStaticFooter(footerText)
}

func (m *TableModel) SetSize(width int, height int) {
	columnWidth := (width - 10) / m.totalUnits

	var tableColumns []table.Column
	for _, column := range m.columns {
		newTableColumn := table.NewColumn(column.Title, column.Title, columnWidth*column.Size).WithStyle(common.TableColumn)
		tableColumns = append(tableColumns, newTableColumn)
	}
	m.table = table.New(tableColumns).Focused(true).WithRows(m.rows).WithPageSize(height - 10)
	m.updateFooter()
	m.width = width - 10
}

func (m TableModel) Update(msg tea.Msg) (TableModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.TableKeys.ChooseRow):
			if m.ManualChanges {
				return m, common.ApplyTableChange(true)
			} else {
				if len(m.rows) > 0 {
					return m, common.EditItem()
				}
			}
		case key.Matches(msg, common.TableKeys.RefreshData):
			return m, common.RefreshData()
		case key.Matches(msg, common.CredTableKeys.SetAsDefault):
			return m, common.SetAsDefault()
		case key.Matches(msg, common.RuleTableKeys.MoveUp, common.RuleTableKeys.MoveDown):
			return m, common.ChangeOrder(key.Matches(msg, common.RuleTableKeys.MoveUp))
		case key.Matches(msg, common.RuleTableKeys.CancelOrder):
			return m, common.ApplyTableChange(false)
		}
	}
	m.table, cmd = m.table.Update(msg)
	m.updateFooter()
	return m, cmd
}

func (m TableModel) View() string {
	var helpStr string
	if m.additionalKeys != nil {
		helpStr = m.help.View(m.additionalKeys)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.table.View(),
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.help.View(common.TableKeys)),
		helpStr,
	)
}
