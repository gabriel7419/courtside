// Package assets provides embedded static assets for the golazo application.
package assets

import (
	_ "embed"
)

// Logo is the golazo logo PNG image, embedded at compile time.
// Used for desktop notifications on Linux and Windows.
//
//go:embed golazo-logo.png
var Logo []byte
