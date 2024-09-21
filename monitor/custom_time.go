package monitor

import "time"

type CustomTime struct {
	time.Time
}

const ctLayout = "2006-01-02T15:04:05"

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	// Remove quotes
	if len(s) > 2 {
		s = s[1 : len(s)-1]
	}
	t, err := time.Parse(ctLayout, s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}
