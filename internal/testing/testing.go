package testing

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"
)

var (
	asciiLetters    = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	japaneseLetters = []rune("あいうえおアイウエオ亜伊卯惠尾")
	mux             sync.Mutex
)

func AsciiLetters() []rune {
	return asciiLetters
}

func JapaneseLetters() []rune {
	return japaneseLetters
}

// RandStr generates a random string of the specified number of characters.
//
// Be sure to specify the seed at rand.Seed(time.Now().UnixNano()) before executing.
// The characters that make up are defined in (part of) Ascii characters.
// Use RandStrWithLetter if you want to specify a different character
func RandStr(n int) string {
	return RandStrWithLetter(n, asciiLetters)
}

// RandStrWithLetter generates a random string of the specified number of characters.
//
// Example:
//
//   fmt.Println(RandStrWithLetter(100, []rune("0123456789")))
//   fmt.Println(RandStrWithLetter(100, JapaneseLetters()))
func RandStrWithLetter(n int, runes []rune) string {
	b := make([]rune, n)
	for i := range b {
		mux.Lock()
		b[i] = runes[rand.Intn(len(runes))]
		mux.Unlock()
	}
	return string(b)
}

// RandChoice randomly selects an element from the given array and returns it.
func RandChoice(srclist interface{}) interface{} {
	v := reflect.ValueOf(srclist)
	if k := v.Kind(); k != reflect.Array && k != reflect.Slice {
		panic("want Array or Slice but " + k.String())
	}
	switch l := v.Len(); l {
	case 0:
		panic("srclist is empty")
	case 1:
		return v.Index(0).Interface()
	default:
		return v.Index(rand.Intn(l)).Interface()
	}
}

// ParallelWithDeadline returns whether or not the execution of fn (number of attempts n)
// finishes within the time limit d.
//
// Example:
//
//	   rand.Seed(time.Now().UnixNano())
//	   fn := func(_ *testing.T, i int) {
//	   	t.Logf("exercise: %d", i)
//	   }
//	   err := ParallelWithDeadline(context.Background(), t, 100, 1*time.Second, fn)
//	   if err != nil {
//	   	t.Error(err)
//	   	t.FailNow()
//	   }
func ParallelWithDeadline(
	ctx context.Context,
	t *testing.T,
	n uint,
	d time.Duration,
	fn func(t *testing.T, i int),
) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(d))
	defer cancel()
	// Exercise
	var wg sync.WaitGroup
	for i := 0; i < int(n); i++ {
		wg.Add(1)
		go func(t *testing.T, i int) {
			defer wg.Done()
			fn(t, i)
		}(t, i)
	}
	// Verify
	allDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(allDone)
	}()
	fmt.Printf("timeout: %v", d)
	select {
	case <-ctx.Done():
		t.Error(ctx.Err())
		return errors.New("fialed")
	case <-allDone:
		return nil
	}
}
