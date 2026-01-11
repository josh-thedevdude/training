package main

import "fmt"

func main() {
	m := map[string]int{"I": 1, "V": 5, "X": 10, "L": 50, "C": 100, "D": 500, "M": 1000}

	input := "MCMXCIV"

	n := len(input)
	currElement := string(input[n-1])
	currValue := m[string(currElement)]

	for i := n - 2; i >= 0; i-- {
		c := string(input[i])
		if m[c] < m[currElement] {
			currValue -= m[c]
		} else {
			currValue += m[c]
		}
		currElement = c
	}
	fmt.Println(currValue)
}
