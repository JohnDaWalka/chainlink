package addressbook

func filterAll[E any](element E, filters []func(E) bool) bool {
	do := true
	for _, f := range filters {
		do = do && f(element)
	}
	return do
}

// transform performs a mapping operation on a slice, producing a slice of the output type of the mapping function.
func transform[E any, O any](slice []E, f func(E) O) []O {
	mapped := make([]O, len(slice))
	for i, v := range slice {
		mapped[i] = f(v)
	}
	return mapped
}
