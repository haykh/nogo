package notionApi

import (
	"context"
	"fmt"

	notion "github.com/jomei/notionapi"
)

func ShowPage(token, pageID string) error {
	client := notion.NewClient(notion.Token(token))
	page, err := client.Page.Get(context.Background(), notion.PageID(pageID))
	if err != nil {
		return err
	}
	showPageTitle(page)
	blocks, err := client.Block.GetChildren(context.Background(), notion.BlockID(pageID), nil)
	if err != nil {
		return err
	}
	for _, block := range blocks.Results {
		ShowBlock(client, block, 0)
	}
	return nil
}

var NumberedListCounter int

func CreatePage(token, parentID string, title, icon string) (string, error) {
	client := notion.NewClient(notion.Token(token))
	parent := notion.Parent{
		Type:   "page_id",
		PageID: notion.PageID(parentID),
	}
	emoji := notion.Emoji(icon)
	pagerequest := notion.PageCreateRequest{
		Parent: parent,
		Properties: map[string]notion.Property{
			"title": notion.TitleProperty{
				Type: "title",
				Title: []notion.RichText{
					{
						Text: notion.Text{
							Content: title,
						},
					},
				},
			},
		},
		Icon: &notion.Icon{
			Type:  "emoji",
			Emoji: &emoji,
		},
	}
	newpage, err := client.Page.Create(context.Background(), &pagerequest)
	if err != nil {
		return "", err
	}
	return string(newpage.ID), nil
}

func ShowBlock(c *notion.Client, b notion.Block, level int) {
	switch b.GetType() {
	case "heading_1", "heading_2", "heading_3":
		showHeading(b, level)
	case "paragraph":
		showParagraph(b, level)
	case "to_do":
		showToDo(b, level)
	case "bulleted_list_item":
		showBulletedListItem(b, level)
	case "numbered_list_item":
		showNumberedListItem(b, level)
	case "toggle":
		showToggle(c, b, true, level)
	case "equation":
		showEquation(b, level)
	case "code":
		showCode(b, level)
	case "divider":
		showDivider(b, level)
	case "column_list":
		showColumnList(c, b, level)
	case "column":
		showColumn(c, b, level)
	case "image":
		showImage(b, level)
	case "child_page":
		showChildPage(b, level)
	default:
		fmt.Println("Unknown block type:", b.GetType())
	}
	if b.GetType() == "numbered_list_item" {
		NumberedListCounter++
	} else {
		NumberedListCounter = 1
	}
}
