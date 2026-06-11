// Package templates embeds the public template trees, one subdirectory per
// preset. The "all:" prefix is required so dot-directories like .context are
// included in the embed.
package templates

import "embed"

//go:embed all:core all:minimal
var FS embed.FS
