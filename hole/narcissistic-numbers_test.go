package hole

import (
	"strings"
	"testing"

	"github.com/code-golf/code-golf/config"
	"github.com/stretchr/testify/assert"
)

func TestNarcissisticNumbers(t *testing.T) {
	holeCfg, ok := config.AllHoleByID["narcissistic-numbers"]
	if !assert.True(t, ok, "narcissistic-numbers hole should be registered in config") {
		return
	}
	if !assert.NotNil(t, holeCfg.AnswerFunc, "narcissistic-numbers AnswerFunc should be registered") {
		return
	}

	answers := holeCfg.AnswerFunc()
	if !assert.Len(t, answers, 1) {
		return
	}

	assert.Empty(t, answers[0].Args)
	assert.Equal(t, holeCfg.Answer, answers[0].Answer)

	lines := strings.Split(answers[0].Answer, "\n")
	assert.Len(t, lines, 219)
	assert.Equal(t, "5", lines[0])
	assert.Equal(t, "153", lines[55])
	assert.Equal(t, "1000", lines[len(lines)-1])
}

func TestNarcissisticInBase(t *testing.T) {
	for _, tt := range []struct {
		n    int
		b    int
		want bool
	}{
		{0, 2, false},
		{1, 1, false},
		{1, 2, false},
		{2, 2, false},
		{3, 2, false},
		{3, 3, false},
		{4, 2, false},
		{4, 3, false},
		{4, 4, false},
		{5, 2, false},
		{5, 3, true},   // 12_3: 1^2 + 2^2 = 5
		{8, 3, true},   // 22_3: 2^2 + 2^2 = 8
		{10, 7, true},  // 13_7: 1^2 + 3^2 = 10
		{17, 3, true},  // 122_3: 1^3 + 2^3 + 2^3 = 17
		{153, 10, true},
		{154, 10, false},
		{370, 10, true},
		{371, 10, true},
		{407, 10, true},
		{1000, 33, true}, // (30, 10)_33: 30^2 + 10^2 = 1000
		{1000, 10, false},
		{10, 11, false},
	} {
		assert.Equal(t, tt.want, isNarcissisticInBase(tt.n, tt.b), "n=%d, b=%d", tt.n, tt.b)
	}
}

func TestNarcissisticAnyBaseEdgeCases(t *testing.T) {
	for _, n := range []int{-5, 0, 1, 2, 3, 4, 6, 7, 9, 11, 12, 999} {
		assert.False(t, isNarcissisticAnyBase(n), "expected %d not to be narcissistic in any base 2..n", n)
	}

	for _, n := range []int{5, 8, 10, 13, 17, 153, 370, 371, 407, 997, 1000} {
		assert.True(t, isNarcissisticAnyBase(n), "expected %d to be narcissistic in at least one base 2..n", n)
	}
}
