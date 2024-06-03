package main

import (
	"crypto/sha1"
	"fmt"
	"strconv"
	"sync"
)

const (
	M         = 8
	RING_SIZE = 1 << M
)

type Node struct {
	id          int
	ipAddress   string
	port        int
	successor   *Node
	predecessor *Node
	fingerTable []*Node
	keys        map[int]string
	mu          sync.Mutex
}

func hashFunction(key string) int {
	h := sha1.New()
	h.Write([]byte(key))
	bs := h.Sum(nil)
	hash := int(bs[0])<<24 | int(bs[1])<<16 | int(bs[2])<<8 | int(bs[3])
	return hash % RING_SIZE
}

func newNode(ipAddress string, port int) *Node {
	node := &Node{
		id:          hashFunction(fmt.Sprintf("%s:%d", ipAddress, port)),
		ipAddress:   ipAddress,
		port:        port,
		successor:   nil,
		predecessor: nil,
		fingerTable: make([]*Node, M),
		keys:        make(map[int]string),
	}
	for i := 0; i < M; i++ {
		node.fingerTable[i] = node
	}
	fmt.Println("NEW NODE:", node.ipAddress, "PORT:", node.port, "HASH:", node.id)
	return node
}

func (n *Node) findSuccessor(id int, shouldPrintHop bool) *Node {
	if shouldPrintHop {
		fmt.Print("->", n.id)
	}
	if n == n.successor {
		return n
	}
	if between(n.id, id, n.successor.id) {
		return n.successor
	} else {
		node := n.closestPrecedingNode(id)
		if node == n {
			return n.successor.findSuccessor(id, shouldPrintHop)
		}
		return node.findSuccessor(id, shouldPrintHop)
	}
}

func (n *Node) closestPrecedingNode(id int) *Node {
	for i := M - 1; i >= 0; i-- {
		if n.fingerTable[i] != nil && between(n.id, n.fingerTable[i].id, id) {
			return n.fingerTable[i]
		}
	}
	return n
}

func (n *Node) join(existingNode *Node) {
	if existingNode != nil {
		n.initFingerTable(existingNode)
		n.updateOthers()
		n.moveKeys()
	} else {
		for i := 0; i < M; i++ {
			n.fingerTable[i] = n
		}
		n.predecessor = n
		n.successor = n
	}
}

func (n *Node) initFingerTable(existingNode *Node) {
	n.fingerTable[0] = existingNode.findSuccessor((n.id+1)%RING_SIZE, false)
	n.successor = n.fingerTable[0]
	n.predecessor = n.successor.predecessor
	n.successor.predecessor = n

	for i := 0; i < M-1; i++ {
		start := (n.id + 1<<i) % RING_SIZE
		if between(n.id, start, n.fingerTable[i].id) {
			n.fingerTable[i+1] = n.fingerTable[i]
		} else {
			n.fingerTable[i+1] = existingNode.findSuccessor(start, false)
		}
	}
}

func (n *Node) updateOthers() {
	for i := 0; i < M; i++ {
		pred := n.findPredecessor((n.id - 1<<i + RING_SIZE) % RING_SIZE)
		pred.updateFingerTable(n, i)
	}
}

func (n *Node) updateFingerTable(node *Node, i int) {
	if node.id == n.id {
		return
	}
	if n.fingerTable[i] == nil || between(n.id, node.id, n.fingerTable[i].id) {
		n.fingerTable[i] = node
		pred := n.predecessor
		if pred != nil && pred != node {
			pred.updateFingerTable(node, i)
		}
	}
}

func (n *Node) findPredecessor(id int) *Node {
	node := n
	for !between(node.id, id, node.successor.id) {
		node = node.closestPrecedingNode(id)
	}
	return node
}

func (n *Node) moveKeys() {
	n.mu.Lock()
	defer n.mu.Unlock()
	for key, value := range n.keys {
		if !between(n.id, hashFunction(strconv.Itoa(key)), n.successor.id) {
			n.successor.keys[key] = value
			delete(n.keys, key)
		}
	}
}

func (n *Node) stabilize() {
	x := n.successor.predecessor
	if x != nil && between(n.id, x.id, n.successor.id) {
		n.successor = x
	}
	n.successor.notify(n)
}

func (n *Node) notify(node *Node) {
	if n.predecessor == nil || between(n.predecessor.id, node.id, n.id) {
		n.predecessor = node
	}
}

func (n *Node) fixFingers() {
	for i := 0; i < M; i++ {
		n.fingerTable[i] = n.findSuccessor((n.id+1<<i)%RING_SIZE, false)
	}
}

func (n *Node) store(key string, value string) {
	keyID := hashFunction(key)
	fmt.Println("FILE KEY ID:", keyID)
	successor := n.findSuccessor(keyID, false)
	if successor != nil {
		successor.mu.Lock()
		defer successor.mu.Unlock()
		fmt.Println("NODE ID STORING THE KEY:", successor.id)
		successor.keys[keyID] = value
	}
}

func (n *Node) lookup(key string) string {
	keyID := hashFunction(key)
	successor := n.findSuccessor(keyID, true)
	if successor == nil {
		return ""
	}
	successor.mu.Lock()
	defer successor.mu.Unlock()
	return successor.keys[keyID]
}

func between(id1, id2, id3 int) bool {
	if id1 < id3 {
		return id1 < id2 && id2 <= id3
	}
	return id1 < id2 || id2 <= id3
}

func (n *Node) printNodeData() {
	fmt.Println("==============================")
	fmt.Println("Node ID", n.id)
	fmt.Println("Node Successor:", n.successor.id)
	fmt.Println("Node Predecessor:", n.predecessor.id)
	fmt.Println("^^^^^^^^^^^^^^^^")
	for i, v := range n.fingerTable {
		fmt.Println("Node FingerTable", i, v.id)
	}
	fmt.Println("==============================")
}

func main() {
	// Create initial node
	node1 := newNode("127.0.0.1", 8000)
	node1.join(nil) // First node in the network

	// Join new nodes
	node2 := newNode("127.0.0.1", 8001)
	node2.join(node1)

	node3 := newNode("127.0.0.1", 8002)
	node3.join(node1)

	node4 := newNode("127.0.0.1", 8003)
	node4.join(node1)

	node5 := newNode("127.0.0.1", 8004)
	node5.join(node1)

	// Run stabilization to update finger tables
	for i := 0; i < 3; i++ {
		node1.stabilize()
		node2.stabilize()
		node3.stabilize()
		node4.stabilize()
		node5.stabilize()
		node1.fixFingers()
		node2.fixFingers()
		node3.fixFingers()
		node4.fixFingers()
		node5.fixFingers()
	}

	node1.printNodeData()
	node2.printNodeData()
	node3.printNodeData()
	node4.printNodeData()
	node5.printNodeData()

	// Store and lookup keys
	node1.store("cnn.com/index.html", "File Content")
	fmt.Println(node2.lookup("cnn.com/index.html")) // Should print "File Content"
	fmt.Println(node3.lookup("cnn.com/index.html"))

	node2.store("bbc.com/index.html", "File Data")
	fmt.Println(node1.lookup("bbc.com/index.html")) // Should print "File Data"
}
