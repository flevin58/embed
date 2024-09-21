package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/flevin58/embed/resources/templates"
	"github.com/flevin58/embed/tools"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func embedCmd(cmd *cobra.Command, args []string) {

	if len(args) != 1 {
		cmd.PrintErrln("You must specify the source folder")
		os.Exit(1)
	}

	startFolder := args[0]

	// Check that -x and -i are not used together
	if len(IncludedExts) > 0 && len(ExcludedExts) > 0 {
		cmd.PrintErrln("Flags -i and -x can't be used together")
		os.Exit(1)
	}

	// Check that -s and -b are not used together
	if len(StringExts) > 0 && len(ByteExts) > 0 {
		cmd.PrintErrln("Flags -b and -s can't be used together")
		os.Exit(1)
	}

	// Fix dots in extensions to make them compatible with path.Ext()
	IncludedExts = tools.FixDot(IncludedExts)
	ExcludedExts = tools.FixDot(ExcludedExts)
	StringExts = tools.FixDot(StringExts)
	ByteExts = tools.FixDot(ByteExts)

	cmd.Printf("Root folder.: %v\n", args[0])
	cmd.Printf("DryRun......: %v\n", DryRun)
	cmd.Printf("IncludedExts: %v\n", IncludedExts)
	cmd.Printf("ExcludedExts: %v\n", ExcludedExts)
	cmd.Printf("ExcludedDirs: %v\n", ExcludedDirs)
	cmd.Printf("ByteExts....: %v\n", ByteExts)
	cmd.Printf("StringExts..: %v\n", StringExts)

	if err := TraverseDir(startFolder); err != nil {
		cmd.PrintErrln(err.Error())
		os.Exit(1)
	}
	if Verbose {
		cmd.Println("Done.")
	}
}

func IsItemIn(item string, list []string) bool {
	// Border cases
	if len(list) == 0 || list[0] == "none" {
		return false
	}
	if list[0] == "all" {
		return true
	}

	// Now check each item
	for _, elem := range list {
		if elem == item {
			return true
		}
	}
	return false
}

func OkToEmbedFile(filename string) bool {

	// If we defined the -x flag
	if len(ExcludedExts) > 0 {
		return !IsItemIn(path.Ext(filename), ExcludedExts)
	}

	// If we defined the -i flag
	if len(IncludedExts) > 0 {
		return IsItemIn(path.Ext(filename), IncludedExts)
	}

	// Neither flag defined, assume all files ok
	return true
}

func OkToEmbedDir(foldername string) bool {
	return !IsItemIn(foldername, ExcludedDirs)
}

// Produces the embed.go file in the current folder
func ProduceEmbedGo(root string, files []string) error {

	// Define the template variables
	type Entry struct {
		File string
		Var  string
		Type string
	}
	type Embed struct {
		Package string
		Entries []Entry
	}

	// Initialize them as empty
	embed := Embed{
		Package: path.Base(root),
		Entries: make([]Entry, 0),
	}

	// Populate with each file
	const (
		TypeByte   = "[]byte"
		TypeString = "string"
	)

	for _, file := range files {
		fvar := strings.Replace(file, ".", "_", -1)
		fvar = cases.Title(language.Und).String(fvar)

		// Set the Type to the default value
		var ftype = TypeByte

		// Flag -s modifies the default if items are included
		if len(StringExts) > 0 && IsItemIn(path.Ext(file), StringExts) {
			ftype = TypeString
		}

		// Flag -b modifies the default if items are excluded
		if len(ByteExts) > 0 && !IsItemIn(path.Ext(file), ByteExts) {
			ftype = TypeString
		}

		embed.Entries = append(embed.Entries, Entry{
			File: file,
			Var:  fvar,
			Type: ftype,
		})
	}

	// Parse the template
	var tmplString = templates.Embed_tmpl
	tmpl, err := template.New("embed").Parse(tmplString)
	if err != nil {
		return fmt.Errorf("error parsing template: %s", err.Error())
	}

	if DryRun {
		fmt.Printf("Package %s: %d asset(s) would have been embedded\n", embed.Package, len(embed.Entries))
	} else {
		// Generate the "embed.go" file
		fh, err := os.Create(path.Join(root, "embed.go"))
		if err != nil {
			return fmt.Errorf("error creating embed.go: %s", err.Error())
		}
		defer fh.Close()

		err = tmpl.Execute(fh, embed)
		if err != nil {
			return fmt.Errorf("error creating embed.go: %s", err.Error())
		}
		if Verbose {
			fmt.Printf("Package %s: %d asset(s) embedded\n", embed.Package, len(embed.Entries))
		}
	}
	return nil
}

// Traverse recursively the specified root folder
// Gather the list of files and the list of subfolders
// If there are files to be processed it calls ProduceEmbedGo()
func TraverseDir(root string) error {
	files := make([]string, 0)
	folders := make([]string, 0)

	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	// First determine files from folders
	// Note: .go files are not considered to be embedded
	for _, entry := range entries {
		if entry.IsDir() {
			gof := path.Join(root, entry.Name(), "embed.go")
			if DryRun {
				fmt.Printf("Would have deleted file %s\n", gof)
			} else {
				os.Remove(gof)
			}
			if OkToEmbedDir(entry.Name()) {
				folders = append(folders, entry.Name())
			}
		} else {
			if OkToEmbedFile(entry.Name()) {
				files = append(files, path.Base(entry.Name()))
			}
		}
	}

	// If there are files, then produce the embed.go file
	if len(files) > 0 {
		err := ProduceEmbedGo(root, files)
		if err != nil {
			return fmt.Errorf("%s", err.Error())
		}
	}

	// Now process all the folders
	for _, folder := range folders {
		if OkToEmbedDir(folder) {
			err := TraverseDir(path.Join(root, folder))
			if err != nil {
				return fmt.Errorf("%s", err.Error())
			}
		}
	}
	return nil
}
