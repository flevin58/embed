package cmd

import "path"

// Returns a slice with all members having a dot as prefix
func FixDot(stringSlice []string) []string {
	outSlice := make([]string, len(stringSlice))
	for i, elem := range stringSlice {
		if elem == "all" || elem == "none" {
			outSlice[i] = elem
			continue
		}
		if len(elem) > 0 && elem[0] != '.' {
			outSlice[i] = "." + elem
		} else {
			outSlice[i] = elem
		}
	}
	return outSlice
}

// Checks if the item is inside the slice of strings
// and returns true if found, false otherwise.
// Special cases:
// first item = "none" returns false
// first item = "all" returns true
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

// Checks if it's ok to parse the given folder
func OkToEmbedDir(foldername string) bool {
	return !IsItemIn(foldername, ExcludedDirs)
}

// Checks if it's ok to embed the given file
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
