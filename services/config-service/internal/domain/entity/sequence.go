package entity

import (
	"fmt"
	"strings"
	"time"
)

type Sequence struct {
	ID            string
	CountryCode   string
	SequenceType  string
	Prefix        string
	CurrentValue  int64
	FormatPattern string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (s Sequence) FormatValue(value int64) string {
	result := s.FormatPattern
	result = strings.ReplaceAll(result, "{prefix}", s.Prefix)
	result = strings.ReplaceAll(result, "{value}", fmt.Sprintf("%06d", value))
	return result
}
