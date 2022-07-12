package utils

import (
  "bufio"
  "fmt"
  "log"
  "os"
  "strings"
  "path/filepath"
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

func PromptBool(label string, def bool) bool {
  r := bufio.NewReader(os.Stdin)
  if def {
    label += " [y]/n "
  } else {
    label += " y/[n] "
  }
  for {
    fmt.Fprint(os.Stderr, label + " ")
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
  return def
}

func Message(msg string) {
  fmt.Println("[ nogo ]:", msg)
}

func Warn(msg string) {
  fmt.Println("[ nogo WARNING ]:", msg)
}

func DashedLine(n int) string {
  line := []rune(strings.Repeat("-", n))
  for k := 0; k < len(line); k += 2 {
    line[k] = ' '
  }
  return string(line)
}

func Hint(msg string) {
  s := strings.Split(msg, "\n")
  // find max length of line
  max := 0
  for _, l := range s {
    if len(l) > max {
      max = len(l)
    }
  }
  max += 4
  fmt.Println()
  for si, i := 0, 0; i < len(s) + 4; i++ {
    var str string
    if i == 0 {
      left := DashedLine(max / 2 - 7)
      right := DashedLine(max - max / 2 - 7)
      str = fmt.Sprintf("+%s [ nogo HINT ]%s+", left, right)
    } else if i == len(s) + 3 {
      mid := DashedLine(max)
      str = fmt.Sprintf("+%s+", mid)
      //fmt.Println(str)
    } else if i == 1 || i == len(s) + 2 {
      str = fmt.Sprintf("|%s|", strings.Repeat(" ", max))
    } else {
      str = fmt.Sprintf("|  %s%s|", s[si], strings.Repeat(" ", max - len(s[si]) - 2))
      si++
    }
    fmt.Println(str)
  }
  fmt.Println()
}

//fmt.Println("[ nogo HINT ]:", msg)
//}

//func Prompt(label string, continueIf func(interface{}) bool, values ...interface{}) string {
//var s string
//if continueIf == nil {
//if len(values) == 0 {
//log.Fatal("Incorrect use of Prompt: no values and no continueIf function")
//} else {
//var ss []string
//for _, v := range values {
//ss = append(ss, fmt.Sprintf("%v", v))
//}
//label += " [" + strings.Join(ss, "/") + "]"
//continueIf = func(v interface{}) bool {
//return isIn(v, values...)
//}
//}
//} else {
//if len(values) > 0 {
//log.Fatal("Incorrect use of Prompt: both values and continueIf function")
//}
//}
//r := bufio.NewReader(os.Stdin)
//for {
//fmt.Fprint(os.Stderr, label + " ")
//if s, err := r.ReadString('\n'); err == nil {
//s = strings.TrimSpace(strings.ToLower(s))
//if continueIf(s) {
//return s
//}
//} else {
//log.Fatal(err)
//}
//}
//return s
//}
