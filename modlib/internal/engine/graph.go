package engine

import "unsafe"

// WorldGraph represents the WorldGraph global variable in HLDLL
var WorldGraph Graph

// GraphConsts contains the offsets and sizes for Graph
var GraphConsts = graphConsts{
	PNodesOffset: 0xc,
	CNodesOffset: 0x18,
	CNodeSize:    0x58,
}

// NodeConsts contains the offsets and sizes for Node
var NodeConsts = nodeConsts{
	Origin: 0x0,
}

type graphConsts struct {
	PNodesOffset uintptr
	CNodesOffset uintptr
	CNodeSize    uintptr
}

type nodeConsts struct {
	Origin uintptr
}

// Node represents CNode
type Node struct {
	ptr unsafe.Pointer
}

// MakeNode creates a new instance of Node
func MakeNode(pointer unsafe.Pointer) Node {
	return Node{ptr: pointer}
}

// Origin returns CNode::m_vecOrigin
func (node Node) Origin() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(uintptr(node.ptr) + NodeConsts.Origin))
}

// Graph represents CGraph
type Graph struct {
	ptr unsafe.Pointer
}

// SetPointer sets the base pointer of WorldGraph
func (graph *Graph) SetPointer(pointer unsafe.Pointer) {
	graph.ptr = pointer
}

// NumNodes returns CGraph::m_cNodes
func (graph Graph) NumNodes() int {
	return *(*int)(unsafe.Pointer(uintptr(graph.ptr) + GraphConsts.CNodesOffset))
}

// Node returns CGraph::m_pNodes[index]
func (graph Graph) Node(index int) Node {
	base := *(*uintptr)(unsafe.Pointer(uintptr(graph.ptr) + GraphConsts.PNodesOffset))
	return MakeNode(unsafe.Pointer(base + uintptr(index)*GraphConsts.CNodeSize))
}
