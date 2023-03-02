package controller

import (
	"engine/app/engine"
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
	m := make(map[int]float64, 6)
	for _, v := range list {
		m[v.Module] += v.Score
	}
	value := make(map[string]interface{})
	for i := 1; i <= 6; i++ {
		value[model.GetModule[i]] = m[i]
	}
	engine := engine.NewEngine(rule.Message)
	err := engine.Calculate(value)
	if err != nil {
		log.Println(err)
		return
	}
	val, _ := engine.GetVal()
	utility.JsonResponse(200, "ok", val, c)
}
