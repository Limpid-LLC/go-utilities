package utilities

import "sort"

var Sorter *SorterUtility

type SorterUtility struct{}

func InitSorter() {
	Sorter = &SorterUtility{}
}

func (util *SorterUtility) Sort(data []map[string]interface{}, key string, direction int) []map[string]interface{} {
	sort.Slice(data, func(i, j int) bool {
		if direction < 0 {
			return data[i][key].(float64) > data[j][key].(float64)
		} else {
			return data[i][key].(float64) <= data[j][key].(float64)
		}
	})

	return data
}
