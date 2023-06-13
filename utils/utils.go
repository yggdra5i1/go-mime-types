package utils

func Map(coll []string, fn func(string) string) []string {
	mapped := make([]string, len(coll))

	for i, e := range coll {
		mapped[i] = fn(e)
	}

	return mapped
}
