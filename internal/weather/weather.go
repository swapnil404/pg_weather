package weather

import "github.com/swapnil404/pg_weather/internal/metrics"

// Condition represents the current weather state
type Condition string

const (
	Sunny     Condition = "sunny"
	Cloudy    Condition = "cloudy"
	Overcast  Condition = "overcast"
	Rain      Condition = "rain"
	Storm     Condition = "storm"
	Fog       Condition = "fog"
	Hurricane Condition = "hurricane"
)

// Result holds the weather condition and the reason for it
type Result struct {
	Condition   Condition
	Reason      string
	Severity    int // 0-100, higher is worse
}

// FromMetrics maps database metrics to a weather condition
func FromMetrics(m metrics.DBMetrics) Result {
	// start with worst condition and work up
	severity := 0
	condition := Sunny
	reason := "All systems healthy"

	// check cache hit rate
	switch {
	case m.CacheHitRate < metrics.CacheHitCloudy:
		condition = Overcast
		reason = "Poor cache hit rate"
		severity += 30
	case m.CacheHitRate < metrics.CacheHitSunny:
		condition = Cloudy
		reason = "Cache hit rate could be better"
		severity += 10
	}

	// check connections
	connPct := m.ConnPercent()
	switch {
	case connPct > metrics.ConnCritical*100:
		condition = Storm
		reason = "Connection pool nearly exhausted"
		severity += 40
	case connPct > metrics.ConnWarning*100:
		condition = Rain
		reason = "High connection count"
		severity += 20
	}

	// lock waits override everything
	if m.LockWaits > 5 {
		condition = Storm
		reason = "Heavy lock contention"
		severity += 30
	} else if m.LockWaits > 0 {
		condition = Rain
		reason = "Lock waits detected"
		severity += 15
	}

	// dead tuples add fog
	if m.DeadTuplesRatio > metrics.DeadTupleWarn && condition == Sunny {
		condition = Fog
		reason = "High dead tuple ratio — run VACUUM"
		severity += 10
	}

	// long queries escalate to storm
	if m.LongestQuerySecs > metrics.LongQueryCrit {
		condition = Storm
		reason = "Query running over 60 seconds"
		severity += 40
	} else if m.LongestQuerySecs > metrics.LongQueryWarn {
		if condition == Sunny || condition == Cloudy {
			condition = Rain
		}
		reason = "Long running query detected"
		severity += 20
	}

	// everything bad at once = hurricane
	if severity >= 80 {
		condition = Hurricane
		reason = "Database under critical stress"
	}

	// cap severity at 100
	if severity > 100 {
		severity = 100
	}

	return Result{
		Condition: condition,
		Reason:    reason,
		Severity:  severity,
	}
}
