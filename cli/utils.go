package main

func maybeAddSuffix(name string) string {
	if name[len(name)-3:] == ".md" {
		return name
	}
	return name + ".md"
}
