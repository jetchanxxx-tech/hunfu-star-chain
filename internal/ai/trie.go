package ai

import "strings"

type Trie struct {
	children map[rune]*Trie
	isEnd    bool
}

func NewTrie(keywords []string) *Trie {
	root := &Trie{children: make(map[rune]*Trie)}
	for _, kw := range keywords {
		root.Insert(strings.ToLower(kw))
	}
	return root
}

func (t *Trie) Insert(word string) {
	node := t
	for _, r := range word {
		if node.children[r] == nil {
			node.children[r] = &Trie{children: make(map[rune]*Trie)}
		}
		node = node.children[r]
	}
	node.isEnd = true
}

func (t *Trie) Search(text string) []string {
	text = strings.ToLower(text)
	runes := []rune(text)
	var found []string
	for i := range runes {
		node := t
		for j := i; j < len(runes); j++ {
			r := runes[j]
			child, ok := node.children[r]
			if !ok {
				break
			}
			node = child
			if node.isEnd {
				found = append(found, string(runes[i:j+1]))
			}
		}
	}
	return found
}
