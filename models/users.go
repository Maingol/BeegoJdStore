package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

/* 获取get请求中的参数（users接口） */
type UsersParams struct {
	//这里不对query参数进行校验
	Pagenum  int `valid:"Required;Min(1)"`
	Pagesize int `valid:"Required;Min(1)"`
}

/*定义login接口返回响应数据中的data中的users结构（users接口）*/
type ResUser struct {
	MgId     int    `json:"id"`
	RoleName string `json:"role_name"`
	MgName   string `json:"username"`
	MgTime   int    `json:"create_time"`
	MgMobile string `json:"mobile"`
	MgEmail  string `json:"email"`
	MgState  bool   `json:"mg_state"`
}

/*定义login接口返回响应数据中的data结构（users接口）*/
type ResUsersData struct {
	Total   interface{} `json:"total"`
	PageNum interface{} `json:"pagenum"`
	Users   interface{} `json:"users"`
}

/*定义login接口返回响应数据中的meta结构（users接口）*/
type ResUsersMeta struct {
	Msg    interface{} `json:"msg"`
	Status interface{} `json:"status"`
}

/*定义login接口返回响应数据的结构*/
type ResUsers struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta"`
}

/* 定义数据库中的sp_role表的模型 */
type SpRole struct {
	RoleId   int    `orm:"pk;auto;size(6)" description:"主键id"`
	RoleName string `orm:"size(20)" description:"角色名称"`
	PsIds    string `orm:"size(512)" description:"权限"`
	PsCa     string `orm:"null" description:"控制器-操作"`
	RoleDesc string `orm:"null" description:"角色描述"`
}

/* 修改用户状态时返回的数据中data格式 */
type ResStateData struct {
	MgId     interface{} `json:"id"`
	RoleId   interface{} `json:"rid"`
	MgName   interface{} `json:"username"`
	MgMobile interface{} `json:"mobile"`
	MgEmail  interface{} `json:"email"`
	MgState  interface{} `json:"mg_state"`
}

/* 修改用户状态时返回的数据格式 */
type ResUserState struct {
	Data *ResStateData `json:"data"`
	Meta *ResUsersMeta `json:"meta"`
}

/* 添加用户时获取参数的数据结构（uers接口post） */
type CreateUser struct {
	UserName   string `valid:"Required;MinSize(3);MaxSize(10)"`
	PassWord   string `valid:"Required;MinSize(6);MaxSize(15)"`
	ConfirmPwd string `valid:"Required;"`
	Email      string `valid:"Required;Email"`
	Mobile     string `valid:"Required;Mobile"`
}

/* 定义返回数据data中的结构体（user接口，post请求） */
type ResAddUserData struct {
	Id         interface{} `json:"id"`
	UserName   interface{} `json:"username"`
	Mobile     interface{} `json:"mobile"`
	Email      interface{} `json:"email"`
	RoleId     interface{} `json:"role_id"`
	CreateTime interface{} `json:"create_time"`
}

/* 定义返回数据中的结构体（user接口，post请求） */
type ResAddUser struct {
	Data *ResAddUserData `json:"data"`
	Meta *ResUsersMeta   `json:"meta"`
}

/* 定义返回数据中data的结构体（user/:id接口，get请求） */
/* 根据Id查询用户信息时返回的数据中data格式 */
type ResUserData struct {
	Id       interface{} `json:"id"`
	Rid      interface{} `json:"rid"`
	UserName interface{} `json:"username"`
	Mobile   interface{} `json:"mobile"`
	Email    interface{} `json:"email"`
}

/* 定义返回数据中的结构体（user/:id接口，get请求） */
type ResGetUser struct {
	Data *ResUserData  `json:"data"`
	Meta *ResUsersMeta `json:"meta"`
}

/* 定义获取请求参数的结构体（users/:id接口，put请求） */
type GetPutParams struct {
	Email  string `valid:"Required;Email"`
	Mobile string `valid:"Required;Mobile"`
}

/* post参数的结构体 接口：users/:id/role，请求：put */
type RoleId struct {
	Rid int `valid:"Required"`
}

func (this *RoleId) Valid(v *validation.Validation) {
	role := SpRole{RoleId: this.Rid}
	o := orm.NewOrm()
	err := o.Read(&role)
	if err != nil {
		v.SetError("Rid", "角色id不存在")
	}
}

// 如果你的 struct 实现了接口 validation.ValidFormer
// 当 StructTag 中的测试都成功时，将会执行 Valid 函数进行自定义验证
func (this *CreateUser) Valid(v *validation.Validation) {
	if this.ConfirmPwd != this.PassWord {
		// 通过 SetError 设置 Name 的错误信息，HasErrors 将会返回 true
		v.SetError("ConfirmPwd", "两次密码输入不一致")
	}
}

/* sql查询管理员列表 */
func GetManagers(query string, pageNum int, pageSize int) ([]ResUser, error) {
	/* 计算LIMIT后的查询起始位置start（从0开始）
	*  start = [ 当前页码(不含0)-1 ] × 每页显示的记录数 */
	start := (pageNum - 1) * pageSize
	sqlStr := fmt.Sprintf(`
		SELECT
			t1.mg_id,
			(CASE WHEN t2.role_name is null && t1.role_id = 0 
				THEN '超级管理员' ELSE t2.role_name END) AS role_name,
			t1.mg_name,
			t1.mg_time,
			t1.mg_mobile,
			t1.mg_email,
			t1.mg_state
		FROM
			sp_manager AS t1 LEFT JOIN sp_role AS t2
		ON
			t1.role_id = t2.role_id
		WHERE
			t1.mg_name LIKE '%%%s%%'
		LIMIT 
			%d,%d
		`, query, start, pageSize)
	var managers []ResUser
	/* 根据菜单等级选择最终执行的相应sql语句 */
	o := orm.NewOrm()
	raws, err := o.Raw(sqlStr).QueryRows(&managers)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	/* 查询到的记录数为0时，返回[]，而不是null */
	if raws == 0 {
		managers = []ResUser{}
	}
	return managers, nil
}

//初始化模型
func init() {
	orm.RegisterModel(new(SpRole))
}
