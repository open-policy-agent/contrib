package esquery

// Source represents the "_source" option which is commonly accepted in ES
// queries. Currently, only the "includes" option is supported.
type Source struct {
	includes []string
	excludes []string
}

// Map returns a map representation of the Source object.
func (source Source) Map() map[string]interface{} {
	m := make(map[string]interface{})
	if len(source.includes) > 0 {
		m["includes"] = source.includes
	}
	if len(source.excludes) > 0 {
		m["excludes"] = source.excludes
	}
	return m
}

// Sort represents a list of keys to sort by.
type Sort []map[string]interface{}

// Order is the ordering for a sort key (ascending, descending).
type Order string

const (
	// OrderAsc represents sorting in ascending order.
	OrderAsc Order = "asc"

	// OrderDesc represents sorting in descending order.
	OrderDesc Order = "desc"
)
