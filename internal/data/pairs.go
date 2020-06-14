package data

import (
	"time"
)

// Pair is the model for pairs of long and short URLs in the DB
type Pair struct {
	Token     Token
	Target    string
	LastUsed  time.Time
	CreatedAt time.Time
}

// CreatePair creates a new pair in the DB
func (db *DB) CreatePair(target string) (p Pair, err error) {
	token, err := GenToken()
	p = Pair{
		Token:     token,
		Target:    target,
		LastUsed:  time.Now(),
		CreatedAt: time.Now(),
	}
	q := `	INSERT INTO pairs
		(token, target, last_used, created_at)
		VALUES (?, ?, ?, ?)`
	_, err = db.pool.Exec(q, p.Token, p.Target, p.LastUsed, p.CreatedAt)
	return

}

// GetTarget returns a target, given a pair's token
func (db *DB) GetTarget(token string) (target string, err error) {
	// TODO: Also update last used
	q := `	SELECT target
		FROM pairs
		WHERE token = ?`
	row := db.pool.QueryRow(q, token)
	err = row.Scan(&target)
	return
}
