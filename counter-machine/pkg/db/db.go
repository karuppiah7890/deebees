package db

var counter = 0

func start(dbChannel chan int, done chan bool) {
	over := false
	for !over {
		select {
		case incrementBy := <-dbChannel:
			counter += incrementBy

		case <-done:
			over = true
		}
	}
}

func Start(dbChannel chan int, done chan bool) {
	go start(dbChannel, done)
}

func GetCounter() int {
	return counter
}
