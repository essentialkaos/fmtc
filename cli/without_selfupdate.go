//go:build !selfupdate
// +build !selfupdate

package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import "github.com/essentialkaos/ek/v13/usage"

// ////////////////////////////////////////////////////////////////////////////////// //

var withSelfUpdate = false

// ////////////////////////////////////////////////////////////////////////////////// //

// updateBinary updates current binary to the latest version
func updateBinary() int {
	return 1
}

// getUpdateChecker returns update checker
func getUpdateChecker() usage.UpdateChecker {
	return usage.UpdateChecker{}
}

// ////////////////////////////////////////////////////////////////////////////////// //
