package custom_type

import (
	"errors"
	"strconv"
)

type StatusMessage int64

const (
	Success StatusMessage = iota
	Fail
	Error
)

func (r StatusMessage) MarshalJSON() ([]byte, error) {

	switch r {
	case Success:
		return []byte(strconv.Quote("success")), nil
	case Fail:
		return []byte(strconv.Quote("fail")), nil
	case Error:
		return []byte(strconv.Quote("error")), nil
	}

	return nil, errors.New("unsupported response status")

}
