package hashs

//一致性哈希(Consistent Hashing)

import (
	"sort"
	"strconv"
	"sync"

	"github.com/spaolacci/murmur3"
)

// DefaultReplicas 默认大小
const DefaultReplicas = 160

//HashRing ..
type HashRing []uint32

func (c HashRing) Len() int {
	return len(c)
}

func (c HashRing) Less(i, j int) bool {
	return c[i] < c[j]
}

func (c HashRing) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

//Node 节点
type Node struct {
	ID       int
	Addr     string
	HostName string
	Weight   int
}

//NewNode new
func NewNode(id int, addr string, name string, weight int) *Node {
	return &Node{
		ID:       id,
		Addr:     addr,
		HostName: name,
		Weight:   weight,
	}
}

//Consistent ..
type Consistent struct {
	Nodes     map[uint32]Node
	numReps   int
	Resources map[int]*Node
	ring      HashRing
	sync.RWMutex
}

//NewConsistent new
func NewConsistent() *Consistent {
	return &Consistent{
		Nodes:     make(map[uint32]Node),
		numReps:   DefaultReplicas,
		Resources: make(map[int]*Node),
		ring:      HashRing{},
	}
}

//Add add
func (c *Consistent) Add(node *Node) bool {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.Resources[node.ID]; ok {
		return false
	}

	count := c.numReps * node.Weight
	for i := 0; i < count; i++ {
		str := c.joinStr(i, node)
		c.Nodes[c.hashStr(str)] = *(node)
	}
	c.Resources[node.ID] = node
	c.sortHashRing()
	return true
}

func (c *Consistent) sortHashRing() {
	c.ring = HashRing{}
	for k := range c.Nodes {
		c.ring = append(c.ring, k)
	}
	sort.Sort(c.ring)
}

func (c *Consistent) joinStr(i int, node *Node) string {
	return node.Addr + "*" + strconv.Itoa(node.Weight) +
		"-" + strconv.Itoa(i) +
		"-" + strconv.Itoa(node.ID)
}

// MurMurHash算法 :https://github.com/spaolacci/murmur3
func (c *Consistent) hashStr(key string) uint32 {
	//	return crc32.ChecksumIEEE([]byte(key))

	return murmur3.Sum32([]byte(key))
}

//Get get
func (c *Consistent) Get(key string) (Node, bool) {
	c.Lock()
	defer c.Unlock()

	if len(c.ring) == 0 {
		return Node{}, false
	}

	hash := c.hashStr(key)
	i := c.search(hash)

	node, exists := c.Nodes[c.ring[i]]
	return node, exists
}

func (c *Consistent) search(hash uint32) int {

	i := sort.Search(len(c.ring), func(i int) bool { return c.ring[i] >= hash })
	if i < len(c.ring) {
		if i == len(c.ring)-1 {
			return 0
		}
		return i
	}
	return len(c.ring) - 1
}

//Remove rem
func (c *Consistent) Remove(node *Node) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Resources[node.ID]; !ok {
		return
	}

	delete(c.Resources, node.ID)

	count := c.numReps * node.Weight
	for i := 0; i < count; i++ {
		str := c.joinStr(i, node)
		delete(c.Nodes, c.hashStr(str))
	}
	c.sortHashRing()
}

//GetNodeByID get
func (c *Consistent) GetNodeByID(id int) (*Node, bool) {
	node, exists := c.Resources[id]
	return node, exists
}
