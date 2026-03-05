package tui

import (
	"fmt"
	"strconv"
	"strings"

	"flashcards/internal/store"
	"flashcards/internal/tui/components"
	"flashcards/internal/tui/keys"
	"flashcards/internal/tui/layout"
	"flashcards/internal/tui/theme"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AdminModel provides CRUD management of flashcards.
type adminView int

const (
	adminList adminView = iota
	adminCreate
	adminEdit
	adminConfirmDelete
	adminConfirmBulkReset
)

type keyMap struct {
	Up        key.Binding
	Down      key.Binding
	PageUp    key.Binding
	PageDown  key.Binding
	Create    key.Binding
	Edit      key.Binding
	Delete    key.Binding
	BulkReset key.Binding
	Reload    key.Binding
	Help      key.Binding
	Quit      key.Binding
	Cancel    key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Cancel}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.PageUp, k.PageDown},
		{k.Create, k.Edit, k.Delete, k.BulkReset, k.Reload},
	}
}

var adminKeys = keyMap{
	Up:        key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k:", "Navigate")),
	Down:      key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j:", "Move Down")),
	PageUp:    key.NewBinding(key.WithKeys("pgup", "pageup"), key.WithHelp("pgup:", "Page Up")),
	PageDown:  key.NewBinding(key.WithKeys("pgdown", "pagedown"), key.WithHelp("pgdown:", "Page Down")),
	Create:    key.NewBinding(key.WithKeys("c"), key.WithHelp("c:", "Create")),
	Edit:      key.NewBinding(key.WithKeys("e"), key.WithHelp("e:", "Edit")),
	Delete:    key.NewBinding(key.WithKeys("d"), key.WithHelp("d:", "Delete")),
	BulkReset: key.NewBinding(key.WithKeys("b"), key.WithHelp("b:", "Bulk Reset")),
	Reload:    key.NewBinding(key.WithKeys("r"), key.WithHelp("r:", "Reload")),
	Help:      key.NewBinding(key.WithKeys("?"), key.WithHelp("?:", "Toggle Help")),
	Quit:      key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q:", "Quit")),
	Cancel:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc:", "Cancel")),
}

type AdminModel struct {
	flashcards []store.Flashcard
	selected   int
	view       adminView
	width      int
	height     int
	help       help.Model
	keys       keyMap

	// table for list view
	table table.Model

	// form fields
	questionInput textinput.Model
	answerInput   textinput.Model
	revisitInput  textinput.Model // days until next review

	status components.StatusMessage

	storeRef *store.Store
}

func NewAdminModel(storeRef *store.Store, flashcards []store.Flashcard) *AdminModel {
	q := textinput.New()
	q.Placeholder = "Question"
	q.Focus()
	a := textinput.New()
	a.Placeholder = "Answer"
	r := textinput.New()
	r.Placeholder = "Days (e.g. 7)"

	// Setup table columns
	columns := []table.Column{
		{Title: "ID", Width: 6},
		{Title: "Question", Width: 40},
		{Title: "Answer", Width: 30},
		{Title: "Revisit In", Width: 12},
	}

	// Convert flashcards to table rows
	rows := makeTableRows(flashcards, columns)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	// Style the table using theme
	t = theme.ApplyTableStyles(t)

	return &AdminModel{
		flashcards:    flashcards,
		selected:      0,
		view:          adminList,
		table:         t,
		questionInput: q,
		answerInput:   a,
		revisitInput:  r,
		storeRef:      storeRef,
		help:          help.New(),
		keys:          adminKeys,
	}
}

// makeTableRows converts flashcards to table rows
func makeTableRows(flashcards []store.Flashcard, columns []table.Column) []table.Row {
	rows := make([]table.Row, len(flashcards))
	for i, fc := range flashcards {
		rows[i] = table.Row{
			fmt.Sprintf("%d", fc.ID),
			truncate(fc.Question, columns[1].Width),
			truncate(fc.Answer, columns[2].Width),
			fmt.Sprintf("%d days", fc.RevisitIn),
		}
	}
	return rows
}

