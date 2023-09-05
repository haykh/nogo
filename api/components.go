package notionApi

import (
	"context"
	"errors"
	"fmt"
	"nogo/utils"
	"strings"
	"unicode/utf8"

	notion "github.com/jomei/notionapi"
)

func padLinebreak(s string, n int) string {
	return strings.Replace(s, "\n", "\n"+strings.Repeat(" ", n), -1)
}

func indent(s string, level int) string {
	return strings.Repeat(" ", level) + padLinebreak(s, level)
}

// func jsonPrint(v interface{}) {
// 	s, _ := json.MarshalIndent(v, "", "  ")
// 	fmt.Println(string(s))
// }

func markdownify(rt notion.RichText) string {
	prefix := ""
	suffix := ""
	if rt.Annotations.Bold {
		prefix += "**"
		suffix = "**" + suffix
	}
	if rt.Annotations.Italic {
		prefix += "*"
		suffix = "*" + suffix
	}
	if rt.Annotations.Strikethrough {
		prefix += "~~"
		suffix = "~~" + suffix
	}
	if rt.Annotations.Code {
		prefix += "`"
		suffix = "`" + suffix
	}
	if rt.Annotations.Underline {
		prefix += "<span style=\"text-decoration: underline;\">"
		suffix = "</span>" + suffix
	}
	switch rt.Annotations.Color {
	case "red":
		prefix += string(utils.ColorRed)
		suffix = string(utils.ColorReset) + suffix
	case "green":
		prefix += string(utils.ColorGreen)
		suffix = string(utils.ColorReset) + suffix
	case "blue":
		prefix += string(utils.ColorBlue)
		suffix = string(utils.ColorReset) + suffix
	case "yellow":
		prefix += string(utils.ColorYellow)
		suffix = string(utils.ColorReset) + suffix
	case "purple":
		prefix += string(utils.ColorPurple)
		suffix = string(utils.ColorReset) + suffix
	case "cyan":
		prefix += string(utils.ColorCyan)
		suffix = string(utils.ColorReset) + suffix
	}

	switch rt.Type {
	case "text":
		trimmed := strings.Trim(rt.PlainText, " ")
		decorated := strings.Replace(rt.PlainText, trimmed, prefix+trimmed+suffix, -1)
		return decorated
	case "equation":
		return fmt.Sprintf("$ %s $", rt.PlainText)
	default:
		return rt.PlainText
	}
}

