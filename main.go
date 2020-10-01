package main

import (
	"fmt"
	"github.com/pycnick/recomendations/ontology"
	"math/rand"
	"sort"
)

var (
	users = func() []string {
		names := []string{"John", "Nick", "Sam", "Dude"}
		var users []string
		for i := 0; i < 4; i++ {
			for j := 0; j < 25; j++ {
				users = append(users, fmt.Sprintf("%s%d", names[i], j))
			}
		}

		return users
	}()

	block = make(map[int][]int)
)

type CollabFilter struct {
	Users []string
	Products []string
	Matrix [][]float64
}

func NewCollabFilter(users []string, products[]string) *CollabFilter {
	filter := &CollabFilter{
		Users: users,
		Products: products,
	}

	matrix := make([][]float64, len(users))
	for i := 0; i < len(users); i++ {
		matrix[i] = make([]float64, len(products))
	}
	for i := range users {
		for j := range products {
			var value float64
			prob := rand.Float64()
			if prob <= 0.4 {
				value = 0 + rand.Float64() * 5
			}
			matrix[i][j] = value
		}
	}

	filter.Matrix = matrix
	return filter
}

func (cF *CollabFilter) getUserID(user string) int {
	for i := range cF.Users {
		if cF.Users[i] == user {
			return i
		}
	}

	return -1
}

func (cF *CollabFilter) getOtherUsers(index int) map[int]float64 {
	aroundUsers := make(map[int]float64)
	for i, row := range cF.Matrix {
		if i == index {
			continue
		}

		var sum float64
		var a float64
		var b float64
		for j, _ := range row {
			sum += cF.Matrix[i][j] * cF.Matrix[index][j]
			a += cF.Matrix[i][j] * cF.Matrix[i][j]
			b += cF.Matrix[index][j] * cF.Matrix[index][j]
		}

		aroundUsers[i] = sum / a / b
	}

	return aroundUsers
}

func (cF *CollabFilter) BlockProduct(user string, product string) {
	index := cF.getUserID(user)

	for i := 0; i < len(cF.Products); i++ {
		if cF.Products[i] == product {
			block[index] = append(block[index], i)
			cF.Matrix[index][i] = -5
			fmt.Println(block[index])
			return
		}
	}
}

func (cF *CollabFilter) UnBlockProduct(user string, product int) {
	index := cF.getUserID(user)

	b := block[index]
	if len(b) == 0 {
		return
	}

	for i := 0; i < len(b); i++ {
		if b[i] == product {
			copy(b[:i], b[i+1:])
			block[index] = b
			return
		}
	}
	fmt.Println(block[index])
}

func (cF *CollabFilter) isProductBlock(userID, productID int) bool {
	userBlocks := block[userID]
	if userBlocks == nil {
		return false
	}

	for _, val := range userBlocks {
		if val == productID {
			return true
		}
	}

	return false
}

func (cF *CollabFilter) GetRecommendation(user string) []string {
	index := cF.getUserID(user)

	aroundUsers := cF.getOtherUsers(index)

	keys := make([]int, 0, len(aroundUsers))
	for key := range aroundUsers {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return aroundUsers[keys[i]] > aroundUsers[keys[j]] })

	var ind int
	var best []int
	for _, key := range keys {
		if ind == 10 {
			break
		}
		best = append(best, key)
		fmt.Printf("%s, %d\n", key, aroundUsers[key])
		ind++
	}

	productsRecommend := make([]struct {Index int; Value float64}, len(cF.Products))
	for j, _ := range cF.Matrix[0] {
		if cF.Matrix[index][j] != 0 || cF.isProductBlock(index, j){
			continue
		}

		var sum float64
		for i, _ := range keys {
			sum += cF.Matrix[keys[i]][j] * aroundUsers[keys[i]]
		}
		productsRecommend[j].Value = sum
		productsRecommend[j].Index = j
	}

	sort.Slice(productsRecommend, func(i, j int) bool { return productsRecommend[i].Value > productsRecommend[j].Value })

	var productNames []string
	for i := 0; i < 5; i++ {
		productNames = append(productNames, cF.Products[productsRecommend[i].Index])
	}

	return productNames
}




func main() {
	//owl, _ := ontology.NewOwl("./ontology/raw/ont.owl")
	//_ = owl.SaveToFile("/ontology/raw/ont.json")
	//
	//obj := owl.GetJsonOntology()
	//_ = owl.SaveJsonToFile("/ontology/raw/true.json", obj)

	ont, _ := ontology.NewJsonOntology("./ontology/raw/true.json")
	sheets := ont.GetAllSheets(ont.Root)

	products := []string{}
	for _, v := range sheets {
		products = append(products, v.Name)
	}

	f := NewCollabFilter(users, products)

	for {
		fmt.Println("Select menu tab:\n1) Block product (id)\n2) Get All Products\n3) GetRecommendations\n")
		var command int
		fmt.Scanf("%d", &command)

		switch command {
		case 1:
			var user, product string
			fmt.Scanf("%s", &user)
			fmt.Scanf("%s", &product)
			fmt.Println("USERPRODUCT: ", user, product)
			f.BlockProduct(user, product)
		case 2:
			fmt.Println(f.Products)
		case 3:
			var user string
			fmt.Scanf("%s", &user)
			fmt.Println(f.GetRecommendation(user))
		}
	}
}