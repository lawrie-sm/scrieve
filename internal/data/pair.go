package data

import (
	"time"
)

// Pair is the model for pairs of long and short URLs in the DB
type Pair struct {
	Token     string
	Target    string
	TimesUsed int
	LastUsed  time.Time
	CreatedAt time.Time
}

// CreatePair creates a new pair in the DB
func (db *DB) CreatePair(target string) (p Pair, err error) {
	token, err := GenToken()
	if err != nil {
		return
	}
	p = Pair{
		Token:     token,
		Target:    target,
		TimesUsed: 0,
		LastUsed:  time.Now(),
		CreatedAt: time.Now(),
	}
	q := `	INSERT INTO pairs
		(token, target, times_used, last_used, created_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(token) DO UPDATE SET
			token=excluded.token,
			target=excluded.target,
			times_used=excluded.times_used,
			last_used=excluded.last_used,
			created_at=excluded.created_at`
	_, err = db.pool.Exec(q,
		p.Token, p.Target, p.TimesUsed, p.LastUsed, p.CreatedAt)
	return

}

// GetTarget returns a target, given a pair's token
func (db *DB) GetTarget(token string) (target string, err error) {
	q := `	SELECT target
		FROM pairs
		WHERE token = ?`
	row := db.pool.QueryRow(q, token)
	err = row.Scan(&target)

	// Update the pair if it exists
	if err == nil {
		q := `	UPDATE pairs
			SET times_used = times_used + 1, last_used = ?
			WHERE token = ?`
		_, err = db.pool.Exec(q, time.Now(), token)
	}

	return
}
