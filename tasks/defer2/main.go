package main

import "fmt"

type Account struct {
	Balance int
}

func main() {
	initBalance := 1000
	account := &Account{Balance: initBalance}

	defer printBalance("Initial balance", account.Balance) // by copy - 1000
	defer printBalance("Current balance", account.Balance) // by copy - 1000
	defer printAccountBalance("Pointer", account)          // by pointer - 1000 + 500 - 200 = 1300

	account.Balance += 500
	updateBalance(account, 200)
	account = &Account{Balance: 300} // a new variable, new pointer, so it's not passed to the printAccountBalance defer
}

func updateBalance(account *Account, amount int) {
	account.Balance -= amount
}

func printAccountBalance(message string, account *Account) {
	fmt.Printf("%s : %d\n", message, account.Balance)
}

func printBalance(message string, balance int) {
	fmt.Printf("%s %d\n", message, balance)
}
