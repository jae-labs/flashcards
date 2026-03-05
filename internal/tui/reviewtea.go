package tui

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"flashcards/internal/store"
	"flashcards/internal/tui/keys"
	"flashcards/internal/tui/layout"
	"flashcards/internal/tui/theme"
)

// ReviewModel manages the state for the review session
// Views: question, answer, correct/incorrect, revisitIn, done

type viewState int

const (
	viewQuestion viewState = iota
	viewAnswer
	viewRevisitIn
	viewDone
	viewTimeout
)

var completionMessages = []string{
	"The Void retreats… for now 🕳️🐾",
	"Knowledge absorbed. The Void purrs in approval 😼",
	"You’ve conquered the deck. The Void whispers… ‘impressive.’ 🌌",
	"Session synced with the Void’s neural core 🧠✨",
	"The Void stares back, but you stand unshaken ⚫",
	"Memory integrated. The Void grows quieter… temporarily 🔮",
	"Another victory against the Void 🐈‍⬛",
}

type ReviewModel struct {
	flashcards    []store.Flashcard
	current       int
	view          viewState
	resultMsg     string
	quitting      bool
	correct       []bool
	revisitIn     []int
	width         int
	height        int
	progress      progress.Model
	timer         timer.Model
	startTime     time.Time
	duration      time.Duration
	interval      time.Duration // add interval for timer ticks
	completionMsg string
}

func NewReviewModel(flashcards []store.Flashcard) *ReviewModel {
	// Use a custom gradient for the progress bar
	d := 30 * time.Second
	interval := 100 * time.Millisecond // smoother animation
	p := progress.New(progress.WithGradient("#ff00e1ff", "#ff00e1ff"))
	p.ShowPercentage = false
	return &ReviewModel{
		flashcards: flashcards,
		current:    0,
		view:       viewQuestion,
		correct:    make([]bool, len(flashcards)),
		revisitIn:  make([]int, len(flashcards)),
		progress:   p,
		timer:      timer.NewWithInterval(d, interval),
		startTime:  time.Now(),
		duration:   d,
		interval:   interval,
	}
}

func (m *ReviewModel) Init() tea.Cmd {
	return m.timer.Init()
}

func (m *ReviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case timer.TimeoutMsg:
		m.view = viewAnswer
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		if msg.String() == keys.Q {
			m.quitting = true
			return m, tea.Quit
		}
		switch m.view {
		case viewQuestion:
			if msg.String() == keys.Enter {
				m.view = viewAnswer
			}
		case viewAnswer:
			if msg.String() == keys.C {
				m.correct[m.current] = true
				m.view = viewRevisitIn
			} else if msg.String() == keys.I {
				m.correct[m.current] = false
				m.resultMsg = "Marked incorrect. Card will not be scheduled for repetition."
				cmd = m.nextCard()
				cmds = append(cmds, cmd)
			}
		case viewRevisitIn:
			switch msg.String() {
			case keys.One:
				m.revisitIn[m.current] = 1
				m.resultMsg = "Revisit in 1 day"
				cmd = m.nextCard()
				cmds = append(cmds, cmd)
			case keys.Three:
				m.revisitIn[m.current] = 3
				m.resultMsg = "Revisit in 3 days"
				cmd = m.nextCard()
				cmds = append(cmds, cmd)
			case keys.Seven:
				m.revisitIn[m.current] = 7
				m.resultMsg = "Revisit in 7 days"
				cmd = m.nextCard()
				cmds = append(cmds, cmd)
			case keys.Nine:
				m.revisitIn[m.current] = 9
				m.resultMsg = "Revisit in 9 days"
				cmd = m.nextCard()
				cmds = append(cmds, cmd)
			}
		case viewDone:
			if msg.String() == keys.Q {
				m.quitting = true
				return m, tea.Quit
			}
		}
	case tea.QuitMsg:
		m.quitting = true
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	case timer.TickMsg:
		// Animate progress bar on every timer tick
		if m.view == viewQuestion {
			elapsed := time.Since(m.startTime)
			percent := float64(elapsed) / float64(m.duration)
			if percent < 0 {
				percent = 0
			}
			if percent > 1.0 {
				percent = 1.0
			}
			cmd = m.progress.SetPercent(percent)
			cmds = append(cmds, cmd)
		}
		m.timer, cmd = m.timer.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}
	return m, tea.Batch(cmds...)
}

