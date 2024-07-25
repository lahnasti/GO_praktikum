/*
Создайте программу, которая симулирует банк и использует мьютексы для защиты операций
депозита и снятия средств.
*/

package main

import (
	"fmt"
	"sync"
)

type SberCard struct {
	balance int
	mutex   sync.Mutex
}

func (card *SberCard) Deposit(amount int) {
	card.mutex.Lock()
	defer card.mutex.Unlock()
	card.balance += amount
	fmt.Printf("A new deposit: %d.\nThe account has been updated: %d\n", amount, card.balance)
}

func (card *SberCard) WithdrawalMoney(amount int) {
	card.mutex.Lock()
	defer card.mutex.Unlock()
	if card.balance >= amount {
		card.balance -= amount
		fmt.Printf("You have withdrawn %d money.\nThe account has been updated: %d\n", amount, card.balance)
	} else {
		fmt.Println("Not enough funds to withdraw")
	}
}

func main() {
	card := &SberCard{balance: 5000}
	wg := sync.WaitGroup{}

	wg.Add(3)

	go func() {
		defer wg.Done()
		card.Deposit(3000)
	}()

	go func() {
		defer wg.Done()
		card.WithdrawalMoney(3000)
	}()

	go func() {
		defer wg.Done()
		card.WithdrawalMoney(1000000)
	}()

	wg.Wait()

	fmt.Printf("Final balance is %d\n", card.balance)
}
