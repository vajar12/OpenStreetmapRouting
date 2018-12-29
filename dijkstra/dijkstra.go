package dijkstra

import (
	"container/heap"
	"errors"
	"math"
	"test/data"

	"github.com/sirupsen/logrus"
)

// CalcDijkstra takes a starting node and returns all edges on the way
// uses edges for the overview of cost and the way to the previous node
func CalcDijkstra(g *data.GraphProd, start *data.Node, target *data.Node) (int64, []*data.Node, error) {

	pq := make(data.PriorityQueue, 0, 10)

	//sets the edge that led to thid node
	prevs := make([]data.Edge, len(g.Nodes))

	for i := range prevs {

		edge := data.Edge{ID: -1, End: start.ID, Start: start.ID, Cost: math.MaxInt64}
		prevs[i] = edge

	}

	heap.Init(&pq)
	//edge for the begining
	edge := data.Edge{ID: -1, End: start.ID, Start: start.ID, Cost: 0}
	heap.Push(&pq, &data.Item{Value: edge, Priority: 0})

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*data.Item)

		currentEdge := item.Value.(data.Edge)

		if item.Priority >= prevs[currentEdge.End].Cost {
			continue
		}

		currentEdge.Cost = item.Priority
		prevs[currentEdge.End] = currentEdge
		// look at all reachable nodes
		edgeBegin := g.Offset[currentEdge.End]
		edgeEnd := g.Offset[currentEdge.End+1]
		for i := edgeBegin; i < edgeEnd; i++ {

			logrus.Debug(i)

			newItem := data.Item{Value: g.Edges[i], Priority: item.Priority + g.Edges[i].Cost}

			// skip if cost is bigger then what we already know
			if newItem.Priority < prevs[g.Edges[i].End].Cost {
				heap.Push(&pq, &newItem)
			}
		}
	}

	// add all nodes that are on the optimal way

	optWay := make([]*data.Node, 0)
	edge = prevs[target.ID]
	minCost := edge.Cost
	if edge.Cost == math.MaxInt64 {
		return math.MaxInt64, nil, errors.New("no way found")
	}
	optWay = append(optWay, &g.Nodes[edge.End])

	for edge.End != start.ID {

		optWay = append(optWay, &g.Nodes[edge.Start])
		edge = prevs[edge.Start]
	}

	return minCost, optWay, nil
}