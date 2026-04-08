package store

import (
	"context"
)

func (s *Store) SessionEventCounts(ctx context.Context) (map[string]int64, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT event_type, COUNT(*)::bigint FROM session_events GROUP BY event_type ORDER BY event_type
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	m := make(map[string]int64)
	for rows.Next() {
		var t string
		var n int64
		if err := rows.Scan(&t, &n); err != nil {
			return nil, err
		}
		m[t] = n
	}
	return m, rows.Err()
}
