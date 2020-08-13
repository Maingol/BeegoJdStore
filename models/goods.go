package models

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/disintegration/imaging"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

// 数据库中sp_goods_attr表的模型
type SpGoodsPics struct {
	PicsId  int `orm:"pk;auto"`
	GoodsId int
	PicsBig string
	PicsMid string
	PicsSma string
}

// 数据库中sp_goods_attr表的模型
type SpGoodsAttr struct {
	Id        int `orm:"pk;auto"`
	GoodsId   int
	AttrId    int
	AttrValue string
	AddPrice  float64
}

// 数据库中sp_category表的模型
type SpCategory struct {
	CatId      int `orm:"pk;auto"`
	CatName    string
	CatPid     int
	CatLevel   int
	CatDeleted int
	CatIcon    string
	CatSrc     string
}

// 数据库中sp_goods表的模型
type SpGoods struct {
	GoodsId        int     `orm:"pk;auto" json:"goods_id"`
	CatId          int     `json:"-"`
	GoodsName      string  `json:"goods_name"`
	GoodsPrice     float64 `json:"goods_price"`
	GoodsNumber    int     `json:"goods_number"`
	GoodsWeight    int     `json:"goods_weight"`
	GoodsIntroduce string  `json:"-"`
	GoodsBigLogo   string  `json:"-"`
	GoodsSmallLogo string  `json:"-"`
	IsDel          string  `json:"-"`
	GoodsState     int     `json:"goods_state"`
	AddTime        int     `json:"add_time"`
	UpdTime        int     `json:"upd_time"`
	DeleteTime     int     `json:"-"`
	CatOneId       int     `json:"-"`
	CatTwoId       int     `json:"-"`
	CatThreeId     int     `json:"-"`
	HotNumber      int     `json:"hot_number"`
	IsPromote      bool    `json:"is_promote"`
}

type PicBody struct {
	Pic string
}

type AttrBody struct {
	Attr_id    int
	Attr_value string
	Attr_name  string
	Attr_sel   string
	Attr_write string
	Attr_vals  string
}

type AddGoodBody struct {
	Goods_name      string
	Goods_cate      string
	Goods_price     float64
	Goods_number    int
	Goods_weight    int
	Goods_introduce string
	Pics            []*PicBody
	Attrs           []*AttrBody
}
type ResAddGoodAttrData struct {
	GoodsId   int     `json:"goods_id"`
	AttrId    int     `json:"attr_id"`
	AttrValue string  `json:"attr_value"`
	AddPrice  float64 `json:"add_price"`
	AttrName  string  `json:"attr_name"`
	AttrSel   string  `json:"attr_sel"`
	AttrWrite string  `json:"attr_write"`
	AttrVals  string  `json:"attr_vals"`
}

type ResAddGoodData struct {
	GoodsId        int                   `json:"goods_id"`
	GoodsName      string                `json:"goods_name"`
	GoodsPrice     float64               `json:"goods_price"`
	CatId          int                   `json:"-"`
	GoodsNumber    int                   `json:"goods_number"`
	GoodsWeight    int                   `json:"goods_weight"`
	GoodsIntroduce string                `json:"-"`
	GoodsBigLogo   string                `json:"-"`
	GoodsSmallLogo string                `json:"-"`
	GoodsState     int                   `json:"goods_state"`
	AddTime        int                   `json:"add_time"`
	UpdTime        int                   `json:"upd_time"`
	HotNumber      int                   `json:"hot_number"`
	IsPromote      bool                  `json:"is_promote"`
	Pics           []*SpGoodsPics        `json:"pics"`
	Attrs          []*ResAddGoodAttrData `json:"attrs"`
}

type ResAddGood struct {
	Data *ResAddGoodData `json:"data"`
	Meta *ResMeta        `json:"meta"`
}

// 获取商品列表时返回结果中Data字段的数据 接口：goods 请求方式：get
type ResGoodsData struct {
	Total   int        `json:"total"`
	Pagenum int        `json:"pagenum"`
	Goods   []*SpGoods `json:"goods"`
}

// 获取商品列表时返回的数据 接口：goods 请求方式：get
type ResGoodsList struct {
	Data *ResGoodsData `json:"data"`
	Meta *ResMeta      `json:"meta"`
}

// 定义商品分类中每个分类的结构
type GoodsCate struct {
	CatId      int    `json:"cat_id"`
	CatName    string `json:"cat_name"`
	CatPid     int    `json:"cat_pid"`
	CatLevel   int    `json:"cat_level"`
	CatDeleted bool   `json:"cat_deleted"`
}

