// Package cpfvalidator contains functions for validating CPF, the Brazilian unique identifacation number for every person in the country.
// CPFs must contain 11 numeric digits, being the last 2 for validation.
// For the purpose of this validation, the first 9 digits can be considered random
// (although there are rules related to the State in which it was emitted, for example).
// The official reference can be found at http://sa.previdencia.gov.br/site/2015/07/rgrv_RegrasValidacao.pdf.
package cpf

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
)

type Cpf struct {
	value  string // 11 numeric digits
	masked string // XXX.XXX.XXX-XX format
}

// knownInvalids pass the algorithm but are considered invalid.
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

var (
	errInvalid   = errors.New("invalid cpf")
	errScan      = errors.New("scan failed")
	errScanEmpty = errors.New("scan returned empty")
)

const (
	unmaskedCpfLength = 11
	module            = 11
	minimumRest       = 2
	maxRandom         = 10
)

// NewCpf creates a Cpf struct with value and masked, or empty and returns error if invalid
func NewCpf(input string) (Cpf, error) {
	var c Cpf

	if ok := c.validate(input); !ok {
		return c, errInvalid
	}

	return c, nil
}

// Random generates a random valid cpf
func Random() Cpf {
	var body string

	for i := 0; i < 9; i++ {
		randomDigit, _ := rand.Int(rand.Reader, big.NewInt(maxRandom))
		body += fmt.Sprintf("%d", int(randomDigit.Int64()))
	}

	body += iterateDigits(body)
	body += iterateDigits(body)

	_, isKnownInvalid := knownInvalids[body]

	for isKnownInvalid {
		body = ""

		for i := 0; i < 9; i++ {
			randomDigit, _ := rand.Int(rand.Reader, big.NewInt(maxRandom))
			body += fmt.Sprintf("%d", int(randomDigit.Int64()))
		}

		body += iterateDigits(body)
		body += iterateDigits(body)

		_, isKnownInvalid = knownInvalids[body]
	}

	generatedCpf, _ := NewCpf(body)

	return generatedCpf
}

// Value returns a cpf with only numeric digits
func (c Cpf) Value() string {
	return c.value
}

// Scan implements database/sql/driver Scanner interface.
// Scan parses a string value to Cpf (if valid) or returns error.
func (c *Cpf) Scan(value interface{}) error {
	if value == nil {
		*c = Cpf{}

		return errScanEmpty
	}

	if value, ok := value.(string); ok {
		cpf, err := NewCpf(value)
		if err != nil {
			return err
		}

		*c = cpf

		return nil
	}

	return errScan
}

// Masked returns a cpf in XXX.XXX.XXX-XX format
func (c Cpf) Masked() string {
	return c.masked
}

// validate runs CPF algorithm to check if CPF is valid and sets Cpf fields. Accepts XXX.XXX.XXX-XX and XXXXXXXXXXX formats.
func (c *Cpf) validate(inputCpf string) bool {
	unmasked, invalidFormat := removeMask(inputCpf)

	if invalidFormat != nil {
		return false
	}

	_, isKnownInvalid := knownInvalids[unmasked]

	if isKnownInvalid {
		return false
	}

	isValid := checkVerifyingDigits(unmasked)

	if isValid {
		c.value = unmasked
		c.masked = mask(inputCpf)

		return true
	}

	return false
}

// removeMask converts a XXX.XXX.XXX-XX or XXXXXXXXXX format CPF to 11 numeric digits.
func removeMask(masked string) (string, error) {
	// The error here is unnecessary because the regex is being passsed directly.
	inputIsNumeric, _ := regexp.MatchString(`^\d{11}$`, masked)

	if inputIsNumeric {
		return masked, nil
	}

	var unmaskedCpf string

	re := regexp.MustCompile(`^(\d{3})\.(\d{3})\.(\d{3})\-(\d{2})`)
	unmaskedCpf = re.ReplaceAllString(masked, "$1$2$3$4")

	if len(unmaskedCpf) != unmaskedCpfLength {
		return "", errInvalid
	}

	return unmaskedCpf, nil
}

// mask converts an 11 long numeric string to XXX.XXX.XXX-XX format, without validating
func mask(cString string) string {
	re := regexp.MustCompile(`^(\d{3})(\d{3})(\d{3})(\d{2})`)
	masked := re.ReplaceAllString(cString, "$1.$2.$3-$4")

	return masked
}

func convertRestToDigit(dividend, divisor int) string {
	rest := dividend % divisor

	if rest < minimumRest {
		return "0"
	}

	return strconv.Itoa(module - rest)
}

func checkVerifyingDigits(cpf string) bool {
	firstVerifyingDigit := iterateDigits(cpf[:9])

	if firstVerifyingDigit != string(cpf[9]) {
		return false
	}

	secondVerifyingDigit := iterateDigits(cpf[:10])

	return secondVerifyingDigit == string(cpf[10])
}

func iterateDigits(cpf string) string {
	var sum int
	factor := len(cpf) + 1

	for i := 0; i < len(cpf); i++ {
		char := string(cpf[i])
		digit, _ := strconv.Atoi(char)
		sum += digit * factor
		factor--
	}

	digit := convertRestToDigit(sum, module)

	return digit
}
