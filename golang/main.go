package main

import (
	"fmt"
	"os"
	"time"
)

const DATETIME_LAYOUT = "2006-01-02"
const START_DATE = "2564-06-01"
const END_DATE = "2564-08-31"
const SENIOR_LOWER_AGE_YEAR = 65
const CHILD_LOWER_AGE_MONTH = 6
const CHILD_UPPER_AGE_YEAR = 2
var thaiMonths = [...]string{
	"มกราคม", "กุมภาพันธ์", "มีนาคม", "เมษายน", "พฤษภาคม", "มิถุนายน",
	"กรกฎาคม", "สิงหาคม", "กันยายน", "ตุลาคม", "พฤศจิกายน", "ธันวาคม",
}


type Gender int64
const (
	MALE Gender = iota
	FEMALE
)

func (gender Gender) String() string {
	switch gender {
		case MALE:
			return "male"
		case FEMALE:
			return "female"
	}
	return "unknown"
}

type Person struct {
	gender Gender
	birthdate time.Time
}

func (p Person) String() string {
	return fmt.Sprintf("{%s: %v}", p.gender, thaiDateFormat(p.birthdate))
}

type Age struct {
	year int
	month int
}

func (a Age) String() string {
	return fmt.Sprintf("[%dy %dm]", a.year, a.month)
}

// returns age in (years, months) from birthdate relative to the set date.
func (p Person) calcAge(setDate time.Time) Age {
	ageYears := setDate.Year() - p.birthdate.Year()
	setDateMonth := int(setDate.Month())
	birthdateMonth := int(p.birthdate.Month())

	// has passed the set date's birthday for that year
	if setDateMonth < birthdateMonth || 
		(setDateMonth == birthdateMonth && setDate.Day() < p.birthdate.Day()) {
		ageYears--
	}
	
	// prefer positive as modulus operands
	ageMonths := ((setDateMonth - birthdateMonth) + 12) % 12
	return Age{ageYears, ageMonths}
}


// returns if the person is able to apply for the service within the sevice period
// with start and end date of the period when the person can apply for the service.
// For senior citizens (65 years old or older) 
// [65, inf)
// For children (between 6 months and 2 years old) 
// [0.6, 2]
// If the person does not meet any of these criteria, they are deemed ineligible.
// This function will modify the age of the person though
func (p Person) eligible(startDate, endDate time.Time) (bool, time.Time , time.Time) {
	ageStart := p.calcAge(startDate)
	ageEnd := p.calcAge(endDate)
	fmt.Println(ageStart,ageEnd)
	inSeniorRange := ageStart.year >= SENIOR_LOWER_AGE_YEAR
	inChildrenRange := ageStart.month >= CHILD_LOWER_AGE_MONTH && ageEnd.year <= CHILD_UPPER_AGE_YEAR && ageEnd.month < 1
	willBe65yo := ageEnd.year == SENIOR_LOWER_AGE_YEAR
	willBe6m := (ageEnd.month >= CHILD_LOWER_AGE_MONTH) && (ageEnd.year < CHILD_UPPER_AGE_YEAR)
	willBe2yo := ageEnd.year == CHILD_UPPER_AGE_YEAR
	
	if inSeniorRange || inChildrenRange {
		return true, startDate, endDate
	} else if willBe65yo {
		return true, p.birthdate.AddDate(65, 0, 0), endDate
	} else if willBe6m {
		return true, p.birthdate.AddDate(0, 6, 0), endDate
	} else if willBe2yo {
		return true, startDate, p.birthdate.AddDate(2, 0, 0)
	}

	// ineligible
	return false, time.Time{}, time.Time{}
}

// Just to return nil
func (p Person) wrapperEligible() (bool, *time.Time, *time.Time) {
	startDate, err  := time.Parse(DATETIME_LAYOUT, START_DATE)
	handleErr(err)
	endDate, err := time.Parse(DATETIME_LAYOUT, END_DATE)
	handleErr(err)
	eligible, start, end := p.eligible(startDate, endDate)
	if start.IsZero() && end.IsZero() {
		return eligible, nil, nil 
	}
	return eligible, &start, &end
}

func main() {
	bdates := [...]string {
		"2499-03-10", // 10 มีนาคม พ.ศ.2499
		"2500-10-08", // 8 ตุลาคม พ.ศ.2500 
		"2562-07-01", // 1 กรกฎาคม พ.ศ.2562 
		"2564-01-05", // 5 มกราคม พ.ศ.2564
	}
	gens := [...]Gender { FEMALE, MALE, FEMALE, FEMALE, }
	var persons []Person
	
	for i := 0; i < len(bdates); i++ {
		date, err := time.Parse(DATETIME_LAYOUT, bdates[i])
		handleErr(err)
		persons = append(persons, Person{gens[i], date})
	}
	
	for i := 0; i < len(persons); i++ {
		fmt.Printf("%d: %s\n", i+1, persons[i])
		eligible, start, end := persons[i].wrapperEligible()
		if start != nil && end != nil {
			fmt.Println(eligible, thaiDateFormat(*start), thaiDateFormat(*end))
		} else {
			fmt.Println(eligible, start, end)
		}
		fmt.Println("----------")
	}
}

func handleErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func thaiDateFormat(date time.Time) string {
	return fmt.Sprintf("%d %s พ.ศ.%d", date.Day(), thaiMonths[date.Month()-1], date.Year())
}