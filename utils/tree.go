package utils

import (
	"BeeDemo/models"
	"encoding/json"
	"fmt"
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
}

// ConvertToINodeArray 将当前数组转换成父类 INode 接口 数组
// 这里没有参数，为啥会(能)传参？
func (s Tress) ConvertToINodeArray() (nodes []INode) {
	for _, v := range s {
		nodes = append(nodes, v)
	}

	return nodes
}

// 未经过封装的
// ---------- tree
// --------------- 示例

type Datanode struct {
	Id    int         `json:"id"`
	PId   int         `json:"pid"`
	Name  string      `json:"name"`
	Child []*Datanode `json:"child"`
}

func main() {

	Data := GetResult()
	//父节点
	pid := 0
	MakeTree(Data, Data[pid]) //调用生成tree
	TransFormJson(Data[pid])  //转化为json
}

func GetResult() []*Datanode {
	return []*Datanode{
		{
			Id:   0,
			PId:  -1,
			Name: "目录",
		},
		{
			Id:   1,
			PId:  0,
			Name: "一、水果",
		},
		{
			Id:   2,
			PId:  1,
			Name: "1.苹果",
		},
		{
			Id:   3,
			PId:  1,
			Name: "2.香蕉",
		},
		{
			Id:   4,
			PId:  0,
			Name: "二、蔬菜",
		},
		{
			Id:   5,
			PId:  4,
			Name: "1.芹菜",
		},
		{
			Id:   6,
			PId:  4,
			Name: "2.黄瓜",
		},
		{
			Id:   7,
			PId:  6,
			Name: "(1)黄瓜特点",
		},
		{
			Id:   8,
			PId:  4,
			Name: "3.西红柿",
		},
		{
			Id:   9,
			PId:  0,
			Name: "三、关系",
		},
		{
			Id:   10,
			PId:  6,
			Name: "(2)黄瓜品质",
		},
		{
			Id:   11,
			PId:  6,
			Name: "(2)黄瓜颜色",
		},
		{
			Id:   12,
			PId:  6,
			Name: "(2)黄瓜产地",
		},
	}
}
func MakeTree(Data []*Datanode, node *Datanode) { //参数为父节点，添加父节点的子节点指针切片
	childs, _ := HasChild(Data, node) //判断节点是否有子节点并返回
	if childs == nil {
		return
	}

	node.Child = append(node.Child, childs[0:]...) //添加子节点

	for _, v := range childs {
		//查询子节点的子节点，并添加到子节点
		_, has := HasChild(Data, v)
		if has {
			MakeTree(Data, v) //递归添加节点
		}
	}
	return
}

func HasChild(Data []*Datanode, node *Datanode) (child []*Datanode, yes bool) {
	for idx, v := range Data {
		if v.PId == node.Id {
			child = append(child, v)
			continue
		}

		if Data[idx].Child == nil {
			Data[idx].Child = []*Datanode{}
		}
	}

	if child != nil {
		yes = true
	}

	return child, yes
}

func TransFormJson(Data *Datanode) { //转为json

	JsonData, _ := json.Marshal(Data)

	fmt.Println(string(JsonData))
}
