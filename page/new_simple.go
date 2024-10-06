package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	frames  []string
	current int
}

func initialModel() model {
	return model{
		frames:  []string{"🌑 ", "🌒 ", "🌓 ", "🌔 ", "🌕 ", "🌖 ", "🌗 ", "🌘 "},
		current: 0,
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	case tickMsg:
		m.current = (m.current + 1) % len(m.frames)
		return m, tickCmd()
	}
	return m, nil
}

func (m model) View() string {
	return m.frames[m.current]
}

type tickMsg struct{}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
	}
}