func (m *ReviewModel) nextCard() tea.Cmd {
	m.current++
	if m.current >= len(m.flashcards) {
		m.view = viewDone
		b := make([]byte, 4)
		if _, err := rand.Read(b); err != nil {
			m.completionMsg = completionMessages[0]
		} else {
			index := int(binary.LittleEndian.Uint32(b)) % len(completionMessages)
			m.completionMsg = completionMessages[index]
		}
		return nil
	}
	m.view = viewQuestion
	m.resultMsg = ""
	m.duration = 30 * time.Second
	m.startTime = time.Now()
	m.timer = timer.NewWithInterval(m.duration, m.interval)
	m.progress = progress.New(progress.WithGradient("#ff00e1ff", "#ff00e1ff"))
	m.progress.ShowPercentage = false
	return m.timer.Init()
}

// Results API for review command
func (m *ReviewModel) FlashcardWasCorrect(idx int) bool {
	if idx < 0 || idx >= len(m.correct) {
		return false
	}
	return m.correct[idx]
}

func (m *ReviewModel) FlashcardRevisitIn(idx int) int {
	if idx < 0 || idx >= len(m.revisitIn) {
		return 0
	}
	return m.revisitIn[idx]
}

func (m *ReviewModel) View() string {
	if m.quitting {
		return "Goodbye!"
	}
	total := len(m.flashcards)
	current := m.current + 1
	// If review is done, show total instead of current+1
	if m.view == viewDone {
		current = total
	}

	// Use layout helper for content width
	width := layout.CalculateContentWidth(m.width)

	// Create frame using layout helper
	frame := layout.CreateFrame(width,
		layout.WithAlignment(lipgloss.Center, lipgloss.Center))

	// Count correct and incorrect answers
	correctCount := 0
	incorrectCount := 0
	for i := range m.flashcards {
		// Only count cards that have been answered (i.e., where correct/incorrect has been set)
		if m.view == viewDone || i < m.current {
			if m.FlashcardWasCorrect(i) {
				correctCount++
			} else {
				incorrectCount++
			}
		}
	}

	// Bottom bar: ✅ left, current/total center, ❌ right
	bottomBar := layout.CreateBottomBar(width,
		fmt.Sprintf("✅ %d", correctCount),
		fmt.Sprintf("%d/%d", current, total),
		fmt.Sprintf("❌ %d", incorrectCount))

	exitMsg := theme.InfoStyle.Render("Enter: Confirm • q: Quit")

	var content string
	switch m.view {
	case viewQuestion:
		// Animated progress bar for countdown
		progressBar := m.progress.View()
		content = fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s", theme.QuestionStyle.Render("Question:"), m.flashcards[m.current].Question, progressBar, bottomBar)
	case viewAnswer:
		content = fmt.Sprintf("%s\n\n%s\n\n%s\n%s", theme.AnswerStyle.Render("Answer:"), m.flashcards[m.current].Answer, theme.InfoStyle.Render("Was your answer correct? [c]orrect / [i]ncorrect\n"), bottomBar)
	case viewRevisitIn:
		content = fmt.Sprintf("%s\n%s\n%s", theme.InfoStyle.Render("\nRevisit in (days): [1]  [3]  [7]  [9]"), m.resultMsg, bottomBar)
	case viewDone:
		content = fmt.Sprintf("\n%s\n%s", theme.SuccessStyle.Render(m.completionMsg+"\n"), bottomBar)
	}
	return layout.CenterContent(m.width, m.height, frame.Render(content)+"\n"+exitMsg)
}
