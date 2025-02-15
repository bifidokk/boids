package main

import "fmt"

func main() {
	m := map[string][]string{
		"group1": {"apple", "banana"},
		"group2": {"carrot"},
	}

	newValues := []string{"banana", "cherry"}
	key := "group1"

	mergeToMap(m, newValues, key)

	fmt.Println(m)
}

func mergeToMap(m map[string][]string, values []string, key string) {
	if _, ok := m[key]; !ok {
		m[key] = values
		return
	}

	unique := make(map[string]struct{})

	for _, v := range m[key] {
		unique[v] = struct{}{}
	}

	for _, v := range values {
		if _, ok := unique[v]; !ok {
			m[key] = append(m[key], v)
		}
	}
}
