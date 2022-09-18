package set

func MakeSet(initialMap map[string]bool) []string {
	list := make([]string, len(initialMap))
	i := 0
	for k := range initialMap {
		list[i] = k
		i++
	}
	return list
}
