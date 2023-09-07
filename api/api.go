package api

import (
	"context"
	"errors"

	"github.com/haykh/nogo/config"
	"github.com/haykh/nogo/utils"

	survey "github.com/AlecAivazis/survey/v2"
	notionapi "github.com/jomei/notionapi"
)

func InitAPI() (*notionapi.Client, string, error) {
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

func NewClient(token string) *notionapi.Client {
	return notionapi.NewClient(notionapi.Token(token))
}

func AddToStack(client *notionapi.Client, pageID string) error {
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
		if _, err := client.Block.AppendChildren(context.Background(), parent.GetID(), &notionapi.AppendBlockChildrenRequest{
			Children: []notionapi.Block{
				&notionapi.ToDoBlock{
					BasicBlock: notionapi.BasicBlock{
						Object: notionapi.ObjectTypeBlock,
						Type:   notionapi.BlockTypeToDo,
					},
					ToDo: notionapi.ToDo{
						RichText: []notionapi.RichText{
							{
								Text: &notionapi.Text{
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

func ModifyStack(client *notionapi.Client, pageID string) error {
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
			if _, err := client.Block.Update(context.Background(), blocks[idx].GetID(), &notionapi.BlockUpdateRequest{
				ToDo: &notionapi.ToDo{
					RichText: []notionapi.RichText{
						{
							Text: &notionapi.Text{
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

func RmFromStack(client *notionapi.Client, pageID string) error {
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

func ToggleStack(client *notionapi.Client, pageID string) error {
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
					request := blocks[mi].(*notionapi.ToDoBlock).ToDo
					request.Checked = isin
					if _, err := client.Block.Update(
						context.Background(),
						blocks[mi].GetID(),
						&notionapi.BlockUpdateRequest{
							ToDo: &request,
						}); err != nil {
						return err
					}
				}
			}
			return nil
		}
	}
}

func CreatePage(client *notionapi.Client, parentID string, title, icon string) (string, error) {
	parent := notionapi.Parent{
		Type:   "page_id",
		PageID: notionapi.PageID(parentID),
	}
	emoji := notionapi.Emoji(icon)
	pagerequest := notionapi.PageCreateRequest{
		Parent: parent,
		Properties: map[string]notionapi.Property{
			"title": notionapi.TitleProperty{
				Type: "title",
				Title: []notionapi.RichText{
					{
						Text: &notionapi.Text{
							Content: title,
						},
					},
				},
			},
		},
		Icon: &notionapi.Icon{
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
