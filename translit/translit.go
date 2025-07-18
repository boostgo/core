package translit

func LatinByCyrillic(text string) string {
	var result string

	for _, char := range text {
		if replacement, exists := cyrillicSchema[string(char)]; exists {
			result += replacement
		} else {
			result += string(char)
		}
	}

	return result
}
