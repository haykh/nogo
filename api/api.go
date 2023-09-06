package notionApi

import (
	"context"
	"fmt"
	"nogo/utils"
	"regexp"
	"strconv"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	notion "github.com/jomei/notionapi"
)

type ToDoStack = []string

func NewClient(token string) *notion.Client {
	return notion.NewClient(notion.Token(token))
}

func ParseStack(client *notion.Client, pageID string) (*ToDoStack, error) {
	if blocks, err := client.Block.GetChildren(context.Background(), notion.BlockID(pageID), nil); err != nil {
		return nil, err
	} else {
		if len(blocks.Results) == 0 {
			return nil, fmt.Errorf("no todo items found")
		}
		if entries, err := client.Block.GetChildren(context.Background(), blocks.Results[0].GetID(), nil); err != nil {
			return nil, err
		} else {
			todo_stack := ToDoStack{}
			for _, entry := range entries.Results {
				if todo_str, err := Block2String(client, entry, 0); err != nil {
					return nil, err
				} else {
					todo_stack = append(todo_stack, utils.Clean(todo_str))
				}
			}
			return &todo_stack, nil
		}
	}
}

func FindInStack(client *notion.Client, pageID, request string) (notion.Block, error) {
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

func AddToStack(client *notion.Client, pageID, request string) error {
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

func ModifyStack(client *notion.Client, pageID, request string) error {
	fmt.Println(request)
	return nil
}

func RmFromStack(client *notion.Client, pageID string) error {
	if stack, err := ParseStack(client, pageID); err != nil {
		return err
	} else {
		entry_idx := -1
		entries := []string{}
		entries = append(entries, *stack...)
		survey.AskOne(&survey.Select{
			Message: "pick to rm:",
			Options: entries,
		}, &entry_idx)
		// if entry_dx >= 0 && entry_idx < len(entries)
		fmt.Printf("selected %d\n", entry_idx)
		return nil
		// if block, err := FindInStack(client, pageID, entry); err != nil {
		// 	return err
		// } else {
		// 	if _, err := client.Block.Delete(context.Background(), block.GetID()); err != nil {
		// 		return err
		// 	} else {
		// 		return ShowBlock(client, block, 1)
		// 	}
		// }
	}
}

func MarkAs(client *notion.Client, pageID, request string, checked bool) error {
	if block, err := FindInStack(client, pageID, request); err != nil {
		return nil
	} else {
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

func CreatePage(client *notion.Client, parentID string, title, icon string) (string, error) {
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
	if newpage, err := client.Page.Create(context.Background(), &pagerequest); err != nil {
		return "", err
	} else {
		return string(newpage.ID), nil
	}
}
