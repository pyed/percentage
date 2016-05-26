package percentage

import "testing"

var testCases = []struct{ input, expectedF, expectedValue string }{
	// +
	{"55 + 55", "55 + 55%", "85.25"},
	{"703.4 + 273.555%", "703.4 + 273.555%", "2,627.59"},
	{"-703.4 + 273.555%", "-703.4 + 273.555%", "-2,627.59"},
	{"-703.4 + -273.555%", "-703.4 + -273.555%", "1,220.79"},
	{"703.4 + -273.555%", "703.4 + -273.555%", "-1,220.79"},
	// -
	{"100 -50%", "100 - 50%", "50"},
	{"77 - 77", "77 - 77%", "17.71"},
	{"-100 --70.3", "-100 - -70.3%", "-170.3"},
	{"102340 - 200", "102,340 - 200%", "-102,340"},
	{"803 - -800%", "803 - -800%", "7,227"},
	// ×
	{"50 * 10", "50 × 10%", "250"},
	{"-30 X-10%", "-30 × -10%", "-90"},
	{"850 x 25.5", "850 × 25.5%", "184,237.5"},
	{"0100000 * 0000300000", "100,000 × 300,000%", "30,000,000,000,000"},
	{"-111*10%", "-111 × 10%", "1,232.1"},
	// ÷
	{"10 / 10", "10 ÷ 10%", "10"},
	{"-80/-10%", "-80 ÷ -10%", "-10"},
	{"10023.23 / 11.2", "10,023.23 ÷ 11.2%", "8.93"},
	{"7/-7", "7 ÷ -7%", "-14.29"},
	{"5/2.5", "5 ÷ 2.5%", "40"},
	// of
	{"50%of100", "50% of 100", "50"},
	{"80of-20", "80% of -20", "-16"},
	{"200% OF 10", "200% of 10", "20"},
	{"500000 oF 30", "500,000% of 30", "150,000"},
	{"-40Of000038", "-40% of 38", "-15.2"},
	// in
	{"15 in 30", "15 in 30", "50%"},
	{"-850IN-730", "-850 in -730", "116.44%"},
	{"88888iN-99999", "88,888 in -99,999", "-88.89%"},
	{"-007In 10", "-7 in 10", "-70%"},
	{"708 in1000", "708 in 1,000", "70.8%"},
}

func TestPercentages(t *testing.T) {
	for _, tcase := range testCases {
		out, err := NewExpr(tcase.input)
		if err != nil {
			t.Fatal(err)
		}
		if out.PrintExpr() != tcase.expectedF {
			t.Fatalf("PrintExpr: expected: %s got: %s", tcase.expectedF, out.PrintExpr())
		}
		if out.PrintValue() != tcase.expectedValue {
			t.Fatalf("PrintValue: expected: %s got: %s", tcase.expectedValue, out.PrintValue())
		}

	}
}

func benchmarkPercentages(exp string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewExpr(exp)
	}
}

func BenchmarkP0(b *testing.B) { benchmarkPercentages(testCases[0].input, b) }
func BenchmarkP1(b *testing.B) { benchmarkPercentages(testCases[5].input, b) }
func BenchmarkP2(b *testing.B) { benchmarkPercentages(testCases[10].input, b) }
func BenchmarkP3(b *testing.B) { benchmarkPercentages(testCases[15].input, b) }
func BenchmarkP4(b *testing.B) { benchmarkPercentages(testCases[20].input, b) }
