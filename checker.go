package main

var Checker *CheckerUtility

type CheckerUtility struct{}

func (util *CheckerUtility) IsStringExistsInSlice(target string, slice []string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
