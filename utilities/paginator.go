package utilities

type PaginatorUtility struct{}

func (util *PaginatorUtility) Paginate(data []map[string]interface{}, skipPassed interface{}, limitPassed interface{}) []map[string]interface{} {
	if skipPassed != nil || limitPassed != nil {
		limit := 20
		from := 0

		if skipPassed != nil {
			from = int(skipPassed.(float64))
		}
		if limitPassed != nil {
			limit = int(limitPassed.(float64))
		}

		to := from + limit

		if len(data) < from {
			return []map[string]interface{}{}
		}

		if len(data) < to {
			return data[from:]
		}

		return data[from:to]
	}

	return data
}
