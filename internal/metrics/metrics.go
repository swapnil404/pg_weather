package metrics

// DBMetrics holds all raw values fetched from the database
type DBMetrics struct {
	CacheHitRate     float64 // percentage 0-100
	ActiveConns      int     // current active connections
	MaxConns         int     // max_connections setting
	LockWaits        int     // queries waiting on locks
	DeadTuplesRatio  float64 // percentage 0-100
	LongestQuerySecs float64 // longest running query in seconds
}

// Thresholds define what values trigger what weather
const (
	CacheHitSunny   = 95.0 // above this = sunny
	CacheHitCloudy  = 80.0 // above this = cloudy
	ConnWarning     = 0.7  // above 70% of max = warning
	ConnCritical    = 0.9  // above 90% of max = critical
	DeadTupleWarn   = 5.0  // above 5% dead tuples = warning
	LongQueryWarn   = 30.0 // query running longer than 30s = warning
	LongQueryCrit   = 60.0 // query running longer than 60s = critical
)

// ConnPercent returns active connections as a percentage of max
func (m DBMetrics) ConnPercent() float64 {
	if m.MaxConns == 0 {
		return 0
	}
	return float64(m.ActiveConns) / float64(m.MaxConns) * 100
}