func (m *AdminModel) Init() tea.Cmd { return nil }

func (m *AdminModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate the effective frame width using layout helper
		frameWidth := layout.CalculateContentWidth(m.width)

		// Calculate column widths using layout helper
		idWidth, questionWidth, answerWidth, revisitInWidth := layout.CalculateTableColumnWidths(frameWidth)

		m.table.SetColumns([]table.Column{
			{Title: "ID", Width: idWidth},
			{Title: "Question", Width: questionWidth},
			{Title: "Answer", Width: answerWidth},
			{Title: "Revisit In", Width: revisitInWidth},
		})

		// Calculate table height using layout helper
		tableHeight := layout.CalculateTableHeight(m.height)
		m.table.SetHeight(tableHeight)

		// Update rows with new column widths
		rows := makeTableRows(m.flashcards, m.table.Columns())
		m.table.SetRows(rows)
		return m, nil

	case tea.KeyMsg:
		switch m.view {
		case adminList:
			return m.handleListView(msg)
		case adminCreate:
			return m.handleCreateView(msg)
		case adminEdit:
			return m.handleEditView(msg)
		case adminConfirmDelete:
			return m.handleDeleteConfirm(msg)
		case adminConfirmBulkReset:
			return m.handleBulkResetConfirm(msg)
		}
	}
	return m, nil
}

func (m *AdminModel) handleListView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.Create):
		m.view = adminCreate
		m.resetForm()
		return m, nil
	case key.Matches(msg, m.keys.Edit):
		if len(m.flashcards) == 0 {
			return m, nil
		}
		m.selected = m.table.Cursor()
		m.loadSelectedIntoForm()
		m.view = adminEdit
		return m, nil
	case key.Matches(msg, m.keys.Delete):
		if len(m.flashcards) == 0 {
			return m, nil
		}
		m.selected = m.table.Cursor()
		m.view = adminConfirmDelete
		return m, nil
	case key.Matches(msg, m.keys.BulkReset):
		if len(m.flashcards) == 0 {
			return m, nil
		}
		m.view = adminConfirmBulkReset
		return m, nil
	case key.Matches(msg, m.keys.Reload):
		m.reload()
		m.status.SetSuccess("Table refreshed")
		return m, nil
	case key.Matches(msg, m.keys.Help):
		m.help.ShowAll = !m.help.ShowAll
		return m, nil
	case key.Matches(msg, m.keys.PageDown):
		newPos := m.table.Cursor() + 10
		if newPos >= len(m.flashcards) {
			newPos = len(m.flashcards) - 1
		}
		m.table.SetCursor(newPos)
		m.selected = newPos
		return m, nil
	case key.Matches(msg, m.keys.PageUp):
		newPos := m.table.Cursor() - 10
		if newPos < 0 {
			newPos = 0
		}
		m.table.SetCursor(newPos)
		m.selected = newPos
		return m, nil
	default:
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		m.selected = m.table.Cursor()
		return m, cmd
	}
}

func (m *AdminModel) handleCreateView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Cancel):
		m.view = adminList
		return m, nil
	}
	if msg.String() == keys.Tab {
		m.cycleFocus()
		return m, nil
	}
	if msg.String() == keys.Enter {
		m.createFlashcard()
		return m, nil
	}
	m.updateInputs(msg)
	return m, nil
}

func (m *AdminModel) handleEditView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Cancel):
		m.view = adminList
		return m, nil
	}
	if msg.String() == keys.Tab {
		m.cycleFocus()
		return m, nil
	}
	if msg.String() == keys.Enter {
		m.updateFlashcard()
		return m, nil
	}
	m.updateInputs(msg)
	return m, nil
}

func (m *AdminModel) handleDeleteConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Cancel):
		m.view = adminList
		return m, nil
	}
	if msg.String() == keys.Y {
		m.deleteFlashcard()
		return m, nil
	}
	if msg.String() == keys.N {
		m.view = adminList
	}
	return m, nil
}

