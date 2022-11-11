package artifacts_table

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/massdriver-cloud/massdriver-cli/pkg/api2"
)

const checked = "✓"
const unchecked = ""

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

type model struct {
	table      table.Model
	selected   map[string]bool
	rows       []table.Row
	quitting   bool
	sourceData []api2.Artifact
}

func (m model) Init() tea.Cmd {
	return nil
}

func buildRows(selected map[string]bool, artifacts []api2.Artifact) []table.Row {
	rows := []table.Row{}

	for _, artifact := range artifacts {
		selectedIndicator := unchecked
		// TODO: DRY up a selectable table
		if _, ok := selected[artifact.Name]; ok {
			selectedIndicator = checked
		}

		// TODO: DRY up a selectable table
		cloudName := artifact.Name
		row := table.Row{selectedIndicator, cloudName, artifact.Name}
		rows = append(rows, row)
	}

	return rows
}

func (m model) toggleSelectedRow() {
	// TODO: DRY up a selectable table
	selectedArtifact := m.table.SelectedRow()[2]

	if m.selected[selectedArtifact] {
		delete(m.selected, selectedArtifact)
	} else {
		m.selected[selectedArtifact] = true
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.quitting = true
			return m, tea.Quit
		case "s":
			// Save and continue
			return m, tea.Quit
		case "enter":
			m.toggleSelectedRow()
		}
	}

	m.table, cmd = m.table.Update(msg)
	m.table.SetRows(buildRows(m.selected, m.sourceData))

	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		promptStyle("Which credential do you want to use for authenticating?\n\n%s\n\n%s"),
		baseStyle.Render(m.table.View()),
		m.helpView(),
	) + "\n\n"
}

func (m model) helpView() string {
	return helpStyle("\n  ↑/↓: Navigate • esc: Quit • enter: Select • s: Save\n")
}

func New(sourceArtifacts []api2.Artifact) ([]api2.Artifact, error) {
	columns := []table.Column{
		{Title: checked, Width: 3},
		{Title: "Cloud", Width: 10},
		{Title: "Type", Width: 40},
	}

	selectedArtifacts := []api2.Artifact{}
	selected := map[string]bool{}
	rows := buildRows(map[string]bool{}, sourceArtifacts)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{
		table:      t,
		selected:   selected,
		rows:       rows,
		sourceData: sourceArtifacts,
	}

	out, err := tea.NewProgram(m).Run()

	if err != nil {
		return selectedArtifacts, err
	}

	if out, ok := out.(model); ok {
		if out.quitting {
			os.Exit(0)
		}

		for _, row := range out.rows {
			// TODO: DRY up a selectable table
			typeName := row[2]
			if m.selected[typeName] {
				// TODO: DRY up a selectable table
				def := api2.Artifact{Name: typeName}
				selectedArtifacts = append(selectedArtifacts, def)
			}
		}

		return selectedArtifacts, nil
	}

	return selectedArtifacts, errors.New("failed to cast model")
}
