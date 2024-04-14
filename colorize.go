/*
Package colorize provides functions for formatting text in true color or Xterm approximations depending on the system support.

When importing this package, it's recommended to use the alias "c" for brevity:

	import c "github.com/dan-almenar/colorize"

Author: Dan Almenar Williams

Version: 0.1.0

License: MIT (https://github.com/dan-almenar/colorize/blob/master/LICENSE)
*/
package colorize

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/* Package specific error type and functions */

/*
colorizeErr represents a non-fatal error specific to the colorize package.

This type is used to encapsulate errors that occur within the colorize package.
It provides a name for categorizing the error and a message describing the error.

Note that whenever an error occurs, the original text string is returned unmodified. This design choice ensures that the formatted text is always displayed, even if there's an issue with the provided options or system support.

During development, it's recommended to handle these errors appropriately to ensure the integrity of the formatted text.
In production environments, omitting error handling or simply logging them out in favor of displaying the unformatted text may be acceptable, depending on the application's requirements.

Fields:

	name string: A name categorizing the error.
	msg  string: A message describing the error.
*/
type colorizeErr struct {
	name string
	msg  string
}

/*
newColorizeErr creates a new instance of colorizeErr with the provided name and message.

Parameters:

	name string: A name categorizing the error.
	msg  string: A message describing the error.

Returns:

	*colorizeErr: A pointer to the newly created colorizeErr instance.
*/
func newColorizeErr(name string, msg string) *colorizeErr {
	return &colorizeErr{name: name, msg: msg}
}

/*
Error returns the string representation of the colorizeErr.

This method formats the error with the following pattern: "<name>: <message>".

Returns:

	string: The string representation of the error.
*/
func (e *colorizeErr) Error() string {
	return fmt.Sprintf("%s: %s", e.name, e.msg)
}

/* The ColorContext type represents the context of the color (background or foreground) */
type ColorContext string

const (
	/* Constants for background and foreground contexts */
	background ColorContext = "background"
	foreground ColorContext = "foreground"
)

/* The Options type represents the options for formatting text */
type Options struct {
	BgColor string   // background color
	FgColor string   // foreground color
	Styles  []string // text style(s): bold, italic, underline, blink, reverse, hidden and stroke
}

/* The color type represents an RGB color */
type color struct {
	r uint8
	g uint8
	b uint8
}

const (
	// escape codes
	fgTrueColor = "\033[38;2;"
	bgTrueColor = "\033[48;2;"
	fgXterm     = "\033[38;5;"
	bgXterm     = "\033[48;5;"
	reset       = "\033[0m"
	Reset       = reset // reset internally refers to the escape code for resetting any formatting

	/* xTerm specific constants */
	scalingFactor = 255 / 5 // 6-bit color scaling factor
	// Xterm color codes
	xTermBlack = 0
	xTermWhite = 15
	grayOffset = 232
	// Xterm color conversion factors
	colorOffset  = 16
	colorFactor1 = 36
	colorFactor2 = 6
)

var (
	/* System color support */
	trueColor = os.Getenv("COLORTERM") == "truecolor"
	xTerm     = os.Getenv("TERM") == "xterm"

	styles = map[string]string{
		"bold":      "\033[1m",
		"italic":    "\033[3m",
		"underline": "\033[4m",
		"blink":     "\033[5m",
		"reverse":   "\033[7m",
		"hidden":    "\033[8m",
		"stroke":    "\033[9m",
	}

	// regex for hex color code
	regex = regexp.MustCompile(`^#?([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})$`)

	// color pointer
	colorPtr *color
)

/*
validateHex validates the provided hexadecimal color code.

If the hex string is invalid, an error is returned.

Parameters:
  - hex: The hexadecimal color code, either with or without the # prefix (e.g., "#RRGGBB").
*/
func validateHex(hex string) error {
	if !regex.MatchString(hex) {
		err := newColorizeErr("HEXERR", fmt.Sprintf("invalid hex code: %s", hex))
		return fmt.Errorf(err.Error())
	}
	return nil
}

