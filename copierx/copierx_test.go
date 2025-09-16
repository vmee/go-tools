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
	Num        uint64
}

type IntCreate struct {
	Title      string
	CreateTime string
	TestTime   sql.NullTime
	Num        sql.NullInt64
}

func TestCopy(t *testing.T) {

	c := &Create{
		Title:      "aa",
		CreateTime: time.Now(),
		TestTime:   "2023-02-02 12:22:33",
		Num:        123,
	}

	t.Log(c)

	i := &IntCreate{}

	Copy(i, c)

	t.Error(i)
}

func TestCopy2(t *testing.T) {

	c := &IntCreate{
		Title:      "aa",
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		TestTime:   sql.NullTime{Time: time.Now(), Valid: true},
		Num:        sql.NullInt64{Int64: 0, Valid: false},
	}

	t.Log(c)

	i := &Create{}

	Copy(i, c)

	t.Error(i)
}
