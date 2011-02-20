package graph

import (
	"github.com/StepLg/go-erx/src/erx"
)

// Mixed graph with matrix as a internal representation.
//
// Doesn't allow duplicate edges and arcs, loops and reversed arcs.
// Graph can't have more than size vertexes, where size set during initialization.
// MixedMatrix use over (size^2/2) * sizeof(MixedConnectionType) bytes.
type MixedMatrix struct {
	nodes []MixedConnectionType
	size int
	VertexIds map[VertexId]int // internal node ids, used in nodes array
	edgesCnt int
	arcsCnt int
}

func NewMixedMatrix(size int) *MixedMatrix {
	if size<=0 {
		panic(erx.NewError("Trying to create mixed matrix graph with zero size"))
	}
	g := new(MixedMatrix)
	g.nodes = make([]MixedConnectionType, size*(size-1)/2)
	g.size = size
	g.VertexIds = make(map[VertexId]int)
	return g
}

///////////////////////////////////////////////////////////////////////////////
// GraphVertexesWriter

// Adding single node to graph
func (gr *MixedMatrix) AddNode(node VertexId) {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Add node to graph.", e)
			err.AddV("node id", node)
			panic(err)
		}
	}()
	
	if _, ok := gr.VertexIds[node]; ok {
		panic(erx.NewError("Node already exists."))
	}
	
	if len(gr.VertexIds) == gr.size {
		panic(erx.NewError("Not enough space to add new node"))
	}
	
	gr.VertexIds[node] = len(gr.VertexIds)
}

///////////////////////////////////////////////////////////////////////////////
// GraphVertexesReader

// Getting nodes count in graph
func (gr *MixedMatrix) Order() int {
	return len(gr.VertexIds)
}

///////////////////////////////////////////////////////////////////////////////
// GraphVertexesRemover

// Removing node from graph
func (gr *MixedMatrix) RemoveNode(node VertexId) {
	panic("Function RemoveNode doesn't implement in MixedMatrix graph yet.")
}
	
///////////////////////////////////////////////////////////////////////////////
// ConnectionsIterable
func (gr *MixedMatrix) ConnectionsIter() <-chan Connection {
	ch := make(chan Connection)
	go func() {
		for from, _ := range gr.VertexIds {
			for to, _ := range gr.VertexIds {
				if from>=to {
					continue
				}
				
				conn := gr.getConnectionId(from, to, false)
				if gr.nodes[conn]!=CT_NONE {
					ch <- Connection{from, to}
				}
			}
		}
		close(ch)
	}()
	return ch
}

///////////////////////////////////////////////////////////////////////////////
// VertexesIterable
func (gr *MixedMatrix) VertexesIter() <-chan VertexId {
	ch := make(chan VertexId)
	go func() {
		for VertexId, _ := range gr.VertexIds {
			ch <- VertexId
		}
		close(ch)
	}()
	return ch
}

///////////////////////////////////////////////////////////////////////////////
// VertexesChecker

func (g *MixedMatrix) CheckNode(node VertexId) (exists bool) {
	_, exists = g.VertexIds[node]
	return
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphEdgesWriter

// Adding new edge to graph
func (gr *MixedMatrix) AddEdge(node1, node2 VertexId) {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Add edge to mixed graph.", e)
			err.AddV("node 1", node1)
			err.AddV("node 2", node2)
			panic(err)
		}
	}()

	conn := gr.getConnectionId(node1, node2, true)
	if gr.nodes[conn]!=CT_NONE {
		err := erx.NewError("Duplicate edge.")
		err.AddV("connection id", conn)
		err.AddV("type", gr.nodes[conn]) 
		panic(err)
	}
	
	gr.nodes[conn] = CT_UNDIRECTED
	gr.edgesCnt++
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphEdgesRemover

// Removing edge, connecting node1 and node2
func (gr *MixedMatrix) RemoveEdge(node1, node2 VertexId) {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Remove edge from mixed graph.", e)
			err.AddV("node 1", node1)
			err.AddV("node 2", node2)
			panic(err)
		}
	}()

	conn := gr.getConnectionId(node1, node2, true)
	if gr.nodes[conn]!=CT_UNDIRECTED {
		err := erx.NewError("Edge doesn't exists.")
		err.AddV("connection id", conn)
		err.AddV("type", gr.nodes[conn])
		panic(err)
	}
	
	gr.nodes[conn] = CT_NONE
	gr.edgesCnt--
}

