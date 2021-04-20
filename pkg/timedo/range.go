package timedo

import "time"

func Today(now time.Time) (from, to time.Time) {
	from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	to = now
	return
}

func Week(now time.Time, turn bool) (from, to time.Time) {
	from = now.AddDate(0, 0, -7)
	if turn {
		_, offset := from.Zone()
		from = from.Truncate(24 * time.Hour).Add(-1 * time.Duration(offset) * time.Second)
	}
	to = now
	return
}

func Month(now time.Time) (from, to time.Time) {
	from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	to = now
	return
}

func LastMonth(now time.Time) (from, to time.Time) {
	from = time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
	to = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Add(-1 * time.Millisecond)
	return
}

func LastDay(now time.Time) (from, to time.Time) {
	from = time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
	to = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(-1 * time.Millisecond)
	return
}

func LastHour(now time.Time) (from, to time.Time) {
	from = now.Add(-1 * time.Hour)
	to = now
	return
}

func DateDuration(formdate, todate string, locs ...*time.Location) (from, to time.Time, err error) {
	var loc *time.Location
	if len(locs) > 0 && locs[0] != nil {
		loc = locs[0]
	} else {
		loc = time.Now().Location()
	}
	from, err = time.ParseInLocation("2006-01-02", formdate, loc)
	if err != nil {
		return
	}
	to, err = time.ParseInLocation("2006-01-02", todate, loc)
	if err != nil {
		return
	}
	to = to.Add(24*time.Hour - time.Millisecond)
	return
}

func SplitDayList(from, to time.Time) (list [][]time.Time) {
	var pair []time.Time
	var f, t time.Time
	for {
		pair = make([]time.Time, 0)
		f, t = Today(to)
		if f.Sub(from) >= 0 {
			pair = append(pair, f)
			pair = append(pair, t)
			list = append(list, pair)
		} else {
			break
		}
		to = f.Add(-1 * time.Millisecond)
	}
	return
}

func Intersection(s1 time.Time, d1 time.Duration, s2 time.Time, d2 time.Duration) (yes bool) {
	if t1 := s2.Sub(s1); t1 >= 0 && t1 <= d1 {
		return true
	}
	if t2 := s1.Sub(s2); t2 >= 0 && t2 <= d2 {
		return true
	}
	return
}
