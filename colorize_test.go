package colorize

import (
	"testing"
)

var (
	badHex = []string{
		"FF00000",
		"FF00",
		"FF000H",
		"#FF00000",
		"#FF00",
		"#FF000H",
	}
	validHex = []string{
		"#FFFFFF",
		"FFFFFF",
		"#000000",
		"000000",
		"#abcdef",
		"ABCDEF",
		"#12abAB",
		"12abAB",
	}
	validOpts = []*Options{
		{FgColor: "#FF0000"},
		{BgColor: "#0000FF"},
		{FgColor: "#FF0000", BgColor: "#0000FF"},
		{Styles: []string{"bold"}},
		{FgColor: "#FF0000", Styles: []string{"bold"}},
		{BgColor: "#0000FF", Styles: []string{"bold"}},
		{FgColor: "#FF0000", BgColor: "#0000FF", Styles: []string{"bold"}},
	}
	invalidOpts = []*Options{
		{FgColor: "#FF00000"},
		{BgColor: "#0000FF0"},
		{FgColor: "#FF00000", BgColor: "#0000FF0"},
		{FgColor: "#FF00000", Styles: []string{"bold-italic"}},
		{BgColor: "#0000FF0", Styles: []string{"bold-italic"}},
		{FgColor: "#FF00000", BgColor: "#0000FF0", Styles: []string{"bold-italic"}},
	}
	prevTrueColor = trueColor
	prevXTerm     = xTerm
)

// defer func
func restore() {
	trueColor = prevTrueColor
	xTerm = prevXTerm
}

/* TestValidateHex tests the validateHex function */
func TestValidateHex(t *testing.T) {
	// valid hex
	for _, hex := range validHex {
		err := validateHex(hex)
		if err != nil {
			t.Error("Expected no error but got", err)
		}
	}

	// invalid hex
	for _, hex := range badHex {
		err := validateHex(hex)
		if err == nil {
			t.Error("Expected an error but got nil")
		}
	}
}

/* TestNewColorizeErr tests the newColorizeErr function */
func TestNewColorizeErr(t *testing.T) {
	name := "test"
	msg := "test message"

	err := newColorizeErr(name, msg)

	if err == nil {
		t.Error("Expected an error but got nil")
	} else if err.Error() != "test: test message" {
		t.Errorf("Expected error message to be 'test: test message' but got '%s'", err.Error())
	}
}

/* TestGetColor tests the GetColor function */
func TestGetColor(t *testing.T) {
	// defer restore
	defer restore()

	// invalid hex
	for _, hex := range badHex {
		_, err := GetColor(hex, foreground)
		if err == nil {
			t.Error("Expected an error but got nil")
		}
	}

	// valid hex, true color support
	trueColor = true
	for _, hex := range validHex {
		_, err := GetColor(hex, foreground)
		if err != nil {
			t.Error("Expected no error but got", err)
		}
	}

	// valid hex, xterm support
	trueColor = false
	xTerm = true
	for _, hex := range validHex {
		_, err := GetColor(hex, foreground)
		if err != nil {
			t.Error("Expected no error but got", err)
		}
	}

	// valid hex, no color support
	xTerm = false
	for _, hex := range validHex {
		_, err := GetColor(hex, foreground)
		if err == nil {
			t.Error("Expected an error but got nil")
		}
	}
}

/* TestInternalGetColor tests the getColor function */
func TestInternalGetColor(t *testing.T) {
	// no options provided
	_, err := getColor("")
	if err == nil {
		t.Error("Expected an error but got nil")
	}

	// invalid hex code
	for _, hex := range badHex {
		_, err = getColor(hex)
		if err == nil {
			t.Error("Expected an error but got nil")
		}
	}

	// valid hex code
	for _, hex := range validHex {
		_, err = getColor(hex)
		if err != nil {
			t.Error("Expected no error but got", err)
		}
	}
}

