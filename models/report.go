package models

import (
	"github.com/astaxie/beego/orm"
)

// 数据库中sp_report_1表的模型
type SpReport_1 struct {
	Id           int `orm:"pk;auto"`
	Rp1UserCount int
	Rp1Area      string
	Rp1Date      string
}

type ResChildData struct {
	Data []interface{} `json:"data"`
}

type ResSeries struct {
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Stack     string                 `json:"stack"`
	AreaStyle map[string]interface{} `json:"areaStyle"`
	Data      interface{}            `json:"data"`
}

type ResReportData struct {
	Legend ResChildData             `json:"legend"`
	YAxis  []map[string]string      `json:"yAxis"`
	XAxis  []map[string]interface{} `json:"xAxis"`
	Series []*ResSeries             `json:"series"`
}

type ResReport struct {
	Data *ResReportData `json:"data"`
	Meta *ResMeta       `json:"meta"`
}

// 获取数据报表 【接口：reports/type/1  请求方式：get】
func GetReport() (*ResReportData, error) {
	report := make([]*SpReport_1, 0)
	o := orm.NewOrm()
	_, err := o.QueryTable("sp_report_1").All(&report)
	if err != nil {
		return nil, err
	}

	resReportData := new(ResReportData)
	arr1 := make([]interface{}, 0)
	arr2 := make([]interface{}, 0)
	// legend、xAxis数据去重
	for _, v := range report {
		arr1 = append(arr1, v.Rp1Area)
		arr2 = append(arr2, v.Rp1Date)
	}
	arr1 = RemoveRepByMap(arr1)
	arr2 = RemoveRepByMap(arr2)

	resReportData.Legend.Data = arr1

	resReportData.YAxis = append(resReportData.YAxis, map[string]string{"type": "value"})
	resReportData.XAxis = append(resReportData.XAxis, map[string]interface{}{"data": arr2})

	for _, area := range resReportData.Legend.Data {
		resSeries := new(ResSeries)
		resSeries.Name = area.(string)
		resSeries.Type = "line"
		resSeries.Stack = "总量"
		resSeries.AreaStyle = map[string]interface{}{"normal": map[int]int{}}

		values := make([]int, 0)
		for _, v := range report {
			if v.Rp1Area == area.(string) {
				values = append(values, v.Rp1UserCount)
			}
		}
		resSeries.Data = values
		resReportData.Series = append(resReportData.Series, resSeries)
	}
	return resReportData, nil
}

// slice去重(保持原顺序)
func RemoveRepByMap(slc []interface{}) []interface{} {
	var result []interface{}          //存放返回的不重复切片
	tempMap := map[interface{}]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0 //当e存在于tempMap中时，再次添加是添加不进去的，，因为key不允许重复
		//如果上一行添加成功，那么长度发生变化且此时元素一定不重复
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e) //当元素不重复时，将元素添加到切片result中
		}
	}
	return result
}

// 初始化模型
func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(SpReport_1))
}
