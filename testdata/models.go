package testdata

import (
	"time"
)

// User.tablename=users
type User struct {
	ID        string
	CreatedOn time.Time `slk:"created_on"`
	UserName  string    `slk:"user_name"`
	Age       int       `slk:"age"`
}

// Notification.tablename=notifications
type Notification struct {
	Seen      bool
	ID        string
	UserID    string    `slk:"user_id"`
	CreatedOn time.Time `slk:"created_on"`
}
