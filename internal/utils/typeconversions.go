package utils

import (
	"log"
	"strconv"
)

func ConvertStrToFloat32(value string) float32 {
	floatValue, err := strconv.ParseFloat(value, 32)
	if err != nil {
		log.Println("Error converting value to float32: ", err)
	}
	return float32(floatValue)
}

func CelsiusToFahrenheit(c float32) float32 {
	return c*9/5 + 32
}
