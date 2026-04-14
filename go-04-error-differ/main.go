package main

import (
    "errors"
    "fmt"
)

// ---- CUSTOM ERROR ----

type InsufficientFundsError struct {
    Requested float64
    Available float64
}

func (e *InsufficientFundsError) Error() string {
    return fmt.Sprintf("insufficient funds: requested %.2f but only %.2f available",
        e.Requested, e.Available)
}

// ---- STRUCT & METHODS ----

type BankAccount struct {
    Owner   string
    Balance float64
}

// VARIADIC — deposit multiple amounts at once
func (b *BankAccount) Deposit(amounts ...float64) {
    for _, amount := range amounts {
        b.Balance += amount
        fmt.Printf("  Deposited %.2f | Balance: %.2f\n", amount, b.Balance)
    }
}

// CUSTOM ERROR — withdraw with proper error type
func (b *BankAccount) Withdraw(amount float64) error {
    if amount > b.Balance {
        return &InsufficientFundsError{
            Requested: amount,
            Available: b.Balance,
        }
    }
    b.Balance -= amount
    return nil
}

// DEFER — prints summary when function exits
func processTransaction(account *BankAccount) {
    // This always runs when processTransaction() returns
    defer fmt.Printf("\n=== Final Balance for %s: %.2f ===\n",
        account.Owner, account.Balance)

    fmt.Println("Starting transactions...")
    account.Deposit(500, 250.50, 100)

    err := account.Withdraw(200)
    if err != nil {
        fmt.Println("Error:", err)
    }

    err = account.Withdraw(10000)
    if err != nil {
        var fundErr *InsufficientFundsError
        if errors.As(err, &fundErr) {
            fmt.Printf("Transaction blocked: you are %.2f short\n",
                fundErr.Requested - fundErr.Available)
        }
    }
}

// CLOSURE — creates a transaction logger with private history
func makeTransactionLogger(owner string) func(string, float64) {
    history := []string{}   // private state captured by closure

    return func(txType string, amount float64) {
        entry := fmt.Sprintf("%s: %s %.2f", owner, txType, amount)
        history = append(history, entry)
        fmt.Println("Logged:", entry)
        fmt.Printf("Total transactions so far: %d\n", len(history))
    }
}

func main() {
    account := &BankAccount{Owner: "Alice", Balance: 1000}

    // Use the closure-based logger
    log := makeTransactionLogger(account.Owner)
    log("DEPOSIT", 500)
    log("WITHDRAW", 200)

    fmt.Println()

    // Run transactions with defer tracking final balance
    processTransaction(account)
}
```

---

## Mental Model Summary
```
ERROR HANDLING
  errors.New("msg")         →  simple sentinel error
  fmt.Errorf("msg: %w", e)  →  wrap with context
  errors.Is(err, target)    →  check if it's a specific error
  errors.As(err, &target)   →  extract a specific error type
  Custom error type         →  struct that implements Error() string

VARIADIC FUNCTIONS
  func fn(args ...int)      →  accepts any number of ints
  Inside function           →  args is just a []int slice
  Calling with a slice      →  fn(mySlice...)

DEFER
  defer fn()                →  runs when the current function exits
  Multiple defers           →  run in reverse order (LIFO)
  Best use                  →  cleanup (close files, unlock, disconnect)

CLOSURES
  func inside a func        →  captures outer variables by reference
  Returned function         →  remembers its captured variables forever
  Common uses               →  counters, factories, middleware, handlers
```

---

## How All Four Connect
```
Variadic     →  flexible inputs  (accept anything)
Defer        →  guaranteed cleanup  (always runs at exit)
Closures     →  stateful functions  (remember their context)
Error types  →  rich error info  (carry data, not just messages)