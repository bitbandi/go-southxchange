package southxchange

import (
	"time"
	"strconv"
)

type Timestamp time.Time

func (t Timestamp) String() string {
	return t.Format("2006-01-02T15:04:05.999999999")
}

func (t Timestamp) Format(layout string) string {
	return time.Time(t).Format(layout)
}

func (t *Timestamp) UnmarshalJSON(body []byte) (err error) {
	s, err := strconv.Unquote(string(body))
	if err != nil {
		return err
	}
	if *(*time.Time)(t), err = time.Parse("2006-01-02T15:04:05.999999999", s); err != nil {
		return err
	}
	return nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(time.Time(t).Format("2006-01-02T15:04:05.999999999"))), nil
}
