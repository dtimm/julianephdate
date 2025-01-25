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

// leapSecondEntry holds a date and the total (TAI - UTC) in effect
// starting right at that instant.
type leapSecondEntry struct {
	effectiveUTC time.Time
	taiMinusUTC  int // in whole seconds
}

// Full table of TAI−UTC from 1972 through the last announced leap second.
// The date/time here is the instant *after* the leap second occurs, i.e.,
// "YYYY-MM-DD 00:00:00 UTC" on the day it first takes effect (except for
// 1972-01-01 which is the moment leap seconds began).
var leapSecondsTable = []leapSecondEntry{
	// The first leap-second adjustment took effect on 1972-01-01, at which
	// time TAI was already 10s ahead of UTC:
	{time.Date(1972, time.January, 1, 0, 0, 0, 0, time.UTC), 10},

	// 1972-07-01 => TAI−UTC = 11
	{time.Date(1972, time.July, 1, 0, 0, 0, 0, time.UTC), 11},

	// 1973-01-01 => 12
	{time.Date(1973, time.January, 1, 0, 0, 0, 0, time.UTC), 12},

	// 1974-01-01 => 13
	{time.Date(1974, time.January, 1, 0, 0, 0, 0, time.UTC), 13},

	// 1975-01-01 => 14
	{time.Date(1975, time.January, 1, 0, 0, 0, 0, time.UTC), 14},

	// 1976-01-01 => 15
	{time.Date(1976, time.January, 1, 0, 0, 0, 0, time.UTC), 15},

	// 1977-01-01 => 16
	{time.Date(1977, time.January, 1, 0, 0, 0, 0, time.UTC), 16},

	// 1978-01-01 => 17
	{time.Date(1978, time.January, 1, 0, 0, 0, 0, time.UTC), 17},

	// 1979-01-01 => 18
	{time.Date(1979, time.January, 1, 0, 0, 0, 0, time.UTC), 18},

	// 1980-01-01 => 19
	{time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC), 19},

	// 1981-07-01 => 20
	{time.Date(1981, time.July, 1, 0, 0, 0, 0, time.UTC), 20},

	// 1982-07-01 => 21
	{time.Date(1982, time.July, 1, 0, 0, 0, 0, time.UTC), 21},

	// 1983-07-01 => 22
	{time.Date(1983, time.July, 1, 0, 0, 0, 0, time.UTC), 22},

	// 1985-07-01 => 23
	{time.Date(1985, time.July, 1, 0, 0, 0, 0, time.UTC), 23},

	// 1988-01-01 => 24
	{time.Date(1988, time.January, 1, 0, 0, 0, 0, time.UTC), 24},

	// 1990-01-01 => 25
	{time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), 25},

	// 1991-01-01 => 26
	{time.Date(1991, time.January, 1, 0, 0, 0, 0, time.UTC), 26},

	// 1992-07-01 => 27
	{time.Date(1992, time.July, 1, 0, 0, 0, 0, time.UTC), 27},

	// 1993-07-01 => 28
	{time.Date(1993, time.July, 1, 0, 0, 0, 0, time.UTC), 28},

	// 1994-07-01 => 29
	{time.Date(1994, time.July, 1, 0, 0, 0, 0, time.UTC), 29},

	// 1996-01-01 => 30
	{time.Date(1996, time.January, 1, 0, 0, 0, 0, time.UTC), 30},

	// 1997-07-01 => 31
	{time.Date(1997, time.July, 1, 0, 0, 0, 0, time.UTC), 31},

	// 1999-01-01 => 32
	{time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC), 32},

	// 2006-01-01 => 33
	{time.Date(2006, time.January, 1, 0, 0, 0, 0, time.UTC), 33},

	// 2009-01-01 => 34
	{time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC), 34},

	// 2012-07-01 => 35
	{time.Date(2012, time.July, 1, 0, 0, 0, 0, time.UTC), 35},

	// 2015-07-01 => 36
	{time.Date(2015, time.July, 1, 0, 0, 0, 0, time.UTC), 36},

	// 2017-01-01 => 37
	{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC), 37},
}

// leapsAtUTC returns TAI-UTC (in whole seconds) at the given UTC time.
func taiMinusUTCAt(t time.Time) int {
	// Find the largest entry in leapSecondsTable whose effectiveUTC <= t
	// The table is in chronological order.
	out := 0
	for _, entry := range leapSecondsTable {
		if !t.Before(entry.effectiveUTC) {
			out = entry.taiMinusUTC
		} else {
			break
		}
	}
	return out
}

// StdTime approximates the UTC time for a given Julian Ephemeris Date (TT).
// We do a small iteration: first guess, then correct leaps, etc.
func StdTime(jed float64) time.Time {
	// 1) Make an initial guess by using the *latest known* TAI - UTC from the table.
	last := leapSecondsTable[len(leapSecondsTable)-1]
	guessTAIminusUTC := last.taiMinusUTC
	// TT - UTC = (TAI - UTC) + 32.184
	// => JD(UTC) = JED - ( (TAI-UTC) + 32.184 ) / 86400
	jdUTCguess := jed - (float64(guessTAIminusUTC)+32.184)/86400.0
	guessTime := jdToCalendarTime(jdUTCguess)

	// 2) Now see what TAI - UTC actually was at guessTime. Then recalc.
	correctTAIminusUTC := taiMinusUTCAt(guessTime)
	jdUTCCorrect := jed - (float64(correctTAIminusUTC)+32.184)/86400.0

	// 3) Convert that corrected JD(UTC) to a time.Time
	// (For many cases, a single correction step is enough. If you want
	//  super-precise alignment around leap-second boundaries, you could
	//  iterate again until stable.)
	return jdToCalendarTime(jdUTCCorrect)
}

// jdToCalendarTime converts JD in UTC to a Go time.Time in UTC
// using a standard Gregorian formula (valid for JD > 2299160.5).
func jdToCalendarTime(jdUTC float64) time.Time {
	// Shift by 0.5 so day boundaries align at midnight
	z := math.Floor(jdUTC + 0.5)
	f := (jdUTC + 0.5) - z

	alpha := math.Floor((z - 1867216.25) / 36524.25)
	a := z + 1 + alpha - math.Floor(alpha/4)
	b := a + 1524
	c := math.Floor((b - 122.1) / 365.25)
	d := math.Floor(365.25 * c)
	e := math.Floor((b - d) / 30.6001)

	dayFloat := b - d - math.Floor(30.6001*e) + f
	dayInt := math.Floor(dayFloat)
	fracDay := dayFloat - dayInt

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

	// Break down fracDay into H:M:S
	hours := fracDay * 24.0
	hourInt := math.Floor(hours)
	fracHours := hours - hourInt

	minutes := fracHours * 60.0
	minInt := math.Floor(minutes)
	fracMinutes := minutes - minInt

	seconds := fracMinutes * 60.0
	secInt := math.Floor(seconds)
	fracSeconds := seconds - secInt

	nano := math.Round(fracSeconds * 1e9)

	return time.Date(
		int(year),
		time.Month(int(month)),
		int(dayInt),
		int(hourInt),
		int(minInt),
		int(secInt),
		int(nano),
		time.UTC,
	)
}
