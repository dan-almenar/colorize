package colorize

import (
	"testing"
)

/* BenchmarkGetColor benchmarks the GetColor function */
func BenchmarkGetColor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, hex := range validHex {
			_, _ = GetColor(hex, foreground)
		}
	}
}

/* BenchmarkValidateHex benchmarks the validateHex function */
func BenchmarkValidateHex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, hex := range validHex {
			_ = validateHex(hex)
		}
	}
}

/* BenchmarkInternalGetColor benchmarks the getColor function */
func BenchmarkInternalGetColor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, hex := range validHex {
			_, _ = getColor(hex)
		}
	}
}

/* BenchmarkFormatText benchmarks the FormatText function */
func BenchmarkFormatText(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, opt := range validOpts {
			_, _ = FormatText("", opt)
		}
	}
}

/* BenchmarkStyleText benchmarks the StyleText function */
func BenchmarkStyleText(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, style := range []string{"bold", "italic", "underline"} {
			_ = StyleText("", []string{style})
		}
	}
}

/* BenchmarkForegroundText benchmarks the ForegroundText function */
func BenchmarkForegroundText(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, color := range []string{"#FF0000", "#00FF00", "#0000FF"} {
			_, _ = ForegroundText("", color)
		}
	}
}

/* BenchmarkBackgroundText benchmarks the BackgroundText function */
func BenchmarkBackgroundText(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, color := range []string{"#FF0000", "#00FF00", "#0000FF"} {
			_, _ = BackgroundText("", color)
		}
	}
}
