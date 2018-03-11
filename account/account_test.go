package account

import (
	"testing"
)

func TestAccounts(t *testing.T) {
	testAccount(t, NewChanLockAccount(4, 6), 20, 7, 3)
	testAccount(t, NewMutexAccount(4, 6), 20, 7, 3)
	testAccount(t, NewChanAccount(4, 6), 20, 7, 3)
}

func testAccount(t *testing.T, a IAccount, exp float64, amounts ...float64) {
	a.Sum(amounts...)
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
		a.Sum(10, 20, 30)
	}
}
