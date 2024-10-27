package utils

import "encoding/json"

func StringP(stringToConvert string) *string {
	return &stringToConvert
}

func MustUnmarshal(data []byte, v any) {
	err := json.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
}

func MustMarshal(v any) []byte {
	result, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return result
}
