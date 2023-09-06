package main

type SliceUtility struct{}

func (util *SliceUtility) AddStringToSlice(elemToAdd string, slice []string) []string {
	return append(slice, elemToAdd)
}

func (util *SliceUtility) RemoveStringFromSlice(elemToRemove string, slice []string) []string {
	var result []string
	for _, elemSlice := range slice {
		if elemSlice != elemToRemove {
			result = append(result, elemSlice)
		}
	}
	return result
}
