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

	for i := range decodeStd {
		decodeStd[i] = math.MaxUint64
	}

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
		chars = append(chars, encodeStd[remainder])
		number = number / length
	}

	chars = reverse(chars)

	return string(chars)
}

func Decode(token string) (number uint64, err error) {
	chars := reverse([]byte(token))
	length := float64(len(encodeStd))

	for i, c := range chars {
		decodeValue := decodeStd[c]
		if decodeValue == math.MaxUint64 {
			return 0, decodeError
		}
		number += decodeValue * uint64(math.Pow(length, float64(i)))
	}
	return number, nil
}

func reverse(array []byte) []byte {
	for i, j := 0, len(array)-1; i < j; i, j = i+1, j-1 {
		array[i], array[j] = array[j], array[i]
	}
	return array
}
