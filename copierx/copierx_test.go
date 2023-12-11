package copierx

import (
	"testing"
	"time"
)

type Create struct {
	Title      string
	CreateTime time.Time
}

type IntCreate struct {
	Title      string
	CreateTime string
}

func TestCopy(t *testing.T) {

	c := &Create{
		Title:      "aa",
		CreateTime: time.Now(),
	}

	t.Log(c)

	i := &IntCreate{}

	Copy(i, c)

	t.Error(i)
}
