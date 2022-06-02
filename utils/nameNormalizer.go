package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type NameNormalizer struct {
	Value string
}

func NormalizeName(value string) string {
	normalizer := NameNormalizer{Value: value}
	return normalizer.Normalize()
}

func (normalizer *NameNormalizer) Normalize() string {
	return normalizer.toLowercase().
		extractUntilStop(".").
		extractUntilStop("(").
		extractUntilStop("[").
		replace(" - ", "_").
		replace(" -", "_").
		replace("- ", "_").
		replace("-", "_").
		replace(",", "_").
		replace("&", "and").
		replace("%", "perc").
		replace(".", "").
		replace(";", "").
		replace(",", "").
		replace("+", "").
		replace("!", "").
		replace(":", "").
		replace("/", "").
		replace("?", "").
		replace("--", "").
		replace("'s", "s").
		replace("’s", "s").
		replace("'", "").
		replace("[", "_").
		replace("]", "_").
		replace(" ", "_").
		replace("__", "_").
		replace("___", "_").
		replace("_–_", "_").
		replace("_—_", "_").
		replace("<", "less").
		replace(">", "greater").
		replace("₀", "0").
		replace("₁", "1").
		replace("₂", "2").
		replace("₃", "3").
		replace("₄", "4").
		replace("₅", "5").
		replace("₆", "6").
		replace("₇", "7").
		replace("₈", "8").
		replace("₉", "9").
		replace("ö", "o").
		replace("ü", "u").
		replace("è̀", "e").
		sanitizeHead().
		sanitizeTail().
		Value
}

func (normalizer *NameNormalizer) extractUntilStop(stopFlag string) *NameNormalizer {
	parts := strings.SplitN(normalizer.Value, stopFlag, -1)
	if len(parts) > 0 {
		normalizer.Value = parts[0]
	}
	return normalizer
}

func (normalizer *NameNormalizer) toLowercase() *NameNormalizer {
	normalizer.Value = strings.ToLower(normalizer.Value)
	return normalizer
}

func (normalizer *NameNormalizer) sanitizeHead() *NameNormalizer {
	if len(normalizer.Value) > 0 {
		if _, err := strconv.Atoi(normalizer.Value[0:1]); err == nil {
			normalizer.Value = fmt.Sprintf("o%s", normalizer.Value)
		}
	}
	return normalizer
}

func (normalizer *NameNormalizer) sanitizeTail() *NameNormalizer {
	if strings.HasSuffix(normalizer.Value, "_") {
		normalizer.Value = normalizer.Value[:len(normalizer.Value)-len("_")]
	}
	return normalizer
}

func (normalizer *NameNormalizer) replace(old string, character string) *NameNormalizer {
	normalizer.Value = strings.Replace(normalizer.Value, old, character, -1)
	return normalizer
}
