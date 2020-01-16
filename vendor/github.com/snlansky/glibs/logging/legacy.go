/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package logging

import (
	"fmt"
	"io"
	"os"
	"strings"

	gologging "github.com/op/go-logging"
)

// These interfaces are used by the chaincode shim at the 1.2 version.
// If we remove them, vendored shims are unlikely to compile against
// newer levels of the peer.

// SetFormat(string) logging.Formatter
// InitBackend(logging.Formatter, io.Writer)
// DefaultLevel() string
// InitFromSpec(string) string

// SetFormat sets the logging format.
func SetFormat(formatSpec string) gologging.Formatter {
	if formatSpec == "" {
		formatSpec = defaultFormat
	}
	return gologging.MustStringFormatter(formatSpec)
}

// InitBackend sets up the logging backend based on
// the provided logging formatter and I/O writer.
func InitBackend(formatter gologging.Formatter, output io.Writer) {
	backend := gologging.NewLogBackend(output, "", 0)
	backendFormatter := gologging.NewBackendFormatter(backend, formatter)
	gologging.SetBackend(backendFormatter).SetLevel(gologging.INFO, "")
}

// DefaultLevel returns the fallback value for loggers to use if parsing fails.
func DefaultLevel() string {
	return strings.ToUpper(Global.DefaultLevel().String())
}

// InitFromSpec initializes the logging based on the supplied spec. It is
// exposed externally so that consumers of the flogging package may parse their
// own logging specification. The logging specification has the following form:
//		[<logger>[,<logger>...]=]<level>[:[<logger>[,<logger>...]=]<logger>...]
func InitFromSpec(spec string) string {
	err := Global.ActivateSpec(spec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to activate logging spec: %s", err)
	}
	return DefaultLevel()
}
