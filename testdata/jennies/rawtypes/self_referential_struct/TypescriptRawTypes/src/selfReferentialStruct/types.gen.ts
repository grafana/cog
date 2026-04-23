// Node represents a node in a singly-linked list.
// The next field points to the following node, or is absent if this is the last node.
export interface Node {
	value: string;
	next?: Node;
}

export const defaultNode = (): Node => ({
	value: "",
});

