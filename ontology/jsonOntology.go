package ontology

import (
	"encoding/json"
	"io/ioutil"
)

type Emotion string

const (
	Sad   Emotion = "sad"
	Happy         = "happy"
	Hard          = "hard"
	Neutral       = "neutral"
)

type Node struct {
	Name      string  `json:"name"`
	Emotional Emotion `json:"emotional"`
	IsModern  bool    `json:"isModern"`
	Volume    int     `json:"volume"`
	Century   int     `json:"century"`
	
	Nodes []Node `json:"node"`
}

type JsonOntology struct {
	Root Node `json:"root"`
}

func NewJsonOntology(filename string) (*JsonOntology, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	ont := &JsonOntology{}

	err = json.Unmarshal(file, ont)
	if err != nil {
		return nil, err
	}

	return ont, nil
}

func (jO *JsonOntology) GetAllSheets(node Node) []Node {
	sheets := []Node{}
	if node.Nodes == nil {
		sheets = append(sheets, node)
		return sheets
	}

	for _, n := range node.Nodes {
		newSheets := jO.GetAllSheets(n)
		for _, nn := range newSheets {
			sheets = append(sheets, nn)
		}
	}

	return sheets
}
