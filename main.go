package main

//go:generate embed -i tmpl -s all resources

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/flevin58/embed/resources/templates"
	"github.com/flevin58/embed/src/args"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	if err := TraverseDir(args.Root); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if args.Verbose {
		fmt.Println("Done.")
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
	if len(args.ExcludedExts) > 0 {
		return !IsItemIn(path.Ext(filename), args.ExcludedExts)
	}

	// If we defined the -i flag
	if len(args.IncludedExts) > 0 {
		return IsItemIn(path.Ext(filename), args.IncludedExts)
	}

	// Neither flag defined, assume all files ok
	return true
}

func OkToEmbedDir(foldername string) bool {
	return !IsItemIn(foldername, args.ExcludedDirs)
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
		var ftype string

		// If flag is -s then the default is []byte unless specified
		if len(args.StringExts) > 0 {
			ftype = TypeByte
			if IsItemIn(path.Ext(file), args.StringExts) {
				ftype = TypeString
			}
		}

		// If flag is -b then the default is string unless specified
		if len(args.ByteExts) > 0 {
			ftype = TypeString
			if IsItemIn(path.Ext(file), args.ByteExts) {
				ftype = TypeByte
			}
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

	if args.DryRun {
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
		if args.Verbose {
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
			if args.DryRun {
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
