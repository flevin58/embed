package main

//go:generate embed -v -s .tmpl resources

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
	err := TraverseDir(args.Root)
	if err == nil && args.Verbose {
		fmt.Println("Done.")
	}
}

func IsFileNameIn(name string, arg string) bool {
	if arg == "all" || arg == "any" {
		return true
	}

	extensions := strings.Split(arg, ",")
	fileExt := path.Ext(name)
	for _, ext := range extensions {
		// Make sure that the extesion always has the '.'
		if len(ext) > 0 && ext[0] != '.' {
			ext = "." + ext
		}

		if ext == fileExt {
			return true
		}
	}
	return false
}

func OkToEmbedFile(filename string) bool {
	return !IsFileNameIn(filename, args.Xe)
}

func OkToEmbedDir(foldername string) bool {
	folders := strings.Split(args.Xd, ",")
	for _, folder := range folders {
		if folder == foldername {
			return false
		}
	}
	return true
}

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
		ftype := TypeByte
		if IsFileNameIn(file, args.StrVal) {
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
	return nil
}

func TraverseDir(root string) error {
	files := make([]string, 0)
	folders := make([]string, 0)

	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	// First determine files from folders
	//Note: .go files are not considered to be embedded
	for _, entry := range entries {
		if entry.IsDir() {
			gof := path.Join(root, entry.Name(), "embed.go")
			os.Remove(gof)
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
		if args.Verbose {
			fmt.Printf("Processing %s\n", folder)
		}
		err := TraverseDir(path.Join(root, folder))
		if err != nil {
			return fmt.Errorf("%s", err.Error())
		}
	}
	return nil
}
