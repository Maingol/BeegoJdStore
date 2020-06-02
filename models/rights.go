package models

import (
	"JDStore/utils"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"strconv"
	"strings"
)

/* 定义获取权限列表（list形式）时data中数据的结构 */
type ListRightsData struct {
	PsId      int    `json:"id"`
	PsName    string `json:"authName"`
	PsLevel   string `json:"level"`
	PsPid     int    `json:"pid"`
	PsApiPath string `json:"path"`
}

/* 定义获取权限列表（list形式）时数据的结构 */
type ResListRights struct {
	Data *[]ListRightsData `json:"data"`
	Meta *ResUsersMeta     `json:"meta"`
}

/* 定义获取权限列表（tree形式）时每条数据的原始结构 */
type TreeRightsData struct {
	PsId      int    `json:"id"`
	PsName    string `json:"authName"`
	PsApiPath string `json:"path"`
	PsPid     int    `json:"pid"`
}

/* 定义获取权限列表（tree形式）时每条数据包含子列表的结构 */
type ResRight struct {
	PsId      int         `json:"id"`
	PsName    string      `json:"authName"`
	PsApiPath string      `json:"path"`
	PsPid     int         `json:"pid"`
	Children  []*ResRight `json:"children"`
}

/* 定义获取权限列表（tree形式）时数据的结构 */
type ResTreeRights struct {
	Data []*ResRight   `json:"data"`
	Meta *ResUsersMeta `json:"meta"`
}

/* 定义获取角色列表时每个角色的结构 */
type Role struct {
	RoleId   int
	RoleName string
	RoleDesc string
}

/* 定义获取角色列表时每个角色返回数据的结构 */
type ResRole struct {
	RoleId   int         `json:"id"`
	RoleName string      `json:"roleName"`
	RoleDesc string      `json:"roleDesc"`
	Children interface{} `json:"children"`
}

/* 定义获取角色列表时每个权限的结构 */
type Right struct {
	PsId      int
	PsName    string
	PsApiPath string
}

/* 定义获取角色列表时每个权限返回数据的结构 */
type ResRoleRight struct {
	RoleId    int         `json:"id"`
	RoleName  string      `json:"authName"`
	PsApiPath string      `json:"path"`
	Children  interface{} `json:"children"`
}

/* 定义获取角色列表时返回数据的结构 */
type ResRoles struct {
	Data []ResRole     `json:"data"`
	Meta *ResUsersMeta `json:"meta"`
}

/* 删除角色指定权限时返回数据的结构（接口：roles/:roleId/rights/:rightId，请求：delete） */
type ResDelRight struct {
	Data []ResRoleRight `json:"data"`
	Meta *ResUsersMeta  `json:"meta"`
}

/* 定义body中请求参数的结构 接口：roles/:roleId/rights 请求方式：post */
type PostRights struct {
	Rids string
}

/* 编辑角色返回的数据中data的结构 接口：roles/:id 请求方式：put */
type ResRoleInfoData struct {
	RoleId   interface{} `json:"roleId"`
	RoleName interface{} `json:"roleName"`
	RoleDesc interface{} `json:"roleDesc"`
}

/* 编辑角色返回的数据的结构 接口：roles/:id 请求方式：put */
type ResRoleInfo struct {
	Data *ResRoleInfoData `json:"data"`
	Meta *ResMeta         `json:"meta"`
}

/* put请求中参数的结构 接口：roles/:id 请求方式：put  */
type GetRoleParams struct {
	RoleName string
	RoleDesc string
}

func (this *PostRights) Valid(v *validation.Validation) {
	if this.Rids == "" {
		return
	}
	RightsSlice := strings.Split(this.Rids, ",")
	for _, right := range RightsSlice {
		/* 将切片中每个元素转换为数字 */
		rightNum, err := strconv.Atoi(right)
		if err != nil {
			v.SetError("Rids", "权限格式错误")
			return
		}
		/* 验证权限列表中是否有权限id不存在于数据库表 */
		if rightNum == 0 {
			/* 不验证根权限 */
			continue
		}
		isNotInTable := true
		var spPermission []SpPermission
		o := orm.NewOrm()
		o.QueryTable("sp_permission").All(&spPermission)
		for _, per := range spPermission {
			if per.PsId == rightNum {
				isNotInTable = false
				break
			}
		}
		if isNotInTable {
			v.SetError("Rids", "包含不存在的权限")
		}
	}
}

/* 查询权限列表（list形式） */
func QueryRightsList() ([]ListRightsData, error) {
	sqlStr := `
		SELECT 
			t1.ps_id, t1.ps_name, t1.ps_level, t1.ps_pid, t2.ps_api_path
		FROM
			sp_permission t1 JOIN sp_permission_api t2
		ON
			t1.ps_id = t2.ps_id`
	var rightsList []ListRightsData
	o := orm.NewOrm()
	_, err := o.Raw(sqlStr).QueryRows(&rightsList)
	if err != nil {
		return nil, err
	}
	return rightsList, nil
}

/* 查询权限列表（tree形式） */
func RightsTree() ([]*ResRight, error) {
	return GetTree(0)
}

