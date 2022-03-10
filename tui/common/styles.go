package common

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

var (
	HeaderTextStyle = lipgloss.NewStyle().Align(lipgloss.Center).Height(1).Foreground(lipgloss.Color("13"))
	BorderredStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

	VerticalMarginHeight = lipgloss.Height(HeaderTextStyle.String())*2 + lipgloss.Height(StatusErrorStyle.String())

	BorderColorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("62"))
	ActiveBorderColorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))

	TableColumn = lipgloss.NewStyle().Align(lipgloss.Left)

	ItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))

	StatusMessageStyle = lipgloss.NewStyle().Height(1).Foreground(lipgloss.Color("145"))
	StatusErrorStyle   = lipgloss.NewStyle().Height(1).Background(lipgloss.Color("1"))

	ValidStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	InvalidStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

	ResourceItemStyle         = lipgloss.NewStyle().Height(1)
	ResourceSelectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	ResourceTitleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))

	BlockStyle = lipgloss.NewStyle().
			PaddingRight(1).
			PaddingLeft(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))

	ActiveBlockBorder = lipgloss.NewStyle().BorderStyle(lipgloss.DoubleBorder()).BorderForeground(lipgloss.Color("42"))

	BottomBlockBorderStyle = lipgloss.NewStyle().
				BorderLeft(false).
				BorderRight(false).
				BorderBottom(true).
				BorderTop(true)

	OcrBlockBorderStyle = lipgloss.NewStyle().
				BorderLeft(false).
				BorderRight(false).
				BorderBottom(false).
				BorderTop(true)

	OcrMatchedTextStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))

	NotSetString = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")).Render("Not Set")

	KeyStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#909090",
		Dark:  "#bdbdbd",
	})

	DescStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#B2B2B2",
		Dark:  "#969696",
	})

	SepStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#DDDADA",
		Dark:  "#828282",
	})
)

func ApplyActiveStyle(style lipgloss.Style) lipgloss.Style {
	return style.Copy().BorderStyle(lipgloss.DoubleBorder()).BorderForeground(lipgloss.Color("42"))
}
func UpdateStyle(style lipgloss.Style, width int, height int) lipgloss.Style {
	return style.Copy().Width(width - 2).MaxWidth(width).Height(height - 2).MaxHeight(height)
}

func UpdateStyleWidth(style lipgloss.Style, width int) lipgloss.Style {
	return style.Copy().Width(width - 2).MaxWidth(width)
}

func GetHelpModel() help.Model {
	helpModel := help.New()
	helpModel.Styles.ShortKey = KeyStyle
	helpModel.Styles.ShortDesc = DescStyle
	helpModel.Styles.ShortSeparator = SepStyle
	helpModel.Styles.FullKey = KeyStyle
	helpModel.Styles.FullDesc = DescStyle
	helpModel.Styles.FullSeparator = SepStyle
	return helpModel
}
