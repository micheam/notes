package testing

import (
	"context"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestParallelWithDeadline(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	exercise := func(_ *testing.T, i int) {
		// t.Logf("exercise: %d", i)
	}
	err := ParallelWithDeadline(context.Background(), t, 100, 1*time.Second, exercise)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestRandChoice(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	want := "one"
	got := RandChoice([]string{want})
	if got.(string) != want {
		t.Errorf("want %s but got %s", want, got)
	}
	log.Println(RandChoice([]int{1, 2, 3, 4, 5, 6}))
}