func (m *AdminModel) handleBulkResetConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Cancel):
		m.view = adminList
		return m, nil
	}
	if msg.String() == keys.Y {
		m.bulkResetRevisitIn()
		return m, nil
	}
	if msg.String() == keys.N {
		m.view = adminList
	}
	return m, nil
}

func (m *AdminModel) View() string {
	// Set max content width using layout helper
	width := layout.CalculateContentWidth(m.width)

	// Set max content height using layout helper
	maxHeight := layout.CalculateMaxFrameHeight(m.height)

	// Create frame with border using layout helper
	frame := layout.CreateFrame(width,
		layout.WithMaxHeight(maxHeight),
		layout.WithAlignment(lipgloss.Left, lipgloss.Top),
		layout.WithPadding(1, 1))

	helpBar := m.help.View(m.keys)

	statusBar := m.status.Render()

	var mainContent string
	var exitMsg string

	switch m.view {
	case adminList:
		var b strings.Builder

		if len(m.flashcards) == 0 {
			b.WriteString("No flashcards. Press 'c' to create.\n")
		} else {
			b.WriteString(m.table.View())
		}
		b.WriteString("\n")
		b.WriteString(helpBar)
		b.WriteString(statusBar)
		mainContent = b.String()
		exitMsg = ""

	case adminCreate:
		formContent := components.RenderFormFields(
			components.FormField{Label: "Question:", Input: m.questionInput},
			components.FormField{Label: "Answer:", Input: m.answerInput},
			components.FormField{Label: "Revisit (days):", Input: m.revisitInput},
		)
		mainContent = formContent + statusBar
		exitMsg = theme.HelpStyle.Render("tab: Next Field • Enter: Confirm • esc: Cancel")

	case adminEdit:
		title := theme.TitleStyle.Render(fmt.Sprintf("Edit Flashcard (ID %d)", m.flashcards[m.selected].ID))
		formContent := components.RenderFormFields(
			components.FormField{Label: "Question:", Input: m.questionInput},
			components.FormField{Label: "Answer:", Input: m.answerInput},
			components.FormField{Label: "Revisit (days):", Input: m.revisitInput},
		)
		mainContent = fmt.Sprintf("%s\n\n%s%s", title, formContent, statusBar)
		exitMsg = theme.HelpStyle.Render("tab: Next Field • Enter: Confirm • esc: Cancel")

	case adminConfirmDelete:
		warning := theme.ErrorStyle.Render(fmt.Sprintf("Delete Flashcard ID %d?", m.flashcards[m.selected].ID))
		mainContent = fmt.Sprintf("%s\n", warning)
		exitMsg = theme.HelpStyle.Render("y: Yes • n: No • esc: Cancel")

	case adminConfirmBulkReset:
		warning := theme.ErrorStyle.Render(fmt.Sprintf("Set RevisitIn to 0 for ALL %d flashcards?", len(m.flashcards)))
		info := theme.InfoStyle.Render("This will make all flashcards due for immediate review.")
		mainContent = fmt.Sprintf("%s\n\n%s\n", warning, info)
		exitMsg = theme.HelpStyle.Render("y: Yes • n: No • esc: Cancel")
	}

	// Render frame with content and exit message below (outside frame)
	framedContent := frame.Render(mainContent)
	if exitMsg != "" {
		framedContent += "\n" + exitMsg
	}

	return layout.CenterContent(m.width, m.height, framedContent)
}

// helpers
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func (m *AdminModel) resetForm() {
	m.questionInput.SetValue("")
	m.answerInput.SetValue("")
	m.revisitInput.SetValue("")
	m.status.Clear()
	m.questionInput.Focus()
	m.answerInput.Blur()
	m.revisitInput.Blur()
}

func (m *AdminModel) loadSelectedIntoForm() {
	fc := m.flashcards[m.selected]
	m.questionInput.SetValue(fc.Question)
	m.answerInput.SetValue(fc.Answer)
	m.revisitInput.SetValue(strconv.Itoa(fc.RevisitIn))
	m.questionInput.Focus()
	m.answerInput.Blur()
	m.revisitInput.Blur()
}

