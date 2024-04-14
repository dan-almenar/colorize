# colorize

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/dan-almenar/colorize/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/dan-almenar/colorize)](https://goreportcard.com/report/github.com/dan-almenar/colorize)

Package colorize provides functions for formatting text with true color or Xterm approximations.

## Overview

The `colorize` package offers developers the ability to format text with various colors and styles, supporting both true color (24-bit) and Xterm (256-color) systems. It provides flexibility in text formatting for terminal-based applications.

**Features:**
- Supports true color (24-bit) and Xterm (256-color) systems
- No dependencies
- Lightweight
- Easy to use

## Installation

```bash
go get github.com/dan-almenar/colorize
```

## API Documentation

### Functions
- **GetColor(hex string, ctx ColorContext) (string, error)**:
  Retrieves the ANSI escape code for setting true color (24-bit) or Xterm (256-color) color based on the provided hexadecimal color code and context (background or foreground).

  Example:
  ```go

  red, err := c.GetColor("#FF0000", c.Foreground)
  if err != nil {
	  fmt.Println("Error:", err)
  }

  // since there was no error, we can safely use red
  fmt.Printf("%sThis text will be red%s\n", red, c.Reset)

  ```

- **ForegroundText(text string, color string) (string, error)**:
  Formats text with the specified foreground color.

  Example:
  ```go

  example, err := c.ForegroundText("Hello, world!", "#FF0000") // Red color
  if err != nil {
      fmt.Println("Error:", err)
  }
  fmt.Println(example)

  ```

- **BackgroundText(text string, color string) (string, error)**:
  Formats text with the specified background color.

  Example:
  ```go

  example, err := c.BackgroundText("Hello, world!", "#FF0000") // Red color
  if err != nil {
	  fmt.Println("Error:", err)
  }
  fmt.Println(example)

  ```

- **StyleText(text string, style string) string**:
  Formats text with the specified style.
  Valid styles are: bold, italic, underline, blink, reverse, hidden, stroke
  Unlike the ForegroundText and BackgroundText functions, the **StyleText** function does not return an error. If an invalid style is provided, it will be ignored.

  Example:
  ```go

  fmt.Println(c.StyleText("Hello, world!", []string{"bold", "underline"})) // Bold and underline text

  ```

- **FormatText(text string, options *Options) (string, error)**:
  Formats text with the specified options.

  Example:
  ```go

  example, err := c.FormatText("Hello, world!", &Options{
	Foreground: "#FF0000",
	Background: "#00FF00",
	Style:      []string{"bold", "underline"},
  })
  if err != nil {
	  fmt.Println("Error:", err)
  }
  fmt.Println(example)

  ```
	
### Types
- **Options**: 
  Represents the options for formatting text.
  Fields:
  - **Foreground**: (string) The foreground color for the text.
  - **Background**: (string) The background color for the text.
  - **Style**: ([]string) The style(s) for the text.
- **ColorContext**:
  Represents the context of the color ("background" or "foreground").
	 
## Test Information
### Tests
Unit tests have been conducted with over 96% code coverage. Detailed test results can be found in [tests_results](https://github.com/dan-almenar/colorize/blob/master/tests_results/tests_results.txt).

To run the tests locally, navigate to the package directory and execute:
```bash
go test -v -cover
```

### Benchmarks
Benchmarks have been performed to evaluate the performance of the package. Results are available in [benchmark_results](https://github.com/dan-almenar/colorize/blob/master/tests_results/benchmarks_results.txt).

To run the benchmarks locally, navigate to the package directory and execute:
```bash
go test -v -bench=.
```

## How to Contribute
Contributions to the colorize package are welcome! Feel free to submit bug reports, feature requests, or pull requests through GitHub.

To contribute, follow these steps:
1. Clone the repository:
  ```bash
   git clone https://github.com/dan-almenar/colorize.git
  ```

2. Create a new branch for your changes:
   ```bash
   git checkout -b feature/new-feature
   ```

3. Make your changes and ensure that tests pass:
   ```bash
   go test
   ```

4. Commit your changes and push to your forked repository:
   ```bash
   git commit -m "Add new feature"
   git push origin feature/new-feature
   ```

5. Create a pull request on GitHub from your forked repository to the main repository.

## Additional Information
Author:
**Dan Almenar Williams**

Version
0.1.0

## License
This project is licensed under the MIT License. See the [LICENSE](https://github.com/dan-almenar/colorize/blob/master/LICENSE.md) file for details.
