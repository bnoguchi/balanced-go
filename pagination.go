package balanced

import "fmt"

type PaginationParams struct {
	Limit, Offset, Total int
	First, Href, Last    string
	Next, Previous       interface{}
}

func NewPaginationParams(meta map[string]interface{}) *PaginationParams {
	return &PaginationParams{
		Limit:  int(meta["limit"].(float64)),
		Offset: int(meta["offset"].(float64)),
		Total:  int(meta["total"].(float64)),

		First: meta["first"].(string),
		Href:  meta["href"].(string),
		Last:  meta["last"].(string),

		Next:     meta["next"],
		Previous: meta["previous"],
	}
}

func paginatedArgsToQuery(args []interface{}) map[string]interface{} {
	params := make(map[string]interface{})
	numInts := 0
	for _, arg := range args {
		switch arg := arg.(type) {
		case int:
			numInts++
			if numInts == 1 {
				params["offset"] = arg
			} else if numInts == 2 {
				params["limit"] = arg
			}
		case map[string]interface{}:
			for k, v := range arg {
				params[k] = v
			}
		default:
			fmt.Printf("unexpected type %T", arg)
		}
	}
	return params
}
