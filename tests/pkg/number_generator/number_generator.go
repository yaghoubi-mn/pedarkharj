package number_generator

import "fmt"

var number = 9000000000

func GetNumber() string {
	number++
	return fmt.Sprintf("+98%v", number)
}