///////////////////////////////////////////////////////////////////////////////
// UndirectedGraphReader

// Edges count in graph
func (gr *MixedMatrix) EdgesCnt() int {
	return gr.edgesCnt
}

// Checking edge existance between node1 and node2
//
// node1 and node2 must exist in graph or error will be returned
func (gr *MixedMatrix) CheckEdge(node1, node2 VertexId) bool {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Check edge in mixed graph.", e)
			err.AddV("node 1", node1)
			err.AddV("node 2", node2)
			panic(err)
		}
	}()

	if node1==node2 {
		return false
	}
	
	return gr.nodes[gr.getConnectionId(node1, node2, false)]==CT_UNDIRECTED
}

///////////////////////////////////////////////////////////////////////////////
// DirectedGraphArcsWriter

// Adding directed arc to graph
func (gr *MixedMatrix) AddArc(tail, head VertexId) {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Add arc to mixed graph.", e)
			err.AddV("tail", tail)
			err.AddV("head", head)
			panic(err)
		}
	}()

	conn := gr.getConnectionId(tail, head, true)
	if gr.nodes[conn]!=CT_NONE {
		err := erx.NewError("Duplicate edge.")
		err.AddV("connection id", conn)
		err.AddV("type", gr.nodes[conn]) 
		panic(err)
	}
	
	if tail<head {
		gr.nodes[conn] = CT_DIRECTED
	} else {
		gr.nodes[conn] = CT_DIRECTED_REVERSED
	}
	
	gr.arcsCnt++
}

///////////////////////////////////////////////////////////////////////////////
// DirectedGraphArcsRemover

// Removding directed arc
func (gr *MixedMatrix) RemoveArc(tail, head VertexId) {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Remove arc from mixed graph.", e)
			err.AddV("tail", tail)
			err.AddV("head", head)
			panic(err)
		}
	}()

	conn := gr.getConnectionId(tail, head, true)
	expectedType := CT_NONE
	if tail<head {
		expectedType = CT_DIRECTED
	} else {
		expectedType = CT_DIRECTED_REVERSED
	}
	
	if gr.nodes[conn]!=expectedType {
		err := erx.NewError("Arc doesn't exists.")
		err.AddV("connection id", conn)
		err.AddV("type", gr.nodes[conn])
		panic(err)
	}
	
	gr.nodes[conn] = CT_NONE
	gr.arcsCnt--
}

///////////////////////////////////////////////////////////////////////////////
// DirectedGraphReader

// Getting arcs count in graph
func (gr *MixedMatrix) ArcsCnt() int {
	return gr.arcsCnt
}

// Checking arrow existance between node1 and node2
//
// node1 and node2 must exist in graph or error will be returned
func (gr *MixedMatrix) CheckArc(tail, head VertexId) bool {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Check arc in mixed graph.", e)
			err.AddV("tail", tail)
			err.AddV("head", head)
			panic(err)
		}
	}()
	
	checkingType := CT_NONE
	if tail < head {
		checkingType = CT_DIRECTED
	} else {
		checkingType = CT_DIRECTED_REVERSED
	}
	
	return gr.nodes[gr.getConnectionId(tail, head, false)]==checkingType
}

///////////////////////////////////////////////////////////////////////////////
// MixedGraphSpecificReader

// Iterate over only undirected edges
func (gr *MixedMatrix) EdgesIter() <-chan Connection {
	ch := make(chan Connection)
	go func() {
		for from, _ := range gr.VertexIds {
			for to, _ := range gr.VertexIds {
				if from>=to {
					continue
				}
				
				if gr.nodes[gr.getConnectionId(from, to, false)]==CT_UNDIRECTED {
					ch <- Connection{from, to}
				}
			}
		}
		close(ch)
	}()
	return ch
}
	