// 定义商品分类返回结果中每个分类的结构
type ResGoodsCate struct {
	CatId      int    `json:"cat_id"`
	CatName    string `json:"cat_name"`
	CatPid     int    `json:"cat_pid"`
	CatLevel   int    `json:"cat_level"`
	CatDeleted bool   `json:"cat_deleted"`
	/* omitempty标签的作用：如果该字段为空，则忽略
	   -标签的作用：无论该字段是否为空，都忽略 */
	Children []*ResGoodsCate `json:"children,omitempty"`
}

// 获取商品分类列表时，包含分页参数的返回结果中data的结构
type ResCatePageData struct {
	Total    int             `json:"total"`
	PageNum  int             `json:"pagenum"`
	PageSize int             `json:"pagesize"`
	Result   []*ResGoodsCate `json:"result"`
}

// 获取商品分类列表时，包含分页参数的返回结果
type ResCatePage struct {
	Data *ResCatePageData `json:"data"`
	Meta *ResMeta         `json:"meta"`
}

// 获取商品分类列表时，不包含分页参数的返回结果
type ResCateList struct {
	Data []*ResGoodsCate `json:"data"`
	Meta *ResMeta        `json:"meta"`
}

// post请求中请求参数的结构体 接口：categories
type AddCateParams struct {
	Cat_pid   int
	Cat_name  string `valid:"Required"`
	Cat_level int    `valid:"Range(0, 2)"`
}

// 添加分类接口返回的数据的结构
type ResAddCate struct {
	Data *GoodsCate `json:"data"`
	Meta *ResMeta   `json:"meta"`
}

// put请求中请求参数的结构体 接口：categories/:id
type UpdateCateParams struct {
	Cat_name string
}

// 数据库中sp_attribute表的模型
type SpAttribute struct {
	AttrId     int    `json:"attr_id" orm:"pk;auto"`
	AttrName   string `json:"attr_name"`
	CatId      int    `json:"cat_id"`
	AttrSel    string `json:"attr_sel"`
	AttrWrite  string `json:"attr_write"`
	AttrVals   string `json:"attr_vals"`
	DeleteTime *int   `json:"-"`
}

// 获取分类参数列表时返回的数据 接口：categories/:id/attributes 请求方式：get
type ResAttrList struct {
	Data []*SpAttribute `json:"data"`
	Meta *ResMeta       `json:"meta"`
}

// post请求中请求参数的结构体 接口：categories/:id/attributes
type AddAttrParams struct {
	Attr_name string
	Attr_sel  string
	Attr_vals string
}

// 验证请求体中的参数是否包含Attr_vals字段 接口：categories/:id/attributes/:attrId 请求方式：put
// 如果包含这个字段，且该字段值为空字符串"",则会报错——类型不匹配；如果不包含这个字段，解析出来后，该字段的值默认为0
type ValidParamExist struct {
	Attr_vals int
}

// 添加分类参数时返回的数据 接口：categories/:id/attributes 请求方式：post
type ResAttr struct {
	Data *SpAttribute `json:"data"`
	Meta *ResMeta     `json:"meta"`
}

// 上传图片接口返回数据中Data字段的结构 接口：upload 请求方式：post
type ResUploadData struct {
	Tmp_path string `json:"tmp_path"`
	Url      string `json:"url"`
}

// 上传图片接口返回的数据 接口：upload 请求方式：post
type ResUpload struct {
	Data *ResUploadData `json:"data"`
	Meta *ResMeta       `json:"meta"`
}

