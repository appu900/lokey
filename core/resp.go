package core

import (
	"errors"
)

func readInt64(data []byte) (int64, int, error) {
	pos := 1
	var value int64 = 0
	for ; data[pos] != '\r'; pos++ {
		value = value*10 + int64(data[pos]-'0')
	}
	return value, pos + 2, nil

}
func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}


// ** function for read the simple string data
// ** start with a pointer and travere untill we get '\r' and incrememt the position
// ** the build this byte array and return this string.
func readSimpleString(data []byte) (string, int, error) {
	pos := 1
	for ; data[pos] != '\r'; pos++ {
	}
	stringData := data[1:pos]
	return string(stringData), pos + 2, nil
}




func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}
	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		return readError(data)
	case ':':
		return readInt64(data)
	case '$':
		return readSimpleString(data)
	case '*':
		return readSimpleString(data)
	}
	return nil, 0, errors.New("no matching founc")
}
func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}
	value, _, err := DecodeOne(data)
	if err != nil {
		return nil, err
	}
	return value, nil

}
