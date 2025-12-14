package testdata

func StringLiterals() {
	a := "// not a comment"
	b := "/* also not a comment */"
	c := `
	// raw string not a comment
	/* neither is this */
	`
	_ = a
	_ = b
	_ = c
}
