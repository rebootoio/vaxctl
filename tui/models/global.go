package models

import (
	"strings"
	"vaxctl/tui/common"

	"github.com/charmbracelet/lipgloss"
)

const (
	title = "vaxctl interactive mode"

	showYamlAction     = "Show YAML"
	saveToFileAction   = "Save To File"
	saveToServerAction = "Save To Server"
	clearFieldsAction  = "Clear Fields"
	backAction         = "Back"
	createRuleAction   = "Create Rule from State"

	mainView    = "main"
	dataView    = "data"
	dynamicView = "dynamic"
	viewerView  = "viewer"

	createMode = "Create"
	editMode   = "Edit"

	yamlViewerOption = "YAML Viewer"
	yamlSaveOption   = "YAML Save"

	topRowMinHeight    = 14
	mainWindowWidth    = 30
	dataWindowMinWidth = 60
)

var (
	mainActions      = []string{showYamlAction, saveToFileAction, saveToServerAction, clearFieldsAction, backAction}
	stateMainActions = []string{saveToServerAction, createRuleAction, clearFieldsAction, backAction}
	views            = []string{mainView, dataView, dynamicView, viewerView}
)

func getTopRowHeight(height int) int {
	topRowHeight := (height - common.VerticalMarginHeight) / 3
	if topRowHeight < topRowMinHeight {
		topRowHeight = topRowMinHeight
	}
	return topRowHeight
}
func getDataWidth(width int) int {
	dataWidth := (width - mainWindowWidth) / 3
	if dataWidth < dataWindowMinWidth {
		dataWidth = dataWindowMinWidth
	}
	return dataWidth
}

func getSizing(width int, height int) (int, int, int, int, int) {
	topRowHeight := getTopRowHeight(height)
	borderWidth, borderHeight := lipgloss.Size(common.ActiveBlockBorder.String())
	rowHeight := topRowHeight - borderHeight
	mainWidth := mainWindowWidth

	dataWidth := getDataWidth(width)
	dynamicWidth := (width - mainWindowWidth - dataWidth)

	viewerHeight := height - common.VerticalMarginHeight - topRowHeight
	return rowHeight, mainWidth - borderWidth, dataWidth - borderWidth, dynamicWidth - borderWidth, viewerHeight
}

func getWindowStyles(width int, height int) (lipgloss.Style, lipgloss.Style, lipgloss.Style, lipgloss.Style) {
	topRowHeight := getTopRowHeight(height)
	dataWidth := getDataWidth(width)
	mainStyle := common.UpdateStyle(common.BlockStyle, mainWindowWidth, topRowHeight)
	dataStyle := common.UpdateStyle(common.BlockStyle, dataWidth, topRowHeight)
	dynamicStyle := common.UpdateStyle(common.BlockStyle, (width - mainWindowWidth - dataWidth), lipgloss.Height(mainStyle.String()))
	viewerStyle := common.UpdateStyle(common.BlockStyle, width, height-common.VerticalMarginHeight-lipgloss.Height(mainStyle.String())).Inherit(common.BottomBlockBorderStyle)
	return mainStyle, dataStyle, dynamicStyle, viewerStyle
}

func getStatusMessage(status string, isError bool) string {
	var statusLine string
	if isError {
		statusLine = common.StatusErrorStyle.Render(strings.ReplaceAll(status, "\n", " "))
	} else {
		statusLine = common.StatusMessageStyle.Render(status)
	}
	return statusLine
}
