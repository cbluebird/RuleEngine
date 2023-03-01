package controller

import (
	"engine/app/compiler"
	"engine/app/model"
	"engine/app/utility"
	"engine/config/database"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func GetScore(c *gin.Context) {
	age_ := c.Query("age")
	age, _ := strconv.Atoi(age_)
	userid_ := c.Query("userid")
	userid, _ := strconv.Atoi(userid_)
	ruleId_ := c.Query("rule")
	ruleid, _ := strconv.Atoi(ruleId_)
	var list []model.Detail
	result := database.DB.Where(
		&model.Detail{
			Grade:  age,
			Userid: userid,
		},
	).Find(&list)
	if result.Error != nil {
		utility.JsonResponseInternalServerError(c)
		return
	}
	var rule model.Rule
	result = database.DB.Where(
		&model.Rule{ID: ruleid},
	).First(&rule)
	log.Println(rule)
	scanner := compiler.NewScanner(rule.Message)
	tokens, _ := scanner.Lexer()
	parser := compiler.NewParser(tokens)
	parser.Print()
	err := parser.CheckBalance()
	if err != nil {
		log.Println(err)
		return
	}
	err = parser.ParseSyntax()
	if err != nil {
		log.Println(err)
		return
	}
	bulider := compiler.NewBuilder(parser)
	node, err_ := bulider.Build()
	if err_ != nil {
		log.Println(err)
		return
	}
	m := make(map[int]float64, 6)
	for _, v := range list {
		m[v.Module] += v.Score
	}
	value := make(map[string]interface{})
	for i := 1; i <= 6; i++ {
		value[model.GetModule[i]] = m[i]
	}
	node.Eval(value)
	log.Println(node.GetVal())
	val, _ := node.GetVal()
	log.Println(val)
	utility.JsonResponse(200, "ok", val, c)
}
