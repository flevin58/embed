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
	ByteExts     []string `kong:"xor='byte',short='b'"`
	StringExts   []string `kong:"xor='byte',short='s'"`
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

// Parse flags and do some checks
// func init() {
// 	rootCmd.Flags().BoolVarP(&DryRun, "dry-run", "d", false, "Sets verbose mode and does not touch the file system (use for testing your flags)")
// 	rootCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "Outputs info to standard output")
// 	rootCmd.Flags().StringSliceVarP(&IncludedExts, "include", "i", []string{}, "Specify extensions to include (all others are excluded)")
// 	rootCmd.Flags().StringSliceVarP(&ExcludedExts, "exclude", "x", []string{}, "Specify extensions to exclude (all others are included)")
// 	rootCmd.Flags().StringSliceVarP(&ExcludedDirs, "exclude-dirs", "f", []string{}, "The folder(s) to be excluded from embedding")
// 	rootCmd.Flags().StringSliceVarP(&ByteExts, "byte", "b", []string{}, "Files with given extension(s) are []byte (all others are string)")
// 	rootCmd.Flags().StringSliceVarP(&StringExts, "string", "s", []string{}, "Files with given extension(s) are string (all others are []byte)")
// }
