package models

import (
	"vaxctl/tui/blocks"
	"vaxctl/tui/common"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	navigationMenu = "Navigation"
	ruleSubMenu    = "Rules"
	actionSubMenu  = "Actions"
	deviceSubMenu  = "Devices"
	stateSubMenu   = "States"
	credSubMenu    = "Creds"
	quitMenu       = "Quit"
)

type NavigationModel struct {
	currentMenu string
	height      int
	width       int
	help        help.Model
	mainList    blocks.MainSelectionModel
	actionModel ActionModel
	ruleModel   RuleModel
	stateModel  StateModel
	deviceModel DeviceModel
	credModel   CredModel
}

func InitialNavigationModel(interactiveData common.InteractiveData) NavigationModel {
	help := common.GetHelpModel()

	mainList := blocks.NewMainSelectionModel(
		"Select Resource",
		[]string{credSubMenu, deviceSubMenu, actionSubMenu, ruleSubMenu, stateSubMenu, quitMenu},
		credSubMenu)

	actionModel := InitialActionModel(interactiveData)
	ruleModel := InitialRuleModel(interactiveData)
	stateModel := InitialStateModel(interactiveData)
	deviceModel := InitialDeviceModel(interactiveData)
	credModel := InitialCredModel(interactiveData)

	var currentMenu string
	if interactiveData.CurrentSubMenu != "" {
		currentMenu = interactiveData.CurrentSubMenu
	} else {
		currentMenu = navigationMenu
	}

	return NavigationModel{
		currentMenu: currentMenu,
		mainList:    mainList,
		help:        help,
		actionModel: actionModel,
		ruleModel:   ruleModel,
		stateModel:  stateModel,
		deviceModel: deviceModel,
		credModel:   credModel,
	}
}

func (m NavigationModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tea.DisableMouse)
}