// Iterate over only directed arcs
func (gr *MixedMatrix) ArcsIter() <-chan Connection {
	ch := make(chan Connection)
	go func() {
		for from, _ := range gr.VertexIds {
			for to, _ := range gr.VertexIds {
				if from>=to {
					continue
				}
				
				conn := gr.getConnectionId(from, to, false)
				if gr.nodes[conn]==CT_DIRECTED {
					ch <- Connection{from, to}
				}
				if gr.nodes[conn]==CT_DIRECTED_REVERSED {
					ch <- Connection{to, from}
				}
			}
		}
		close(ch)
	}()
	return ch
}

func (gr *MixedMatrix) CheckEdgeType(tail VertexId, head VertexId) MixedConnectionType {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Check edge type in mixed graph.", e)
			err.AddV("tail", tail)
			err.AddV("head", head)
			panic(err)
		}
	}()
	
	conn := gr.getConnectionId(tail, head, false)
	connType := gr.nodes[conn]
	if connType==CT_DIRECTED && tail>head {
		return CT_DIRECTED_REVERSED
	}
	return connType
}

func (g *MixedMatrix) ConnectionsCnt() int {
	return g.arcsCnt + g.edgesCnt
}

func (gr *MixedMatrix) TypedConnectionsIter() <-chan TypedConnection {
	ch := make(chan TypedConnection)
	go func() {
		for from, _ := range gr.VertexIds {
			for to, _ := range gr.VertexIds {
				if from>=to {
					continue
				}
				
				conn := gr.getConnectionId(from, to, false)
				switch gr.nodes[conn] {
					case CT_NONE:
					case CT_UNDIRECTED:
						ch <- TypedConnection{Connection:Connection{Tail: from, Head:to}, Type:CT_UNDIRECTED} 
					case CT_DIRECTED:
						ch <- TypedConnection{Connection:Connection{Tail: from, Head:to}, Type:CT_DIRECTED}
					case CT_DIRECTED_REVERSED:
						ch <- TypedConnection{Connection:Connection{Tail: to, Head:from}, Type:CT_DIRECTED}
					default:
						err := erx.NewError("Internal error: wrong connection type in mixed graph matrix")
						err.AddV("connection type", gr.nodes[conn])
						err.AddV("connection id", conn)
						err.AddV("tail node", from)
						err.AddV("head node", to)
						panic(err)
				}
			}
		}
		close(ch)
	}()
	return ch
}

func (gr *MixedMatrix) getConnectionId(node1, node2 VertexId, create bool) int {
	defer func() {
		if e := recover(); e!=nil {
			err := erx.NewSequent("Calculating connection id.", e)
			err.AddV("node 1", node1)
			err.AddV("node 2", node2)
			panic(err)
		}
	}()
	
	var id1, id2 int
	node1Exist := false
	node2Exist := false
	id1, node1Exist = gr.VertexIds[node1]
	id2, node2Exist = gr.VertexIds[node2]
	
	// checking for errors
	{
		if node1==node2 {
			panic(erx.NewError("Equal nodes."))
		}
		if !create {
			if !node1Exist {
				panic(erx.NewError("First node doesn't exist in graph"))
			}
			if !node2Exist {
				panic(erx.NewError("Second node doesn't exist in graph"))
			}
		} else if !node1Exist || !node2Exist {
			if node1Exist && node2Exist {
				if gr.size - len(gr.VertexIds) < 2 {
					panic(erx.NewError("Not enough space to create two new nodes."))
				}
			} else {
				if gr.size - len(gr.VertexIds) < 1 {
					panic(erx.NewError("Not enough space to create new node."))
				}
			}
		}
	}
	
	if !node1Exist {
		id1 = int(len(gr.VertexIds))
		gr.VertexIds[node1] = id1
	}

	if !node2Exist {
		id2 = int(len(gr.VertexIds))
		gr.VertexIds[node2] = id2
	}
	
	// switching id1, id2 in order to id1 < id2
	if id1>id2 {
		id1, id2 = id2, id1
	}
	
	// id from upper triangle matrix, stored in vector
	connId := id1*(gr.size-1) + id2 - 1 - id1*(id1+1)/2
	return connId 
}
