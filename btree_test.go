package btree

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tree := New()
	for i := uint64(0); i <= 30; i += 2 {
		tree.Insert(i)
		val, err := tree.Search(i)
		if err != nil {
			log.Fatal(err)
		}
		if val != i {
			t.Fatalf("want: %d, got: %s", i, val)
		}
	}

	for i := 29; i >= 1; i -= 2 {
		v := uint64(i)
		tree.Insert(v)
		val, err := tree.Search(v)
		if err != nil {
			log.Fatal(err)
		}
		if val != v {
			t.Fatalf("want: %d, got: %s", i, val)
		}
	}

	//fmt.Printf("%s\n", tree)
}

func TestRand(t *testing.T) {
	seed := time.Now().UnixNano()
	fmt.Printf("seed: %d\n", seed)
	rand.Seed(seed)

	tree := New()
	for i := 0; i < 50; i++ {
		num := uint64(rand.Intn(998) + 1)
		if _, err := tree.Search(num); err == nil {
			continue
		}
		tree.Insert(num)

		val, err := tree.Search(num)
		if err != nil {
			log.Fatal(err)
		}
		if val != num {
			t.Fatalf("want: %d, got: %s", num, val)
		}
	}

	fmt.Printf("%s\n", tree)
}
