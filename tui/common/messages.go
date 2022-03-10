package common

import tea "github.com/charmbracelet/bubbletea"

type ExitDynamicViewMsg struct {
	Save bool
}

func ExitDynamicView(save bool) tea.Cmd {
	return func() tea.Msg {
		return ExitDynamicViewMsg{save}
	}
}

type ApplyMainActionMsg struct {
	Action string
}

func ApplyMainAction(action string) tea.Cmd {
	return func() tea.Msg {
		return ApplyMainActionMsg{action}
	}
}

type SetDynamicViewMsg struct {
	View string
}

func SetDynamicView(view string) tea.Cmd {
	return func() tea.Msg {
		return SetDynamicViewMsg{view}
	}
}

type UpdateToggleValueMsg struct {
	Name  string
	Value bool
}

func UpdateToggleValue(name string, value bool) tea.Cmd {
	return func() tea.Msg {
		return UpdateToggleValueMsg{Name: name, Value: value}
	}
}

type BackToNavigationMsg struct{}

func BackToNavigation() tea.Cmd {
	return func() tea.Msg {
		return BackToNavigationMsg{}
	}
}

type EditItemMsg struct{}

func EditItem() tea.Cmd {
	return func() tea.Msg {
		return EditItemMsg{}
	}
}

type CreateRuleMsg struct {
	StateId string
}

func CreateRule(stateId string) tea.Cmd {
	return func() tea.Msg {
		return CreateRuleMsg{StateId: stateId}
	}
}

type SetAsDefaultMsg struct{}

func SetAsDefault() tea.Cmd {
	return func() tea.Msg {
		return SetAsDefaultMsg{}
	}
}

type UpdateCredNamesMsg struct{}

func UpdateCredNames() tea.Cmd {
	return func() tea.Msg {
		return UpdateCredNamesMsg{}
	}
}

type UpdateActionNamesMsg struct{}

func UpdateActionNames() tea.Cmd {
	return func() tea.Msg {
		return UpdateActionNamesMsg{}
	}
}

type RefreshDataMsg struct{}

func RefreshData() tea.Cmd {
	return func() tea.Msg {
		return RefreshDataMsg{}
	}
}

type ApplyTableChangeMsg struct {
	Save bool
}

func ApplyTableChange(save bool) tea.Cmd {
	return func() tea.Msg {
		return ApplyTableChangeMsg{Save: save}
	}
}

type ChangeOrderMsg struct {
	Up bool
}

func ChangeOrder(up bool) tea.Cmd {
	return func() tea.Msg {
		return ChangeOrderMsg{Up: up}
	}
}
