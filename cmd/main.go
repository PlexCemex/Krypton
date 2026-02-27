package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Автор: Фесун Семён
// 1961 - год первого полета человека в космос

const year = 1961

type SafeMap struct {
	sync.Mutex
	m         map[int]int
	accesses  int
	additions int
}

func main() {
	rand.Seed(time.Now().UnixNano())

	sm := &SafeMap{m: make(map[int]int)}

	var keys []int
	for i := 0; i < 3; i++ {
		for k := 1; k <= year; k++ {
			keys = append(keys, k)
		}
	}

	rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })

	ch := make(chan int, len(keys))
	for _, k := range keys {
		ch <- k
	}
	close(ch)

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := range ch {
				sm.Lock()
				sm.accesses++
				val, ok := sm.m[k]
				if !ok {
					sm.additions++
				}
				sm.m[k] = val + 1
				sm.Unlock()
			}
		}()
	}
	wg.Wait()

	fmt.Printf("Обращения: %d (ожидалось %d)\n", sm.accesses, year*3)
	fmt.Printf("Добавления: %d (ожидалось %d)\n", sm.additions, year)

	for k := 1; k <= year; k++ {
		if sm.m[k] != 3 {
			panic(fmt.Sprintf("Ошибка: по ключу %d значение %d, а не 3", k, sm.m[k]))
		}
	}
	fmt.Println("Тест пройден успешно!")
}
