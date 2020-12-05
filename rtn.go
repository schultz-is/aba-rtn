// Copyright (c) 2020 Matt Schultz <matt@schultz.is>. All rights reserved.
// Use of this source code is governed by an ISC license that can be found in
// the LICENSE file.

// Package rtnutil provides utilities for working with American Bankers
// Association routing transit numbers; also known as ABA RTNs.
//
// For more information about the structure and composition of ABA RTNs, please
// see [1] [2].
//
// [1] https://en.wikipedia.org/wiki/ABA_routing_transit_number
// [2] http://s3-eu-west-1.amazonaws.com/cjp-rbi-accuity/wp-content/uploads/2016/09/21222732/ROUTING_NUMBER_POLICY-2016.pdf
package rtnutil

import (
	"errors"
)

// ErrIncorrectLength indicates that an RTN is not the correct length of 9
// characters.
var ErrIncorrectLength = errors.New("incorrect length")

// ErrInvalidCharacter indicates that an RTN contains a character which is not
// valid for the current context.
var ErrInvalidCharacter = errors.New("invalid character")

// ErrChecksumMismatch indicates that the check digit of an RTN does not match
// the remaining digits.
var ErrChecksumMismatch = errors.New("checksum mismatch")

// ErrTooManyMissingDigits indicates that a provided RTN has more than a single
// missing digit.
var ErrTooManyMissingDigits = errors.New("too many missing digits")

// ErrNoMissingDigits indicates that a provided RTN contains no missing digits
// when one is expected.
var ErrNoMissingDigits = errors.New("no missing digits")

// checksumMultipliers is a set of numbers that multiply RTN digits to
// calculate a checksum.
var checksumMultipliers = []int{3, 7, 1}

// Validate determins whether a provided RTN is in valid MICR format with a
// correct check digit.
func Validate(rtn string) (err error) {
	// MICR RTNs are 9 digits
	if len(rtn) != 9 {
		return ErrIncorrectLength
	}

	var (
		i         int
		digitRune rune
		digit     int
		ok        bool
		checksum  int
	)

	// Iterate over each character in the string
	for i, digitRune = range rtn {
		// Attempt to convert the character to a digit
		digit, ok = runeToDigit(digitRune)
		if !ok {
			return ErrInvalidCharacter
		}

		// Multiply the digit by its respective multiplier and add to the checksum
		checksum += digit * checksumMultipliers[i%3]
	}

	// If the checksum is not evenly divisible by 10, the RTN is invalid
	if checksum%10 != 0 {
		return ErrChecksumMismatch
	}

	return nil
}

// GetMissingDigit calculates a single unknown digit within the provided RTN
// Input must be an RTN in MICR format with a single digit replaced by the
// character 'X'.
func GetMissingDigit(rtn string) (digit int, err error) {
	if len(rtn) != 9 {
		return 0, ErrIncorrectLength
	}

	var (
		i                 int
		digitRune         rune
		missingMultiplier int
		ok                bool
		checksum          int
	)

	// Iterate over each character in the string
	for i, digitRune = range rtn {
		// Check for the "missing digit" rune
		if digitRune == 'X' {
			// If the missing multiplier has already been set, there are too many
			// digits missing from the provided RTN
			if missingMultiplier > 0 {
				return 0, ErrTooManyMissingDigits
			}

			// Set the multiplier for the missing digit based on its index
			missingMultiplier = checksumMultipliers[i%3]
			continue
		}

		// Attempt to convert the character to a digit
		digit, ok = runeToDigit(digitRune)
		if !ok {
			return 0, ErrInvalidCharacter
		}

		// Multiply the digit by its respective multiplier and add to the checksum
		checksum += digit * checksumMultipliers[i%3]
	}

	// If the missing multiplier was never set, no digits were missing from the
	// provided RTN
	if missingMultiplier == 0 {
		return 0, ErrNoMissingDigits
	}

	// Check digits 0-8 to see if they satisfy the checksum
	for i = 0; i < 9; i++ {
		if (checksum+(missingMultiplier*i))%10 == 0 {
			return i, nil
		}
	}

	// If it's not 0-8, it can only be 9
	return 9, nil
}

// runeToDigit attempts to convert the provided rune into a digit.
func runeToDigit(r rune) (digit int, ok bool) {
	switch r {
	case '0':
		return 0, true
	case '1':
		return 1, true
	case '2':
		return 2, true
	case '3':
		return 3, true
	case '4':
		return 4, true
	case '5':
		return 5, true
	case '6':
		return 6, true
	case '7':
		return 7, true
	case '8':
		return 8, true
	case '9':
		return 9, true
	}

	return 0, false
}
