package saver

import (
	"fmt"
	"ova-exam-api/internal/domain/entity/user"
	"ova-exam-api/internal/flusher"
	"sync"
	"time"
)

type Saver interface {
	Save(entity user.User)
	Close()
}

// NewSaver возвращает Saver с поддержкой переодического сохранения
func NewSaver(
	capacity uint,
	flusher flusher.Flusher,
// ...
) Saver {
	saver := &saver{
		mutex: sync.Mutex{},
		users: make([]user.User, 0, capacity),
		capacity: capacity,
		flusher:  flusher,
	}

	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
			saver.Close()
		}
	}()

	// interruptCh := make(chan os.Signal, 1)
	// signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)
	// fmt.Printf("Got %v...\n", <-interruptCh)

	return saver
}

type saver struct{
	users []user.User
	capacity uint
	flusher flusher.Flusher
	mutex sync.Mutex
}

func (s *saver) Save(entity user.User) {
	fmt.Println("Save on saver called")
	s.users = append(s.users, entity)
	if len(s.users) < int(s.capacity){
		return
	}
}

func (s *saver) Close() {
	fmt.Println("Close on saver called")
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.users) == 0{
		return
	}

	res := s.flusher.Flush(s.users)
	if res != nil {
		s.users = res
		return
	}

	s.users = s.users[:0]
}



