package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Ошибка: передайте параметр K (шаг тикера в секундах)")
		os.Exit(1)
	}

	kSeconds, err := strconv.ParseUint(os.Args[1], 10, 32)
	if err != nil {
		fmt.Printf("Ошибка парсинга параметра K: %v\n", err)
		os.Exit(1)
	}
	k := uint(kSeconds)

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	doneChan := make(chan struct{})

	go func() {
		var tickNum int = 1
		var timeElapsed uint = 0

		for {

			time.Sleep(time.Duration(k) * time.Second)

			select {
			case <-doneChan:
				return
			default:
				timeElapsed += k
				fmt.Printf("Tick %d since %d\n", tickNum, timeElapsed)
				tickNum++
			}
		}
	}()
	//t.Println("Termination")
	<-sigChan

	close(doneChan)
	fmt.Println("Termination")
}

/*
   "А close(doneChan) успеет закрыть горутину прежде чем main завершится?"

   Нет, в данном коде горутина гарантированно НЕ успеет завершиться до того, как закроется main.

   Как только в главной горутине срабатывает close(doneChan), управление моментально переходит
   к следующей строчке: fmt.Println("Termination"). Сразу после вывода этой строки функция main завершает свою работу.
   В Go завершение функции main мгновенно уничтожает всё приложение вместе со всеми фоновыми горутинами,
   независимо от того, что они делали. Фоновая горутина просто не успеет проснуться от time.Sleep и выполнить return.




   Как сделать так, чтобы main дождался горутину?

   Если вам критически важно, чтобы фоновая горутина успела корректно завершиться
   (например, записать данные в файл или закрыть сетевые соединения), нужно использовать sync.WaitGroup.

      Использование sync.WaitGroup подразумевает, что у нас будет счетчик, который мы увеличим на 1 Add(1),
      а потом после close(doneChan) напишем wg.Wait(), который будет ждать обнуления счетчика запущенных горутин
        Внутри запускаемой горутины напишем defer wg.Done(), чтобы она при завершении вычла единицу и обнулила счетчик запущенных горутин

   	 А main будет ждать и не завершаться пока счетчик не обнулится: wg.Wait() - это обеспечит корректное завершение горутин и обработку их результатов, если это необходимо
*/