// 根据图片路径，将该图在相同目录下生成大、中、小三张图片，并返回这三张图片的路径
func ResizeImg(oriPicPath string) ([]string, error) {
	//读取本地文件
	imgData, err := ioutil.ReadFile("." + oriPicPath)
	if err != nil {
		logs.Error("图片读取错误：", err)
		return nil, err
	}
	buf := bytes.NewBuffer(imgData)
	image, err := imaging.Decode(buf)
	if err != nil {
		logs.Error("图片解码错误", err)
		return nil, err
	}
	// 生成缩略图，尺寸800*800，并保存文件
	largeImage := imaging.Resize(image, 800, 800, imaging.Lanczos)
	largePicPath := strings.Replace(oriPicPath, ".", "_800×800.", -1)
	err = imaging.Save(largeImage, "."+largePicPath)
	if err != nil {
		logs.Error("图片(800×800)保存错误", err)
		return nil, err
	}

	// 生成缩略图，尺寸400*400，并保存文件
	midImage := imaging.Resize(image, 400, 400, imaging.Lanczos)
	midPicPath := strings.Replace(oriPicPath, ".", "_400×400.", -1)
	err = imaging.Save(midImage, "."+midPicPath)
	if err != nil {
		logs.Error("图片(400×400)保存错误", err)
		return nil, err
	}

	// 生成缩略图，尺寸200*200，并保存文件
	smallImage := imaging.Resize(image, 200, 200, imaging.Lanczos)
	smallPicPath := strings.Replace(oriPicPath, ".", "_200×200.", -1)
	err = imaging.Save(smallImage, "."+smallPicPath)
	if err != nil {
		logs.Error("图片(200×200)保存错误", err)
		return nil, err
	}

	return []string{largePicPath, midPicPath, smallPicPath}, nil
}

// 添加商品
func AddGood(addGoodBody AddGoodBody) (*ResAddGoodData, error) {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return nil, err
	}

	// 获取分类id
	goodsCateArr := strings.Split(addGoodBody.Goods_cate, ",")
	if len(goodsCateArr) != 3 {
		return nil, errors.New("分类参数错误")
	}
	var catIds []int
	for _, v := range goodsCateArr {
		catId, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		catIds = append(catIds, catId)
	}

	// 当前时间的时间戳
	current := int(time.Now().Unix())

	good := &SpGoods{
		GoodsName:      addGoodBody.Goods_name,
		GoodsPrice:     addGoodBody.Goods_price,
		GoodsNumber:    addGoodBody.Goods_number,
		GoodsWeight:    addGoodBody.Goods_weight,
		CatId:          catIds[2],
		GoodsIntroduce: addGoodBody.Goods_introduce,
		IsDel:          "0",
		AddTime:        current,
		UpdTime:        current,
		CatOneId:       catIds[0],
		CatTwoId:       catIds[1],
		CatThreeId:     catIds[2],
	}

	_, err = o.Insert(good)
	if err != nil {
		o.Rollback()
		return nil, err
	}

	var resPics []*SpGoodsPics
	for _, v := range addGoodBody.Pics {
		resizePathArr, err := ResizeImg(v.Pic)
		if err != nil {
			o.Rollback()
			return nil, err
		}
		baseUrl := "http://127.0.0.1:8700"
		goodPic := &SpGoodsPics{
			GoodsId: good.GoodsId,
			PicsBig: baseUrl + resizePathArr[0],
			PicsMid: baseUrl + resizePathArr[1],
			PicsSma: baseUrl + resizePathArr[2],
		}
		_, err = o.Insert(goodPic)
		if err != nil {
			o.Rollback()
			return nil, err
		}
		resPics = append(resPics, goodPic)
	}

	var resAttrs []*ResAddGoodAttrData
	for _, v := range addGoodBody.Attrs {
		goodAttr := &SpGoodsAttr{
			GoodsId:   good.GoodsId,
			AttrId:    v.Attr_id,
			AttrValue: v.Attr_value,
		}
		_, err = o.Insert(goodAttr)
		if err != nil {
			o.Rollback()
			return nil, err
		}
		resGoodAttr := &ResAddGoodAttrData{
			GoodsId:   goodAttr.GoodsId,
			AttrId:    goodAttr.AttrId,
			AttrValue: goodAttr.AttrValue,
			AddPrice:  goodAttr.AddPrice,
			AttrName:  v.Attr_name,
			AttrSel:   v.Attr_sel,
			AttrWrite: v.Attr_write,
			AttrVals:  v.Attr_vals,
		}
		resAttrs = append(resAttrs, resGoodAttr)
	}

	resAddGoodData := &ResAddGoodData{
		GoodsId:        good.GoodsId,
		GoodsName:      good.GoodsName,
		GoodsPrice:     good.GoodsPrice,
		CatId:          good.CatId,
		GoodsNumber:    good.GoodsNumber,
		GoodsWeight:    good.GoodsWeight,
		GoodsIntroduce: good.GoodsIntroduce,
		GoodsBigLogo:   good.GoodsBigLogo,
		GoodsSmallLogo: good.GoodsSmallLogo,
		GoodsState:     good.GoodsState,
		AddTime:        good.AddTime,
		UpdTime:        good.UpdTime,
		HotNumber:      good.HotNumber,
		IsPromote:      good.IsPromote,
		Pics:           resPics,
		Attrs:          resAttrs,
	}
	o.Commit()
	return resAddGoodData, nil
}

