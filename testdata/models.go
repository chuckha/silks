package testdata

import (
	"time"
)

// User.tablename=users

type User struct {
	ID        string
	CreatedOn time.Time `slk:"created_on"`
	UserName  string    `slk:"user_name"`
}

// Notification.tablename=notifications

type Notification struct {
	ID        string
	UserID    string    `slk:"user_id"`
	CreatedOn time.Time `slk:"created_on"`
	Seen      bool
}
