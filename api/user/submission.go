package apiUser

import (
	"time"
)

// Submission 表示一次flag提交记录
type Submission struct {
	ID          int
	UserID      int
	ChallengeID int
	Flag        string
	Time        time.Time
}

// Solve 表示一次正确的flag提交记录
type Solve struct {
	ID           int
	SubmissionID int
}
