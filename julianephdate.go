package julianephdate

import (
	"math"
	"time"
)

func Date(t time.Time) float64 {
	// Ensure we're working in UTC
	t = t.UTC()

	// Extract date/time components
	year, month, day := t.Date()
	hour := t.Hour()
	min := t.Minute()
	sec := t.Second()
	nsec := t.Nanosecond()

	// Convert time-of-day to fraction of a day
	dayFraction := (float64(hour) +
		float64(min)/60.0 +
		(float64(sec)+float64(nsec)/1e9)/3600.0) / 24.0

	// If month is January or February, treat them as months 13, 14 of previous year
	y := year
	m := int(month)
	if m <= 2 {
		y -= 1
		m += 12
	}

	// The Julian Date (JD) formula (for the start of the day)
	A := math.Floor(float64(y) / 100.0)
	B := 2.0 - A + math.Floor(A/4.0)

	// JD at 0h UT
	jd0 := math.Floor(365.25*float64(y+4716)) +
		math.Floor(30.6001*float64(m+1)) +
		float64(day) + B - 1524.5

	// Add fraction of the day to get precise JD in UTC
	jdUTC := jd0 + dayFraction

	// ---------------------------------------------------------
	// Convert from JD(UTC) -> JD(TT) -> Julian Ephemeris Date.
	// TT = UTC + (leapSeconds) + 32.184 seconds
	// As of this writing, the leap-second offset is 37 seconds.
	// If this changes in the future or for historical times,
	// adjust accordingly or look up the correct offset.
	// ---------------------------------------------------------
	leapSeconds := 37.0
	ttOffset := 32.184 + leapSeconds // total offset in seconds
	jed := jdUTC + ttOffset/86400.0

	return jed
}

func StdTime(jed float64) time.Time {
	// 1. Subtract the offset between UTC and TT from JED
	//    to get the Julian Date (JD) in UTC.
	//
	//    TT = UTC + leapSeconds + 32.184 seconds
	//    => JD(UTC) = JED - (leapSeconds + 32.184)/86400
	leapSeconds := 37.0                 // Hardcoded leap seconds (as of now)
	ttOffsetSec := 32.184 + leapSeconds // in seconds
	jdUTC := jed - ttOffsetSec/86400.0

	// 2. Convert JD(UTC) to a calendar date/time in UTC.
	//    We'll use a standard algorithm that works for JD >= 0.
	//    If you need to handle negative JDs (dates before 4713 BCE),
	//    additional caution is required.

	// Use the common "JD -> Gregorian" algorithm:
	//   a) Shift by 0.5 so day starts at midnight
	//   b) Split integer and fractional part
	//   c) Convert integer part to year, month, day
	//   d) Convert fractional part to hour, minute, second
	//
	// NOTE: The formula changes slightly if JD < 2299160.5 (Julian vs. Gregorian),
	// but for modern epochs, we can safely assume Gregorian.

	// Shift by 0.5 to align day boundaries
	z := math.Floor(jdUTC + 0.5)
	f := (jdUTC + 0.5) - z // fractional part

	// Depending on historical date, you may need to handle the Julian calendar.
	// We'll assume everything is post-Gregorian reform for simplicity.
	alpha := math.Floor((z - 1867216.25) / 36524.25)
	a := z + 1 + alpha - math.Floor(alpha/4)
	b := a + 1524
	c := math.Floor((b - 122.1) / 365.25)
	d := math.Floor(365.25 * c)
	e := math.Floor((b - d) / 30.6001)

	// Extract integer day
	dayFloat := b - d - math.Floor(30.6001*e) + f
	dayInt := math.Floor(dayFloat)
	dayFrac := dayFloat - dayInt

	// Derive month and year
	var month, year float64
	if e < 14 {
		month = e - 1
	} else {
		month = e - 13
	}

	if month > 2 {
		year = c - 4716
	} else {
		year = c - 4715
	}

	// 3. Convert the fractional day into hours, minutes, seconds
	fracDays := dayFrac
	hours := fracDays * 24.0
	hourInt := math.Floor(hours)
	fracHours := hours - hourInt

	minutes := fracHours * 60.0
	minInt := math.Floor(minutes)
	fracMinutes := minutes - minInt

	seconds := fracMinutes * 60.0
	secInt := math.Floor(seconds)
	fracSeconds := seconds - secInt

	// Convert fractional seconds into nanoseconds
	nanoSeconds := math.Round(fracSeconds * 1e9)

	// 4. Construct a Go time.Time in UTC.
	//    We have year, month, day, hour, minute, second, nanosecond.
	//    Note that Go's time.Month is 1-based (January=1, etc.).
	//    We'll clamp the nanoSeconds to avoid floating rounding issues.
	return time.Date(
		int(year),
		time.Month(int(month)),
		int(dayInt),
		int(hourInt),
		int(minInt),
		int(secInt),
		int(nanoSeconds),
		time.UTC,
	)
}
