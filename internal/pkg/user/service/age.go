package service

import "time"

func calcAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()
	if now.Month() <= dob.Month() && now.Day() < dob.Day() {
		return age - 1
	}
	return age
}
