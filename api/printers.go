package notionApi

import (
	"context"
	"fmt"
	"nogo/utils"

	notion "github.com/jomei/notionapi"
)

func ShowPage(client *notion.Client, pageID string) error {
	if page, err := client.Page.Get(context.Background(), notion.PageID(pageID)); err != nil {
		return err
	} else {
		if err := ShowPageTitle(page); err != nil {
			return err
		}
		if blocks, err := client.Block.GetChildren(context.Background(), notion.BlockID(pageID), nil); err != nil {
			return err
		} else {
			for _, block := range blocks.Results {
				if err := ShowBlock(client, block, 0); err != nil {
					return err
				}
			}
			return nil
		}
	}
}

func ShowBlock(c *notion.Client, b notion.Block, level int) error {
	if str, err := Block2String(c, b, level); err != nil {
		return err
	} else {
		fmt.Print(str)
		return nil
	}
}

func ShowTitle(title *notion.TitleProperty) error {
	fmt.Print(Title2String(title))
	return nil
}

func ShowPageTitle(page *notion.Page) error {
	fmt.Print(PageTitle2String(page))
	return nil
}

func ShowRichText(rts []notion.RichText, prefix string, level int, style ...utils.TextStyle) error {
	fmt.Print(RichText2String(rts, prefix, level, style...))
	return nil
}

func ShowParagraph(b notion.Block, level int) error {
	fmt.Print(Paragraph2String(b, level))
	return nil
}

func ShowHeading(b interface{}, level int) error {
	switch b.(type) {
	case *notion.Heading1Block, *notion.Heading2Block, *notion.Heading3Block:
		fmt.Print(Heading2String(b, level))
		return nil
	default:
		return fmt.Errorf("unknown heading type")
	}
}

func ShowToDo(b notion.Block, level int) error {
	fmt.Print(ToDo2String(b, level))
	return nil
}

func ShowBulletedListItem(b notion.Block, level int) error {
	fmt.Print(BulletedListItem2String(b, level))
	return nil
}

func ShowNumberedListItem(b notion.Block, level int) error {
	fmt.Print(NumberedListItem2String(b, level))
	return nil
}

func ShowToggle(c *notion.Client, b notion.Block, open bool, level int) error {
	if str, err := Toggle2String(c, b, open, level); err != nil {
		return err
	} else {
		fmt.Print(str)
		return nil
	}
}

func ShowEquation(b notion.Block, level int) error {
	fmt.Print(Equation2String(b, level))
	return nil
}

func ShowCode(b notion.Block, level int) error {
	fmt.Print(Code2String(b, level))
	return nil
}

func ShowDivider(b notion.Block, level int) error {
	fmt.Print(Divider2String(b, level))
	return nil
}

func ShowColumn(c *notion.Client, b notion.Block, level int) error {
	if str, err := Column2String(c, b, level); err != nil {
		return err
	} else {
		fmt.Print(str)
		return nil
	}
}

func ShowColumnList(c *notion.Client, b notion.Block, level int) error {
	if str, err := ColumnList2String(c, b, level); err != nil {
		return err
	} else {
		fmt.Print(str)
		return nil
	}
}

func ShowImage(b notion.Block, level int) error {
	fmt.Print(Image2String(b, level))
	return nil
}

func ShowChildPage(b notion.Block, level int) error {
	fmt.Print(ChildPage2String(b, level))
	return nil
}
