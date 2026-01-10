package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func roundToTwo(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

/*
First-line has the comma-separated values of the Principal, rate and time (in years) respective
*constraints: Round simple interest float value to 2 decimal places
*/
func getSimpleInterest(p, n, r float64) (si float64) {
	si = (p * n * r) / 100
	return
}

/*
calculate the area of the circle, First line has a value of the radius of the circle
constraint
1. Use const PI from the package math package
2. Use the Pow function from the math package
3. Round area float value to 2 decimal places
*/
func getAreaOfCircle(r float64) float64 {
	return math.Pi * math.Pow(r, 2)
}

func main() {
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
