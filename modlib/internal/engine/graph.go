package engine

import "unsafe"

// WorldGraph represents the WorldGraph global variable in HLDLL
var WorldGraph Graph

// GraphConsts contains the offsets and sizes for Graph
var GraphConsts = graphConsts{
	PNodesOffset:    0xc,
	PLinkPoolOffset: 0x10,
	CNodesOffset:    0x18,
	CLinksOffset:    0x1c,
	CNodeSize:       0x58,
	CLinkSize:       0x18,
}

// NodeConsts contains the offsets and sizes for Node
var NodeConsts = nodeConsts{
	Origin: 0x0,
}

// LinkConsts contains the offsets and sizes for Link
var LinkConsts = linkConsts{
	SrcNode:  0x0,
	DestNode: 0x4,
	LinkEnt:  0x8,
}

type graphConsts struct {
	PNodesOffset    uintptr
	PLinkPoolOffset uintptr
	CNodesOffset    uintptr
	CLinksOffset    uintptr
	CNodeSize       uintptr
	CLinkSize       uintptr
}

type nodeConsts struct {
	Origin uintptr
}

type linkConsts struct {
	SrcNode  uintptr
	DestNode uintptr
	LinkEnt  uintptr
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

// Link represents CLink
type Link struct {
	ptr unsafe.Pointer
}

// MakeLink creates a new instance of Link
func MakeLink(pointer unsafe.Pointer) Link {
	return Link{ptr: pointer}
}

// Source returns a node pointed to by CLink::m_iSrcNode
func (link Link) Source() Node {
	idx := *(*int)(unsafe.Pointer(uintptr(link.ptr) + LinkConsts.SrcNode))
	return WorldGraph.Node(idx)
}

// Destination returns a node pointed to by CLink::m_iDestNode
func (link Link) Destination() Node {
	idx := *(*int)(unsafe.Pointer(uintptr(link.ptr) + LinkConsts.DestNode))
	return WorldGraph.Node(idx)
}

// LinkEnt returns CLink::m_pLinkEnt
func (link *Link) LinkEnt() EntVars {
	ptr := *(*unsafe.Pointer)(unsafe.Pointer(uintptr(link.ptr) + LinkConsts.LinkEnt))
	return MakeEntVars(ptr)
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
	base := *(*unsafe.Pointer)(unsafe.Pointer(uintptr(graph.ptr) + GraphConsts.PNodesOffset))
	return MakeNode(unsafe.Pointer(uintptr(base) + uintptr(index)*GraphConsts.CNodeSize))
}

// NumLinks returns CGraph::m_cLinks
func (graph Graph) NumLinks() int {
	return *(*int)(unsafe.Pointer(uintptr(graph.ptr) + GraphConsts.CLinksOffset))
}

// Link returns CGraph::m_pLinkPool[index]
func (graph Graph) Link(index int) Link {
	base := *(*unsafe.Pointer)(unsafe.Pointer(uintptr(graph.ptr) + GraphConsts.PLinkPoolOffset))
	return MakeLink(unsafe.Pointer(uintptr(base) + uintptr(index)*GraphConsts.CLinkSize))
}
