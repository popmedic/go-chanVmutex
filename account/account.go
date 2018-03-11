package account

import (
	"sync"
)

type IAccount interface {
	Sum(amounts ...float64)
	Balance() float64
}

type ChanLockAccount struct {
	balance float64
	lock    chan byte
}

type MutexAccount struct {
	balance float64
	lock    *sync.RWMutex
}

type ChanAccount struct {
	sumChannel     chan []float64
	balanceChannel chan float64
}

func NewChanLockAccount(amounts ...float64) *ChanLockAccount {
	var balance float64 = 0
	for _, amount := range amounts {
		balance = balance + amount
	}
	return &ChanLockAccount{
		balance: balance,
		lock:    make(chan byte, 1),
	}
}

func (a *ChanLockAccount) Sum(amounts ...float64) {
	a.lock <- 0
	defer func(a *ChanLockAccount) { <-a.lock }(a)
	for _, amount := range amounts {
		a.balance = a.balance + amount
	}
}

func (a *ChanLockAccount) Balance() (balance float64) {
	a.lock <- 0
	defer func(a *ChanLockAccount) { <-a.lock }(a)
	balance = a.balance
	return
}

func NewMutexAccount(amounts ...float64) *MutexAccount {
	var balance float64
	for _, amount := range amounts {
		balance = balance + amount
	}
	return &MutexAccount{
		balance: balance,
		lock:    &sync.RWMutex{},
	}
}

func (a *MutexAccount) Sum(amounts ...float64) {
	a.lock.Lock()
	defer a.lock.Unlock()
	for _, amount := range amounts {
		a.balance = a.balance + amount
	}
}

func (a *MutexAccount) Balance() (balance float64) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	balance = a.balance
	return
}

func NewChanAccount(amounts ...float64) *ChanAccount {
	chanAccount := &ChanAccount{
		sumChannel:     make(chan []float64),
		balanceChannel: make(chan float64, 1),
	}
	chanAccount.balanceChannel <- 0.0
	go chanAccount.worker()
	chanAccount.Sum(amounts...)
	return chanAccount
}

func (a *ChanAccount) Sum(amounts ...float64) {
	a.sumChannel <- amounts
}

func (a *ChanAccount) Balance() float64 {
	balance := <-a.balanceChannel
	defer func(a *ChanAccount, b float64) { a.balanceChannel <- b }(a, balance)
	return balance
}

func (a *ChanAccount) worker() {
	for {
		balance := <-a.balanceChannel
		amounts := <-a.sumChannel
		for _, amount := range amounts {
			balance = balance + amount
		}
		a.balanceChannel <- balance
	}
}
