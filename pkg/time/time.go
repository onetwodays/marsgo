package time

import (
	"context"
	"database/sql/driver"
	"strconv"
	xtime "time"
)

// Time be used to MySql timestamp converting.
type Time int64

// Scan scan time. 把time.Time和string 表示的时间戳转整数时间戳
// 把src的值读到jt里面
func (jt *Time) Scan(src interface{}) (err error) {
	switch sc := src.(type) {
	case xtime.Time:
		*jt = Time(sc.Unix()) //linux时间1970的秒数
	case string:
		var i int64
		i, err = strconv.ParseInt(sc, 10, 64)
		*jt = Time(i)
	}
	return
}

// Value get time value.
func (jt Time) Value() (driver.Value, error) {
	return xtime.Unix(int64(jt), 0), nil
}

// Time get time.time.Time =A Time represents an instant in time with nanosecond precision
func (jt Time) Time() xtime.Time {
	return xtime.Unix(int64(jt), 0)
}

// Duration be used toml unmarshal string time, like 1s, 500ms.
type Duration xtime.Duration

// UnmarshalText unmarshal text to duration.ParseDuration parses a duration string
func (d *Duration) UnmarshalText(text []byte) error {
	tmp, err := xtime.ParseDuration(string(text))
	if err == nil {
		*d = Duration(tmp) //类型显示转化
	}
	return err
}

// Shrink(收缩) will decrease(减少) the duration by comparing with context's timeout duration
// and return new timeout\context\CancelFunc. Until = It is shorthand for t.Sub(time.Now())
// 假设context的过期时间>d,新的过期时间缩减为d,否则不变.
// 将c的过期时间缩减为d
func (d Duration) Shrink(c context.Context) (Duration, context.Context, context.CancelFunc) {
	if deadline, ok := c.Deadline(); ok {
		if ctimeout := xtime.Until(deadline); ctimeout < xtime.Duration(d) {
			// deliver small timeout
			return Duration(ctimeout), c, func() {} //不要管
		}
	}
	ctx, cancel := context.WithTimeout(c, xtime.Duration(d)) // 提前到期
	return d, ctx, cancel
}
