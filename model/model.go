package model

/**
* 
* @author willian
* @created 2017-04-25 20:29
* @email 18702515157@163.com  
**/

type UserVisit struct {
	Date               int64
	User_id            int
	Session_id         string
	Page_id            int
	Action_time        int64
	Search_keyword     string
	Click_category_id  int
	Click_product_id   int
	Order_category_ids string
	Order_product_ids  string
	Pay_category_ids   string
	Pay_product_ids    string
}