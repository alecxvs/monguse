package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var defaultStyle = lipgloss.NewStyle().
	Foreground(lipgloss.CompleteColor{TrueColor: "#FAFAFA", ANSI256: "254", ANSI: "15"}).
	Background(lipgloss.CompleteColor{TrueColor: "#1F1F1F", ANSI256: "234", ANSI: "0"})

var screenStyle = defaultStyle.
	Bold(true).
	PaddingTop(2).
	PaddingLeft(4).
	PaddingBottom(2).
	PaddingRight(40).
	// Width(80).
	Height(16).
	BorderStyle(lipgloss.DoubleBorder()).
	BorderForeground(lipgloss.CompleteColor{TrueColor: "#1B913E", ANSI256: "65", ANSI: "2"})

var headerStyle = defaultStyle.
	Bold(true).
	Foreground(lipgloss.CompleteColor{TrueColor: "#CACACA", ANSI256: "251", ANSI: "7"})

var hotStyle = defaultStyle.
	Bold(true).
	Foreground(lipgloss.CompleteColor{TrueColor: "#87FA87", ANSI256: "120", ANSI: "10"})

var completedStyle = defaultStyle.
	Bold(false).
	Foreground(lipgloss.CompleteColor{TrueColor: "#4A5A4A", ANSI256: "65", ANSI: "8"})

var errorStyle = defaultStyle.
	Bold(false).
	Foreground(lipgloss.CompleteColor{TrueColor: "#CF6565", ANSI256: "197", ANSI: "9"}).
	Background(lipgloss.CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"})

type loadingStep int

const (
	step_LoadConfig loadingStep = iota
	step_ParseConfig
	step_FindGamePath
	step_CheckGamePath
	step_LoadModProfiles
	step_GetFunky
)

type LoaderModel struct {
	spinner spinner.Model
	step    loadingStep
	err     error
}

func NewLoaderModel() LoaderModel {
	return LoaderModel{
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Line),
			spinner.WithStyle(defaultStyle.Foreground(lipgloss.ANSIColor(15))),
		),
		step: step_LoadModProfiles,
	}
}

func (m LoaderModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m LoaderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		default:
			return m, nil
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case error:
		m.err = msg
		return m, nil

	default:
		return m, nil
	}
}

func (m LoaderModel) View() string {

	build := strings.Builder{}
	build.WriteString(headerStyle.Render("Pre-flight check") + "\n")

	if m.err != nil {
		build.WriteString(errorStyle.Render("*** Unhandled exception during loading sequence ***\n"))
		build.WriteString(fmt.Sprintf("Error: %s\n", m.err.Error()))
		build.WriteString(fmt.Sprintf("Additional Info: m.step == `%d`\n", m.step))
		return screenStyle.Render(build.String())
	}

	strChecklist := []string{}

	switch m.step {
	case step_GetFunky:
		strChecklist = append(strChecklist, " * Loading Completed")
		fallthrough
	case step_LoadModProfiles:
		strChecklist = append(strChecklist, " * Loading Mod Profiles")
		fallthrough
	case step_CheckGamePath:
		strChecklist = append(strChecklist, " * Verifying Among Us Installation Path")
		fallthrough
	case step_FindGamePath:
		strChecklist = append(strChecklist, " * Finding Among Us Installation Path")
		fallthrough
	case step_ParseConfig:
		strChecklist = append(strChecklist, " * Parsing Configuration")
		fallthrough
	case step_LoadConfig:
		strChecklist = append(strChecklist, " * Loading Configuration")
		fallthrough
	default:
		if len(strChecklist) <= 0 {
			build.WriteString(errorStyle.Render(fmt.Sprintf("\n\n%s\n%s\n%s\n\n%s\n",
				"*** Fatal error: Wrong state step in loading sequence ***",
				"This should never happen! Please submit a bug report!",
				fmt.Sprintf("Additional Info: m.step == `%d`", m.step),
				"Please quit the application and try again.",
			)))
		}
	}

	for i := len(strChecklist) - 1; i >= 1; i-- {
		build.WriteString(completedStyle.Render(strChecklist[i]) + "\n")
	}
	build.WriteString(hotStyle.Render(strChecklist[0], m.spinner.View()) + "\n")

	// str := fmt.Sprintf("\n\n   %s Loading forever...press q to quit\n\n", m.spinner.View())
	return screenStyle.Render(build.String())
}
