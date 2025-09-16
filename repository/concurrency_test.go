package repository

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestConcurrentCreateSameSymbol(t *testing.T) {
	repo := NewMemoryCryptoRepo()

	var successes int32
	var errs int32
	var goroutines = 100
	wg := sync.WaitGroup{}
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_, err := repo.Create("btc")
			if err != nil {
				atomic.AddInt32(&errs, 1)
			} else {
				atomic.AddInt32(&successes, 1)
			}
		}()
	}

	wg.Wait()

	s := atomic.LoadInt32(&successes)
	e := atomic.LoadInt32(&errs)
	if s != 1 {
		t.Fatalf("expected exactly 1 success, got %d", s)
	}
	if int(e) != goroutines-1 {
		t.Fatalf("expected %d errors, got %d", goroutines-1, e)
	}
}

func TestConcurrentRefreshPrice(t *testing.T) {
	repo := NewMemoryCryptoRepo()

	c, err := repo.Create("eth")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	initial := len(c.History)

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	var errs int32
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			if _, err := repo.RefreshPrice("eth"); err != nil {
				atomic.AddInt32(&errs, 1)
			}
		}()
	}

	wg.Wait()
	if e := atomic.LoadInt32(&errs); e != 0 {
		t.Fatalf("RefreshPrice errors: %d", e)
	}

	updated, err := repo.Get("eth")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	// History is capped at 100 latest records
	want := initial + goroutines
	if want > 100 {
		want = 100
	}
	if len(updated.History) != want {
		t.Fatalf("history length: got %d want %d", len(updated.History), want)
	}

	last := updated.History[len(updated.History)-1]
	if updated.CurrentPrice != last.Price {
		t.Errorf("CurrentPrice != last history price: %v vs %v", updated.CurrentPrice, last.Price)
	}
	if !updated.LastUpdated.Equal(last.Timestamp) {
		t.Errorf("LastUpdated != last history timestamp: %v vs %v", updated.LastUpdated, last.Timestamp)
	}

	for i := 0; i+1 < len(updated.History); i++ {
		a := updated.History[i].Timestamp
		b := updated.History[i+1].Timestamp
		if b.Before(a) {
			t.Errorf("timestamps not non-decreasing at %d: %v then %v", i, a, b)
			break
		}
	}
}

func TestReadersVsWriter(t *testing.T) {
	repo := NewMemoryCryptoRepo()

	c, err := repo.Create("eth")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	initial := len(c.History)

	const writers = 100
	const readers = 10

	done := make(chan struct{})

	var wgReaders sync.WaitGroup
	wgReaders.Add(readers)
	for i := 0; i < readers; i++ {
		go func() {
			defer wgReaders.Done()
			for {
				select {
				case <-done:
					return
				default:
					if _, err := repo.Get("eth"); err != nil {
						t.Errorf("Get failed: %v", err)
						return
					}
					if _, err := repo.List(); err != nil {
						t.Errorf("List failed: %v", err)
						return
					}
				}
			}
		}()
	}

	var wgWriters sync.WaitGroup
	wgWriters.Add(writers)
	for i := 0; i < writers; i++ {
		go func() {
			defer wgWriters.Done()
			if _, err := repo.RefreshPrice("eth"); err != nil {
				t.Errorf("RefreshPrice failed: %v", err)
			}
		}()
	}

	wgWriters.Wait()
	close(done)
	wgReaders.Wait()

	updated, err := repo.Get("eth")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	// History is capped at 100 latest records
	want := initial + writers
	if want > 100 {
		want = 100
	}
	if len(updated.History) != want {
		t.Fatalf("history length: got %d want %d", len(updated.History), want)
	}

	last := updated.History[len(updated.History)-1]
	if updated.CurrentPrice != last.Price {
		t.Errorf("CurrentPrice != last history price: %v vs %v", updated.CurrentPrice, last.Price)
	}
	if !updated.LastUpdated.Equal(last.Timestamp) {
		t.Errorf("LastUpdated != last history timestamp: %v vs %v", updated.LastUpdated, last.Timestamp)
	}

	for i := 0; i+1 < len(updated.History); i++ {
		a := updated.History[i].Timestamp
		b := updated.History[i+1].Timestamp
		if b.Before(a) {
			t.Errorf("timestamps not non-decreasing at %d: %v then %v", i, a, b)
			break
		}
	}
}

