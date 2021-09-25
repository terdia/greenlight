package custom_type

import (
	"errors"
	"strconv"

	"github.com/speps/go-hashids/v2"
)

var ErrInvalidIDFormat = errors.New("invalid ID format")

type ID int64

func NewIdHasher() (*hashids.HashID, error) {
	hd := hashids.NewData()
	hd.Salt = "salt is a salt salt"
	hd.MinLength = 30

	return hashids.NewWithData(hd)
}

func (id ID) MarshalJSON() ([]byte, error) {

	hashId, err := NewIdHasher()
	if err != nil {
		panic(err)
	}

	encodedId, err := hashId.Encode([]int{int(id)})
	if err != nil {
		return nil, ErrInvalidIDFormat
	}

	quotedJSONValue := strconv.Quote(encodedId)

	return []byte(quotedJSONValue), nil
}

func (id *ID) UnmarshalJSON(jsonValue []byte) error {

	unQuotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidIDFormat
	}

	hashId, err := NewIdHasher()
	if err != nil {
		panic(err)
	}

	decodedId, err := hashId.DecodeWithError(unQuotedJSONValue)
	if err != nil {
		return ErrInvalidIDFormat
	}

	// Convert the int32 to a ID type and assign this to the receiver. Note that we
	// use the * operator to deference the receiver (which is a pointer to a ID
	// type) in order to set the underlying value of the pointer.
	*id = ID(decodedId[0])

	return nil
}
