package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

func incrementByTen(p *int) {
	*p += 10
}

func testForIncrementByTen() {
	num := 15
	fmt.Printf("修改前的值：%d\n", num)
	incrementByTen(&num)
	fmt.Printf("修改后的值：%d\n", num)
}

func safeDoubleEachEle(p *[]int) {
	if p == nil {
		fmt.Println("传入切片指针为nil")
		return
	}
	if *p == nil {
		fmt.Println("传入切片指针指向切片为nil")
		return
	}
	for i := range *p {
		(*p)[i] *= 2
	}
}
func testForSafeDoubleEachEle() {
	nums := []int{1, 2, 3, 5, 4, 6}
	fmt.Printf("修改前的切片：%v\n", nums)
	safeDoubleEachEle(&nums)
	fmt.Printf("修改后的切片：%v\n", nums)
}
func goroutinePractice() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("奇数打印协程启动")
		for i := 1; i < 10; i += 2 {
			fmt.Printf("奇数：%d\n", i)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("偶数打印协程启动")
		for i := 2; i <= 10; i += 2 {
			fmt.Printf("偶数：%d\n", i)
		}
	}()

	wg.Wait()
	fmt.Println("所有协程执行完毕")
}

type Task func() error
type TaskResult struct {
	taskIndex int
	costTime  time.Duration
	err       error
}

func scheduler(taskList []Task, maxWorkers int) []TaskResult {
	results := make([]TaskResult, len(taskList))
	taskCh := make(chan int, len(taskList))
	resultCh := make(chan TaskResult, len(taskList))

	var wg sync.WaitGroup

	wg.Add(maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		go func(workId int) {
			defer wg.Done()
			for taskIdx := range taskCh {
				fmt.Printf("协程%d开始执行任务：%d\n", i, taskIdx)
				start := time.Now()
				err := taskList[taskIdx]()
				costTime := time.Since(start)

				resultCh <- TaskResult{
					taskIndex: taskIdx,
					costTime:  costTime,
					err:       err,
				}
				fmt.Printf("协程%d完成任务%d,耗时:%s\n", i, taskIdx, costTime)
			}
		}(i)
	}
	go func() {
		for idx := range taskList {
			taskCh <- idx
		}
		close(taskCh)
	}()
	go func() {
		wg.Wait()
		close(resultCh)
	}()
	for res := range resultCh {
		results[res.taskIndex] = res
	}
	return results
}
func testForSchedual() {
	taskList := []Task{
		// 任务0：执行200ms，成功
		func() error {
			time.Sleep(200 * time.Millisecond)
			return nil
		},
		// 任务1：执行100ms，成功
		func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		},
		// 任务2：执行300ms，模拟错误
		func() error {
			time.Sleep(300 * time.Millisecond)
			return fmt.Errorf("任务执行失败：模拟业务错误")
		},
		// 任务3：执行150ms，成功
		func() error {
			time.Sleep(150 * time.Millisecond)
			return nil
		},
		// 任务4：执行250ms，成功
		func() error {
			time.Sleep(250 * time.Millisecond)
			return nil
		},
	}
	fmt.Println("===任务调度器启动===")
	startTime := time.Now()
	results := scheduler(taskList, 4)
	totalCost := time.Since(startTime)
	fmt.Println("===所有任务执行完毕===")

	fmt.Printf("总执行时间：%s\n", totalCost)
	fmt.Println("\n各子任务执行详情")
	for _, res := range results {
		if res.err != nil {
			fmt.Sprintf("任务%d失败（%s）", res.taskIndex, res.err.Error())
		}
		fmt.Sprintf("任务%d成功耗时：%s", res.taskIndex, res.costTime)
	}
}

type Shape interface {
	Area() float64
	Perimeter() float64
	Name() string
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r Rectangle) Name() string {
	return "长方形"
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}
func (c Circle) Name() string {
	return "圆形"
}

func TestForShape() {
	rect := Rectangle{Width: 5, Height: 3.3}
	circle := Circle{Radius: 4.5}

	fmt.Println("=== 长方形 ===")
	fmt.Printf("宽：%.2f, 高：%.2f\n", rect.Width, rect.Height)
	fmt.Printf("面积：%.2f\n", rect.Area())
	fmt.Printf("周长：%.2f\n", rect.Perimeter())

	fmt.Println("\n=== 圆形 ===")
	fmt.Printf("半径：%.2f\n", circle.Radius)
	fmt.Printf("面积：%.2f\n", circle.Area())
	fmt.Printf("周长：%.2f\n", circle.Perimeter())

	fmt.Println("\n====接口多态====")
	shapes := []Shape{rect, circle}
	for _, s := range shapes {
		fmt.Printf("%s:面积：%.2f，周长：%.2f\n", s.Name(), s.Area(), s.Perimeter())
	}
}

type Person struct {
	Name string
	Age  int
}
type Employee struct {
	Person
	EmployeeId string
}

func (e Employee) PrintInfo() {
	// 访问组合的 Person 字段：可直接访问（匿名嵌入的特性），也可通过 e.Person.Name 显式访问
	fmt.Println("=== 员工信息 ===")
	fmt.Printf("员工ID：%s\n", e.EmployeeId)
	fmt.Printf("姓名：%s\n", e.Name)       // 直接访问组合字段（推荐，简洁）
	fmt.Printf("年龄：%d\n", e.Person.Age) // 显式访问组合字段（等价，适合字段名冲突场景）
}
func testForEmplyee() {
	emp1 := Employee{
		Person: Person{
			Name: "第一名员工",
			Age:  25,
		},
		EmployeeId: "employee1",
	}
	emp2 := Employee{
		Person:     Person{Name: "第二名员工", Age: 29},
		EmployeeId: "employee2",
	}
	emp1.PrintInfo()
	fmt.Println()
	emp2.PrintInfo()
}
func channelStudy() {
	numChan := make(chan int)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(numChan)
		for i := 1; i <= 10; i++ {
			fmt.Printf("发送：%d\n", i)
			numChan <- i
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for num := range numChan {
			fmt.Printf("接收并打印：%d\n", num)
		}
	}()
	wg.Wait()
	fmt.Print("程序结束")
}

func channelWithBuffer() {
	bufferChan := make(chan int, 10)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(bufferChan)
		for i := 1; i <= 100; i++ {
			fmt.Printf("[生产者] 发送：%d，通道当前缓冲数：%d\n", i, len(bufferChan))
			bufferChan <- i
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for num := range bufferChan {
			fmt.Printf("接收并打印：%d\n", num)
		}
	}()
	wg.Wait()
	fmt.Print("程序结束")
}

// 锁机制练习
type Counter struct {
	mu    sync.Mutex // 互斥锁
	count int
}

func (c *Counter) increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}
func (c *Counter) Get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func testForMutex() {
	var wg sync.WaitGroup
	counter := &Counter{}

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				counter.increment()
			}
		}()
	}
	wg.Wait()
	fmt.Printf("最终计数器值：%d\n", counter.Get())
}

// 原子操作练习
func testForAtomic() {
	var wg sync.WaitGroup
	var counter int64
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Printf("最终计数器值：%d\n", atomic.LoadInt64(&counter))
}
func main() {
	// testForIncrementByTen()
	// testForSafeDoubleEachEle()
	// goroutinePractice()
	// testForSchedual()
	// TestForShape()
	//testForEmplyee()
	//channelStudy()
	//channelWithBuffer()
	//testForMutex()
	testForAtomic()
}
