package utils

import (
	"BeeDemo/models"
	"github.com/beego/beego/v2/client/orm"
	"time"
)

// Tree 统一定义菜单树的数据结构，也可以自定义添加其他字段
type Trees struct {
	Title    string      `json:"title"`    //节点名字
	Data     interface{} `json:"data"`     //自定义对象
	Children []Trees     `json:"children"` //子节点
}

// ConvertToINodeArray 其他的结构体想要生成菜单树，直接实现这个接口
type INode interface {
	// GetId获取id
	GetId() int
	// Pid 获取父id
	GetPid() int
	// GetTitle 获取显示名字
	GetTitle() string
	// GetData 获取附加数据
	GetData() interface{}
	// IsRoot 判断当前节点是否是顶层根节点
	IsRoot() bool
}

type INodes []INode

func (nodes INodes) Len() int {
	return len(nodes)
}
func (nodes INodes) Swap(i, j int) {
	nodes[i], nodes[j] = nodes[j], nodes[i]
}
func (nodes INodes) Less(i, j int) bool {
	return nodes[i].GetId() < nodes[j].GetId()
}

// GenerateTree 自定义的结构体实现 INode 接口后调用此方法生成树结构
// nodes 需要生成树的节点
// selectedNode 生成树后选中的节点
// menuTrees 生成成功后的树结构对象
func GenerateTree(nodes []INode) (trees []Trees) {
	trees = []Trees{}
	// 定义顶层根(roots)和子节点(children)
	var roots, children []INode
	for _, v := range nodes {

		if v.IsRoot() {
			// 判断顶层根节点
			roots = append(roots, v)
		}
		children = append(children, v)
	}

	for _, v := range roots {
		childTree := &Trees{
			Title:    v.GetTitle(),
			Data:     v.GetData(),
			Children: []Trees{},
		}
		// 递归
		recursiveTree(childTree, children)

		trees = append(trees, *childTree)
	}

	return trees
}

// recursiveTree 递归生成树结构
// tree 递归的树对象
// nodes 递归的节点
// selectedNodes 选中的节点
func recursiveTree(tree *Trees, nodes []INode) {
	//类型断言
	data := tree.Data.(INode)

	for _, v := range nodes {
		if v.IsRoot() {
			// 如果当前节点是顶层根节点就跳过
			continue
		}

		if data.GetId() == v.GetPid() {
			childTree := &Trees{
				Title:    v.GetTitle(),
				Data:     v.GetData(),
				Children: []Trees{},
			}

			recursiveTree(childTree, nodes)

			tree.Children = append(tree.Children, *childTree)
		}
	}
}

// =================================================================== //
// =================================================================== //
// =================================================================== //
// 定义我们自己的菜单对象
type SystemMenu struct {
	Id       int    `json:"id"`        //id
	FatherId int    `json:"father_id"` //上级菜单id
	Name     string `json:"name"`      //菜单名
	Route    string `json:"route"`     //页面路径
	Icon     string `json:"icon"`      //图标路径
}

type CommentTree struct {
	Id         int       `json:"id"`
	Pid        int       `json:"pid"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"create_time"`
}

func (s CommentTree) GetTitle() string {
	return s.Content
}

func (s CommentTree) GetId() int {
	return s.Id
}

func (s CommentTree) GetPid() int {
	return s.Pid
}

func (s CommentTree) GetData() interface{} {
	return s
}

func (s CommentTree) IsRoot() bool {
	// 这里通过FatherId等于0 或者 FatherId等于自身Id表示顶层根节点
	return s.Pid == 0 || s.Pid == s.Id
}

//多维节点Tree数据
type Tress []CommentTree

func TestGenerateTree() []Trees {
	// 模拟获取数据库中所有菜单，在其它所有的查询中，也是首先将数据库中所有数据查询出来放到数组中，
	// 后面的遍历递归，都在这个 allMenu中进行，而不是在数据库中进行递归查询，减小数据库压力。
	comments := make([]models.Comment, 0)
	_, err := orm.NewOrm().
		QueryTable(new(models.Comment)).
		Filter("post_id", 10).
		All(&comments)

	if err != nil {
		return nil
	}

	var result Tress
	for _, comment := range comments {
		result = append(result, CommentTree{
			Id:         comment.Id,
			Pid:        comment.Pid,
			Content:    comment.Content,
			CreateTime: comment.CreateTime,
		})
	}

	return GenerateTree(Tress.ConvertToINodeArray(result))

	//return GenerateTree(ConvertToINodeArrayBak(result))
}

// ConvertToINodeArray 将当前数组转换成父类 INode 接口 数组
// 这里没有参数，为啥会(能)传参？
func (s Tress) ConvertToINodeArray() (nodes []INode) {
	//fmt.Println(ids)
	for _, v := range s {
		nodes = append(nodes, v)
	}

	return nodes
}

func ConvertToINodeArrayBak(s Tress) (nodes []INode) {
	for _, v := range s {
		nodes = append(nodes, v)
	}

	return nodes
}
