// Contains functions for validating CPF, the Brazilian unique identifacation number for every person in the country.
// CPFs must contain 11 numeric digits, being the last 2 for validation.
// For the purpose of this validation, the first 9 digits can be considered random
// (although there are rules related to the State in which it was emitted, for example).
// The official refference can be found at http://sa.previdencia.gov.br/site/2015/07/rgrv_RegrasValidacao.pdf.
package cpfvalidator

import (
	"fmt"
	"regexp"
	"strconv"
)

type Cpf string

// These values pass the algorithm but are considered invalid.
var knownInvalids = map[string]struct{}{
	"00000000000": {},
	"11111111111": {},
	"22222222222": {},
	"33333333333": {},
	"44444444444": {},
	"55555555555": {},
	"66666666666": {},
	"77777777777": {},
	"88888888888": {},
	"99999999999": {},
}

func IsValid(cpf Cpf) (bool, error) {
	unmasked, err := RemoveMask(cpf)

	if err != nil {

		return false, err
	}

	_, isKnownInvalid := knownInvalids[string(unmasked)]

	if isKnownInvalid {
		return false, nil
	}

	return checkVerifyingDigits(unmasked), nil
}

// Receives a XXX.XXX.XXX-XX or XXXXXXXXXX format CPF and returns always 11 numeric digits.
func RemoveMask(input Cpf) (Cpf, error) {
	inputString := string(input)

	// The error here is unnecessary because the regex is being passsed directly.
	inputIsNumeric, _ := regexp.MatchString(`^\d{11}$`, inputString)

	if inputIsNumeric {

		return Cpf(input), nil
	}

	var unmaskedCpf Cpf

	re := regexp.MustCompile(`^(\d{3})\.(\d{3})\.(\d{3})\-(\d{2})`)
	trimmed := re.ReplaceAllString(inputString, "$1$2$3$4")

	if len(trimmed) == len(inputString) {

		return unmaskedCpf, fmt.Errorf("invalid input")
	}

	unmaskedCpf = Cpf(trimmed)

	return unmaskedCpf, nil
}

func Mask(input Cpf) (Cpf, error) {
	inputString := string(input)
	var maskedCpf Cpf

	// The error here is unnecessary because the regex is being passsed directly.
	inputIsNumeric, _ := regexp.MatchString(`^\d{11}$`, inputString)

	if !inputIsNumeric {
		return maskedCpf, fmt.Errorf("invalid format")
	}

	re := regexp.MustCompile(`^(\d{3})(\d{3})(\d{3})(\d{2})`)
	maskedInput := re.ReplaceAllString(inputString, "$1.$2.$3-$4")

	maskedCpf = Cpf(maskedInput)

	return maskedCpf, nil
}

func convertRestToDigit(dividend, divisor int) string {
	rest := dividend % divisor

	if rest < 2 {

		return "0"
	}

	return strconv.Itoa(11 - rest)
}

func checkVerifyingDigits(cpf Cpf) bool {
	firstVerifyingDigit := iterateDigits(cpf[:9])

	if firstVerifyingDigit != string(cpf[9]) {
		return false
	}

	secondVerifyingDigit := iterateDigits(cpf[:10])

	return secondVerifyingDigit == string(cpf[10])
}

func iterateDigits(cpf Cpf) string {
	cpfString := string(cpf)

	var sum int
	var factor int = len(cpfString) + 1

	for i := 0; i < len(cpfString); i++ {
		var char string = string(cpfString[i])
		digit, _ := strconv.Atoi(char)
		sum += digit * factor
		factor--

	}

	digit := convertRestToDigit(sum, 11)

	return digit
}
