package utilities

import "sort"

var Sorter *sorterUtility

type sorterUtility struct{}

func InitSorter() {
	Sorter = &sorterUtility{}
}

func (util *sorterUtility) Sort(data []map[string]interface{}, key string, direction int) []map[string]interface{} {
	sort.Slice(data, func(i, j int) bool {
		if direction < 0 {
			return data[i][key].(float64) > data[j][key].(float64)
		} else {
			return data[i][key].(float64) <= data[j][key].(float64)
		}
	})

	return data
}
