package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swapnil404/pg_weather/internal/metrics"
)

// Connect opens a connection pool with prepared statements disabled
func Connect(connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("could not parse connection string: %w", err)
	}

	// disable prepared statements — required for PgBouncer/Supabase pooler
	config.ConnConfig.DefaultQueryExecMode = 2

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	return pool, nil
}

// FetchMetrics runs all health queries and returns a DBMetrics struct
func FetchMetrics(ctx context.Context, pool *pgxpool.Pool) (metrics.DBMetrics, error) {
	var m metrics.DBMetrics

	err := pool.QueryRow(ctx, `
		SELECT coalesce(
			sum(heap_blks_hit) / nullif(sum(heap_blks_hit) + sum(heap_blks_read), 0) * 100,
		0) FROM pg_statio_user_tables
	`).Scan(&m.CacheHitRate)
	if err != nil {
		return m, fmt.Errorf("cache hit query failed: %w", err)
	}

	err = pool.QueryRow(ctx, `
		SELECT count(*),
		(SELECT setting::int FROM pg_settings WHERE name = 'max_connections')
		FROM pg_stat_activity WHERE state = 'active'
	`).Scan(&m.ActiveConns, &m.MaxConns)
	if err != nil {
		return m, fmt.Errorf("connections query failed: %w", err)
	}

	err = pool.QueryRow(ctx, `
		SELECT count(*) FROM pg_stat_activity
		WHERE wait_event_type = 'Lock'
	`).Scan(&m.LockWaits)
	if err != nil {
		return m, fmt.Errorf("lock waits query failed: %w", err)
	}

	err = pool.QueryRow(ctx, `
		SELECT coalesce(
			sum(n_dead_tup)::float / nullif(sum(n_live_tup + n_dead_tup), 0) * 100,
		0) FROM pg_stat_user_tables
	`).Scan(&m.DeadTuplesRatio)
	if err != nil {
		return m, fmt.Errorf("dead tuples query failed: %w", err)
	}

	err = pool.QueryRow(ctx, `
		SELECT coalesce(
			max(extract(epoch FROM (now() - query_start))),
		0) FROM pg_stat_activity
		WHERE state = 'active' AND query_start IS NOT NULL
	`).Scan(&m.LongestQuerySecs)
	if err != nil {
		return m, fmt.Errorf("longest query failed: %w", err)
	}

	return m, nil
}
