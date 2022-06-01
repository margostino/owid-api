package utils

import (
	"fmt"
	"testing"
)

func normalizeName(name string) string {
	normalizer := NameNormalizer{Value: name}
	return normalizer.Normalize()
}

func assert(t *testing.T, expected string, result string) {
	if result != expected {
		t.Fatalf(fmt.Sprintf("expected %s but it was %s", expected, result))
	}
}

func TestNameNormalizer(t *testing.T) {
	var original, expected, result string

	original = "Top Net Personal Wealth Shares – WID (2018)"
	expected = "top_net_personal_wealth_shares_wid"
	result = normalizeName(original)
	assert(t, expected, result)

	original = "Trade – Giovanni and Tena-Junguito (2016)"
	expected = "trade_giovanni_and_tena_junguito"
	result = normalizeName(original)
	assert(t, expected, result)

	//TODO
}
