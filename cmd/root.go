/*
Copyright © 2024 Fernando Levin <flevin58@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	DryRun       bool
	Verbose      bool
	IncludedExts []string
	ExcludedExts []string
	ExcludedDirs []string
	ByteExts     []string
	StringExts   []string
)

var rootCmd = &cobra.Command{
	Use:   "embed",
	Short: "Generates go:embed instruction in your assets folder",
	Long: `Parses the given assets folder and generates inside each subfolder
an 'embed.go' file with package name = name of the subfolder and
several go:embed instructions according to file content and given flags.
Example: a file named explosion.ogg inside the sounds folder will generate:

package sound

include _ "embed"

//go:embed explosion.ogg
var Explosion_ogg []byte
`,
	Run: embedCmd,
}

// Entry point from main()
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Parse flags and do some checks
func init() {
	rootCmd.Flags().BoolVarP(&DryRun, "dry-run", "d", false, "Sets verbose mode and does not touch the file system (use for testing your flags)")
	rootCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "Outputs info to standard output")
	rootCmd.Flags().StringSliceVarP(&IncludedExts, "include", "i", []string{}, "Specify extensions to include (all others are excluded)")
	rootCmd.Flags().StringSliceVarP(&ExcludedExts, "exclude", "x", []string{}, "Specify extensions to exclude (all others are included)")
	rootCmd.Flags().StringSliceVarP(&ExcludedDirs, "exclude-dirs", "f", []string{}, "The folder(s) to be excluded from embedding")
	rootCmd.Flags().StringSliceVarP(&ByteExts, "byte", "b", []string{}, "Files with given extension(s) are []byte (all others are string)")
	rootCmd.Flags().StringSliceVarP(&StringExts, "string", "s", []string{}, "Files with given extension(s) are string (all others are []byte)")
}
