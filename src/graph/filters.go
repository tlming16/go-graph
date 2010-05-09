package graph

// Arcs filter
//
// This is arcs filter for DirectedGraphReader. Initialize it with arcs, which need to be filtered
// and they never appeared in GetAccessors, GetPredecessors, CheckArc and Iter functions.
//
// Be careful! Filter doesn't affect GetSources and GetSinks functions. Also it doesn't recalculate
// dangling vertexes.
type DirectedGraphArcsFilter struct {
	DirectedGraphReader
	arcs []Connection
}

// Create arcs filter with array of filtering arcs
func NewArcsFilter(g DirectedGraphReader, arcs []Connection) *DirectedGraphArcsFilter {
	filter := &DirectedGraphArcsFilter{
		DirectedGraphReader: g,
		arcs: arcs,
	}
	return filter
}

// Create arcs filter with single arc
func NewArcFilter(g DirectedGraphReader, tail, head NodeId) *DirectedGraphArcsFilter {
	filter := &DirectedGraphArcsFilter{
		DirectedGraphReader: g,
		arcs: make([]Connection, 1),
	}
	filter.arcs[0].Tail = tail
	filter.arcs[0].Head = head
	return filter	
}

// Getting node accessors
func (filter *DirectedGraphArcsFilter) GetAccessors(node NodeId) Nodes {
	accessors := filter.DirectedGraphReader.GetAccessors(node)
	newAccessorsLen := len(accessors)
	for _, filteringConnection := range filter.arcs {
		if node == filteringConnection.Tail {
			// need to remove filtering arc
			k := 0
			for k=0; k<newAccessorsLen; k++ {
				if accessors[k]==filteringConnection.Head {
					break
				}
			}
			if k<newAccessorsLen {
				copy(accessors[k:newAccessorsLen-1], accessors[k+1:newAccessorsLen])
				newAccessorsLen--
			}
		}
	}
	return accessors[0:newAccessorsLen]
}

// Getting node predecessors
func (filter *DirectedGraphArcsFilter) GetPredecessors(node NodeId) Nodes {
	accessors := filter.DirectedGraphReader.GetAccessors(node)
	newAccessorsLen := len(accessors)
	for _, filteringConnection := range filter.arcs {
		if node == filteringConnection.Head {
			// need to remove filtering arc
			k := 0
			for k=0; k<newAccessorsLen; k++ {
				if accessors[k]==filteringConnection.Tail {
					break
				}
			}
			if k<newAccessorsLen {
				copy(accessors[k:newAccessorsLen-1], accessors[k+1:newAccessorsLen])
				newAccessorsLen--
			}
		}
	}
	return accessors[0:newAccessorsLen]
}

// Checking arrow existance between node1 and node2
//
// node1 and node2 must exist in graph or error will be returned
func (filter *DirectedGraphArcsFilter) CheckArc(node1, node2 NodeId) bool {
	res := filter.DirectedGraphReader.CheckArc(node1, node2)
	if res {
		for _, filteringConnection := range filter.arcs {
			if filteringConnection.Tail==node1 && filteringConnection.Head==node2 {
				res = false
				break
			}
		}
	}
	return res
}

func (filter *DirectedGraphArcsFilter) ConnectionsIter() <-chan Connection {
	ch := make(chan Connection)
	go func() {
		for conn := range filter.DirectedGraphReader.ConnectionsIter() {
			for _, filteringConnection := range filter.arcs {
				if filteringConnection.Head==conn.Head && filteringConnection.Tail==conn.Tail {
					continue
				}
			}
		}
		close(ch)
	}()
	return ch
}