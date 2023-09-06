package utilities

var Checker *CheckerUtility

type CheckerUtility struct{}

func InitChecker() {
	Checker = &CheckerUtility{}
}

func (util *CheckerUtility) IsStringExistsInSlice(target string, slice []string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
