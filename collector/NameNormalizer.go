package main

import (
	"strings"
)

type NameNormalizer struct {
	value string
}

func (normalizer *NameNormalizer) Normalize() string {
	return normalizer.toLowercase().
		replace(" - ", "").
		replace("- ", "").
		replace(" -", "").
		replace("-", "_").
		replace(" (", "_").
		replace(") ", "_").
		replace(")", "_").
		replace("&", "and").
		replace("%", "perc").
		replace(".", "").
		replace(" ", "_").
		cleanTail().
		value
}

func (normalizer *NameNormalizer) toLowercase() *NameNormalizer {
	normalizer.value = strings.ToLower(normalizer.value)
	return normalizer
}

func (normalizer *NameNormalizer) cleanTail() *NameNormalizer {
	if strings.HasSuffix(normalizer.value, "_") {
		normalizer.value = normalizer.value[:len(normalizer.value)-len("_")]
	}
	return normalizer
}

func (normalizer *NameNormalizer) replace(old string, character string) *NameNormalizer {
	normalizer.value = strings.Replace(normalizer.value, old, character, -1)
	return normalizer
}
