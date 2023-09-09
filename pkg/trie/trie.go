package trie

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type Trie struct {
	root *trie_Node
}

type trie_Node struct {
	children   [26]*trie_Node
	wordEnds   bool
	childCount int
}

func New() *Trie {
	t := new(Trie)
	t.root = new(trie_Node)
	return t
}

func NewFromFile(path string) (*Trie, error) {
	words, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	t := New()
	for _, word := range strings.Split(string(words), "\n") {
		t.Insert(word)
	}
	return t, nil
}

func (t *Trie) Insert(word string) {
	current := t.root
	current.childCount++
	for _, wr := range word {
		index := wr - 'a'
		if current.children[index] == nil {
			current.children[index] = new(trie_Node)
		}
		current = current.children[index]
		current.childCount++
	}
	current.wordEnds = true
}

func (t *Trie) Search(word string) bool {
	current := t.root
	for _, wr := range word {
		index := wr - 'a'
		if current.children[index] == nil {
			return false
		}
		current = current.children[index]
	}
	return current.wordEnds
}

func (t *Trie) Random() string {
	index := rand.Intn(t.root.childCount)
	return t.getWordAtIndex(index)
}

func (t *Trie) getWordAtIndex(index int) string {
	if index >= t.root.childCount {
		return ""
	}
	current := 0
	node := t.root

	word := ""
	for {
		for c, n := range node.children {
			if n != nil {
				if current+n.childCount > index+1 {
					word += fmt.Sprintf("%c", 'a'+c)
					node = n
					break
				} else if current+n.childCount == index+1 {
					word += fmt.Sprintf("%c", 'a'+c)
					if n.wordEnds {
						return word
					}

					node = n
					break
				}
				current += n.childCount
			}
		}
	}
}
