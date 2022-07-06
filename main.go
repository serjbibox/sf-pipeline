package main

import (
	"bufio"
	"fmt"
	"log"
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

var ilog = log.New(os.Stdout, "Main INFO\t", log.Ldate|log.Ltime)
var elog = log.New(os.Stderr, "Main ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	done, source := dataSource()
	negativeFilter := pipeline.UnitInt(negativeFilterInt)
	divisibleFilter := pipeline.UnitInt(divisibleFilterInt)
	buffering := pipeline.UnitInt(bufferingInt)

	pipeLine := pipeline.PipelineInt(divisibleFilter, negativeFilter, buffering).Setup(done)

	cunsumer(done, pipeLine.Run(source))
	<-time.After(1 * time.Second)
}

func cunsumer(done <-chan struct{}, c <-chan int) {
	for {
		select {
		case data := <-c:
			fmt.Printf("Обработаны данные: %d\n", data)
			ilog.Printf("Обработаны данные: %d", data)
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
				ilog.Printf("Программа завершила работу")
				return
			}
			i, err := strconv.Atoi(data)
			if err != nil {
				fmt.Println("Программа обрабатывает только целые числа!")
				elog.Printf("Введены неверные данные: %v", data)
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
						ilog.Printf("Юнит divisibleFilterInt завершил работу")
						return
					}
				}
				ilog.Printf("Данные попали в фильтр divisibleFilterInt: %v", data)
			case <-done:
				ilog.Printf("Юнит divisibleFilterInt завершил работу")
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
						ilog.Printf("Юнит negativeFilterInt завершил работу")
						return
					}
				}
				ilog.Printf("Данные попали в фильтр negativeFilterInt: %v", data)
			case <-done:
				ilog.Printf("Юнит negativeFilterInt завершил работу")
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
				ilog.Printf("Юнит bufferingInt завершил работу")
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
							ilog.Printf("Модуль задержки юнита bufferingInt завершил работу")
							return
						}
					})
				}
			case <-done:
				ilog.Printf("Модуль задержки юнита bufferingInt завершил работу")
				return
			}
		}
	}()
	return bufferChan
}
