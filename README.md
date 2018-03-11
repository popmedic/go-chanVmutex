# Benchmarking Channels Vs. Mutexes in Go.

## Results\*

| Lock Mechanism | Results (ns/op) |
| --- | --- |
| Channel Lock | 104 ns/op |
| Mutex Lock | 75.8 ns/op |
| Channel Worker | 309 ns/op |

\*_MacBook Pro, 2.8 GHz Intel Core i7, 16 GB 1600 MHz DDR3_

## Objective

The objective of this project was to benchmark the difference in using a Go Channel against a Go Mutex when used to guard a critical section.

## Running for yourself

To get these results yourself from the terminal 

1. Get the project 

``` bash
go get github.com/popmedic/go-chanVmutex/...
```

2. Change directory into the account directory 

``` bash
cd $GOPATH/src/github.com/popmedic/go-chanVmutex/account
```

3. Run the benchmarks 

``` bash
go test -bench=.
```

## Scenario

For the scenario I selected a Bank Account object because summing a bank account balance is the classic example of a "race condition."  The "race condition" occurs when we have an Account object that stores a users bank balance.  If multiple "threads" try to "sum" on this bank account without locking we can end up with the problem of:

1. Thread 1 gets the balance (say 10.00)
2. Thread 2 gets the balance before thread 1 has finished with the balance (so they get 10.00 also)
3. Thread 1 adds its amount (say 10.00, totaling 20.00)
4. Thread 2 adds its amount (say 10.00, totaling 20.00)
5. Gosh darn it the total is 20.00 BUT WE ADDED 20.00 to 10.00, it should be 30.00!!!

If instead we lock the "critical section" (adding to the balance) so that Thread 2 will not get the balance until Thread 1 is finished with the "critical section," we can avoid this "race condition."  Both [ChanAccount](#chanlockaccount) and [MutexAccount](#mutexaccount) use this technique.

Go language introduced a way to avoid locks by communication through channels, instead of relying on the condition of a lock ("communication instead of conditions.")  One can do this by adding a Go worker routine to the class that tries getting a value off a channel, and if it can, use the value, if no value is on the channel, do nothing. The [ChanAccount](#chanaccount) uses this technique.

## ChanLockAccount

One way I have seen people use channels is in place of a Mutex by having a buffered channel of 1, and putting a value on the channel to lock before the critical section and removing a value off the channel after the critical section to unlock.  

``` Go
a.chlock <- 0
defer func(a *ChanLockAccount) { <-a.chlock }(a)
for _, amount := range amounts {
    a.balance = a.balance + amount
}
```

This works because by definition a buffered channel of 1 will pause if a value is in the channel until the channel is empty again, effectively making it a lock.

This is not the intended way to use a channel for "communication over conditions" like Go designed them, so I put this here for all the newer developers that seem to use this technique for locking.

## MutexAccount

The time tested, user approved mutex is the standard lock used most by developers for decades.  This guards the critical section by using a kernel mutex that all go routines can use to lock before the critical section and unlock after the critical section.

``` Go
a.lock.Lock()
defer a.lock.Unlock()
for _, amount := range amounts {
    a.balance = a.balance + amount
}
```

## ChanAccount

For this I implement the standard "worker" routine that will pause until something is put on a channel, retrieve the value on the channel, preform the critical section, and communicate back the value.  Go would like us to use channels for "communication over conditions" and avoids using any locks.

``` Go
for {
    balance := <-a.balanceChannel
    amounts := <-a.sumChannel
    for _, amount := range amounts {
        balance = balance + amount
    }
    a.balanceChannel <- balance
}
```

## Conclusion

I found through this scenario that it is best to do access control for thread safety by using a Mutex.  I decided this based on the results of benchmarking, and also on the fact that it seems more readable and maintainable to use the common pattern of locking then the concept of "communication over conditions."
