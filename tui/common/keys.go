package common

import "github.com/charmbracelet/bubbles/key"

type footerKeyMap struct {
	NextView      key.Binding
	PreviousView  key.Binding
	Quit          key.Binding
	ForceQuit     key.Binding
	Suspend       key.Binding
	SwitchView    key.Binding
	GoBack        key.Binding
	PageUp        key.Binding
	PageDown      key.Binding
	SwitchToOne   key.Binding
	SwitchToTwo   key.Binding
	SwitchToThree key.Binding
	SwitchToFour  key.Binding
}

func (k footerKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NextView, k.PreviousView, k.SwitchView, k.PageUp, k.PageDown, k.GoBack, k.Quit, k.ForceQuit}
}
func (k footerKeyMap) FullHelp() [][]key.Binding { return nil }

type credFooterKeyMap struct {
	ShowPasswords key.Binding
}

func (k credFooterKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.ShowPasswords}
}
func (k credFooterKeyMap) FullHelp() [][]key.Binding { return nil }

type confirmKeyMap struct {
	ApplyData key.Binding
	ExitMode  key.Binding
}

func (k confirmKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.ApplyData, k.ExitMode}
}
func (k confirmKeyMap) FullHelp() [][]key.Binding { return nil }

type orderedListKeyMap struct {
	MoveUp     key.Binding
	MoveDown   key.Binding
	SelectItem key.Binding
}

func (k orderedListKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.SelectItem, k.MoveUp, k.MoveDown}
}
func (k orderedListKeyMap) FullHelp() [][]key.Binding { return nil }

type filteredListKeyMap struct {
	SelectItem key.Binding
}

func (k filteredListKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.SelectItem}
}
func (k filteredListKeyMap) FullHelp() [][]key.Binding { return nil }

type resourceDataKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Edit     key.Binding
	Toggle   key.Binding
	PageUp   key.Binding
	PageDown key.Binding
}

func (k resourceDataKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Edit, k.Toggle}
}
func (k resourceDataKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Edit, k.Toggle},
	}
}

type numberKeyMap struct {
	Number key.Binding
}

type keycomboListKeyMap struct {
	MoveUp          key.Binding
	MoveDown        key.Binding
	EditItem        key.Binding
	RemoveItem      key.Binding
	AddSpecialItem  key.Binding
	AddSequenceItem key.Binding
}

func (k keycomboListKeyMap) ShortHelp() []key.Binding { return nil }
func (k keycomboListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.EditItem, k.AddSpecialItem, k.AddSequenceItem, k.RemoveItem},
		{k.MoveUp, k.MoveDown},
	}
}

type tableKeyMap struct {
	RowDown key.Binding
	RowUp   key.Binding

	PageDown key.Binding
	PageUp   key.Binding

	ChooseRow  key.Binding
	ReorderRow key.Binding

	RefreshData key.Binding
}

func (k tableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.RowDown,
		k.RowUp,
		k.PageDown,
		k.PageUp,
		k.ChooseRow,
		k.RefreshData,
	}
}
func (k tableKeyMap) FullHelp() [][]key.Binding { return nil }

type credTableKeyMap struct {
	SetAsDefault  key.Binding
	ApplyDefault  key.Binding
	CancelDefault key.Binding
}

func (k credTableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.SetAsDefault, k.ApplyDefault, k.CancelDefault}
}
func (k credTableKeyMap) FullHelp() [][]key.Binding { return nil }

type ruleTableKeyMap struct {
	MoveUp      key.Binding
	MoveDown    key.Binding
	ApplyOrder  key.Binding
	CancelOrder key.Binding
}

func (k ruleTableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.MoveUp, k.MoveDown, k.ApplyOrder, k.CancelOrder}
}
func (k ruleTableKeyMap) FullHelp() [][]key.Binding { return nil }

