package opennms

import "iter"

// Paginate yields every item from a paginated list endpoint,
// handling limit/offset pagination transparently and stopping when
// totalCount is reached or a partial page is returned.
//
// fetch is called with successive (limit, offset) pairs and must
// return the raw page object; key is the response key holding the
// list of items (e.g. "alarm" for alarms, "node" for nodes).
//
// Example — fetch every critical alarm without a pagination loop:
//
//	pages := func(limit, offset int) (map[string]any, error) {
//		return client.GetAlarms(ctx,
//			&opennms.ListOptions{Limit: limit, Offset: offset},
//			map[string]string{"severity": "CRITICAL"})
//	}
//	for alarm, err := range opennms.Paginate(pages, "alarm", 100) {
//		if err != nil {
//			return err
//		}
//		fmt.Println(alarm.(map[string]any)["id"])
//	}
func Paginate(fetch func(limit, offset int) (map[string]any, error), key string, pageSize int) iter.Seq2[any, error] {
	return func(yield func(any, error) bool) {
		offset := 0
		for {
			page, err := fetch(pageSize, offset)
			if err != nil {
				yield(nil, err)
				return
			}
			items, _ := page[key].([]any)
			for _, item := range items {
				if !yield(item, nil) {
					return
				}
			}
			offset += len(items)
			if total, ok := page["totalCount"].(float64); ok && offset >= int(total) {
				return
			}
			if len(items) < pageSize {
				return
			}
		}
	}
}
