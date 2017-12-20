package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"sort"
	"time"

	"github.com/RoaringBitmap/roaring"
	"github.com/zhenjl/bitmap/ewah"
)

type BitPackIntArray struct {
	values []int
	data   []uint64
}

func (a BitPackIntArray) Get(index int) int {
	bucket, shift := index>>4, uint(index&0xf)*4
	// fmt.Printf("%d=%d %d=%d %#x\n", index/16, bucket, index%16, pos, (a.data[index/16] >> uint(index%16*4)))
	return a.values[(a.data[bucket]>>shift)&0x0f]
}

func (a BitPackIntArray) Index(index int) int {
	bucket, shift := index>>4, uint(index&0xf)*4
	return int((a.data[bucket] >> shift) & 0x0f)
}

func (a BitPackIntArray) Set(index int, value int) {
	v := uint64(0)
	for a.values[v] != value {
		v++
	}

	bucket, shift := index>>4, uint(index&0xf)*4
	// fmt.Printf("%d %#x\n", index/16, ^(uint64(0x0f) << uint(index%16*4)))
	a.data[bucket] = a.data[bucket] & ^(uint64(0x0f)<<shift) | (v << shift)
}

func main() {
	// a := BitPackIntArray{[]int{500, 100, 200, 400}, []uint64{0, 0, 0, 0}}
	// a.Set(0, 100)
	// a.Set(1, 200)
	// a.Set(2, 200)
	// a.Set(20, 200)
	// a.Set(3, 500)
	// a.Set(4, 400)
	// a.Set(2, 100)

	// for i := 0; i < 30; i++ {
	// 	fmt.Println(a.Get(i))
	// }
	// return
	data, err := ioutil.ReadFile("cod.clm")
	if err != nil {
		panic(err)
	}

	var v []int
	if err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(&v); err != nil {
		panic(err)
	}

	count := len(v)

	fmt.Println(len(v), cap(v))

	values := []int{}
	f := make(map[int]int)
	for _, c := range v {
		f[c] = f[c] + 1
	}

	for a, b := range f {
		values = append(values, a)
		fmt.Println(a, b)
	}

	sort.Ints(values)
	sort.Ints(v)
	fmt.Println(values)

	runtime.GC()
	fmt.Println(len(v), cap(v))
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("\nAlloc = %v\nTotalAlloc = %v\nSys = %v\nNumGC = %v\n\n", m.Alloc, m.TotalAlloc, m.Sys, m.NumGC)

	a := BitPackIntArray{values, make([]uint64, len(v)>>4+1)}
	for i := range v {
		a.Set(i, v[i])
	}

	// check
	for i := range v {
		if a.Get(i) != v[i] {
			panic(i)
		}
	}

	runtime.GC()
	fmt.Println(len(a.data))
	runtime.ReadMemStats(&m)
	log.Printf("\nAlloc = %v\nTotalAlloc = %v\nSys = %v\nNumGC = %v\n\n", m.Alloc, m.TotalAlloc, m.Sys, m.NumGC)

	q := []*roaring.Bitmap{}
	for range values {
		q = append(q, roaring.NewBitmap())
	}

	for i := 0; i < count; i++ {
		q[a.Index(i)].AddInt(i)
		// q[a.Index(i)].AddInt(2 * i)
		// q[a.Index(i)].AddInt(3 * i)
		// q[a.Index(i)].AddInt(4 * i)
	}
	// for i := 0; i < len(q)-1; i++ {
	// 	q[i+1].Or(q[i])
	// }

	runtime.GC()
	runtime.ReadMemStats(&m)
	log.Printf("\nAlloc = %v\nTotalAlloc = %v\nSys = %v\nNumGC = %v\n\n", m.Alloc, m.TotalAlloc, m.Sys, m.NumGC)

	for i := range values {
		fmt.Println(values[i], q[i].GetCardinality(), q[i].GetSizeInBytes())
	}

	z := []*ewah.Ewah{}

	for range values {
		z = append(z, ewah.New().(*ewah.Ewah))
	}

	for i := range z {
		it := q[i].Iterator()
		for it.HasNext() {
			z[i].Set(int64(it.Next()))
		}
	}

	// for i := 0; i < count; i++ {
	// 	z[a.Index(i)].Set(int64(i))
	// 	// z[a.Index(i)].AddInt(2 * i)
	// 	// z[a.Index(i)].AddInt(3 * i)
	// 	// z[a.Index(i)].AddInt(4 * i)
	// }

	runtime.GC()
	runtime.ReadMemStats(&m)
	log.Printf("\nAlloc = %v\nTotalAlloc = %v\nSys = %v\nNumGC = %v\n\n", m.Alloc, m.TotalAlloc, m.Sys, m.NumGC)

	for i := range values {
		fmt.Println(values[i], z[i].Cardinality(), z[i].SizeInBytes())
	}

	time.Sleep(1 * time.Minute)

}
