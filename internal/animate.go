package internal

import (
	"fmt"
	"time"
)

// Функция спиннера
func spinner(done chan bool) {
	chars := []rune{'|', '/', '-', '\\'}

	i := 0
	for {
		select {
		case <-done:
			fmt.Print("\r          \r") // Очистка строки
			return
		default:
			fmt.Printf("\r[%c] Выполняется...", chars[i%len(chars)])
			time.Sleep(100 * time.Millisecond)
			i++
		}
	}
}
