package user

import "time"

type Record struct {
	ID             string     `db:"id"`
	Email          string     `db:"email"`
	Nickname       string     `db:"nickname"`
	HashedPassword string     `db:"hashed_password"`
	Verified       bool       `db:"verified"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at"`
}
