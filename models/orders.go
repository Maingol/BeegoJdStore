package models

import (
	"github.com/astaxie/beego/orm"
	"strings"
)

// 数据库中sp_order表的模型
type SpOrder struct {
	OrderId            int     `orm:"pk;auto" json:"order_id"`
	UserId             int     `json:"user_id"`
	OrderNumber        string  `json:"order_number"`
	OrderPrice         float64 `json:"order_price"`
	OrderPay           string  `json:"order_pay"`
	IsSend             string  `json:"is_send"`
	TradeNo            string  `json:"trade_no"`
	OrderFapiaoTitle   string  `json:"order_fapiao_title"`
	OrderFapiaoContent string  `json:"order_fapiao_content"`
	ConsigneeAddr      string  `json:"consignee_addr"`
	PayStatus          string  `json:"pay_status"`
	CreateTime         int     `json:"create_time"`
	UpdateTime         int     `json:"update_time"`
}

// 获取订单列表时返回的数据 【接口：orders 请求方式：get】
type ResOrdersList struct {
	Data *ResOrdersData `json:"data"`
	Meta *ResMeta       `json:"meta"`
}

// 获取订单列表时返回结果中Data字段的数据 【接口：orders 请求方式：get】
type ResOrdersData struct {
	Total   int        `json:"total"`
	Pagenum int        `json:"pagenum"`
	Orders  []*SpOrder `json:"orders"`
}

// 修改订单地址时返回结果中Data字段的数据 【接口：orders/:order_id 请求方式：put】
type ResAddrData struct {
	OrderId       int    `json:"order_id"`
	ConsigneeAddr string `json:"consignee_addr"`
}

// 修改订单地址时返回结果的数据 【接口：orders/:order_id 请求方式：put】
type ResUpdateAddr struct {
	Data *ResAddrData `json:"data"`
	Meta *ResMeta     `json:"meta"`
}

// 根据订单id验证订单是否存在 【接口：orders/:order_id 请求方式：put】
func OrderExists(id int) bool {
	o := orm.NewOrm()
	return o.QueryTable("sp_order").Filter("order_id", id).Exist()
}

// 修改订单地址 【接口：orders/:order_id 请求方式：put】
func UpdateAddr(order *SpOrder, id int) (*ResAddrData, error) {
	order.OrderId = id
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return nil, err
	}
	_, err = o.Update(order, "consignee_addr")
	if err != nil {
		o.Rollback()
		return nil, err
	}
	o.Commit()
	resAddrData := &ResAddrData{order.OrderId, order.ConsigneeAddr}
	return resAddrData, nil
}

// 获取订单数据列表 【接口：orders 请求方式：get】
func GetOrdersList(query string, pagenum, pagesize int) (int, []*SpOrder, error) {
	o := orm.NewOrm()
	orderList := make([]*SpOrder, 0)
	offset := (pagenum - 1) * pagesize
	qs := o.QueryTable("sp_order")
	if strings.Trim(query, " ") != "" {
		// 按条件模糊查询
		qs = qs.Filter("order_number__icontains", query)
	}
	total, err := qs.Count()
	if err != nil {
		return 0, nil, err
	}
	// OrderBy的参数前使用减号“-”意味着倒叙排列
	_, err = qs.OrderBy("-order_id").Limit(pagesize, offset).All(&orderList)
	if err != nil {
		return 0, nil, err
	}
	return int(total), orderList, nil
}

// 初始化模型
func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(SpOrder))
}