/*
getColor converts a hexadecimal color code to RGB representation.

Parameters:
  - hex: The hexadecimal color code (e.g., "#RRGGBB").

Return:
  - *color: A pointer to the color struct representing the RGB color.
  - error: An error if the provided hex code is invalid.
*/
func getColor(hex string) (*color, error) {
	err := validateHex(hex)
	if err != nil {
		return nil, err
	}

	// errors are omitted due to regex
	match := regex.FindStringSubmatch(hex)
	r, _ := strconv.ParseUint(match[1], 16, 8)
	g, _ := strconv.ParseUint(match[2], 16, 8)
	b, _ := strconv.ParseUint(match[3], 16, 8)

	colorPtr = &color{uint8(r), uint8(g), uint8(b)}

	return colorPtr, nil
}

/*
The GetColor function is a convenience wrapper around internal package functions, that returns the
ANSI escape code for setting true color (24-bit) or Xterm (256-color) color (depending on the system support)
for the provided ColorContext (background or foreground).

It's purpose is to offer an easy-to-use API that helps avoid repetitive code when the same color is ment
to be used in multiple places.

Parameters:
  - hex: The hexadecimal color code (e.g., "#RRGGBB").
  - ctx: The color context (background or foreground).

Return:
  - string: The ANSI escape code for setting true color.
  - error: An error if the provided hex code is invalid or the system does not support true color or xterm.

Example:

	// Save the code for red foreground in a variable
	red, err := c.GetColor("#FF0000", c.foreground)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Use the red foreground code
	warningMessage := red + "Warning: This text is red" + c.Reset
	redHeart := red + "\u2665" + Reset

Note: Append the Reset constant to the end of the code to reset the color.
*/
func GetColor(hex string, ctx ColorContext) (string, error) {
	var code string = ""

	// get color
	colorPtr, err := getColor(hex)
	if err != nil {
		return code, err
	}

	// set code based on system support
	if trueColor {
		code = getTCCode(colorPtr, ctx)
	} else if xTerm {
		code = getXTCode(colorPtr, ctx)
	} else {
		err = newColorizeErr("SYSNOCOLOR", "System does not support true color or xterm")
	}

	return code, err
}

/*
getTCCode returns the ANSI escape code for setting true color (24-bit) in the terminal.

Parameters:
  - col: A pointer to the color struct representing the RGB color.
  - ctx: The color context (background or foreground).

Return:
  - string: The ANSI escape code for setting true color.
*/
func getTCCode(col *color, ctx ColorContext) string {
	if ctx == background {
		return fmt.Sprintf("%s%d;%d;%dm", bgTrueColor, col.r, col.g, col.b)
	} else {
		return fmt.Sprintf("%s%d;%d;%dm", fgTrueColor, col.r, col.g, col.b)
	}
}

/*
getXTCode returns the ANSI escape code for setting Xterm (256-color) color in the terminal.

Parameters:
  - col: A pointer to the color struct representing the RGB color.
  - ctx: The color context (background or foreground).

Return:
  - string: The ANSI escape code for setting Xterm color.
*/
func getXTCode(col *color, ctx ColorContext) string {
	if ctx == background {
		return fmt.Sprintf("%s%dm", bgXterm, rgbToXterm(col))
	} else {
		return fmt.Sprintf("%s%dm", fgXterm, rgbToXterm(col))
	}
}

/*
rgbToXterm converts an RGB color to the closest Xterm (256-color) approximation.

Parameters:
  - col: A pointer to the color struct representing the RGB color.

Return:
  - uint8: The Xterm color code.
*/
func rgbToXterm(col *color) uint8 {
	xtCode := uint8(0)

	// Convert RGB values to basee-6
	// This process involves several type conversions in order to guarantee that the result is
	// the closest approximation in the Xterm table.
	// These type conversions may affect performance.
	rInt := uint8(math.Round((float64(col.r) / scalingFactor) + 0.4))
	gInt := uint8(math.Round((float64(col.g) / scalingFactor) + 0.4))
	bInt := uint8(math.Round((float64(col.b) / scalingFactor) + 0.4))

	// Calculate Xterm color code
	if rInt == gInt && gInt == bInt {
		// Grayscale
		if rInt == 0 {
			// Black
			xtCode = xTermBlack
		} else if rInt == 5 {
			// White
			xtCode = xTermWhite
		} else {
			// Shade of gray
			xtCode = grayOffset + (rInt-1)*5
		}
	} else {
		// Color
		xtCode = colorOffset + colorFactor1*rInt + colorFactor2*gInt + bInt
	}

	return xtCode
}

