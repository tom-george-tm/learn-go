// Package main is the entry point of the Go program.
// This file demonstrates basic Go concepts including:
//   - Variable declaration (explicit and short-hand)
//   - Functions with single and multiple return values
//   - Error handling using the (result, error) pattern
//   - Formatted output using the fmt package
package main

import (
	"fmt"
)

// add takes two integers and returns their sum.
//
// Parameters:
//   - a: first integer
//   - b: second integer
//
// Returns:
//   - int: the sum of a and b
func add(a int, b int) int {
	return a + b
}

// greet prints a simple greeting message to the console.
// It takes no parameters and returns nothing.
func greet() {
	fmt.Println("Hello!")
}

// divide performs division of two float64 numbers.
// It uses Go's multiple return value feature to return
// both the result and an error (if any).
//
// Parameters:
//   - a: the dividend (number to be divided)
//   - b: the divisor (number to divide by)
//
// Returns:
//   - float64: result of the division (0 if error)
//   - error:   non-nil if b is 0 (cannot divide by zero), nil otherwise
//
// Example:
//
//	result, err := divide(10, 2)  // result = 5.0, err = nil
//	result, err := divide(10, 0)  // result = 0,   err = "cannot divide by zero"
func divide(a, b float64) (float64, error) {
	if b == 0 {
		// Return zero value and a descriptive error
		return 0, fmt.Errorf("cannot divide by zero")
	}
	// Return the result and nil (nil = no error)
	return a / b, nil
}

// main is the entry point of the program.
// Go's runtime automatically calls this function when the program starts.
// All program logic begins here.
func main() {
	fmt.Println("Hello, Go!")

	// -------------------------
	// Variable Declarations
	// -------------------------

	// Explicit declaration using the 'var' keyword.
	// Format: var <name> <type> = <value>
	// Use this style when declaring variables outside functions
	// or when you want to be explicit about the type.
	var x int = 10

	// Short-hand declaration using ':='
	// Go automatically infers the type (int in this case).
	// This is the most common style inside functions.
	y := 20

	// -------------------------
	// Using the add() Function
	// -------------------------

	// Call add() with x and y, store the result in 'sum'
	sum := add(x, y)
	fmt.Printf("The sum of %d and %d is %d\n", x, y, sum)
	// Output: The sum of 10 and 20 is 30

	// -------------------------
	// Using the greet() Function
	// -------------------------

	// greet() has no parameters and no return value.
	// It simply prints a message.
	greet()
	// Output: Hello!

	// -------------------------
	// Using the divide() Function
	// -------------------------

	// divide() returns TWO values: the result and an error.
	// We capture both using ':=' with two variables.
	result, err := divide(1, 2.0)

	// Always check the error before using the result.
	// This is the standard Go error handling pattern.
	if err != nil {
		// Something went wrong — print the error and stop
		fmt.Println("Error:", err)
	} else {
		// No error — safe to use the result
		fmt.Println("Result:", result)
		// Output: Result: 0.5
	}
}