func (m *AdminModel) updateInputs(msg tea.Msg) {
	// Update the focused input with the key message
	switch {
	case m.questionInput.Focused():
		m.questionInput, _ = m.questionInput.Update(msg)
	case m.answerInput.Focused():
		m.answerInput, _ = m.answerInput.Update(msg)
	case m.revisitInput.Focused():
		m.revisitInput, _ = m.revisitInput.Update(msg)
	}
}

func (m *AdminModel) parseRevisitDays() (int, error) {
	val := strings.TrimSpace(m.revisitInput.Value())
	if val == "" {
		return 0, fmt.Errorf("revisit days required")
	}
	d, err := strconv.Atoi(val)
	if err != nil || d < 0 {
		return 0, fmt.Errorf("invalid days")
	}
	return d, nil
}

func (m *AdminModel) createFlashcard() {
	days, err := m.parseRevisitDays()
	if err != nil {
		m.status.SetError(err.Error())
		return
	}
	fc := store.Flashcard{Question: m.questionInput.Value(), Answer: m.answerInput.Value(), File: "manual", RevisitIn: days}
	if err := m.storeRef.InsertFlashcard(fc); err != nil {
		m.status.SetError(err.Error())
		return
	}
	m.status.SetSuccess("Flashcard created")
	m.reload()
	m.view = adminList
}

func (m *AdminModel) updateFlashcard() {
	days, err := m.parseRevisitDays()
	if err != nil {
		m.status.SetError(err.Error())
		return
	}
	fc := m.flashcards[m.selected]
	fc.Question = m.questionInput.Value()
	fc.Answer = m.answerInput.Value()
	fc.RevisitIn = days
	if err := m.storeRef.UpdateFlashcardFull(fc); err != nil {
		m.status.SetError(err.Error())
		return
	}
	m.status.SetSuccess("Flashcard updated")
	m.reload()
	m.view = adminList
}

func (m *AdminModel) deleteFlashcard() {
	id := m.flashcards[m.selected].ID
	if err := m.storeRef.DeleteFlashcard(id); err != nil {
		m.status.SetError(err.Error())
		return
	}
	m.status.SetSuccess(fmt.Sprintf("Deleted flashcard %d", id))
	m.reload()
	if m.selected >= len(m.flashcards) {
		m.selected = len(m.flashcards) - 1
	}
	m.view = adminList
}

func (m *AdminModel) reload() {
	list, err := m.storeRef.GetAllFlashcards()
	if err != nil {
		m.status.SetError(err.Error())
		return
	}
	m.flashcards = list

	// Update table with new data
	rows := makeTableRows(m.flashcards, m.table.Columns())
	m.table.SetRows(rows)

	// Adjust cursor if needed
	if m.selected >= len(m.flashcards) && len(m.flashcards) > 0 {
		m.selected = len(m.flashcards) - 1
		m.table.SetCursor(m.selected)
	}
	if m.selected < 0 {
		m.selected = 0
		m.table.SetCursor(0)
	}
}

// bulkResetRevisitIn resets RevisitIn to 0 for all flashcards
func (m *AdminModel) bulkResetRevisitIn() {
	count := 0
	for _, fc := range m.flashcards {
		fc.RevisitIn = 0
		if err := m.storeRef.UpdateFlashcard(fc); err != nil {
			m.status.SetError(fmt.Sprintf("Error updating flashcard %d: %v", fc.ID, err))
			m.view = adminList
			return
		}
		count++
	}
	m.status.SetSuccess(fmt.Sprintf("Reset RevisitIn to 0 for %d flashcards", count))
	m.reload()
	m.view = adminList
}

// cycleFocus switches focus Question -> Answer -> Revisit -> Question
func (m *AdminModel) cycleFocus() {
	if m.questionInput.Focused() {
		m.questionInput.Blur()
		m.answerInput.Focus()
		m.revisitInput.Blur()
		return
	}
	if m.answerInput.Focused() {
		m.answerInput.Blur()
		m.revisitInput.Focus()
		m.questionInput.Blur()
		return
	}
	// default or revisit focused -> go to question
	m.revisitInput.Blur()
	m.questionInput.Focus()
	m.answerInput.Blur()
}
