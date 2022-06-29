package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/serjbibox/SF20.2.1/pipeline"
)

const (
	bufferDrainInterval time.Duration = 15 * time.Second
	bufferSize          int           = 5
	divisibleBy         int           = 3
)

func main() {
	done, source := dataSource()

	negativeFilter := pipeline.UnitInt(negativeFilterInt)
	divisibleFilter := pipeline.UnitInt(divisibleFilterInt)
	buffering := pipeline.UnitInt(bufferingInt)

	pipeLine := pipeline.PipelineInt(divisibleFilter, negativeFilter, buffering).Setup(done)

	cunsumer(done, pipeLine.Run(source))
}

func cunsumer(done <-chan struct{}, c <-chan int) {
	for {
		select {
		case data := <-c:
			fmt.Printf("Обработаны данные: %d\n", data)
		case <-done:
			return
		}
	}
}

func dataSource() (<-chan struct{}, <-chan int) {
	c := make(chan int)
	done := make(chan struct{})
	go func() {
		defer close(done)
		scanner := bufio.NewScanner(os.Stdin)
		var data string
		fmt.Println("Вводите целые числа. Для выхода наберите exit:")
		for {
			scanner.Scan()
			data = scanner.Text()
			if strings.EqualFold(data, "exit") {
				fmt.Println("Программа завершила работу!")
				return
			}
			i, err := strconv.Atoi(data)
			if err != nil {
				fmt.Println("Программа обрабатывает только целые числа!")
				continue
			}
			c <- i
		}
	}()
	return done, c
}

func divisibleFilterInt(done <-chan struct{}, c <-chan int) <-chan int {
	filterChan := make(chan int)
	go func() {
		for {
			select {
			case data := <-c:
				if data != 0 && data%divisibleBy == 0 {
					select {
					case filterChan <- data:
					case <-done:
						return
					}
				}
			case <-done:
				return
			}
		}
	}()
	return filterChan
}

func negativeFilterInt(done <-chan struct{}, c <-chan int) <-chan int {
	filterChan := make(chan int)
	go func() {
		for {
			select {
			case data := <-c:
				if data > 0 {
					select {
					case filterChan <- data:
					case <-done:
						return
					}
				}
			case <-done:
				return
			}
		}
	}()
	return filterChan
}

func bufferingInt(done <-chan struct{}, c <-chan int) <-chan int {
	bufferChan := make(chan int)
	buffer := pipeline.NewRingBuffer(bufferSize)
	go func() {
		for {
			select {
			case data := <-c:
				buffer.Push(data)
			case <-done:
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-time.After(bufferDrainInterval):
				if buffer.Pos() >= 0 {
					buffer.Get().Do(func(p interface{}) {
						select {
						case bufferChan <- p.(int):
						case <-done:
							return
						}
					})
				}
			case <-done:
				return
			}
		}
	}()
	return bufferChan
}
