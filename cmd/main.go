package main

import (
	"ova-exam-api/internal/domain/entity/user"
	"ova-exam-api/internal/flusher"
	"ova-exam-api/internal/repo"
	"ova-exam-api/internal/saver"
	"time"
)

func main() {
	repo := repo.NewRepo("data.txt")
	flusher := flusher.NewFlusher(2, repo)
	saver := saver.NewSaver(3, flusher)
	user1 := user.User{UserId: 1}
	user2 := user.User{UserId: 2}
	user3 := user.User{UserId: 3}
	user4 := user.User{UserId: 4}
	user5 := user.User{UserId: 5}
	user6 := user.User{UserId: 6}

	saver.Save(user1)
	saver.Save(user2)
	time.Sleep(time.Millisecond * 5000)
	saver.Save(user3)
	time.Sleep(time.Millisecond * 1000)
	saver.Save(user4)
	saver.Save(user5)
	saver.Save(user6)
	saver.Close()
}
