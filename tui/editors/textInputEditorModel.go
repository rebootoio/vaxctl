package editors

import (
	"fmt"
	"vaxctl/tui/common"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInputEditorModel struct {
	title    string
	ti       textinput.Model
	help     help.Model
	validate func(string) bool
}

func NewTextInputEditor(title string, initialValue string) TextInputEditorModel {
	ti := textinput.New()
	ti.SetValue(initialValue)
	ti.CharLimit = 64
	ti.Width = 48

	help := common.GetHelpModel()

	return TextInputEditorModel{
		title: title,
		ti:    ti,
		help:  help,
	}
}

func (m *TextInputEditorModel) SetValidationFunc(validate func(string) bool) {
	m.validate = validate
}

func (m *TextInputEditorModel) SetValue(value string) {
	m.ti.SetValue(value)
	m.ti.SetCursor(len(value))
}

func (m *TextInputEditorModel) SetEchoMode(echoMode textinput.EchoMode) {
	m.ti.EchoMode = echoMode
}

func (m TextInputEditorModel) Value() string {
	return m.ti.Value()
}

func (m *TextInputEditorModel) Focus() tea.Cmd {
	return m.ti.Focus()
}

func (m *TextInputEditorModel) Blur() {
	m.ti.Blur()
}

func (m TextInputEditorModel) Update(msg tea.Msg) (TextInputEditorModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.ConfirmKeys.ApplyData):
			if m.ti.Value() != "" && (m.validate == nil || m.validate(m.ti.Value())) {
				cmd = common.ExitDynamicView(true)
				cmds = append(cmds, cmd)
			}
		case key.Matches(msg, common.ConfirmKeys.ExitMode):
			cmd = common.ExitDynamicView(false)
			cmds = append(cmds, cmd)
		}
	}
	m.ti, cmd = m.ti.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m TextInputEditorModel) View() string {
	helpText := m.help.View(common.ConfirmKeys)
	var title, validationStatus string
	if m.validate != nil {
		if m.validate(m.ti.Value()) {
			validationStatus = common.ValidStyle.Render("✔")
		} else {
			validationStatus = common.InvalidStyle.Render("✘")
		}
		title = fmt.Sprintf("%s %s", m.title, validationStatus)
	} else {
		title = m.title
	}
	return fmt.Sprintf("%s\n%s\n\n%s", title, m.ti.View(), helpText)
}
