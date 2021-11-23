package client

import "time"

type Log struct {
	rows []LogRow
}

func NewLog() *Log {
	return &Log{}
}

type LogRow struct {
	Time    time.Time
	Payload Payload
	Success bool
}

func (l *Log) TotalAttempts() int {
	return len(l.rows)
}

func (l *Log) SuccessCount() int {
	res := 0
	for _, row := range l.rows {
		if row.Success {
			res += 1
		}
	}

	return res
}

func (l *Log) OldestTime() time.Time {
	oldest := time.Date(9999, 12, 31, 9, 0, 0, 0, time.UTC)

	for _, row := range l.rows {
		if row.Time.Before(oldest) {
			oldest = row.Time
		}
	}

	return oldest
}

func (l *Log) LatestTime() time.Time {
	latest := time.Time{}

	for _, row := range l.rows {
		if row.Time.After(latest) {
			latest = row.Time
		}
	}

	return latest
}

func AggregateLog(logs []Log) *Log {
	var rows []LogRow
	for _, log := range logs {
		rows = append(rows, log.rows...)
	}

	return &Log{rows: rows}
}
