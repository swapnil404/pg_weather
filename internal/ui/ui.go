package ui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swapnil404/pg_weather/internal/db"
	"github.com/swapnil404/pg_weather/internal/metrics"
	"github.com/swapnil404/pg_weather/internal/render"
	"github.com/swapnil404/pg_weather/internal/weather"
)

type tickMsg time.Time

type model struct {
	conn     *pgxpool.Pool
	metrics  metrics.DBMetrics
	result   weather.Result
	connStr  string
	interval time.Duration
	err      error
}

func New(conn *pgxpool.Pool, connStr string, interval time.Duration) model {
	return model{
		conn:     conn,
		connStr:  connStr,
		interval: interval,
	}
}

func tick(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tick(m.interval)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case tickMsg:
		metrics, err := db.FetchMetrics(context.Background(), m.conn)
		if err != nil {
			m.err = err
			return m, tick(m.interval)
		}
		m.metrics = metrics
		m.result = weather.FromMetrics(metrics)
		m.err = nil
		return m, tick(m.interval)
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return "error fetching metrics: " + m.err.Error() + "\n"
	}

	if m.metrics.MaxConns == 0 {
		return "pgweather — connecting...\n"
	}

	return render.Layout(m.result, m.metrics, m.connStr)
}
