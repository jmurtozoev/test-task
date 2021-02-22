package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

func buildSearchQuery(query string, filters map[string]interface{}, allowedFilters map[string]string) (string, []interface{}, error) {
	filterString := "WHERE 1=1"
	var inputArgs []interface{}

	// filter key is the url name of the filter used as the lookup for the allowed filters list
	for filterKey, filterValue := range filters {
		if realFilterName, ok := allowedFilters[filterKey]; ok {
			if filterValue == nil {
				continue
			}

			switch filterKey {
			case "cost_min":
				filterString = fmt.Sprintf("%s AND %s>=?", filterString, realFilterName)
				inputArgs = append(inputArgs, filterValue)
			case "cost_max":
				filterString = fmt.Sprintf("%s AND %s<=?", filterString, realFilterName)
				inputArgs = append(inputArgs, filterValue)
			case "name":
				filterString = fmt.Sprintf("%s AND similarity(%s, ?) >= 0.7 OR %s ILIKE  ? || %s", filterString, realFilterName, realFilterName, "'%'")
				inputArgs = append(inputArgs, filterValue, filterValue)
			default:
				filterString = fmt.Sprintf("%s AND %s ILIKE  ? || %s", filterString, realFilterName, "'%'")
				inputArgs = append(inputArgs, filterValue)
			}
		}
	}

	// template the where clause into the original query and then expand the IN clauses with sqlx
	query, args, err := sqlx.In(fmt.Sprintf(query, filterString), inputArgs...)
	if err != nil {
		return "", nil, err
	}
	// using postgres means we need to rebind the ? bindvars that sqlx.IN creates by default to $ bindvars
	// you can omit this if you are using mysql
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	return query, args, nil
}