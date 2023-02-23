package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}
var workersCount = 1000

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	return output
}

func main() {
	rand.Seed(time.Now().Unix())
	startTime := time.Now()
	var wg sync.WaitGroup
	usersCount := 100
	usersChan := generateUsers(usersCount)
	wg.Add(usersCount)
	for i := 0; i < workersCount; i++ {
		go saveUserInfo(usersChan, &wg)
	}
	wg.Wait()

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func generateUsersWorker(usersIdChan <-chan int, outputChan chan<- User) {
	for id := range usersIdChan {
		time.Sleep(time.Millisecond * 100)
		user := User{
			id:    id,
			email: fmt.Sprintf("user%d@company.com", id),
			logs:  generateLogs(rand.Intn(1000)),
		}

		fmt.Printf("user %d was generated, saving user info in file\n", user.id)
		outputChan <- user
	}
	close(outputChan)
}

func saveUserInfo(inputChan <-chan User, wg *sync.WaitGroup) {
	for user := range inputChan {
		time.Sleep(time.Second)

		filename := fmt.Sprintf("users/uid%d.txt", user.id)
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}

		file.WriteString(user.getActivityInfo())
		fmt.Printf("user %d data was written in file\n", user.id)
		wg.Done()
	}
}

func generateUsers(count int) chan User {
	usersIdChan, outputChan := make(chan int, count), make(chan User, count)
	for i := 0; i < count; i++ {
		usersIdChan <- i + 1
	}
	for i := 0; i < workersCount; i++ {
		go generateUsersWorker(usersIdChan, outputChan)
	}

	return outputChan
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}

	return logs
}
