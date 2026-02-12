package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	input := "MCMXCIV"
	ans := romanToInt(input)
	fmt.Printf("Integer value of %s is %d", input, ans)

	/*------------------------------------------------------------------------------------------*/

	var input1 string
	fmt.Print("Input1: ")
	fmt.Scanln(&input1)

	char := strings.Split(input1, ",")
	p, err := strconv.ParseFloat(char[0], 64)
	if err != nil {
		fmt.Println("Principal should be a numeric value")
	}

	n, err := strconv.ParseFloat(char[1], 64)
	if err != nil {
		fmt.Println("Period should be a numeric value")
	}

	r, err := strconv.ParseFloat(char[2], 64)
	if err != nil {
		fmt.Println("Rate should be a numeric value")
	}

	si := roundToTwo(getSimpleInterest(p, n, r))
	fmt.Printf("Simple Interest is " + si + "\n")

	/*-------------------------------------------------------------------------------------*/

	var input2 string
	fmt.Print("Enter Radius: ")
	fmt.Scanln(&input2)

	radius, err := strconv.ParseFloat(input2, 64)
	if err != nil {
		fmt.Println("Radius should be a numeric value")
	}

	area := roundToTwo(getAreaOfCircle(radius))
	fmt.Printf("Area of circle with radius " + roundToTwo(radius) + " is " + area + "\n")
}
