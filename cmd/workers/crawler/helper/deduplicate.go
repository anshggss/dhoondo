package workerhelper

func DeDuplicate(inp []string) []string {
	if len(inp) < 2 {
		return inp
	}
	var res []string
	hmap := make(map[string]bool)
	for _, v := range inp {
		hmap[v] = true
	}
	for key := range hmap {
		res = append(res, key)
	}
	return res
}