// 请求参数CatPid的自定义校验规则
func (this *AddCateParams) Valid(v *validation.Validation) {
	var category SpCategory
	o := orm.NewOrm()
	category = SpCategory{CatId: this.Cat_pid}
	if this.Cat_pid != 0 {
		err := o.Read(&category, "CatId")
		if err != nil {
			v.SetError("CatPid", "分类父ID不存在")
		}
	}
	category = SpCategory{CatName: this.Cat_name}
	err := o.Read(&category, "CatName")
	if err == nil {
		v.SetError("CatName", "分类名称已存在")
	}
}

// 获取商品分类列表，不包含分页参数
func GetCateList(level int) ([]*ResGoodsCate, error) {
	return GetCate(0, level)
}

// 获取商品分类列表，包含分页参数
func GetCatePage(level int, pageNum int, pageSize int) (*ResCatePageData, error) {
	o := orm.NewOrm()
	// 过滤掉被假删的分类
	//total, err := o.QueryTable("sp_category").
	//	Filter("CatPid", 0).Exclude("CatDeleted", 1).Count()

	// 保留被假删的分类
	total, err := o.QueryTable("sp_category").
		Filter("CatPid", 0).Count()
	if err != nil {
		return nil, err
	}
	resCatePageData := ResCatePageData{Total: int(total), PageNum: pageNum, PageSize: pageSize}

	start := (pageNum - 1) * pageSize
	sqlStr := fmt.Sprintf(`
		SELECT
			cat_id, cat_name, cat_pid, cat_level, cat_deleted
		FROM
			sp_category
		WHERE
			cat_pid = 0
		LIMIT
			%d,%d`, start, pageSize)
	/* 存储从数据库中查询到的结果集 */
	var goodsCate []GoodsCate
	/* 从数据库中查询出所有pid为id的记录，例如id为0，查询结果就是所有一级权限 */
	_, err = o.Raw(sqlStr).QueryRows(&goodsCate)
	if err != nil {
		return nil, err
	}
	/* 存储包含children成员的权限列表 */
	var resGoodsCats []*ResGoodsCate
	/* 遍历查询结果，结果中的每个字段都赋值给新的结构体，并且都添加一个children字段赋给这个新结构体 */
	for _, v := range goodsCate {
		children, err := GetCate(v.CatId, level)
		if err != nil {
			return nil, err
		}
		resGoodsCate := ResGoodsCate{
			CatId: v.CatId, CatName: v.CatName, CatPid: v.CatPid,
			CatLevel: v.CatLevel, CatDeleted: v.CatDeleted}
		resGoodsCate.Children = children
		resGoodsCats = append(resGoodsCats, &resGoodsCate)
	}
	resCatePageData.Result = resGoodsCats
	return &resCatePageData, nil
}

// 递归获取商品分类列表。第一个参数是父分类的id，第二个参数是递归截止到哪一级分类，第三个参数表示是否分页展示
func GetCate(id int, level int) ([]*ResGoodsCate, error) {
	query := map[int]string{1: " && cat_level = 0", 2: " && cat_level != 2", 3: "", 0: ""}

	sqlStr := fmt.Sprintf(`
		SELECT
			cat_id, cat_name, cat_pid, cat_level, cat_deleted
		FROM
			sp_category
		WHERE
			cat_pid = %d%s`, id, query[level])
	/* 存储从数据库中查询到的结果集 */
	var goodsCate []GoodsCate
	o := orm.NewOrm()
	/* 从数据库中查询出所有pid为id的记录，例如id为0，查询结果就是所有一级权限 */
	num, err := o.Raw(sqlStr).QueryRows(&goodsCate)
	if err != nil {
		return nil, err
	}
	/* 存储包含children成员的权限列表 */
	var resGoodsCats []*ResGoodsCate
	/* 递归终止条件 */
	if num == 0 {
		/* 返回一个空切片 */
		return []*ResGoodsCate{}, nil
	}
	/* 遍历查询结果，结果中的每个字段都赋值给新的结构体，并且都添加一个children字段赋给这个新结构体 */
	for _, v := range goodsCate {
		children, err := GetCate(v.CatId, level)
		if err != nil {
			return nil, err
		}
		resGoodsCate := ResGoodsCate{
			CatId: v.CatId, CatName: v.CatName, CatPid: v.CatPid,
			CatLevel: v.CatLevel, CatDeleted: v.CatDeleted}
		resGoodsCate.Children = children
		resGoodsCats = append(resGoodsCats, &resGoodsCate)
	}
	return resGoodsCats, nil
}

