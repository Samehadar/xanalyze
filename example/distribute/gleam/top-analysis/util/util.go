package util

import (
	"fmt"
	"github.com/sniperkit/xanalyze/plugin/distribute/gleam/flow"
	"github.com/snluu/uuid"
	"github.com/zhangweilun/session/model"
	"math/rand"
	"strconv"
	"time"
)

func Mock(context *flow.FlowContext) (ret *flow.Dataset) {

	s1 := rand.NewSource(time.Now().Unix())

	r1 := rand.New(s1)

	userList := make([]*model.UserVisit, 1000)

	actions := []string{
		"search", "click", "order", "pay",
	}
	searchKeywords := []string{
		"Hot pot", "cake", "Chongqing spicy chicken", "Chongqing facet", "suckling feed", "new spicy fish pot", "International Trade Building", "shopping arcades", "Japanese cuisine", "spa",
	}
	var count int
	for i := 0; i < 100; i++ {
		userId := r1.Intn(100)
		var id uuid.UUID
		var sessionId string
		dups := make(map[string]bool)
		for i := 0; i < 1024; i++ {
			id = uuid.Rand()
			sessionId = id.Hex()
			if dups[sessionId] {
				fmt.Errorf("Duplicates after %d iterations", i+1)
			}
			dups[sessionId] = true
		}
		for j := 0; j < 10; j++ {
			pageId := r1.Intn(100)
			now := time.Now()
			year, month, day := now.Date()
			date := time.Date(year, month, day, 0, 0, 0, 0, time.Local).Unix()
			actionTime := now.Unix()
			searchKeyword := ""
			clickCategoryId := 0
			clickProductId := 0
			orderCategoryIds := ""
			orderProductIds := ""
			payCategoryIds := ""
			payProductIds := ""
			action := actions[r1.Intn(4)]
			switch action {
			case actions[0]:
				searchKeyword = searchKeywords[r1.Intn(10)]
			case actions[1]:
				clickCategoryId = r1.Intn(100)
				clickProductId = r1.Intn(100)
			case actions[2]:
				orderCategoryIds = strconv.Itoa(r1.Intn(100))
				orderProductIds = strconv.Itoa(r1.Intn(100))
			case actions[3]:
				payCategoryIds = strconv.Itoa(r1.Intn(100))
				payProductIds = strconv.Itoa(r1.Intn(100))
			}
			visit := &model.UserVisit{
				Date:               date,
				User_id:            userId,
				Session_id:         sessionId,
				Page_id:            pageId,
				Action_time:        actionTime,
				Click_category_id:  clickCategoryId,
				Click_product_id:   clickProductId,
				Order_category_ids: orderCategoryIds,
				Order_product_ids:  orderProductIds,
				Pay_category_ids:   payCategoryIds,
				Pay_product_ids:    payProductIds,
				Search_keyword:     searchKeyword,
			}
			userList[count] = visit
			count++
		}
	}
	input := make(chan interface{})

	go func() {
		for _, data := range userList {
			input <- data
		}
		close(input)
	}()

	ret = context.Channel(input)
	return
}

//// Strings begins a flow with an []string
//func (fc *flow.FlowContext) Strings(lines []string) (ret *Dataset) {
//	inputChannel := make(chan interface{})
//
//	go func() {
//		for _, data := range lines {
//			inputChannel <- []byte(data)
//		}
//		close(inputChannel)
//	}()
//
//	return fc.Channel(inputChannel)
//}
