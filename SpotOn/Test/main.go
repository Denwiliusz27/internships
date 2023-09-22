package main

import (
	"context"
	"fmt"
	"time"
)

func ticker() <-chan time.Time {
	now := time.Now()
	_, _, s := now.Clock()

	testT := time.NewTicker(time.Microsecond)
	defer testT.Stop()

	for {
		select {
		case t := <-testT.C:
			_, _, s2 := t.Clock()
			if s != s2 {
				return time.NewTicker(time.Second).C
			}
		}
	}
}

func display(name string, utcOffset int) chan<- time.Time {
	c := make(chan time.Time)

	go func() {
		for {
			select {
			case t := <-c:
				h, m, s := t.UTC().Clock()
				h = abs((h + utcOffset) % 24)
				fmt.Printf("%s - %02d:%02d:%02d\n", name, h, m, s)
			}
		}
	}()

	return c
}

func fanOut[T any](ctx context.Context, trigger <-chan T, receivers ...chan<- T) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("--- Terimnating FanOut ---")
			return

		case evnt := <-trigger:
			fmt.Println()
			for _, rec := range receivers {
				go func(r chan<- T) {
					select {
					case <-ctx.Done():
						return
					case r <- evnt:
					}
				}(rec)
			}
		}
	}
}

func abs(val int) int {
	if val < 0 {
		return -val
	}
	return val
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	krk := display("Kraków", 2)
	irk := display("Irkutsk", 8)
	sp := display("Sao Paulo", -3)
	vcv := display("Vancouver", -7)

	fanOut[time.Time](ctx, ticker(), krk, irk, sp, vcv)
	time.Sleep(3 * time.Second)
}
