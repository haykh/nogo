package notionApi

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	notion "github.com/jomei/notionapi"
)

func ShowPage(token, pageID string) error {
	client := notion.NewClient(notion.Token(token))
	if page, err := client.Page.Get(context.Background(), notion.PageID(pageID)); err != nil {
		return err
	} else {
		if err := showPageTitle(page); err != nil {
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

func FindInStack(token, pageID, request string) (notion.Block, error) {
	client := notion.NewClient(notion.Token(token))
	if blocks, err := client.Block.GetChildren(context.Background(), notion.BlockID(pageID), nil); err != nil {
		return nil, err
	} else {
		if len(blocks.Results) == 0 {
			return nil, fmt.Errorf("no blocks found")
		}
		if entries, err := client.Block.GetChildren(context.Background(), blocks.Results[0].GetID(), nil); err != nil {
			return nil, err
		} else {
			block_idx := -1
			if block_i, err := strconv.Atoi(request); err == nil {
				if block_i < 0 {
					block_i += len(entries.Results) + 1
				}
				block_idx = block_i - 1
			} else {
				nblocks := 0
				for i, block := range entries.Results {
					if block.GetType() == "to_do" {
						plain := ""
						for _, rt := range block.(*notion.ToDoBlock).ToDo.RichText {
							plain += markdownify(rt)
						}
						find := regexp.MustCompile(strings.ToLower(request))
						if find.MatchString(strings.ToLower(plain)) {
							nblocks++
							block_i = i
						}
					}
				}
				if nblocks > 1 {
					return nil, fmt.Errorf("more than one block satisfies the criterion")
				} else if block_i != -1 {
					block_idx = block_i
				} else {
					return nil, fmt.Errorf("could not find the stack entry")
				}
			}
			if block_idx >= 0 && block_idx < len(entries.Results) {
				return entries.Results[block_idx], nil
			} else {
				return nil, fmt.Errorf("invalid block index")
			}
		}
	}
}

func AddToStack(token, pageID, request string) error {
	client := notion.NewClient(notion.Token(token))
	if blocks, err := client.Block.GetChildren(context.Background(), notion.BlockID(pageID), nil); err != nil {
		return err
	} else {
		if len(blocks.Results) == 0 {
			return fmt.Errorf("no blocks found")
		}
		if response, err := client.Block.AppendChildren(context.Background(), blocks.Results[0].GetID(), &notion.AppendBlockChildrenRequest{
			Children: []notion.Block{
				&notion.ToDoBlock{
					BasicBlock: notion.BasicBlock{
						Object: notion.ObjectTypeBlock,
						Type:   notion.BlockTypeToDo,
					},
					ToDo: notion.ToDo{
						RichText: []notion.RichText{
							{
								Text: notion.Text{
									Content: request,
								},
							},
						},
						Checked: false,
					},
				},
			},
		}); err != nil {
			return err
		} else {
			for _, block := range response.Results {
				fmt.Printf("added new stack entry:\n")
				if err := ShowBlock(client, block, 1); err != nil {
					return err
				}
			}
			fmt.Println()
			return nil
		}
	}
}

func RmFromStack(token, pageID, request string) error {
	if block, err := FindInStack(token, pageID, request); err != nil {
		return err
	} else {
		client := notion.NewClient(notion.Token(token))
		if _, err := client.Block.Delete(context.Background(), block.GetID()); err != nil {
			return err
		} else {
			return ShowBlock(client, block, 1)
		}
	}
}

func MarkAs(token, pageID, request string, checked bool) error {
	if block, err := FindInStack(token, pageID, request); err != nil {
		return nil
	} else {
		client := notion.NewClient(notion.Token(token))
		if blocknew, err := client.Block.Update(context.Background(), block.GetID(), &notion.BlockUpdateRequest{
			ToDo: &notion.ToDo{
				Checked: checked,
			},
		}); err != nil {
			return err
		} else {
			as := "not done"
			if checked {
				as = "done"
			}
			fmt.Printf("marked the stack entry as %s:\n", as)
			return ShowBlock(client, blocknew, 1)
		}
	}
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

func ShowBlock(c *notion.Client, b notion.Block, level int) error {
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
		return showHeading(b, level)
	case "paragraph":
		return showParagraph(b, level)
	case "to_do":
		return showToDo(b, level)
	case "bulleted_list_item":
		return showBulletedListItem(b, level)
	case "numbered_list_item":
		return showNumberedListItem(b, level)
	case "toggle":
		return showToggle(c, b, true, level)
	case "equation":
		return showEquation(b, level)
	case "code":
		return showCode(b, level)
	case "divider":
		return showDivider(b, level)
	case "column_list":
		return showColumnList(c, b, level)
	case "column":
		return showColumn(c, b, level)
	case "image":
		return showImage(b, level)
	case "child_page":
		return showChildPage(b, level)
	case "synced_block":
		blocks, err := c.Block.GetChildren(context.Background(), notion.BlockID(b.(*notion.SyncedBlock).ID), nil)
		if err != nil {
			return err
		}
		for _, block := range blocks.Results {
			if err := ShowBlock(c, block, level); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unknown block type: %s", b.GetType())
	}
	return nil
}
