package utilities

var Slice *sliceUtility

type sliceUtility struct{}

func InitSlice() {
	Slice = &sliceUtility{}
}

func (util *sliceUtility) AddStringToSlice(elemToAdd string, slice []string) []string {
	return append(slice, elemToAdd)
}

func (util *sliceUtility) RemoveStringFromSlice(elemToRemove string, slice []string) []string {
	var result []string
	for _, elemSlice := range slice {
		if elemSlice != elemToRemove {
			result = append(result, elemSlice)
		}
	}
	return result
}
