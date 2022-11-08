package openttd

const (
	ticksInSecond  = 74
	daysInYear     = 365
	daysInLeapYear = 366
	monthsInYear   = 12

	originalBaseYear = year(1920)
)

type date int32

//type dateFract uint16
//type ticks int32

type year int32

//type month byte
//type day byte

func leapYearsUntil(y year) int {
	if y == 0 {
		return 0
	}
	return int(((y)-1)/4 - ((y)-1)/100 + ((y)-1)/400 + 1)
}

func daysUntil(y year) int {
	return daysInYear*int(y) + leapYearsUntil(y)
}

func daysUntilOriginalBaseYear() int {
	return daysUntil(originalBaseYear)
}
