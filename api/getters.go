package api

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/haykh/nogo/utils"

	notion "github.com/jomei/notionapi"
)

func GetStack(client *notion.Client, pageID string) (notion.Block, error) {
	if parent, err := client.Block.GetChildren(context.Background(), notion.BlockID(pageID), nil); err != nil {
		return nil, err
	} else if len(parent.Results) == 0 {
		return nil, fmt.Errorf("no parent found")
	} else {
		return parent.Results[0], nil
	}
}

func GetStackEntries(client *notion.Client, pageID string) (notion.Blocks, error) {
	if stack, err := GetStack(client, pageID); err != nil {
		return nil, err
	} else {
		if entries, err := client.Block.GetChildren(context.Background(), stack.GetID(), nil); err != nil {
			return nil, err
		} else {
			return entries.Results, nil
		}
	}
}

func ParseStack(client *notion.Client, pageID string) (*[]string, *[]string, *[]bool, error) {
	if entries, err := GetStackEntries(client, pageID); err != nil {
		return nil, nil, nil, err
	} else {
		return ParseStackFromBlocks(client, entries, pageID)
	}
}

func ParseStackFromBlocks(client *notion.Client, blocks notion.Blocks, pageID string) (*[]string, *[]string, *[]bool, error) {
	rich := []string{}
	plain := []string{}
	marked := []bool{}
	for _, entry := range blocks {
		if todo_str, err := Block2String(client, entry, 0); err != nil {
			return nil, nil, nil, err
		} else {
			rt := utils.Clean(todo_str)
			pl := utils.Clean(regexp.MustCompile(`\[.*?\]`).ReplaceAllString(rt, ""))
			rich = append(rich, rt)
			plain = append(plain, pl)
			marked = append(marked, strings.Contains(rt, "✓"))
		}
	}
	return &rich, &plain, &marked, nil
}
