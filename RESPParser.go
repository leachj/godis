package parser

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

// ParseRESP parses a Redis RESP message from the given reader.
func ParseRESP(reader io.Reader) (interface{}, error) {
	bufReader := bufio.NewReader(reader)
	messageType, err := bufReader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch messageType {
	case '+':
		return parseSimpleString(bufReader)
	case '-':
		return parseError(bufReader)
	case ':':
		return parseInteger(bufReader)
	case '$':
		return parseBulkString(bufReader)
	case '*':
		return parseArray(bufReader)
	default:
		return nil, errors.New("invalid message type")
	}
}

func parseSimpleString(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return line[:len(line)-2], nil
}

func parseError(reader *bufio.Reader) (error, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	return errors.New(line[:len(line)-2]), nil
}

func parseInteger(reader *bufio.Reader) (int64, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(line[:len(line)-2], 10, 64)
}

func parseBulkString(reader *bufio.Reader) (*string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return new(string), err
	}
	length, err := strconv.Atoi(line[:len(line)-2])
	if err != nil {
		return new(string), err
	}
	if length == -1 {
		return nil, nil
	}
	buf := make([]byte, length+2)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return new(string), err
	}
	result := string(buf[:len(buf)-2])
	return &result, nil
}

func parseArray(reader *bufio.Reader) (*[]interface{}, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	length, err := strconv.Atoi(line[:len(line)-2])
	if err != nil {
		return nil, err
	}
	if length == -1 {
		return nil, nil
	}
	array := make([]interface{}, length)
	for i := 0; i < length; i++ {
		value, err := ParseRESP(reader)
		if err != nil {
			return nil, err
		}
		array[i] = value
	}
	result := array
	return &result, nil
}