func (m NavigationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Handle global keys, size changes, custom messages
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, common.FooterKeys.ForceQuit):
			return m, tea.Quit

		case key.Matches(msg, common.FooterKeys.Quit):
			switch m.currentMenu {
			case deviceSubMenu:
				if m.deviceModel.CurrentView != dynamicView || m.deviceModel.DynamicView == "" {
					return m, tea.Quit
				}
			case actionSubMenu:
				if m.actionModel.CurrentView != dynamicView || m.actionModel.DynamicView == "" {
					return m, tea.Quit
				}
			case stateSubMenu:
				if m.stateModel.CurrentView != dynamicView || m.stateModel.DynamicView == "" {
					return m, tea.Quit
				}
			case ruleSubMenu:
				if m.ruleModel.CurrentView != dynamicView || m.ruleModel.DynamicView == "" {
					return m, tea.Quit
				}
			case credSubMenu:
				if m.credModel.CurrentView != dynamicView || m.credModel.DynamicView == "" {
					return m, tea.Quit
				}
			default:
				return m, tea.Quit
			}

		case key.Matches(msg, common.ConfirmKeys.ApplyData):
			if m.currentMenu == navigationMenu {
				m.currentMenu = m.mainList.Value()
				if m.currentMenu == quitMenu {
					return m, tea.Quit
				}
				return m, nil
			}

		case key.Matches(msg, common.FooterKeys.GoBack):
			switch m.currentMenu {
			case deviceSubMenu:
				if m.deviceModel.CurrentView != dynamicView || m.deviceModel.DynamicView == "" {
					m.currentMenu = navigationMenu
					m.deviceModel.CurrentView = viewerView
					m.deviceModel.StatusMessage = ""
				}
			case actionSubMenu:
				if m.actionModel.CurrentView != dynamicView || m.actionModel.DynamicView == "" {
					m.currentMenu = navigationMenu
					m.actionModel.CurrentView = viewerView
					m.actionModel.StatusMessage = ""
				}
			case stateSubMenu:
				if m.stateModel.CurrentView != dynamicView || m.stateModel.DynamicView == "" {
					m.currentMenu = navigationMenu
					m.stateModel.CurrentView = viewerView
					m.stateModel.StatusMessage = ""
				}
			case ruleSubMenu:
				if m.ruleModel.CurrentView != dynamicView || m.ruleModel.DynamicView == "" {
					if m.ruleModel.Editting {
						m.ruleModel.Editting = false
						m.ruleModel.CurrentView = viewerView
						m.ruleModel.StatusMessage = ""
					} else {
						m.currentMenu = navigationMenu
						m.ruleModel.CurrentView = viewerView
						m.ruleModel.StatusMessage = ""
					}
				}
			case credSubMenu:
				if m.credModel.CurrentView != dynamicView || m.credModel.DynamicView == "" {
					m.currentMenu = navigationMenu
					m.credModel.CurrentView = viewerView
					m.credModel.StatusMessage = ""
				}
			}

		case key.Matches(msg, common.FooterKeys.NextView, common.FooterKeys.PreviousView):
			var currentView string
			switch m.currentMenu {
			case deviceSubMenu:
				currentView = m.deviceModel.CurrentView
			case actionSubMenu:
				currentView = m.actionModel.CurrentView
			case stateSubMenu:
				currentView = m.stateModel.CurrentView
			case ruleSubMenu:
				currentView = m.ruleModel.CurrentView
			case credSubMenu:
				currentView = m.credModel.CurrentView
			}

			var currentIndex int
			for idx, view := range views {
				if view == currentView {
					currentIndex = idx
					break
				}
			}
			if key.Matches(msg, common.FooterKeys.NextView) {
				currentIndex++
			} else {
				currentIndex--
			}

			if currentIndex == len(views) {
				currentIndex = 0
			} else if currentIndex < 0 {
				currentIndex = len(views) - 1
			}
			currentView = views[currentIndex]
			switch m.currentMenu {
			case deviceSubMenu:
				m.deviceModel.CurrentView = currentView
				m.deviceModel.StatusMessage = ""
			case actionSubMenu:
				m.actionModel.CurrentView = currentView
				m.actionModel.StatusMessage = ""
			case stateSubMenu:
				m.stateModel.CurrentView = currentView
				m.stateModel.StatusMessage = ""
			case ruleSubMenu:
				m.ruleModel.CurrentView = currentView
				m.ruleModel.StatusMessage = ""
			case credSubMenu:
				m.credModel.CurrentView = currentView
				m.credModel.StatusMessage = ""
			}
			return m, nil

		case key.Matches(msg, common.FooterKeys.SwitchToOne, common.FooterKeys.SwitchToTwo, common.FooterKeys.SwitchToThree, common.FooterKeys.SwitchToFour):
			var currentView string
			if key.Matches(msg, common.FooterKeys.SwitchToOne) {
				currentView = views[0]
			} else if key.Matches(msg, common.FooterKeys.SwitchToTwo) {
				currentView = views[1]
			} else if key.Matches(msg, common.FooterKeys.SwitchToThree) {
				currentView = views[2]
			} else if key.Matches(msg, common.FooterKeys.SwitchToFour) {
				currentView = views[3]
			}
			switch m.currentMenu {
			case deviceSubMenu:
				if m.deviceModel.CurrentView != dynamicView || m.deviceModel.DynamicView == "" {
					m.deviceModel.CurrentView = currentView
					m.deviceModel.StatusMessage = ""
				}
			case actionSubMenu:
				if m.actionModel.CurrentView != dynamicView || m.actionModel.DynamicView == "" {
					m.actionModel.CurrentView = currentView
					m.actionModel.StatusMessage = ""
				}
			case stateSubMenu:
				if m.stateModel.CurrentView != dynamicView || m.stateModel.DynamicView == "" {
					m.stateModel.CurrentView = currentView
					m.stateModel.StatusMessage = ""
				}
			case ruleSubMenu:
				if m.ruleModel.CurrentView != dynamicView || m.ruleModel.DynamicView == "" {
					m.ruleModel.CurrentView = currentView
					m.ruleModel.StatusMessage = ""
				}
			case credSubMenu:
				if m.credModel.CurrentView != dynamicView || m.credModel.DynamicView == "" {
					m.credModel.CurrentView = currentView
					m.credModel.StatusMessage = ""
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		rowHeight, mainWindowWidth, _, _, _ := getSizing(m.width, m.height)
		m.mainList.SetSize(mainWindowWidth, rowHeight)

		m.actionModel, cmd = m.actionModel.Update(msg)
		cmds = append(cmds, cmd)
		m.ruleModel, cmd = m.ruleModel.Update(msg)
		cmds = append(cmds, cmd)
		m.stateModel, cmd = m.stateModel.Update(msg)
		cmds = append(cmds, cmd)
		m.deviceModel, cmd = m.deviceModel.Update(msg)
		cmds = append(cmds, cmd)
		m.credModel, cmd = m.credModel.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case common.BackToNavigationMsg:
		m.currentMenu = navigationMenu
		return m, nil

	case common.CreateRuleMsg:
		m.ruleModel.UpdateStateData(msg.StateId)
		m.currentMenu = ruleSubMenu

	case common.UpdateCredNamesMsg:
		m.deviceModel, cmd = m.deviceModel.Update(msg)
		return m, cmd

	case common.UpdateActionNamesMsg:
		m.ruleModel, cmd = m.ruleModel.Update(msg)
		return m, cmd
	}

	// Send message to the appropriate component
	switch m.currentMenu {
	case navigationMenu:
		m.mainList, cmd = m.mainList.Update(msg)
		cmds = append(cmds, cmd)
	case deviceSubMenu:
		m.deviceModel, cmd = m.deviceModel.Update(msg)
		cmds = append(cmds, cmd)
	case actionSubMenu:
		m.actionModel, cmd = m.actionModel.Update(msg)
		cmds = append(cmds, cmd)
	case ruleSubMenu:
		m.ruleModel, cmd = m.ruleModel.Update(msg)
		cmds = append(cmds, cmd)
	case stateSubMenu:
		m.stateModel, cmd = m.stateModel.Update(msg)
		cmds = append(cmds, cmd)
	case credSubMenu:
		m.credModel, cmd = m.credModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m NavigationModel) View() string {
	var str string

	switch m.currentMenu {
	case navigationMenu:
		mainStyle, dataStyle, dynamicStyle, viewerStyle := getWindowStyles(m.width, m.height)

		mainStyle = common.ApplyActiveStyle(mainStyle)

		top := lipgloss.JoinHorizontal(
			lipgloss.Top,
			mainStyle.Render(m.mainList.View()),
			dataStyle.String(),
			dynamicStyle.String())

		str = lipgloss.JoinVertical(
			lipgloss.Left,
			top,
			common.StatusMessageStyle.String(),
			lipgloss.PlaceHorizontal(m.width, lipgloss.Center, viewerStyle.String()),
		)

	case actionSubMenu:
		str = m.actionModel.View()
	case ruleSubMenu:
		str = m.ruleModel.View()
	case stateSubMenu:
		str = m.stateModel.View()
	case deviceSubMenu:
		str = m.deviceModel.View()
	case credSubMenu:
		str = m.credModel.View()
	}

	var footerStr string
	if m.currentMenu == credSubMenu {
		footerStr = m.help.View(common.FooterKeys) + common.SepStyle.Render(m.help.ShortSeparator) + m.help.View(common.CredFooterKeys)
	} else {
		footerStr = m.help.View(common.FooterKeys)
	}
	header := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, common.HeaderTextStyle.Render(title))
	footer := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, footerStr)
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		str,
		footer,
	)
}
