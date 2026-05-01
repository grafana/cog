package self_referential_struct;


// Node represents a node in a singly-linked list.
// The next field points to the following node, or is absent if this is the last node.
public class Node {
    public String value;
    public Node next;
    public Node() {
        this.value = "";
    }
    public Node(String value,Node next) {
        this.value = value;
        this.next = next;
    }
}
