package test

import (
	"math/rand"
	"search_engine/utils"
	"strconv"
	"sync"
	"testing"
)

var conMap = utils.NewConcurrentMap(8, 1000)
var synMap = sync.Map{}

func writeConMap() {
	for i := 0; i < 10000; i++ {
		key := strconv.Itoa(int(rand.Int63()))
		conMap.Set(key, 1)
	}
}

func readConMap() {
	for i := 0; i < 10000; i++ {
		key := strconv.Itoa(int(rand.Int63()))
		conMap.Get(key)
	}
}

func writeSynMap() {
	for i := 0; i < 10000; i++ {
		key := strconv.Itoa(int(rand.Int63()))
		synMap.Store(key, 1)
	}
}

func readSynMap() {
	for i := 0; i < 10000; i++ {
		key := strconv.Itoa(int(rand.Int63()))
		synMap.Load(key)
	}
}

func BenchmarkConMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		const c = 500

		wg := sync.WaitGroup{}
		wg.Add(2 * c)
		for i := 0; i < c; i++ {
			go func() {
				defer wg.Done()
				writeConMap()
			}()
		}
		for i := 0; i < c; i++ {
			go func() {
				defer wg.Done()
				readConMap()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkSynMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		const c = 500

		wg := sync.WaitGroup{}
		wg.Add(2 * c)
		for i := 0; i < c; i++ {
			go func() {
				defer wg.Done()
				writeSynMap()
			}()
		}
		for i := 0; i < c; i++ {
			go func() {
				defer wg.Done()
				readSynMap()
			}()
		}
		wg.Wait()
	}
}

//go test ./utils/test -bench=Map -run=^$ -count=1 -benchmem -benchtime=3s
//BenchmarkConMap-24             4         929151625 ns/op        890201162 B/op  10154559 allocs/op
//BenchmarkSynMap-24             1        4457198800 ns/op        1131694048 B/op 30154641 allocs/op