func TestCreateVsRefreshPriceRace(t *testing.T) {
	repo := NewMemoryCryptoRepo()

	const preWriters = 50
	const postWriters = 50

	var created int32 // 0 until Create returns successfully
	var preErrs int32
	var preLateOK int32
	var badSuccessBefore int32
	var badErrorAfter int32
	var postOK int32

	startPre := make(chan struct{})
	startPost := make(chan struct{})
	initLenCh := make(chan int, 1)
	createErrCh := make(chan error, 1)

	var wg sync.WaitGroup

	// Pre-create writers
	for i := 0; i < preWriters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-startPre
			if _, err := repo.RefreshPrice("race"); err != nil {
				if atomic.LoadInt32(&created) == 1 {
					atomic.AddInt32(&badErrorAfter, 1)
				} else {
					atomic.AddInt32(&preErrs, 1)
				}
			} else {
				if atomic.LoadInt32(&created) == 0 {
					atomic.AddInt32(&badSuccessBefore, 1)
				} else {
					atomic.AddInt32(&preLateOK, 1)
				}
			}
		}()
	}

	// Create in parallel, then release post writers
	go func() {
		// small delay to ensure some pre-writers run before Create
		time.Sleep(5 * time.Millisecond)
		c, err := repo.Create("race")
		if err != nil {
			// Report error to main goroutine and unblock post writers
			// to avoid deadlocks in case of failure.
			createErrCh <- err
			close(startPost)
			return
		}
		// Publish creation state before any other signalling to avoid
		// misclassifying late pre-writer successes as "before create".
		atomic.StoreInt32(&created, 1)
		initLenCh <- len(c.History)
		close(startPost)
		createErrCh <- nil
	}()

	// Start pre writers
	close(startPre)

	// Post-create writers
	for i := 0; i < postWriters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-startPost
			if _, err := repo.RefreshPrice("race"); err != nil {
				// After Create all RefreshPrice must succeed
				atomic.AddInt32(&badErrorAfter, 1)
			} else {
				if atomic.LoadInt32(&created) == 0 {
					atomic.AddInt32(&badSuccessBefore, 1)
				} else {
					atomic.AddInt32(&postOK, 1)
				}
			}
		}()
	}

	wg.Wait()

	// Ensure Create succeeded (and avoid reading from initLenCh on failure)
	if err := <-createErrCh; err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if s := atomic.LoadInt32(&badSuccessBefore); s != 0 {
		t.Fatalf("RefreshPrice succeeded before Create: %d cases", s)
	}
	if e := atomic.LoadInt32(&badErrorAfter); e != 0 {
		t.Fatalf("RefreshPrice errored after Create: %d cases", e)
	}
	if e := atomic.LoadInt32(&preErrs); e == 0 {
		t.Fatalf("expected at least one pre-create error")
	}
	if ok := atomic.LoadInt32(&postOK); ok != postWriters {
		t.Fatalf("post-create successes: got %d want %d", ok, postWriters)
	}

	initial := <-initLenCh

	updated, err := repo.Get("race")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	// History is capped at 100 latest records
	want := initial + int(atomic.LoadInt32(&preLateOK)) + int(atomic.LoadInt32(&postOK))
	if want > 100 {
		want = 100
	}
	if len(updated.History) != want {
		t.Fatalf("history length: got %d want %d", len(updated.History), want)
	}

	last := updated.History[len(updated.History)-1]
	if updated.CurrentPrice != last.Price {
		t.Errorf("CurrentPrice != last history price: %v vs %v", updated.CurrentPrice, last.Price)
	}
	if !updated.LastUpdated.Equal(last.Timestamp) {
		t.Errorf("LastUpdated != last history timestamp: %v vs %v", updated.LastUpdated, last.Timestamp)
	}

	for i := 0; i+1 < len(updated.History); i++ {
		a := updated.History[i].Timestamp
		b := updated.History[i+1].Timestamp
		if b.Before(a) {
			t.Errorf("timestamps not non-decreasing at %d: %v then %v", i, a, b)
			break
		}
	}
}