var (
	FooterKeys = footerKeyMap{
		NextView: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("Tab", "next view"),
		),
		PreviousView: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("Shift+Tab", "previous view"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		ForceQuit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("Ctrl+c", "force quit"),
		),
		Suspend: key.NewBinding(
			key.WithKeys("ctrl+z"),
			key.WithHelp("Ctrl+z", "suspend"),
		),
		SwitchView: key.NewBinding(
			key.WithHelp("Shift+1/2/3/4", "switch to view"),
		),
		GoBack: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "go back one level"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown"),
			key.WithHelp("pgdown", "page down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("pgup", "page up"),
		),
		SwitchToOne: key.NewBinding(
			key.WithKeys("!"),
			key.WithHelp("shift+1", "Switch to 1"),
		),
		SwitchToTwo: key.NewBinding(
			key.WithKeys("@"),
			key.WithHelp("shift+2", "Switch to 2"),
		),
		SwitchToThree: key.NewBinding(
			key.WithKeys("#"),
			key.WithHelp("shift+3", "Switch to 3"),
		),
		SwitchToFour: key.NewBinding(
			key.WithKeys("$"),
			key.WithHelp("shift+4", "Switch to 4"),
		),
	}
	CredFooterKeys = credFooterKeyMap{
		ShowPasswords: key.NewBinding(
			key.WithKeys("ctrl+p"),
			key.WithHelp("ctrl+p", "show/hide passwords"),
		),
	}
	ConfirmKeys = confirmKeyMap{
		ExitMode: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("Esc", "exit"),
		),
		ApplyData: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("Enter", "apply"),
		),
	}
	ListKeys = orderedListKeyMap{
		MoveUp: key.NewBinding(
			key.WithKeys("+"),
			key.WithHelp("+", "move up"),
		),
		MoveDown: key.NewBinding(
			key.WithKeys("-"),
			key.WithHelp("-", "move down"),
		),
		SelectItem: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "select"),
		),
	}
	FilteredListKeys = filteredListKeyMap{
		SelectItem: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "select"),
		),
	}
	ResourceDataKeys = resourceDataKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		Edit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("Enter", "edit value"),
		),
		Toggle: key.NewBinding(
			key.WithKeys(" ", "left", "right"),
			key.WithHelp("space/left/right", "toggle value"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown"),
			key.WithHelp("pgdown", "page down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("pgup", "page up"),
		),
	}
	NumberKeys = numberKeyMap{
		Number: key.NewBinding(
			key.WithKeys("0", "1", "2", "3", "4", "5", "6", "7", "8", "9"),
		),
	}
	KeycomboListKeys = keycomboListKeyMap{
		MoveUp: key.NewBinding(
			key.WithKeys("+"),
			key.WithHelp("+", "move up"),
		),
		MoveDown: key.NewBinding(
			key.WithKeys("-"),
			key.WithHelp("-", "move down"),
		),
		EditItem: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit item"),
		),
		AddSpecialItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add key combo"),
		),
		AddSequenceItem: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "add string"),
		),
		RemoveItem: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "remove item"),
		),
	}
	SpecialKeyListKeys = keycomboListKeyMap{
		AddSpecialItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add special key"),
		),
		AddSequenceItem: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "add character"),
		),
		RemoveItem: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "remove item"),
		),
	}
	TableKeys = tableKeyMap{
		RowDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("down", "row down"),
		),
		RowUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("up", "row up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("right", "pgdown"),
			key.WithHelp("right,pgdown", "page down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("left", "pgup"),
			key.WithHelp("left,pgup", "page up"),
		),
		ChooseRow: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose item"),
		),
		RefreshData: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "refresh table data"),
		),
	}
	CredTableKeys = credTableKeyMap{
		SetAsDefault: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "set as default"),
		),
		CancelDefault: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel default change"),
		),
		ApplyDefault: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "apply default change"),
		),
	}
	RuleTableKeys = ruleTableKeyMap{
		MoveUp: key.NewBinding(
			key.WithKeys("+"),
			key.WithHelp("+", "move up"),
		),
		MoveDown: key.NewBinding(
			key.WithKeys("-"),
			key.WithHelp("-", "move down"),
		),
		CancelOrder: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel reorder"),
		),
		ApplyOrder: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "apply reorder"),
		),
	}
)
