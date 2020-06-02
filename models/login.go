package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type SpManager struct {
	/*pk以为当前字段为主键,auto意为自动增长,数据库表默认为 NOT NULL，设置 null 代表 ALLOW NULL*/
	MgId     int    `json:"mg_id" orm:"column(mg_id);pk;auto;size(11)" description:"主键id"`
	MgName   string `json:"mg_name" orm:"column(mg_name);size(32)" description:"管理员名称"`
	MgPwd    string `json:"mg_pwd" orm:"column(mg_pwd);size(64);type(char)" description:"管理员密码"`
	MgTime   int    `json:"mg_time" orm:"column(mg_time);size(10)" description:"注册时间"`
	RoleId   int8   `json:"role_id" orm:"column(role_id);size(11)" description:"角色id"`
	MgMobile string `json:"mg_mobile" orm:"column(mg_mobile);size(32);null" description:"管理员手机号"`
	MgEmail  string `json:"mg_email" orm:"column(mg_email);size(64);null" description:"管理员邮箱地址"`
	MgState  int8   `json:"mg_state" orm:"column(mg_state);size(2);null" description:"1:表示启用 0:表示禁用"`
}

/*定义login接口返回响应数据中的data结构*/
type ResData struct {
	Id       int    `json:"id"`
	Rid      int8   `json:"rid"`
	Username string `json:"username"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

/* 定义获取post参数的结构体（login接口） */
type GetParams struct {
	/* 这里使用StructTag进行表单验证。详细介绍参见：https://www.jianshu.com/p/37abab5808bb */
	UserName string `valid:"Required;MinSize(3);MaxSize(10)"`
	PassWord string `valid:"Required;MinSize(6);MaxSize(15)"`
}

/*定义响应数据中的meta结构*/
type ResMeta struct {
	Msg    string `json:"msg"`
	Status int    `json:"status"`
}

/*定义login接口返回响应数据的结构*/
type ResLogin struct {
	Data *ResData `json:"data"`
	Meta *ResMeta `json:"meta"`
}

/*定义数据库中sp_permission表的模型*/
type SpPermission struct {
	PsId    int    `json:"ps_id" orm:"column(ps_id);pk;auto;size(6)" description:"主键id"`
	PsName  string `json:"ps_name" orm:"column(ps_name);size(20)" description:"权限名称"`
	PsPid   int    `json:"ps_pid" orm:"column(ps_pid);size(6)" description:"父id"`
	PsC     string `json:"ps_c" orm:"column(ps_c);size(32)" description:"控制器"`
	PsA     string `json:"ps_a" orm:"column(ps_a);size(32)" description:"操作方法"`
	PsLevel string `json:"ps_level" orm:"column(ps_level)" description:"权限等级"`
}

/*定义数据库中sp_permission_api表的模型*/
type SpPermissionApi struct {
	Id           int    `json:"id" orm:"column(id);pk;auto;size(11)" description:"主键id"`
	PsId         int    `json:"ps_id" orm:"column(ps_id);size(11)"`
	PsApiService string `json:"ps_api_service" orm:"column(ps_api_service);size(255);null"`
	PsApiAction  string `json:"ps_api_action" orm:"column(ps_api_action);size(255);null"`
	PsApiPath    string `json:"ps_api_path" orm:"column(ps_api_path);size(255);null"`
	PsApiOrder   int    `json:"ps_api_order" orm:"column(ps_api_order);size(4);null"`
}

/*定义菜单接口menus的返回数据data中的结构*/
type Menus struct {
	PsId      int     `json:"id" orm:"column(ps_id)"`
	PsName    string  `json:"authName" orm:"column(ps_name)"`
	PsApiPath string  `json:"path" orm:"column(ps_api_path)"`
	Children  []Menus `json:"children"`
	/* order字段为空时，返回给前台null */
	PsApiOrder *int `json:"order" orm:"column(ps_api_order)"`
}

/*定义菜单接口menus的返回数据中的结构*/
type ResMenus struct {
	Data []Menus  `json:"data"`
	Meta *ResMeta `json:"meta"`
}

/* 查询出所有的一级菜单和二级菜单，query是查询参数 */
func TestOneToOne(level int, query int) ([]Menus, error) {
	o := orm.NewOrm()
	/* sql语句1：找出所有的一级菜单 */
	sqlStr1 := `
		SELECT
			t1.ps_id, t1.ps_name, t2.ps_api_path, t2.ps_api_order
		FROM
			sp_permission t1 JOIN sp_permission_api t2 
		ON 
			t1.ps_id = t2.ps_id 
		WHERE
			t1.ps_pid = ?
		ORDER BY
			t2.ps_api_order`

	/* sql语句2：根据一级菜单的ps_id找出相应的二级菜单（只需更改一下查询条件即可） */
	sqlStr2 := `
		SELECT
			t1.ps_id, t1.ps_name, t2.ps_api_path, t2.ps_api_order
		FROM
			sp_permission t1 JOIN sp_permission_api t2 
		ON 
			t1.ps_id = t2.ps_id 
		WHERE
			t1.ps_pid = ? AND t1.ps_level = "1"
		ORDER BY
			t2.ps_api_order`

	sqlMap := map[int]string{1: sqlStr1, 2: sqlStr2}
	var menu []Menus
	/* 根据菜单等级选择最终执行的相应sql语句 */
	_, err := o.Raw(sqlMap[level], query).QueryRows(&menu)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	/* children字段为空时，返回给前台空切片 */
	for i := range menu {
		if menu[i].Children == nil {
			menu[i].Children = []Menus{}
		}
	}
	return menu, nil
}

//初始化模型
func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(SpManager), new(SpPermission), new(SpPermissionApi))
}
