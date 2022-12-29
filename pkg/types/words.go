package types

var lowNames = []string{"Zero", "One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen"}
var tensNames = []string{"Twenty", "Thirty", "Forty", "Fifty", "Sixty", "Seventy", "Eighty", "Ninety"}
var bigNames = []string{"Thousand", "Million", "Billion"}

func convert999(num int) string {
	s1 := convert99(num/100) + " Hundred"
	s2 := convert99(num % 100)
	if num%100 == 0 {
		return s1
	} else {
		return s1 + " " + s2
	}
}

func convert99(num int) string {
	if num < 20 {
		return lowNames[num]
	}
	s := tensNames[num/10-2]
	if num%10 == 0 {
		return s
	}
	return s + "-" + lowNames[num%10]
}

// Adapted from https://socketloop.com/tutorials/golang-how-to-convert-a-number-to-words
func Num2Words(num int) string {
	switch {
	case num < 0:
		return "Negative " + Num2Words(-num)
	case num < 100:
		return convert99(num)
	case num < 1_000:
		return convert999(num)
	case num < 1_000_000:
		s1 := convert999(num / 1_000)
		s2 := convert999(num % 1_000)
		return s1 + bigNames[1] + s2
	case num < 1_000_000_000:
		s1 := Num2Words(num / 1_000_000)
		s2 := Num2Words(num % 1_000_000)
		return s1 + bigNames[2] + s2
	default:
		return "More than a trillion"
	}
}
