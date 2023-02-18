package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}
var workersCount = 100

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

	users := generateUsers(100)

	saveUsersInfo(users)

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func generateUsersWorker(inputChan <-chan bool, outputChan chan<- User, index int) {
	for _ = range inputChan {
		time.Sleep(time.Millisecond * 100)
		outputChan <- User{
			id:    index,
			email: fmt.Sprintf("user%d@company.com", index+1),
			logs:  generateLogs(rand.Intn(1000)),
		}
	}
}

func saveUsersInfo(users []User) {
	inputChan, outputChan := make(chan User, len(users)), make(chan User, len(users))
	for _, user := range users {
		inputChan <- user
	}
	close(inputChan)
	for i := 0; i < workersCount; i++ {
		go saveUserInfo(inputChan, outputChan)
	}

	for i := 0; i < workersCount; i++ {
		user := <-outputChan
		fmt.Printf("WRITING FILE FOR UID %d\n", user.id)
	}
}

func saveUserInfo(inputChan <-chan User, outputChan chan<- User) {
	for user := range inputChan {
		time.Sleep(time.Second)

		filename := fmt.Sprintf("users/uid%d.txt", user.id)
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}

		file.WriteString(user.getActivityInfo())
		outputChan <- user
	}
}

func generateUsers(count int) []User {
	users := make([]User, count)
	inputChan, outputChan := make(chan bool, count), make(chan User, count)
	for i := 0; i < count; i++ {
		inputChan <- true
	}
	for i := 0; i < workersCount; i++ {
		go generateUsersWorker(inputChan, outputChan, i+1)
	}
	close(inputChan)
	for i := 0; i < count; i++ {
		users[i] = <-outputChan
		fmt.Printf("generated user %d\n", users[i].id)
	}

	return users
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
