package api

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/haykh/nogo/utils"

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

var NumberedListCounter int

func Block2String(c *notion.Client, b notion.Block, level int) (string, error) {
	if b.GetType() == "numbered_list_item" {
		defer func() {
			NumberedListCounter++
		}()
	} else {
		defer func() {
			NumberedListCounter = 1
		}()
	}
	switch b.GetType() {
	case "heading_1", "heading_2", "heading_3":
		return Heading2String(b, level), nil
	case "paragraph":
		return Paragraph2String(b, level), nil
	case "to_do":
		return ToDo2String(b, level), nil
	case "bulleted_list_item":
		return BulletedListItem2String(b, level), nil
	case "numbered_list_item":
		return NumberedListItem2String(b, level), nil
	case "toggle":
		return Toggle2String(c, b, true, level)
	case "equation":
		return Equation2String(b, level), nil
	case "code":
		return Code2String(b, level), nil
	case "divider":
		return Divider2String(b, level), nil
	case "column_list":
		return ColumnList2String(c, b, level)
	case "column":
		return Column2String(c, b, level)
	case "image":
		return Image2String(b, level), nil
	case "child_page":
		return ChildPage2String(b, level), nil
	case "synced_block":
		if blocks, err := c.Block.GetChildren(context.Background(), notion.BlockID(b.(*notion.SyncedBlock).ID), nil); err != nil {
			return "", err
		} else {
			result := ""
			for _, block := range blocks.Results {
				if str, err := Block2String(c, block, level); err != nil {
					return "", err
				} else {
					result += str
				}
			}
			return result, nil
		}
	default:
		return "", fmt.Errorf("unknown block type: %s", b.GetType())
	}
}

func RichText2String(rts []notion.RichText, prefix string, level int, style ...utils.TextStyle) string {
	reset_all := utils.ColorReset + utils.HiReset
	result := ""
	if len(style) > 0 {
		for s := range style {
			result = result + style[s]
		}
	}
	if len(rts) > 0 {
		plain := prefix
		for _, rt := range rts {
			plain += markdownify(rt)
		}
		nprefix := utf8.RuneCountInString(prefix)
		result += indent(padLinebreak(plain, nprefix), level)
	} else {
		result += indent(prefix, level)
	}
	if len(style) > 0 {
		result += reset_all
	}
	return result + "\n"
}

func Title2String(title *notion.TitleProperty) string {
	richtext := title.Title[0]
	return RichText2String([]notion.RichText{richtext}, "▓ ", 0, utils.ColorCyan) + "\n"
}

func PageTitle2String(page *notion.Page) string {
	title := page.Properties["title"].(*notion.TitleProperty)
	if (page.Icon != nil) && (page.Icon.Type == "emoji") {
		title.Title[0].PlainText = fmt.Sprintf("%s  %s", string(*page.Icon.Emoji), title.Title[0].PlainText)
	}
	return Title2String(title)
}

func Paragraph2String(b notion.Block, level int) string {
	return RichText2String(b.(*notion.ParagraphBlock).Paragraph.RichText, "", level)
}

func ToDo2String(b notion.Block, level int) string {
	todo := b.(*notion.ToDoBlock).ToDo
	var check string
	if todo.Checked {
		check = string(utils.ColorGreen) + "✓" + string(utils.ColorReset)
	} else {
		check = " "
	}
	return RichText2String(todo.RichText, fmt.Sprintf("[%s] ", check), level)
}

func Heading2String(b interface{}, level int) string {
	switch b := b.(type) {
	case *notion.Heading1Block:
		return RichText2String(b.Heading1.RichText, "# ", level)
	case *notion.Heading2Block:
		return RichText2String(b.Heading2.RichText, "## ", level)
	case *notion.Heading3Block:
		return RichText2String(b.Heading3.RichText, "### ", level)
	default:
		return ""
	}
}

func BulletedListItem2String(b notion.Block, level int) string {
	return RichText2String(b.(*notion.BulletedListItemBlock).BulletedListItem.RichText, "* ", level)
}

func NumberedListItem2String(b notion.Block, level int) string {
	num := b.(*notion.NumberedListItemBlock).NumberedListItem
	prefix := fmt.Sprintf("%d. ", NumberedListCounter)
	return RichText2String(num.RichText, prefix, level)
}

func Toggle2String(c *notion.Client, b notion.Block, open bool, level int) (string, error) {
	tblock := b.(*notion.ToggleBlock)
	toggle := tblock.Toggle
	var icon string
	if open {
		icon = "▼"
	} else {
		icon = "▶"
	}
	result := RichText2String(toggle.RichText, fmt.Sprintf("%s ", icon), level)
	if open && tblock.HasChildren {
		children, err := c.Block.GetChildren(context.Background(), notion.BlockID(b.(*notion.ToggleBlock).ID), nil)
		if err != nil {
			return "", err
		}
		for _, child := range children.Results {
			if bl, err := Block2String(c, child, level+2); err != nil {
				return "", err
			} else {
				result += bl
			}
		}
	}
	return result, nil
}

func Equation2String(b notion.Block, level int) string {
	eqblock := b.(*notion.EquationBlock)
	equation := eqblock.Equation
	return indent(fmt.Sprintf("$$ %s $$", equation.Expression), level)
}

func Code2String(b notion.Block, level int) string {
	code := b.(*notion.CodeBlock).Code
	lang := code.Language
	result := indent(fmt.Sprintf("```%s", lang), level)
	result += RichText2String(code.RichText, "", level)
	result += indent("```", level)
	return result
}

func Divider2String(b notion.Block, level int) string {
	return indent("---", level)
}

func Column2String(c *notion.Client, b notion.Block, level int) (string, error) {
	col := b.(*notion.ColumnBlock)
	result := ""
	if col.HasChildren {
		blockID := notion.BlockID(col.ID)
		if children, err := c.Block.GetChildren(context.Background(), blockID, nil); err != nil {
			return "", err
		} else {
			for _, child := range children.Results {
				if bl, err := Block2String(c, child, level+2); err != nil {
					return "", err
				} else {
					result += bl
				}
			}
		}
	} else {
		result += "\n"
	}
	return result, nil
}

func ColumnList2String(c *notion.Client, b notion.Block, level int) (string, error) {
	clist := b.(*notion.ColumnListBlock)
	result := ""
	if clist.HasChildren {
		blockID := notion.BlockID(clist.ID)
		if children, err := c.Block.GetChildren(context.Background(), blockID, nil); err != nil {
			return "", err
		} else {
			for _, child := range children.Results {
				if bl, err := Block2String(c, child, level); err != nil {
					return "", err
				} else {
					result += bl
				}
			}
		}
	} else {
		result += "\n"
	}
	return result, nil
}

func Image2String(b notion.Block, level int) string {
	img := b.(*notion.ImageBlock).Image
	url := ""
	if img.Type == "external" {
		url = img.External.URL
	} else if img.Type == "file" {
		url = img.File.URL
	} else {
		url = string(img.Type)
	}
	return indent(fmt.Sprintf("![](%s)", url), level)
}

func ChildPage2String(b notion.Block, level int) string {
	child := b.(*notion.ChildPageBlock).ChildPage
	return indent("░ "+child.Title, level)
}
