package account

import (
	"sync"
	"testing"
)

func TestAccounts(t *testing.T) {
	testAccount(t, NewChanLockAccount(4, 6), 20, 4, 3, 2, 1)
	testAccount(t, NewMutexAccount(4, 6), 20, 4, 3, 2, 1)
	testAccount(t, NewChanAccount(4, 6), 20, 4, 3, 2, 1)
}

func testAccount(t *testing.T, a IAccount, exp float64, amounts ...float64) {
	sumAmounts(a, amounts...)
	got := a.Balance()
	if got != exp {
		t.Errorf("falure for type %T: expected %f got %f", a, exp, got)
	}
}

func BenchmarkChanLockAccount(b *testing.B) {
	benchmarkAccount(b, NewChanLockAccount())
}

func BenchmarkMutexAccount(b *testing.B) {
	benchmarkAccount(b, NewMutexAccount())
}

func BenchmarkChanAccount(b *testing.B) {
	benchmarkAccount(b, NewChanAccount())
}

func benchmarkAccount(b *testing.B, a IAccount) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sumAmounts(a, float64(i))
	}
}

func sumAmounts(a IAccount, amounts ...float64) {
	wg := &sync.WaitGroup{}
	for _, amount := range amounts {
		wg.Add(1)
		go func(amount float64) {
			a.Sum(amount)
			wg.Done()
		}(amount)
	}
	wg.Wait()
}
