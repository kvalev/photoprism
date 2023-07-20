package txt

import (
	"strings"
	"time"
)

// DateFromFilePath returns a string as time or the zero time instant in case it can not be converted.
func DateFromFilePath(s string) (result time.Time) {
	defer func() {
		if r := recover(); r != nil {
			result = DefaultTime
		}
	}()

	if len(s) < 6 {
		return DefaultTime
	}

	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}

	b := []byte(s)
	handlers := []func([]byte) time.Time{
		DateTimeHandler,
		DateHandler,
		DatePathHandler,
		WhatsAppDateHandler,
	}

	for _, handler := range handlers {
		result = handler(b)
		if !result.IsZero() {
			break
		}
	}

	return result.UTC()
}

// WhatsAppDateHandler checks whether the data contains a datetime like "VID-20191120-WA0001.jpg"
func WhatsAppDateHandler(b []byte) time.Time {
	found := DateWhatsAppRegexp.Find(b)

	if len(found) == 0 {
		return DefaultTime
	}

	match := DateWhatsAppRegexp.FindSubmatch(found)

	if len(match) != 4 {
		return DefaultTime
	}


	matchMap := make(map[string]string)
	for i, name := range DateWhatsAppRegexp.SubexpNames() {
		if i != 0 {
			matchMap[name] = string(match[i])
		}
	}

	year := ExpandYear(matchMap["year"])
	month := Int(matchMap["month"])
	day := Int(matchMap["day"])

	// Perform date plausibility check.
	if year < YearMin || year > YearMax || month < MonthMin || month > MonthMax || day < DayMin || day > DayMax {
		return DefaultTime
	}

	return time.Date(
		year,
		time.Month(month),
		day,
		0,
		0,
		0,
		0,
		time.UTC)
}

// DateTimeHandler checks whether the data contains a datetime like "2020-01-30_09-57-18"
func DateTimeHandler(b []byte) time.Time {
	found := DateTimeRegexp.Find(b)

	if len(found) == 0 {
		return DefaultTime
	}

	n := DateIntRegexp.FindAll(found, -1)

	if len(n) < 6 {
		return DefaultTime
	}

	year := ExpandYear(string(n[0]))
	month := Int(string(n[1]))
	day := Int(string(n[2]))
	hour := Int(string(n[3]))
	min := Int(string(n[4]))
	sec := Int(string(n[5]))

	if !IsValidDate(year, month, day) || !IsValidTime(hour, min, sec) {
		return DefaultTime
	}

	return time.Date(
		year,
		time.Month(month),
		day,
		hour,
		min,
		sec,
		0,
		time.UTC)
}

// DateHandler checks whether the data contains a date like "2020-01-30"
func DateHandler(b []byte) time.Time {
	found := DateRegexp.Find(b)

	if len(found) == 0 {
		return DefaultTime
	}

	n := DateIntRegexp.FindAll(found, -1)

	if len(n) != 3 {
		return DefaultTime
	}

	year := ExpandYear(string(n[0]))
	month := Int(string(n[1]))
	day := Int(string(n[2]))

	if !IsValidDate(year, month, day) {
		return DefaultTime
	}

	return time.Date(
		year,
		time.Month(month),
		day,
		0,
		0,
		0,
		0,
		time.UTC)
}

// DatePathHandler checks whether the data contains a date like "2020/01/03"
func DatePathHandler(b []byte) time.Time {
	found := DatePathRegexp.Find(b)

	if len(found) == 0 {
		return DefaultTime
	}

	n := DateIntRegexp.FindAll(found, -1)

	if len(n) < 2 || len(n) > 3 {
		return DefaultTime
	}

	year := ExpandYear(string(n[0]))
	month := Int(string(n[1]))

	if !IsValidYearAndMonth(year, month) {
		return DefaultTime
	}

	if len(n) == 2 {
		return time.Date(
			year,
			time.Month(month),
			1,
			0,
			0,
			0,
			0,
			time.UTC)
	} else if day := Int(string(n[2])); IsValidDay(day) {
		return time.Date(
			year,
			time.Month(month),
			day,
			0,
			0,
			0,
			0,
			time.UTC)
	}

	return DefaultTime
}

// IsValidTime tests if the given time is plausible
func IsValidTime(hour int, min int, sec int) bool {
	return IsValidHour(hour) && IsValidMinute(min) && IsValidSecond(sec)
}

// IsValidHour tests if the given hour value is plausible
func IsValidHour(hour int) bool {
	return hour >= HourMin && hour <= HourMax
}

// IsValidMinute tests if the given minutes value is plausible
func IsValidMinute(minute int) bool {
	return minute >= MinMin && minute <= MinMax
}

// IsValidSecond tests if the given second value is plausible
func IsValidSecond(second int) bool {
	return second >= SecMin && second <= SecMax
}

// IsValidDate tests if the given date is plausible
func IsValidDate(year int, month int, day int) bool {
	return IsValidYearAndMonth(year, month) && IsValidDay(day)
}

// IsValidYearAndMonth tests if the given year and month are plausible
func IsValidYearAndMonth(year int, month int) bool {
	return IsValidYear(year) && IsValidMonth(month)
}

// IsValidYear tests if the given year is plausible
func IsValidYear(year int) bool {
	return year >= YearMin && year <= YearMax
}

// IsValidMonth tests if the given month is plausible
func IsValidMonth(month int) bool {
	return month >= MonthMin && month <= MonthMax
}

// IsValidDay tests if the given day is plausible
func IsValidDay(day int) bool {
	return day >= DayMin && day <= DayMax
}