// 检验分类id是否存在
func CateIdExist(id int) bool {
	o := orm.NewOrm()
	return o.QueryTable("sp_category").Filter("cat_id", id).Exist()
}

// 验证某一分类下的分类名称是否存在
func AttrExist(id int, attrName, attrSel string) bool {
	o := orm.NewOrm()
	return o.QueryTable("sp_attribute").
		Filter("attr_name", attrName).
		Filter("cat_id", id).
		Filter("attr_sel", attrSel).
		Exist()
}

// 添加动态参数或者静态属性
func AddAttr(id int, attrName, attrSel, attrVals string) (*SpAttribute, error) {
	o := orm.NewOrm()
	mapAttrWrite := map[string]string{"only": "manual", "many": "list"}
	attr := &SpAttribute{
		AttrName:  attrName,
		CatId:     id,
		AttrSel:   attrSel,
		AttrWrite: mapAttrWrite[attrSel],
		AttrVals:  attrVals,
	}
	_, err := o.Insert(attr)
	if err != nil {
		return attr, err
	}
	return attr, nil
}

// 检验当前分类id下的属性id是否存在
func AttrIdExist(catId, attrId int) bool {
	o := orm.NewOrm()
	return o.QueryTable("sp_attribute").Filter("cat_id", catId).Filter("attr_id", attrId).Exist()
}

// 修改参数时验证某一分类下的分类名称是否存在
func UpdateAttrExist(id, attrId int, attrName, attrSel string) bool {
	o := orm.NewOrm()
	return o.QueryTable("sp_attribute").
		Exclude("attr_id", attrId).
		Filter("attr_name", attrName).
		Filter("cat_id", id).
		Filter("attr_sel", attrSel).
		Exist()
}

// 修改参数时验证属性类型是否正确
func AttrSelExist(attrId int, attrSel string) bool {
	o := orm.NewOrm()
	return o.QueryTable("sp_attribute").
		Filter("attr_id", attrId).
		Filter("attr_sel", attrSel).
		Exist()
}

// 修改参数
func UpdateAttr(attrId int, attrName, attrVals string, valsInBody bool) (*SpAttribute, error) {
	o := orm.NewOrm()
	attr := &SpAttribute{AttrId: attrId}
	err := o.Read(attr)
	if err != nil {
		return nil, err
	}
	attr.AttrName = attrName
	if valsInBody {
		attr.AttrVals = attrVals
	}
	_, err = o.Update(attr, "attr_name", "attr_vals")
	if err != nil {
		return nil, err
	}
	return attr, nil
}

// 删除参数
func DeleteAttr(attrId int) error {
	// 执行的是假删操作
	o := orm.NewOrm()
	attr := &SpAttribute{AttrId: attrId}
	err := o.Read(attr, "attrId")
	if err != nil {
		return err
	}
	currentTime := int(time.Now().Unix())
	attr.DeleteTime = &currentTime
	_, err = o.Update(attr)
	if err != nil {
		return err
	}
	return nil
}

// 获取商品数据列表 【接口：goods 请求方式：get】
func GetGoodsList(query string, pagenum, pagesize int) (int, []*SpGoods, error) {
	o := orm.NewOrm()
	goodsList := make([]*SpGoods, 0)
	offset := (pagenum - 1) * pagesize
	qs1 := o.QueryTable("sp_goods").Filter("is_del", "0")
	// OrderBy的参数前使用减号“-”意味着倒叙排列
	qs2 := o.QueryTable("sp_goods").Filter("is_del", "0")
	if strings.Trim(query, " ") != "" {
		// 按条件模糊查询
		qs1 = qs1.Filter("goods_name__icontains", query)
		qs2 = qs2.Filter("goods_name__icontains", query)
	}
	total, err := qs1.Count()
	if err != nil {
		return 0, nil, err
	}
	_, err = qs2.OrderBy("-add_time").Limit(pagesize, offset).All(&goodsList)
	if err != nil {
		return 0, nil, err
	}
	return int(total), goodsList, nil
}

//初始化模型
func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(SpCategory), new(SpAttribute), new(SpGoods), new(SpGoodsPics), new(SpGoodsAttr))
}
