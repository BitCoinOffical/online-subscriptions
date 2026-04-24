package rules

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidatePrice(fl validator.FieldLevel) bool {
	val := fl.Field().Float()
	valstr := strconv.FormatFloat(val, 'f', -1, 64)
	sides := strings.Split(valstr, ".")

	Lside := len(sides[0])
	Rside := len(sides[1])

	return Lside <= 12 && Rside <= 2
}