/* 递归获取权限列表 */
func GetTree(id int) ([]*ResRight, error) {
	sqlStr := `
		SELECT 
			t1.ps_id, t1.ps_name, t1.ps_level, t1.ps_pid, t2.ps_api_path
		FROM
			sp_permission t1 JOIN sp_permission_api t2
		ON
			t1.ps_id = t2.ps_id
		WHERE
			t1.ps_pid = ?`
	/* 存储从数据库中查询到的结果集 */
	var treeRightsData []TreeRightsData
	o := orm.NewOrm()
	/* 从数据库中查询出所有pid为id的记录，例如id为0，查询结果就是所有一级权限 */
	num, err := o.Raw(sqlStr, id).QueryRows(&treeRightsData)
	if err != nil {
		return nil, err
	}
	/* 存储包含children成员的权限列表 */
	var resRights []*ResRight
	/* 递归终止条件 */
	if num == 0 {
		/* 返回一个空切片 */
		return []*ResRight{}, nil
	}
	/* 遍历查询结果，结果中的每个字段都赋值给新的结构体，并且都添加一个children字段赋给这个新结构体 */
	for _, v := range treeRightsData {
		children, err := GetTree(v.PsId)
		if err != nil {
			return nil, err
		}
		resRight := ResRight{PsId: v.PsId, PsName: v.PsName, PsApiPath: v.PsApiPath, PsPid: v.PsPid}
		resRight.Children = children
		resRights = append(resRights, &resRight)
	}
	return resRights, nil
}

/* 查询角色列表 */
func QueryRoleList() ([]ResRole, error) {
	var roles []SpRole
	o := orm.NewOrm()
	/* 查询出所有角色信息 */
	_, err := o.QueryTable("sp_role").All(&roles)
	if err != nil {
		return nil, err
	}

	var resRoles []ResRole
	for _, v := range roles {
		resRole := ResRole{RoleId: v.RoleId, RoleName: v.RoleName, RoleDesc: v.RoleDesc}
		rightsSlice := strings.Split(v.PsIds, ",")
		var children []ResRoleRight
		/* 校验根权限0。这里调用了自定义的函数 */
		if !utils.Contains(rightsSlice, "0") || v.PsIds == "" {
			children = []ResRoleRight{}
		} else {
			children, err = GetRoleRights(0, rightsSlice)
			if err != nil {
				return nil, err
			}
		}
		/* 当角色没有任何权限时，children字段不返回null，而是空切片 */
		if children == nil {
			children = []ResRoleRight{}
		}
		resRole.Children = children
		resRoles = append(resRoles, resRole)
	}
	return resRoles, nil
}

/* 递归获取每个角色的权限列表 */
func GetRoleRights(id int, rightsSlice []string) ([]ResRoleRight, error) {
	/* 查询出rid为id的那些权限。例如id为0时，返回的是所有一级权限以及所有一级权限的子权限 */
	/* 数据库中的order字段可能是渲染左侧菜单列表的时候才有用，
	*  这里渲染权限列表时应该是不必按order字段排序的，所以这里不必使用order by子句 */
	sqlStr := `
		SELECT 
			t1.ps_id, t1.ps_name, t2.ps_api_path
		FROM
			sp_permission t1 JOIN sp_permission_api t2
		ON
			t1.ps_id = t2.ps_id
		WHERE
			t1.ps_pid = ?`
	/* 存储从数据库中查询到的结果集 */
	var rights []Right
	o := orm.NewOrm()
	num, err := o.Raw(sqlStr, id).QueryRows(&rights)
	if err != nil {
		return nil, err
	}
	/* 存储包含children成员的权限列表 */
	var resRoleRights []ResRoleRight
	/* 递归终止条件 */
	if num == 0 {
		/* 返回一个空切片 */
		return []ResRoleRight{}, nil
	}

	for _, v := range rights {
		/* 如果查询到的权限id不包含在角色的权限内，则中断本次循环 */
		if !utils.Contains(rightsSlice, strconv.FormatInt(int64(v.PsId), 10)) {
			continue
		}
		resRight := ResRoleRight{RoleId: v.PsId, RoleName: v.PsName, PsApiPath: v.PsApiPath}
		children, err := GetRoleRights(v.PsId, rightsSlice)
		if err != nil {
			return nil, err
		}
		/* 当对象没有任何子权限时，children字段不返回null，而是空切片 */
		if children == nil {
			children = []ResRoleRight{}
		}
		resRight.Children = children
		resRoleRights = append(resRoleRights, resRight)
	}
	return resRoleRights, nil
}

/* 对用户访问资源进行权限验证 */
func ValidateRight(ctx *context.Context, rid int) bool {
	// 获取当前登陆的用户名
	token := ctx.Input.Header("Authorization")
	token = strings.Split(token, " ")[1]
	user := utils.CheckToken(token)
	// 判断用户是否是超级管理员
	manager := SpManager{MgName: user}
	o := orm.NewOrm()
	o.Read(&manager, "MgName")
	if manager.RoleId == 0 {
		// 当前用户是超级管理员，直接返回true
		return true
	}
	// 获取当前用户的权限列表
	role := SpRole{RoleId: int(manager.RoleId)}
	o.Read(&role, "RoleId")
	userRights := strings.Split(role.PsIds, ",")

	// 判断接口对应的权限id是否存在于用户的权限列表
	ridStr := strconv.FormatInt(int64(rid), 10)
	if utils.Contains(userRights, ridStr) {
		return true
	} else {
		return false
	}
}
