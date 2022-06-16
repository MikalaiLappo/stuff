package bytes

import (
	"strconv"
)


func Btoi(bs []byte) int {
	r, err := strconv.Atoi(string(bs))
	if err != nil {
		panic(err)
	}
	return r
}