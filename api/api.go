package notionApi

import (
	"context"
	"errors"
	"nogo/config"
	"nogo/utils"

	survey "github.com/AlecAivazis/survey/v2"
	notion "github.com/jomei/notionapi"
)

func InitAPI() (*notion.Client, string, error) {
	if loc_config, err := config.CreateOrReadLocalConfig(true); err != nil {
		return nil, "", err
	} else {
		if token, err := loc_config.GetSecret("api_token"); err != nil {
			return nil, "", err
		} else {
			if stackID, err := loc_config.GetSecret("stack_page_id"); err != nil {
				return nil, "", err
			} else {
				return NewClient(token), stackID, nil
			}
		}
	}
}

func NewClient(token string) *notion.Client {
	return notion.NewClient(notion.Token(token))
}

func AddToStack(client *notion.Client, pageID string) error {
	if parent, err := GetStack(client, pageID); err != nil {
		return err
	} else {
		new_item := ""
		survey.AskOne(&survey.Input{
			Message: "new entry:",
		}, &new_item)
		if new_item == "" {
			return errors.New("empty entry")
		}
		if _, err := client.Block.AppendChildren(context.Background(), parent.GetID(), &notion.AppendBlockChildrenRequest{
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
									Content: new_item,
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
			return nil
		}
	}
}

func ModifyStack(client *notion.Client, pageID string) error {
	if blocks, err := GetStackEntries(client, pageID); err != nil {
		return err
	} else {
		if stack, plain, _, err := ParseStackFromBlocks(client, blocks, pageID); err != nil {
			return err
		} else {
			idx := -1
			survey.AskOne(
				&survey.Select{
					Message: "modify:",
					Options: *stack,
				},
				&idx,
				survey.WithPageSize(10),
			)
			if idx == -1 {
				return errors.New("no selection")
			}
			new_item := ""
			survey.AskOne(&survey.Input{
				Message: "new entry:",
				Suggest: func(string) []string {
					return []string{(*plain)[idx]}
				},
			}, &new_item)
			if new_item == "" {
				return errors.New("empty entry")
			}
			if _, err := client.Block.Update(context.Background(), blocks[idx].GetID(), &notion.BlockUpdateRequest{
				ToDo: &notion.ToDo{
					RichText: []notion.RichText{
						{
							Text: notion.Text{
								Content: new_item,
							},
						},
					},
				},
			}); err != nil {
				return err
			} else {
				return nil
			}
		}
	}
}

func RmFromStack(client *notion.Client, pageID string) error {
	if blocks, err := GetStackEntries(client, pageID); err != nil {
		return err
	} else {
		if stack, _, _, err := ParseStackFromBlocks(client, blocks, pageID); err != nil {
			return err
		} else {
			torm := []int{}
			survey.AskOne(
				&survey.MultiSelect{
					Message: "pick to rm:",
					Options: *stack,
				},
				&torm,
				survey.WithPageSize(10),
				survey.WithIcons(func(icons *survey.IconSet) {
					icons.MarkedOption.Text = "✖"
					icons.MarkedOption.Format = "red"
					icons.UnmarkedOption.Text = " "
				}),
			)
			for _, idx := range torm {
				if _, err := client.Block.Delete(context.Background(), blocks[idx].GetID()); err != nil {
					return err
				}
			}
			return nil
		}
	}
}

func ToggleStack(client *notion.Client, pageID string) error {
	if blocks, err := GetStackEntries(client, pageID); err != nil {
		return err
	} else {
		if _, stack, marked, err := ParseStackFromBlocks(client, blocks, pageID); err != nil {
			return err
		} else {
			preselect := []string{}
			for i, m := range *marked {
				if m {
					preselect = append(preselect, (*stack)[i])
				}
			}
			selected := []int{}
			survey.AskOne(
				&survey.MultiSelect{
					Message: "toggle:",
					Options: *stack,
					Default: preselect,
				},
				&selected,
				survey.WithPageSize(10),

				survey.WithIcons(func(icons *survey.IconSet) {
					icons.MarkedOption.Text = "[✓]"
					icons.MarkedOption.Format = "green"
					icons.UnmarkedOption.Text = "[ ]"
				}),
			)
			for mi, m := range *marked {
				isin := utils.IsIn(mi, selected)
				if (!m && isin) || (m && !isin) {
					if _, err := client.Block.Update(context.Background(), blocks[mi].GetID(), &notion.BlockUpdateRequest{
						ToDo: &notion.ToDo{
							Checked: isin,
						},
					}); err != nil {
						return err
					}
				}
			}
			return nil
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
