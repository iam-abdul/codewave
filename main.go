package main

// An example program demonstrating the pager component from the Bubbles
// component library.

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// You generally won't need this unless you're processing stuff with
// complicated ANSI escape sequences. Turn it on if you notice flickering.
//
// Also keep in mind that high performance rendering only works for programs
// that use the full size of the terminal. We're enabling that below with
// tea.EnterAltScreen().
const useHighPerformanceRenderer = false

type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Millisecond * 150)
	return tickMsg{}
}

var (
	titleStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

var style = lipgloss.NewStyle().
	Bold(true).Align(lipgloss.Center)

type heroAnimation struct {
	content             string
	heroContent         []string
	backwards           bool
	index               int
	charAt              int
	cursor              string
	desiredCursorBlinks int
	currentCursorBlinks int
}

type model struct {
	heroSection      heroAnimation
	ready            bool
	viewport         viewport.Model
	selectedMenuItem int
	menuItems        []string
	debug            string
}

func (m model) Init() tea.Cmd {
	return tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

		if k := msg.String(); k == "l" || k == "right" {
			// move right
			if m.selectedMenuItem >= len(m.menuItems)-1 || m.selectedMenuItem < 0 {
				m.selectedMenuItem = 0
			} else {
				m.selectedMenuItem += 1
			}
		}

		if k := msg.String(); k == "h" || k == "left" {
			// move left
			if m.selectedMenuItem == 0 {
				m.selectedMenuItem = len(m.menuItems) - 1
			} else {
				m.selectedMenuItem -= 1
			}
		}

	case tickMsg:
		newModel := m.updateHeroMenuAnimation()
		return newModel, tick

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(m.headerView())
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}
	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	// return fmt.Sprintf("%s\n%s\n%s\n%s", m.headerView(), style.MarginLeft((m.viewport.Width-len(m.content))/2).Render(m.viewport.View()), m.heroMenuAnimation(), m.footerView())
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.heroMenuAnimation(), m.footerView())
}

func (m model) headerView() string {
	titleContent := "CODEWAVE"
	title := titleStyle.Border(lipgloss.HiddenBorder()).Render(titleContent)

	otherOptions := ""
	for i, item := range m.menuItems {
		highlighted := lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#fefefe"))
		if i == m.selectedMenuItem {
			otherOptions = otherOptions + "  " + highlighted.Render(item)
		} else {
			otherOptions = otherOptions + "  " + item
		}
	}

	line := strings.Repeat(" ", max(0, m.viewport.Width-lipgloss.Width(title)-lipgloss.Width(otherOptions)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line, otherOptions)
}

func (m model) updateHeroMenuAnimation() model {
	currentWord := m.heroSection.heroContent[m.heroSection.index]

	if m.heroSection.backwards {
		if m.heroSection.currentCursorBlinks < m.heroSection.desiredCursorBlinks {
			cursorPresent := strings.Contains(m.heroSection.content, m.heroSection.cursor)
			if cursorPresent {
				// this will not have cursor
				m.heroSection.content = m.heroSection.heroContent[m.heroSection.index]
			} else {
				blinkCursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
				m.heroSection.content = m.heroSection.heroContent[m.heroSection.index] + blinkCursorStyle.Render(m.heroSection.cursor)
			}
			m.heroSection.currentCursorBlinks += 1
			return m
		}
		// this is for removing the text
		if m.heroSection.charAt <= 0 {
			// backward is done, move to next character
			m.heroSection.backwards = false
			m.heroSection.charAt = 0
			m.heroSection.index += 1
			if m.heroSection.index >= len(m.heroSection.heroContent) {
				m.heroSection.index = 0
			}
		} else {
			m.heroSection.content = currentWord[:m.heroSection.charAt-1] + m.heroSection.cursor
			m.heroSection.charAt -= 1
		}

	} else {
		if m.heroSection.charAt >= len(currentWord) {
			m.heroSection.backwards = true
			m.heroSection.currentCursorBlinks = 0
		} else {
			m.heroSection.content = currentWord[:m.heroSection.charAt+1] + m.heroSection.cursor
			m.heroSection.charAt += 1
		}
	}

	return m
}

func (m model) centerHorizontally(content string) string {
	paddingLeftAndRight := strings.Repeat(" ", (m.viewport.Width-lipgloss.Width(m.heroSection.content))/2)
	return lipgloss.JoinHorizontal(lipgloss.Center, paddingLeftAndRight, content, paddingLeftAndRight)
}

func (m model) centerVertically(content string) string {
	paddingTopAndBottom := strings.Repeat("\n", m.viewport.Height)
	return lipgloss.JoinHorizontal(lipgloss.Center, paddingTopAndBottom, content, paddingTopAndBottom)
}

func (m model) heroMenuAnimation() string {
	mainContent := m.centerVertically(lipgloss.JoinVertical(lipgloss.Center, m.centerHorizontally(m.heroSection.content)))
	return lipgloss.JoinHorizontal(lipgloss.Center, mainContent)
}

func (m model) footerView() string {
	info := infoStyle.Render(m.debug)
	// info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.selectedMenuItem))
	return lipgloss.JoinHorizontal(lipgloss.Center, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	// Load some text for our viewport
	content := []string{"BRING CHANGE", "TAKE CHANCES", "CHOOSE DELIGHT", "CREATE VALUE", "THINK FREE", "PILOT IDEAS", "SPARK ACTION"}

	p := tea.NewProgram(
		model{heroSection: heroAnimation{heroContent: content, charAt: 0, index: 0, backwards: false, cursor: "▋", desiredCursorBlinks: 3, currentCursorBlinks: 0}, selectedMenuItem: 0, menuItems: []string{"SERVICES", "DONE IN A WEEK", "INDUSTRIES", "WORKS", "INSIGHTS", "CULTURE", "CONTACT"}},
		tea.WithAltScreen(),      // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseAllMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
