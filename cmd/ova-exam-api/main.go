package main

import (
	"fmt"
	utils "ova-exam-api/cmd/ova-exam-api/internal"
)

func main() {
	source1 := []int{1,2,3,4,5,6,7,8,9,10,11,12}
	fmt.Println(source1)
	result1 := utils.Div(source1, 5)
	fmt.Println(result1)

	source2 := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
		"key5": "value5",
		"key6": "value6",
	}
	fmt.Println(source2)
	result2 := utils.Invert(source2)
	fmt.Println(result2)


	source3 := []string{"слово1", "слово2", "слово3", "слово4", "слово5"}
	fmt.Println(source3)
	result3 := utils.Filter(source3)
	fmt.Println(result3)
}
