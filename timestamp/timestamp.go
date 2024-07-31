package timestamp

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"strconv"
	"time"
)

type Timestamp time.Time

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t).Unix()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	*t = Timestamp(time.Unix(int64(ts), 0))

	return nil
}

func (t *Timestamp) GetBSON() (interface{}, error) {
	if time.Time(*t).IsZero() {
		return nil, nil
	}

	return time.Time(*t), nil
}

func (t *Timestamp) SetBSON(raw bson.Raw) error {
	var tm time.Time

	if err := raw.Unmarshal(&tm); err != nil {
		return err
	}

	*t = Timestamp(tm)

	return nil
}

func (t *Timestamp) String() string {
	return time.Time(*t).String()
}
