package main

import (
	"fmt"
	"net/url"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/swapnil404/pg_weather/internal/db"
	"github.com/swapnil404/pg_weather/internal/ui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: pgweather <connection-string>")
		fmt.Println("example: pgweather postgresql://user:pass@localhost:5432/mydb")
		os.Exit(1)
	}

	connStr := os.Args[1]
	interval := 3 * time.Second

	// parse just the host for display, hide credentials
	displayConn := connStr
	if u, err := url.Parse(connStr); err == nil {
		displayConn = u.Host
	}

	conn, err := db.Connect(connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	m := ui.New(conn, displayConn, interval)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
