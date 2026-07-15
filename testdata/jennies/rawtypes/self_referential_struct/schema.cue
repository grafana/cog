package self_referential_struct

// Node represents a node in a singly-linked list.
// The next field points to the following node, or is absent if this is the last node.
Node: {
	value: string
	next?: Node
}