func showTitle(title *notion.TitleProperty) error {
	richtext := title.Title[0]
	if err := showRichText([]notion.RichText{richtext}, "▓ ", 0, utils.ColorCyan); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

func showPageTitle(page *notion.Page) error {
	title := page.Properties["title"].(*notion.TitleProperty)
	if (page.Icon != nil) && (page.Icon.Type == "emoji") {
		title.Title[0].PlainText = fmt.Sprintf("%s  %s", string(*page.Icon.Emoji), title.Title[0].PlainText)
	}
	return showTitle(title)
}

func showRichText(rts []notion.RichText, prefix string, level int, color ...utils.ColorType) error {
	c := utils.ColorReset
	if len(color) > 0 {
		c = color[0]
	}
	fmt.Print(string(c))
	if len(rts) > 0 {
		plain := prefix
		for _, rt := range rts {
			plain += markdownify(rt)
		}
		nprefix := utf8.RuneCountInString(prefix)
		fmt.Println(indent(padLinebreak(plain, nprefix), level))
	} else {
		fmt.Println(indent(prefix, level))
	}
	fmt.Print(string(utils.ColorReset))
	return nil
}

func showParagraph(b notion.Block, level int) error {
	par := b.(*notion.ParagraphBlock).Paragraph
	return showRichText(par.RichText, "", level)
}

func showHeading(b interface{}, level int) error {
	switch b := b.(type) {
	case *notion.Heading1Block:
		h1 := b.Heading1
		return showRichText(h1.RichText, "# ", level)
	case *notion.Heading2Block:
		h2 := b.Heading2
		return showRichText(h2.RichText, "## ", level)
	case *notion.Heading3Block:
		h3 := b.Heading3
		return showRichText(h3.RichText, "### ", level)
	default:
		return errors.New("unknown heading type")
	}
}

func showToDo(b notion.Block, level int) error {
	todo := b.(*notion.ToDoBlock).ToDo
	var check string
	if todo.Checked {
		check = string(utils.ColorGreen) + "✓" + string(utils.ColorReset)
	} else {
		check = " "
	}
	return showRichText(todo.RichText, fmt.Sprintf("[%s] ", check), level)
}

func showBulletedListItem(b notion.Block, level int) error {
	bullet := b.(*notion.BulletedListItemBlock).BulletedListItem
	return showRichText(bullet.RichText, "* ", level)
}

func showNumberedListItem(b notion.Block, level int) error {
	num := b.(*notion.NumberedListItemBlock).NumberedListItem
	prefix := fmt.Sprintf("%d. ", NumberedListCounter)
	return showRichText(num.RichText, prefix, level)
}

func showToggle(c *notion.Client, b notion.Block, open bool, level int) error {
	tblock := b.(*notion.ToggleBlock)
	toggle := tblock.Toggle
	var icon string
	if open {
		icon = "▼"
	} else {
		icon = "▶"
	}
	if err := showRichText(toggle.RichText, fmt.Sprintf("%s ", icon), level); err != nil {
		return err
	}
	if open && tblock.HasChildren {
		children, err := c.Block.GetChildren(context.Background(), notion.BlockID(b.(*notion.ToggleBlock).ID), nil)
		if err != nil {
			return err
		}
		for _, child := range children.Results {
			if err := ShowBlock(c, child, level+2); err != nil {
				return err
			}
		}
	}
	return nil
}

func showEquation(b notion.Block, level int) error {
	eqblock := b.(*notion.EquationBlock)
	equation := eqblock.Equation
	fmt.Println(indent(fmt.Sprintf("$$ %s $$", equation.Expression), level))
	return nil
}

func showCode(b notion.Block, level int) error {
	code := b.(*notion.CodeBlock).Code
	lang := code.Language
	fmt.Println(indent(fmt.Sprintf("```%s", lang), level))
	if err := showRichText(code.RichText, "", level); err != nil {
		return err
	}
	fmt.Println(indent("```", level))
	return nil
}

func showDivider(b notion.Block, level int) error {
	fmt.Println(indent("---", level))
	return nil
}

func showColumn(c *notion.Client, b notion.Block, level int) error {
	col := b.(*notion.ColumnBlock)
	if col.HasChildren {
		blockID := notion.BlockID(col.ID)
		children, err := c.Block.GetChildren(context.Background(), blockID, nil)
		if err != nil {
			return err
		}
		for _, child := range children.Results {
			if err := ShowBlock(c, child, level+2); err != nil {
				return err
			}
		}
	} else {
		fmt.Println()
	}
	return nil
}

func showColumnList(c *notion.Client, b notion.Block, level int) error {
	clist := b.(*notion.ColumnListBlock)
	if clist.HasChildren {
		blockID := notion.BlockID(clist.ID)
		children, err := c.Block.GetChildren(context.Background(), blockID, nil)
		if err != nil {
			return err
		}
		for _, child := range children.Results {
			if err := ShowBlock(c, child, level); err != nil {
				return err
			}
		}
	} else {
		fmt.Println()
	}
	return nil
}

func showImage(b notion.Block, level int) error {
	img := b.(*notion.ImageBlock).Image
	var url string
	if img.Type == "external" {
		url = img.External.URL
	} else if img.Type == "file" {
		url = img.File.URL
	} else {
		return errors.New("unknown image type")
	}
	fmt.Println(indent(fmt.Sprintf("![](%s)", url), level))
	return nil
}

func showChildPage(b notion.Block, level int) error {
	child := b.(*notion.ChildPageBlock).ChildPage
	fmt.Println(indent("░ "+child.Title, level))
	return nil
}
