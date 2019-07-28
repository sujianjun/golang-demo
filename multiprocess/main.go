package main

/*
	利用goruntinue 进行多协程并发计算
*/
import (
	"fmt"
	"sync"
	"time"
)

// Data 计算结果结构体
type Data struct {
	Name     string
	Total    int
	Increase int
}

// Record 计算前记录结构体
type Record struct {
	ID    int
	Name  string
	Start int
	End   int
}

/* Cal 计算函数:
d map[string]*Data 为共享变量，多线程并发计算的结果会根据key插入到d指针参数中,
使用读写锁及waitgroup保证数据正确性及全部计算完成再返回主线程
*/
func Cal(r Record, d map[string]*Data, mu *sync.RWMutex, wg *sync.WaitGroup) {
	res := fmt.Sprintf("Name: %s, ID: %d, Start:%d, End:%d", r.Name, r.Start, r.End, r.ID)
	fmt.Println(res)
	mu.Lock()
	key := r.Name
	d[key].Increase = r.End - r.Start + d[key].Increase
	d[key].Total = r.End + r.Start + d[key].Total
	mu.Unlock()

	defer wg.Done()
}

func main() {
	var l sync.RWMutex
	var res = map[string]*Data{}
	var D = []Record{
		{ID: 1, Name: "A", Start: 3, End: 8},
		{ID: 2, Name: "A", Start: 4, End: 12},
		{ID: 3, Name: "B", Start: 1, End: 9},
		{ID: 4, Name: "C", Start: 10, End: 13},
		{ID: 4, Name: "B", Start: 4, End: 6},
		{ID: 6, Name: "A", Start: 21, End: 32},
	}
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	var wg sync.WaitGroup
	wg.Add(len(D))
	for _, v := range D {
		if _, ok := res[v.Name]; !ok {
			res[v.Name] = &Data{
				Name:     v.Name,
				Increase: 0,
				Total:    0,
			}
		}
	}
	for _, v := range D {
		go Cal(v, res, &l, &wg)
	}
	wg.Wait()
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	for _, r := range res {
		fmt.Println(r.Name, r.Increase, r.Total)
	}
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
