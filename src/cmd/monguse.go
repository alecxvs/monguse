package main

import (
	"fmt"
	"os"

	"github.com/alecxvs/monguse/src/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

var APP_VERSION string = `v1.0.0`

func main() {
	// s := shell.New()
	// s.Navstack.Push(navstack.NavigationItem{Model: m, Title: "Colors"})

	p := tea.NewProgram(models.NewLoaderModel(), tea.WithAltScreen())

	_, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
