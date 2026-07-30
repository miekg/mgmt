// Mgmt
// Copyright (C) James Shubin and the project contributors
// Written by James Shubin <james@shubin.ca> and the project contributors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// Additional permission under GNU GPL version 3 section 7
//
// If you modify this program, or any covered work, by linking or combining it
// with embedded mcl code and modules (and that the embedded mcl code and
// modules which link with this program, contain a copy of their source code in
// the authoritative form) containing parts covered by the terms of any other
// license, the licensors of this program grant you additional permission to
// convey the resulting work. Furthermore, the licensors of this program grant
// the original author, James Shubin, additional permission to update this
// additional permission if he deems it necessary to achieve the goals of this
// additional permission.

package util

import (
	"errors"
	"strconv"
)

// Error is a constant error type that implements error.
type Error string

// Error fulfills the error interface of this type.
func (obj Error) Error() string { return string(obj) }

// ExitCodeError is an error that carries a process exit code.
type ExitCodeError struct {
	Code int
}

// Error fulfills the error interface of this type.
func (obj ExitCodeError) Error() string {
	return "exit " + strconv.Itoa(obj.Code)
}

// ExitCode returns the process exit code.
func (obj ExitCodeError) ExitCode() int {
	return obj.Code
}

// ExitCode returns a process exit code from an error chain. If the error is nil
// it returns 0. If it has an embedded code, it returns that, otherwise it
// returns 1.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}

	// TODO: can we do this without defining this interface?
	type exitCoder interface {
		ExitCode() int
	}

	var exitErr exitCoder
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}

	return 1 // some uncoded error
}
