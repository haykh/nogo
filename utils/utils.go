package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type MessageType int64

const (
	Normal MessageType = iota
	Hint
	Warning
	Error
)

type ColorType string

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

func isIn(v interface{}, list ...interface{}) bool {
	for _, l := range list {
		if l == v {
			return true
		}
	}
	return false
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

func Message(msg string, msgtype MessageType, newline bool) {
	if msgtype == Warning {
		fmt.Println("  WARNING")
	} else if msgtype == Hint {
		fmt.Println("   HINT")
	}
	fmt.Println(strings.Replace(" : "+msg, "\n", "\n : ", -1))
	if newline {
		fmt.Printf("\n")
	}
}

func PromptBool(msg string, def bool, msgtype MessageType) bool {
	r := bufio.NewReader(os.Stdin)
	for {
		Message(msg, msgtype, false)
		if def {
			fmt.Fprint(os.Stderr, " :> [y]/n ")
		} else {
			fmt.Fprint(os.Stderr, " :> y/[n] ")
		}
		if s, err := r.ReadString('\n'); err == nil {
			s = strings.TrimSpace(strings.ToLower(s))
			if s == "y" || s == "n" || s == "" {
				if s == "" {
					return def
				} else {
					return s == "y"
				}
			}
		} else {
			log.Fatal(err)
		}
	}
}

func PromptString(msg string, def string, msgtype MessageType) string {
	r := bufio.NewReader(os.Stdin)
	for {
		Message(msg, msgtype, false)
		if def != "" {
			fmt.Fprint(os.Stderr, fmt.Sprintf(" :> [%s] ", def))
		} else {
			fmt.Fprint(os.Stderr, " :> ")
		}
		if s, err := r.ReadString('\n'); err == nil {
			return strings.TrimSpace(s)
		} else {
			log.Fatal(err)
		}
	}
}

func DashedLine(n int) string {
	line := []rune(strings.Repeat("-", n))
	for k := 0; k < len(line); k += 2 {
		line[k] = ' '
	}
	return string(line)
}
