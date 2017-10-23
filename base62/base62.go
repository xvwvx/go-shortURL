package base62

import (
	"errors"
	"math"
)

const (
	encodeStd = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	decodeStd   = [256]uint64{}
	decodeError = errors.New("Decode error")
)

func init() {
	encodeStdChar := []byte(encodeStd)
	for i, c := range encodeStdChar {
		decodeStd[int(c)] = uint64(i)
	}
}

func Encode(number uint64) string {
	if number == 0 {
		return string(encodeStd[0])
	}

	var chars = make([]byte, 0, 11)

	length := uint64(len(encodeStd))
	for number > 0 {
		remainder := number % length
		result := number / length
		chars = append(chars, encodeStd[remainder])
		number = result
	}

	return string(chars)
}

func Decode(token string) (number uint64, err error) {
	chars := []byte(token)
	length := float64(len(encodeStd))

	for i, c := range chars {
		decodeValue := decodeStd[c]
		if encodeStd[decodeValue] != c {
			return 0, decodeError
		}
		number += decodeValue * uint64(math.Pow(length, float64(i)))
	}
	return number, nil
}
