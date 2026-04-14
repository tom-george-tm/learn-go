package main

import "fmt"

func main() {

    // ---- IF/ELSE ----
    temperature := 35
    if temperature > 30 {
        fmt.Println("It's hot outside!")
    } else if temperature > 20 {
        fmt.Println("Nice weather")
    } else {
        fmt.Println("It's cold")
    }

    // ---- FOR LOOP ----
    fmt.Println("\n-- Counting --")
    for i := 1; i <= 5; i++ {
        fmt.Printf("Count: %d\n", i)
    }

    // ---- SWITCH ----
    fmt.Println("\n-- Day Check --")
    day := "Saturday"
    switch day {
    case "Saturday", "Sunday":
        fmt.Println("It's the weekend!")
    case "Monday":
        fmt.Println("Back to work...")
    default:
        fmt.Println("Regular weekday")
    }

    // ---- SLICE ----
    fmt.Println("\n-- Fruits --")
    fruits := []string{"apple", "banana", "mango"}
    fruits = append(fruits, "orange")

    for _, fruit := range fruits {
        fmt.Println("-", fruit)
    }

    // ---- MAP ----
    fmt.Println("\n-- Student Scores --")
    scores := map[string]int{
        "Alice": 95,
        "Bob":   87,
        "Carol": 91,
    }

    for name, score := range scores {
        if score >= 90 {
            fmt.Printf("%s: %d ⭐ (Excellent)\n", name, score)
        } else {
            fmt.Printf("%s: %d (Good)\n", name, score)
        }
    }
}
```

---

## Quick Reference Summary
```
IF/ELSE     →  if condition { } else if { } else { }
            →  if x := fn(); x != nil { }   ← init + check in one line

FOR LOOP    →  for i := 0; i < n; i++ { }   ← classic
            →  for condition { }              ← while style
            →  for { }                        ← infinite
            →  for i, v := range slice { }   ← range

SWITCH      →  switch variable { case x: case y: default: }
            →  switch { case condition: }     ← no variable needed

ARRAYS      →  [5]int{1,2,3,4,5}             ← fixed size
SLICES      →  []int{1,2,3}                  ← flexible, use this
            →  append(slice, newItem)         ← adding items
            →  slice[1:3]                     ← slicing

MAPS        →  map[string]int{"key": value}
            →  m["key"]                       ← read
            →  m["key"] = value               ← write
            →  delete(m, "key")               ← delete
            →  value, ok := m["key"]          ← safe check