package utilities

var Checker *checkerUtility

type checkerUtility struct{}

func InitChecker() {
	Checker = &checkerUtility{}
}

func (util *checkerUtility) IsStringExistsInSlice(target string, slice []string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