/*
FormatText formats the given text with the specified options.

Parameters:
  - text: The text to be formatted.
  - options: The formatting options including background color, foreground color, and styles.

Return:
  - string: The formatted text.
  - error: An error if the provided options are invalid or the system does not support true color or Xterm.

Example:

	// Format text with red foreground color and bold underline styles
	formattedText, err := c.FormatText("Hello, world!", &c.Options{FgColor: "#FF0000", Styles: []string{"bold", "underline"}})
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(formattedText)
	}

Note: Valid styles include: bold, italic, underline, blink, reverse, hidden and stroke.
*/
func FormatText(text string, options *Options) (string, error) {
	builder := strings.Builder{}

	// no options provided
	if options == nil || (options.BgColor == "" && options.FgColor == "" && len(options.Styles) == 0) {
		err := fmt.Errorf("No options provided")
		return text, err
	}

	// no system support
	if !trueColor && !xTerm {
		err := newColorizeErr("SYSNOCOLOR", "System does not support true color or xterm")
		return text, fmt.Errorf(err.Error())
	}

	// options provided
	if len(options.Styles) > 0 {
		for _, s := range options.Styles {
			builder.WriteString(styles[s])
		}
	}
	if trueColor {
		if options.BgColor != "" {
			bgColor, err := getColor(options.BgColor)
			if err != nil {
				// HEXERR
				return text, err
			}
			builder.WriteString(getTCCode(bgColor, background))
		}
		if options.FgColor != "" {
			fgColor, err := getColor(options.FgColor)
			if err != nil {
				return text, err
			}
			builder.WriteString(getTCCode(fgColor, foreground))
		}
	} else {
		// xTerm
		if options.BgColor != "" {
			bgColor, err := getColor(options.BgColor)
			if err != nil {
				return text, err
			}
			builder.WriteString(getXTCode(bgColor, background))
		}
		if options.FgColor != "" {
			fgColor, err := getColor(options.FgColor)
			if err != nil {
				return text, err
			}
			builder.WriteString(getXTCode(fgColor, foreground))
		}
	}

	builder.WriteString(text)

	if len(builder.String()) == len(text) {
		return builder.String(), nil
	}
	builder.WriteString(reset)

	return builder.String(), nil
}

/*
ForegroundText formats the given text with the specified foreground color.

Parameters:
  - text: The text to be formatted.
  - color: The foreground color (in hexadecimal format).

Return:
  - string: The formatted text.
  - error: An error if the provided color is invalid or the system does not support true color or Xterm.

Example:

	// Format text with red foreground color
	formattedText, err := c.ForegroundText("Hello, world!", "#FF0000")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(formattedText)
	}
*/
func ForegroundText(text string, color string) (string, error) {
	return FormatText(text, &Options{FgColor: color})
}

/*
BackgroundText formats the given text with the specified background color.

Parameters:
  - text: The text to be formatted.
  - color: The background color (in hexadecimal format).

Return:
  - string: The formatted text.
  - error: An error if the provided color is invalid or the system does not support true color or Xterm.

Example:

	// Format text with blue background color
	formattedText, err := c.BackgroundText("Hello, world!", "#0000FF")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(formattedText)
	}
*/
func BackgroundText(text string, color string) (string, error) {
	return FormatText(text, &Options{BgColor: color})
}

/*
StyleText formats the given text with the specified styles.

Unlike ForegroundText and BackgroundText, StyleText, if the system does not support any of the
provided styles, no error is returned, since no escape sequences are generated for the invalid styles.

Parameters:
  - text: The string to be formatted.
  - styles: A string slice containing the text styles (e.g., bold, italic, underline).

Return:
  - string: The formatted text.

Example:

	// Format text with bold style
	formattedText := c.StyleText("Hello, world!", []string{"bold"}) // assuming the package alias "c" is used
	fmt.Println(formattedText)
*/
func StyleText(text string, styles []string) string {
	t, _ := FormatText(text, &Options{Styles: styles})
	return t
}
