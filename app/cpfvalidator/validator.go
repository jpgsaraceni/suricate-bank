// Contains functions for validating CPF, the Brazilian unique identifacation number for every person in the country.
// CPFs must contain 11 numeric digits, being the last 2 for validation.
// For the purpose of this validation, the first 9 digits can be considered random
// (although there are rules related to the State in which it was emitted, for example).
// The official refference can be found at http://sa.previdencia.gov.br/site/2015/07/rgrv_RegrasValidacao.pdf.
package cpfvalidator

import (
	"errors"
	"regexp"
	"strconv"
)

type Cpf string

// These values pass the algorithm but are considered invalid.
var knownInvalids = map[Cpf]struct{}{
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

// Runs CPF algorithm to check if CPF is valid. Accepts XXX.XXX.XXX-XX and XXXXXXXXXXX formats.
func (c Cpf) IsValid() (bool, error) {
	unmasked, err := c.RemoveMask()

	if err != nil {

		return false, err
	}

	_, isKnownInvalid := knownInvalids[unmasked]

	if isKnownInvalid {
		return false, nil
	}

	return checkVerifyingDigits(unmasked), nil
}

var errInput = errors.New("invalid input")

// Converts a XXX.XXX.XXX-XX or XXXXXXXXXX format CPF to 11 numeric digits.
func (c Cpf) RemoveMask() (Cpf, error) {
	cString := string(c)

	// The error here is unnecessary because the regex is being passsed directly.
	inputIsNumeric, _ := regexp.MatchString(`^\d{11}$`, cString)

	if inputIsNumeric {

		return c, nil
	}

	var unmaskedCpf Cpf

	re := regexp.MustCompile(`^(\d{3})\.(\d{3})\.(\d{3})\-(\d{2})`)
	trimmed := re.ReplaceAllString(cString, "$1$2$3$4")

	if len(trimmed) == len(cString) {

		return unmaskedCpf, errInput
	}

	unmaskedCpf = Cpf(trimmed)

	return unmaskedCpf, nil
}

var errFormat = errors.New("invalid format")

// Converts an 11 long numeric string to XXX.XXX.XXX-XX format
func (c Cpf) Mask() (Cpf, error) {
	cString := string(c)
	var maskedCpf Cpf

	// The error here is unnecessary because the regex is being passsed directly.
	inputIsNumeric, _ := regexp.MatchString(`^\d{11}$`, cString)

	if !inputIsNumeric {
		return maskedCpf, errFormat
	}

	re := regexp.MustCompile(`^(\d{3})(\d{3})(\d{3})(\d{2})`)
	maskedInput := re.ReplaceAllString(cString, "$1.$2.$3-$4")

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
	var sum int
	var factor int = len(cpf) + 1

	for i := 0; i < len(cpf); i++ {
		var char = string(cpf[i])
		digit, _ := strconv.Atoi(char)
		sum += digit * factor
		factor--

	}

	digit := convertRestToDigit(sum, 11)

	return digit
}