/* TestFormatText tests the FormatText function */
func TestFormatText(t *testing.T) {
	// defer restore
	defer restore()

	// no options provided
	_, err := FormatText("", nil)
	if err == nil {
		t.Error("Expected an error but got nil")
	}

	// valid options
	for _, opt := range validOpts {
		_, err = FormatText("", opt)
		if err != nil {
			t.Error("Expected no error but got", err)
		}
	}

	// invalid options
	for _, opt := range invalidOpts {
		_, err = FormatText("", opt)
		if err == nil {
			t.Error("Expected an error but got nil")
		}
	}

	// test for non-supported true color
	xTerm = true
	for _, opt := range validOpts {
		_, err = FormatText("", opt)
		if err != nil {
			t.Error("Expected no error but got", err)
		}
	}

	// test for non-supported true color and xterm
	trueColor = false
	xTerm = false
	for _, opt := range validOpts {
		_, err = FormatText("", opt)
		if err == nil {
			t.Error("Expected an error but got nil")
		}
	}
}

/* TestStyleText tests the StyleText function */
func TestStyleText(t *testing.T) {
	testString := "test"
	validStyles := []string{
		"bold",
		"italic",
		"underline",
	}
	invalidStyles := []string{
		"invalid",
		"bold-italic",
	}

	// valid styles
	for _, style := range validStyles {
		ret := StyleText(testString, []string{style})
		if len(ret) <= len(testString) {
			t.Error("No style was applied")
		}
	}
	ret := StyleText(testString, validStyles)
	if len(ret) <= len(testString) {
		t.Error("No style was applied")
	}

	// invalid styles
	for _, style := range invalidStyles {
		ret := StyleText(testString, []string{style})
		if len(ret) != len(testString) {
			t.Error("Invalid style was applied")
		}
	}
	ret = StyleText(testString, invalidStyles)
	if len(ret) != len(testString) {
		t.Error("Invalid styles were applied")
	}

}

/* TestForegroundText tests the ForegroundText function */
func TestForegroundText(t *testing.T) {
	// defer restore
	defer restore()

	validColors := []string{
		"#FF0000",
		"#00FF00",
		"#0000FF",
	}
	invalidColors := []string{
		"#FF00000",
		"#FF00",
		"#FF000H",
	}
	// invalid colors
	for _, color := range invalidColors {
		_, err := ForegroundText("", color)
		if err == nil {
			t.Error("Expected an error but got nil")
		}
	}

	// valid colors
	for _, color := range validColors {
		_, err := ForegroundText("", color)
		if err != nil {
			t.Error("Expected no error but got", err)
		}
	}

	// valid colors with no true colors support
	trueColor = false
	xTerm = true
	for _, color := range validColors {
		_, err := ForegroundText("", color)
		if err != nil {
			t.Error("Expected no error but got", err)
		}
	}

	// valid colors with no xterm support
	xTerm = false
	for _, color := range validColors {
		_, err := ForegroundText("", color)
		if err == nil {
			t.Error("Expected an error but got nil")
		}
	}
}

/* TestBackgroundText tests the ForegroundText function */
func TestBackgroundText(t *testing.T) {
	// defer restore
	defer restore()

	validColors := []string{
		"#FF0000",
		"#00FF00",
		"#0000FF",
	}
	invalidColors := []string{
		"#FF00000",
		"#FF00",
		"#FF000H",
	}

	// invalid colors
	for _, color := range invalidColors {
		_, err := BackgroundText("", color)
		if err == nil {
			t.Error("Expected an error but got nil")
		}
	}

	// valid colors
	for _, color := range validColors {
		_, err := BackgroundText("", color)
		if err != nil {
			t.Error("Expected no error but got", err)
		}
	}

	// valid colors with no true colors support
	trueColor = false
	xTerm = true
	for _, color := range validColors {
		_, err := BackgroundText("", color)
		if err != nil {
			t.Error("Expected no error but got", err)
		}
	}

	// valid colors with no xterm support
	xTerm = false
	for _, color := range validColors {
		_, err := BackgroundText("", color)
		if err == nil {
			t.Error("Expected an error but got nil")
		}
	}
}
