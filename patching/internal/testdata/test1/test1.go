package test1

import _ "embed"

//go:embed a.txt
var A string

//go:embed b.txt
var B string
