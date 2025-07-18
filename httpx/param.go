package httpx

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/boostgo/core/convert"

	"github.com/google/uuid"
)

func NewParseIntParamError(err error, value string) error {
	return ErrParseIntParam.
		SetError(err).
		SetData(parseContext{
			Value: value,
		})
}

func NewParseFloatParam(err error, value string) error {
	return ErrParseFloatParam.
		SetError(err).
		SetData(parseContext{
			Value: value,
		})
}

func NewParseUUIDParam(err error, value string) error {
	return ErrParseUUIDParam.
		SetError(err).
		SetData(parseContext{
			Value: value,
		})
}

type Param struct {
	value string
}

func NewParam(value string) Param {
	return Param{
		value: value,
	}
}

func EmptyParam() Param {
	return NewParam("")
}

func IsEmptyParam(param Param) bool {
	return param.IsEmpty()
}

func ParamEquals(p1, p2 Param) bool {
	return p1.value == p2.value
}

func (param Param) IsEmpty() bool {
	return param.value == ""
}

func (param Param) Equals(compare Param) bool {
	return ParamEquals(param, compare)
}

func (param Param) String(defaultValue ...string) string {
	if param.value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return param.value
}

func (param Param) Strings() []string {
	return strings.Split(param.value, ",")
}

func (param Param) IntArray() []int {
	integers := make([]int, 0)
	split := strings.Split(param.value, ",")
	for _, value := range split {
		if value == "" {
			continue
		}

		integers = append(integers, convert.Int(value))
	}
	return integers
}

func (param Param) Int() (int, error) {
	intValue, err := strconv.Atoi(param.value)
	if err != nil {
		return 0, NewParseIntParamError(err, param.value)
	}

	return intValue, nil
}

func (param Param) MustInt(defaultValue ...int) int {
	intValue, err := strconv.Atoi(param.value)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return 0
	}

	return intValue
}

func (param Param) Int64() (int64, error) {
	intValue, err := strconv.ParseInt(param.value, 10, 64)
	if err != nil {
		return 0, NewParseIntParamError(err, param.value)
	}

	return intValue, nil
}

func (param Param) MustInt64(defaultValue ...int64) int64 {
	intValue, err := strconv.ParseInt(param.value, 10, 64)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return 0
	}

	return intValue
}

func (param Param) Int32() (int32, error) {
	intValue, err := strconv.ParseInt(param.value, 10, 64)
	if err != nil {
		return 0, NewParseIntParamError(err, param.value)
	}

	return int32(intValue), nil
}

func (param Param) MustInt32(defaultValue ...int32) int32 {
	intValue, err := strconv.ParseInt(param.value, 10, 32)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return 0
	}

	return int32(intValue)
}

func (param Param) Float32() (float32, error) {
	floatValue, err := strconv.ParseFloat(param.value, 32)
	if err != nil {
		return 0, NewParseFloatParam(err, param.value)
	}

	return float32(floatValue), nil
}

func (param Param) Float64() (float64, error) {
	floatValue, err := strconv.ParseFloat(param.value, 64)
	if err != nil {
		return 0, NewParseFloatParam(err, param.value)
	}

	return floatValue, nil
}

func (param Param) Bool() bool {
	return strings.ToLower(param.value) == "true"
}

func (param Param) UUID() (uuid.UUID, error) {
	uuidValue, err := uuid.Parse(param.value)
	if err != nil {
		return uuid.UUID{}, NewParseUUIDParam(err, param.value)
	}

	return uuidValue, nil
}

func (param Param) MustUUID() uuid.UUID {
	uuidValue, err := uuid.Parse(param.value)
	if err != nil {
		return uuid.UUID{}
	}

	return uuidValue
}

func (param Param) Bytes() []byte {
	return convert.BytesFromString(param.value)
}

func (param Param) Parse(export any) error {
	return json.Unmarshal(param.Bytes(), export)
}
