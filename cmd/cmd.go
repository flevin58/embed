package cmd

import (
	"github.com/alecthomas/kong"
)

type CLI struct {
	DryRun       bool     `kong:"short='d',help='Sets verbose mode and does not touch the file system (use for testing your flags)'"`
	Verbose      bool     `kong:"short='v',help='Outputs info to standard output'"`
	Folder       string   `kong:"arg,required,type='foldermustexist',help='The folder to process'"`
	IncludedExts []string `kong:"xor='exts',short='i',help='Extensions to include (all others are excluded)'"`
	ExcludedExts []string `kong:"xor='exts',short='x',help='Extensions to exclude (all others are included)'"`
	ExcludedDirs []string `kong:"short='f',help='The subfolder(s) to be excluded from embedding'"`
	ByteExts     []string `kong:"xor='byte',short='b',help='Files with given extension(s) are []byte (all others are string)'"`
	StringExts   []string `kong:"xor='byte',short='s',help='Files with given extension(s) are string (all others are []byte)'"`
}

var cli CLI

func ParseAndRun() {
	ctx := kong.Parse(&cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}

func (c *CLI) Help() string {
	return `Parses the given assets folder and generates inside each subfolder
an 'embed.go' file with package name = name of the subfolder and
several go:embed instructions according to file content and given flags.
Example: a file named explosion.ogg inside the sounds folder will generate:

package sound

include _ "embed"

//go:embed explosion.ogg
var Explosion_ogg []byte
`
}
