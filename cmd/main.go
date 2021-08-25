package main

import (
	"fmt"
	"log"
	"os"
	"ova-exam-api/internal/domain/entity"
	"ova-exam-api/internal/domain/entity/user"
	"ova-exam-api/internal/utils"
	"path"
)

func OpenAndCloseFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("open file", file)
	defer func() {
		fmt.Println("close file", file)
		file.Close()
	}()
}

func main() {
	// subtask a
	pwd, _ := os.Getwd()
	fileName := path.Join(pwd, "Makefile")

	for i := 0; i < 10; i++ {
		OpenAndCloseFile(fileName)
	}

	// subtask b
	us := user.User{
		UserId:   1,
		Entity:   entity.Entity{},
		Email:    "12",
		Password: "12",
	}

	fmt.Println(us.String())

	// subtask c
	source1 := []user.User{
		{UserId: 1},
		{UserId: 2},
		{UserId: 3},
		{UserId: 4},
		{UserId: 5},
		{UserId: 6},
	}
	fmt.Println(source1)
	result1 := utils.SplitToBulks(source1, 2)
	fmt.Println(result1)

	if result2, err := utils.UsersToMap(source1); err == nil{
		fmt.Println(result2)
	}
}
