package main

import "fmt"

// ---- STRUCTS ----

type Address struct {
    City    string
    Country string
}

type Employee struct {
    Name    string
    Age     int
    Salary  float64
    Address Address
}

// ---- METHODS ----

// Value receiver — just reading
func (e Employee) Describe() string {
    return fmt.Sprintf("%s (%d) from %s", e.Name, e.Age, e.Address.City)
}

// Pointer receiver — modifying data
func (e *Employee) GiveRaise(percent float64) {
    e.Salary += e.Salary * (percent / 100)
}

// ---- INTERFACE ----

type Worker interface {
    Describe() string
    GiveRaise(percent float64)
}

func processRaise(w Worker, percent float64) {
    fmt.Printf("Before raise: %s\n", w.Describe())
    w.GiveRaise(percent)
    fmt.Printf("After raise:  %s\n", w.Describe())
}

func main() {
    emp := &Employee{   // & because GiveRaise uses pointer receiver
        Name:   "Alice",
        Age:    28,
        Salary: 50000,
        Address: Address{
            City:    "Kochi",
            Country: "India",
        },
    }

    fmt.Println(emp.Describe())

    emp.GiveRaise(10)
    fmt.Printf("New Salary: %.2f\n", emp.Salary)
}
```

---

## Mental Model Summary
```
STRUCT      →  Blueprint for grouping related data
               type Person struct { Name string; Age int }

METHODS     →  Functions attached to a struct
               func (p Person) Greet() string { ... }

POINTERS    →  &variable  →  get the memory address
               *pointer   →  get the value at that address

RECEIVER    →  Value    func (p Person) ...      →  works on a copy (read only)
               Pointer  func (p *Person) ...     →  works on original (read/write)

INTERFACE   →  Defines behaviour (methods), not data
               Any type with the required methods automatically satisfies it
               Lets you write functions that work with MULTIPLE types
```

---

## How They All Connect
```
Struct      →  defines the DATA  (what it has)
Methods     →  defines the BEHAVIOUR (what it can do)
Pointers    →  lets methods MODIFY the original data
Interfaces  →  groups types by BEHAVIOUR (what they can do)