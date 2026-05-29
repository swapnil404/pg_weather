package render

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/swapnil404/pg_weather/internal/metrics"
	"github.com/swapnil404/pg_weather/internal/weather"
)

var art = map[weather.Condition]string{
	weather.Sunny: `
    \   /
     .-.
  ― (   ) ―
     '-'
    /   \  `,

	weather.Cloudy: `
    .--.
 .-(    ).
(___.__)__)`,

	weather.Overcast: `
    .--.
 .-(    ).
(___.__)__)
  __ __ __`,

	weather.Rain: `
    .--.
 .-(    ).
(___.__)__)
 ' ' ' ' '
 ' ' ' ' '`,

	weather.Storm: `
    .--.
 .-(    ).
(___.__)__)
 ⚡' '⚡'
 ' ' ' ' '`,

	weather.Fog: `
 _ - _ - _
  _ - _ -
 _ - _ - _
  _ - _ - `,

	weather.Hurricane: `
 @ @ @ @ @
@ .--.  @
@(    ).@
@(___.)@@
 @ @ @ @ @`,
}

var conditionColors = map[weather.Condition]string{
	weather.Sunny:     "#FFD700",
	weather.Cloudy:    "#888888",
	weather.Overcast:  "#666666",
	weather.Rain:      "#4499ff",
	weather.Storm:     "#aa44ff",
	weather.Fog:       "#aaaaaa",
	weather.Hurricane: "#ff4444",
}

// Layout builds the full terminal output string
func Layout(result weather.Result, m metrics.DBMetrics, connStr string) string {
	color := conditionColors[result.Condition]

	artStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#e8e8e8"))

	accentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true)

	conditionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true)

	header := fmt.Sprintf("pgweather — %s\n\n", connStr)

	artStr := artStyle.Render(art[result.Condition])

	stats := fmt.Sprintf(
		"%s %s\n\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s",
		conditionStyle.Render(string(result.Condition)),
		"",
		labelStyle.Render("cache hit:    "),
		valueStyle.Render(fmt.Sprintf("%.1f%%", m.CacheHitRate)),
		labelStyle.Render("connections:  "),
		valueStyle.Render(fmt.Sprintf("%d/%d", m.ActiveConns, m.MaxConns)),
		labelStyle.Render("lock waits:   "),
		valueStyle.Render(fmt.Sprintf("%d", m.LockWaits)),
		labelStyle.Render("dead tuples:  "),
		valueStyle.Render(fmt.Sprintf("%.1f%%", m.DeadTuplesRatio)),
		labelStyle.Render("longest query:"),
		valueStyle.Render(fmt.Sprintf("%.1fs", m.LongestQuerySecs)),
	)

	reason := fmt.Sprintf("\n  %s", accentStyle.Render(result.Reason))

	combined := lipgloss.JoinHorizontal(
		lipgloss.Top,
		artStr,
		"     ",
		stats,
	)

	return header + combined + reason + "\n"
}
