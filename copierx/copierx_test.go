package copierx

import (
	"database/sql"
	"testing"
	"time"
)

type Create struct {
	Title      string
	CreateTime time.Time
	TestTime   string
}

type IntCreate struct {
	Title      string
	CreateTime string
	TestTime   sql.NullTime
}

func TestCopy(t *testing.T) {

	c := &Create{
		Title:      "aa",
		CreateTime: time.Now(),
		TestTime:   "2023-02-02 12:22:33",
	}

	t.Log(c)

	i := &IntCreate{}

	Copy(i, c)

	t.Error(i)
}
