package tools

func FixDot(stringSlice []string) []string {
	outSlice := make([]string, len(stringSlice))
	for i, elem := range stringSlice {
		if len(elem) > 0 && elem[0] != '.' {
			outSlice[i] = "." + elem
		} else {
			outSlice[i] = elem
		}
	}
	return outSlice
}
