package bookworm

func Compact(list []string) []string {
	var compactedList []string
	for _, str := range list {
		if str != "" {
			compactedList = append(compactedList, str)
		}
	}

	return compactedList
}
