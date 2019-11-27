/*
A simple skip list
*/
package skiplist

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Node struct {
	score           int
	obj             interface{}
	level           int
	nextNodeByLevel []*Node
}

func NewNode(score int, obj interface{}, level int) (node *Node) {
	nodeLevel := make([]*Node, level+1)
	node = &Node{
		score:           score,
		obj:             obj,
		level:           level,
		nextNodeByLevel: nodeLevel,
	}
	return
}

type SkipList struct {
	head     *Node
	tail     *Node
	maxLevel int
	prob     float64
}

func NewSkipList(maxLevel int) *SkipList {
	head := NewNode(0, nil, maxLevel)
	skipList := &SkipList{
		head:     head,
		tail:     nil,
		maxLevel: maxLevel,
		prob:     0.5,
	}
	return skipList
}

func (s *SkipList) Insert(score int, val interface{}) {
	preNodes := make([]*Node, s.maxLevel+1)
	preNodes[s.maxLevel] = s.head
	for i := s.maxLevel - 1; i >= 0; i-- {
		p := preNodes[i+1]
		if p.nextNodeByLevel[i] == nil {
			preNodes[i] = p
			continue
		}
		for p.nextNodeByLevel[i] != nil && p.nextNodeByLevel[i].score < score {
			p = p.nextNodeByLevel[i]
		}
		preNodes[i] = p
	}
	if preNodes[0].nextNodeByLevel[0] != nil && preNodes[0].nextNodeByLevel[0].score == score {
		preNodes[0].nextNodeByLevel[0].obj = val
	} else {
		level := s.randLevel()
		node := NewNode(score, val, level)
		for i := 0; i <= level; i++ {
			preNodes[i].nextNodeByLevel[i], node.nextNodeByLevel[i] = node, preNodes[i].nextNodeByLevel[i]
		}
	}
	return
}

func (s *SkipList) Delete(score int) {
	var level = s.Level()
	preNodes := make([]*Node, level+1)
	for p := s.head; level >= 0; {
		if p.nextNodeByLevel[level] != nil && p.nextNodeByLevel[level].score == score {
			// got target
			preNodes[level] = p
			level--
		} else if p.nextNodeByLevel[level] == nil || p.nextNodeByLevel[level].score > score {
			//target is not in this level
			preNodes[level] = nil
			level--
		} else if p.nextNodeByLevel[level].score < score {
			p = p.nextNodeByLevel[level]
		}

	}
	for i := range preNodes {
		if preNodes[i] != nil {
			preNodes[i].nextNodeByLevel[i] = preNodes[i].nextNodeByLevel[i].nextNodeByLevel[i]
		}
	}
}

func (s *SkipList) Find(score int) interface{} {
	level := s.Level()
	for p := s.head; level >= 0; {
		if p.nextNodeByLevel[level] == nil || p.nextNodeByLevel[level].score > score {
			//cant find target in this level , goto next level
			level--
		} else if p.nextNodeByLevel[level].score == score {
			return p.nextNodeByLevel[level].obj
		} else if p.nextNodeByLevel[level].score < score {
			p = p.nextNodeByLevel[level]
		}
	}
	return nil
}

func (s *SkipList) String() string {
	strs := make([]string, s.maxLevel+1)
	for p := s.head.nextNodeByLevel[0]; p != nil; p = p.nextNodeByLevel[0] {
		for i := 0; i < s.maxLevel; i++ {
			if i <= p.level {
				strs[s.maxLevel-i] += fmt.Sprintf("%v\t", p.obj)
			} else {
				strs[s.maxLevel-i] += fmt.Sprintf("%v\t", "-")
			}
		}
	}
	return strings.Join(strs, "\n")
}

func (s *SkipList) Level() int {
	for i := s.maxLevel - 1; i > 0; i-- {
		if s.head.nextNodeByLevel[i] != nil {
			return i
		}
	}
	return 0
}

func (s *SkipList) Len() int {
	p := s.head.nextNodeByLevel[0]
	l := 0
	for ; p != nil; l++ {
		p = p.nextNodeByLevel[0]
	}
	return l
}

func (s *SkipList) randLevel() int {
	return randLevel(s.maxLevel, func() bool {
		if rand.Float64() < s.prob {
			return true
		}
		return false
	})

}
func randLevel(max int, rand func() bool) int {
	l := 0
	for l < max {
		if rand() {
			l++
		} else {
			return l
		}
	}
	return l
}
