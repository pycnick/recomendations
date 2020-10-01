package ontology

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"os"
)

type Class struct {
	Name string `xml:"name,attr"`
}

type Declaration struct {
	Class *Class `json:"class"`
}

type SubClassOf struct {
	Classes []*Class `xml:"Class", json:"classes"`
}

type Owl struct {
	Declarations []*Declaration `xml:"Declaration", json:"declarations"`
	Relations    []*SubClassOf  `xml:"SubClassOf", json:"relations"`
}

func NewOwl(filename string) (*Owl, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	owl := &Owl{}

	err = xml.Unmarshal(file, owl)
	if err != nil {
		return nil, err
	}

	return owl, nil
}

func (o *Owl) RawJson() ([]byte, error) {
	buf, err := json.Marshal(o)
	if err != nil {
		return []byte{}, err
	}

	return buf, nil
}

func (o *Owl) SaveToFile(filename string) error {
	file, err := os.Create(os.Getenv("PWD") + filename)
	defer file.Close()
	if err != nil {
		return err
	}

	buf, err := o.RawJson()
	if err != nil {
		return err
	}

	_, err = file.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func (o *Owl) SaveJsonToFile(filename string, obj *JsonOntology) error {
	file, err := os.Create(os.Getenv("PWD") + filename)
	defer file.Close()
	if err != nil {
		return err
	}

	buf, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = file.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func buildCategories(ids []string, relations map[string][]string) []Node {
	categories := make([]Node, len(ids))
	for i, id := range ids {
		c := Node{Name: id}
		if childIDs, ok := relations[id]; ok {
			c.Nodes = buildCategories(childIDs, relations)
		}
		categories[i] = c
	}
	return categories
}

func (o *Owl) GetJsonOntology() *JsonOntology {
	m := make(map[string][]string)
	for _, decl := range o.Relations {
		m[decl.Classes[1].Name] = append(m[decl.Classes[1].Name], decl.Classes[0].Name)
	}

	res := buildCategories([]string{"Direction"}, m)
	return &JsonOntology{Root: res[0]}
}
