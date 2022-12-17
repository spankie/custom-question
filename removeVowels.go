package main

func removeVowels(strings []string) []string {
	var result []string
	for _, s := range strings {
		var newString string
		for _, c := range s {
			if c != 'a' && c != 'e' && c != 'i' && c != 'o' && c != 'u' && c != 'A' && c != 'E' && c != 'I' && c != 'O' && c != 'U' {
				newString += string(c)
			}
		}
		result = append(result, newString)
	}
	return result
}
