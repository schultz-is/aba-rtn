// Copyright (c) 2020 Matt Schultz <matt@schultz.is>. All rights reserved.
// Use of this source code is governed by an ISC license that can be found in
// the LICENSE file.

package rtnutil

import "testing"

func TestValidate(t *testing.T) {
	tests := []struct {
		input    string
		expected error
	}{
		{"asdf", ErrIncorrectLength},
		{"1234", ErrIncorrectLength},
		{"0123456789", ErrIncorrectLength},
		{"R00000000", ErrInvalidCharacter},
		{"123456789", ErrChecksumMismatch},
		{"322286188", nil},
		{"021200025", nil},
		{"111000025", nil},
		{"026014601", nil},
	}

	var actual error
	for _, test := range tests {
		t.Run(
			test.input,
			func(t *testing.T) {
				actual = Validate(test.input)
				if actual != test.expected {
					t.Fatalf(
						"input \"%s\" generated actual output \"%t\" (expected \"%t\")",
						test.input,
						actual,
						test.expected,
					)
				}
			},
		)
	}
}

func TestGetMissingDigit(t *testing.T) {
	tests := []struct {
		input         string
		expectedDigit int
		expectedError error
	}{
		{"asdf", 0, ErrIncorrectLength},
		{"1234", 0, ErrIncorrectLength},
		{"0123456789", 0, ErrIncorrectLength},
		{"XX2286188", 0, ErrTooManyMissingDigits},
		{"R22286188", 0, ErrInvalidCharacter},
		{"322286188", 0, ErrNoMissingDigits},
		{"X22286188", 3, nil},
		{"3X2286188", 2, nil},
		{"32X286188", 2, nil},
		{"322X86188", 2, nil},
		{"3222X6188", 8, nil},
		{"32228X188", 6, nil},
		{"322286X88", 1, nil},
		{"3222861X8", 8, nil},
		{"32228618X", 8, nil},
		{"03110064X", 9, nil},
	}

	var (
		actualDigit int
		actualError error
	)
	for _, test := range tests {
		actualDigit, actualError = GetMissingDigit(test.input)
		if actualDigit != test.expectedDigit {
			t.Fatalf(
				"input \"%s\" generated actual digit \"%d\" (expected \"%d\")",
				test.input,
				actualDigit,
				test.expectedDigit,
			)
		}

		if actualError != test.expectedError {
			t.Fatalf(
				"input \"%s\" generated actual error \"%s\" (expected \"%s\")",
				test.input,
				actualError,
				test.expectedError,
			)
		}
	}
}
