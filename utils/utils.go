package utils

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
)

type MessageType int64

const (
	Normal MessageType = iota
	Hint
	Warning
	Error
)

type TextStyle = string
type ColorType = string
type HighlightType = string

const (
	ColorReset  ColorType = "\033[0m"
	ColorRed    ColorType = "\033[31m"
	ColorGreen  ColorType = "\033[32m"
	ColorBlue   ColorType = "\033[34m"
	ColorCyan   ColorType = "\033[36m"
	ColorYellow ColorType = "\033[33m"
	ColorPurple ColorType = "\033[35m"
	ColorGray   ColorType = "\033[30m"
)

const (
	HiStrike HighlightType = "\033[9m"
	HiReset  HighlightType = "\033[0m"
)

func IsIn[T cmp.Ordered](element T, array []T) bool {
	for _, e := range array {
		if e == element {
			return true
		}
	}
	return false
}

func Clean(str string) string {
	colors := []ColorType{ColorReset, ColorRed, ColorGreen, ColorBlue, ColorCyan, ColorYellow, ColorPurple, ColorGray}
	highlights := []HighlightType{HiStrike, HiReset}
	for _, c := range colors {
		str = strings.ReplaceAll(str, c, "")
	}
	for _, h := range highlights {
		str = strings.ReplaceAll(str, h, "")
	}
	return strings.Trim(str, " \n")
}

func CreateFile(fname string) error {
	f, err := func(p string) (*os.File, error) {
		if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
			return nil, err
		}
		return os.Create(p)
	}(fname)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer f.Close()
	if _, exists := os.Stat(fname); os.IsNotExist(exists) {
		log.Fatal("Failed to create file")
		return exists
	}
	return nil
}

func Message(msg string, msgtype MessageType, newline bool, color ...ColorType) {
	if len(color) > 0 {
		fmt.Printf("%s", color[0])
	} else if msgtype == Warning {
		fmt.Printf("%s", ColorRed)
		fmt.Println("  WARNING")
	} else if msgtype == Hint {
		fmt.Printf("%s", ColorYellow)
		fmt.Println("   HINT")
	}
	fmt.Println(msg)
	if newline {
		fmt.Printf("\n")
	}
	if len(color) > 0 {
		fmt.Printf("%s", ColorReset)
	}
}

func Prompt(msg string, canskip bool, def interface{}) (interface{}, error) {
	switch def.(type) {
	case string:
		return PromptString(msg, canskip, def.(string))
	case bool:
		return PromptBool(msg, canskip, def.(bool))
	default:
		return nil, fmt.Errorf("unsupported type")
	}
}

func PromptBool(msg string, canskip, def bool) (bool, error) {
	val := false
	var prompt survey.Prompt
	msg += " pick a new value"
	if canskip {
		msg += " or leave blank for default/existing"
		prompt = &survey.Confirm{
			Message: msg,
			Default: def,
		}
	} else {
		prompt = &survey.Confirm{
			Message: msg,
		}
	}
	if err := survey.AskOne(prompt, &val); err != nil {
		return false, err
	} else {
		return val, nil
	}
}

func PromptString(msg string, canskip bool, def string) (string, error) {
	val := ""
	var prompt survey.Prompt
	msg += " enter a new value"
	if canskip {
		msg += " or leave blank for default/existing"
		prompt = &survey.Input{
			Message: msg,
			Default: def,
		}
	} else {
		prompt = &survey.Input{
			Message: msg,
		}
	}
	if err := survey.AskOne(prompt, &val); err != nil {
		return "", err
	} else {
		return val, nil
	}
}

func DashedLine(n int) string {
	line := []rune(strings.Repeat("-", n))
	for k := 0; k < len(line); k += 2 {
		line[k] = ' '
	}
	return string(line)
}
