obj Node {
    key: int
    left: Node?
    right: Node?
    height: int

    fn get_balance_factor = left?.height - right?.height
    fn max_height = math.max left?.height right?.height
}

obj AVLTree {
    root: Node!?

    fn insert_node(node: Node!?, key: int) => Node {
        if !node then => Node{key: key, height: 1}

        if key < node.key
            node.left = insert_node node.left key
        else if key > node.key
            node.right = insert_node node.right key

        node.height = 1 + node.max_height
        let balance = node.left?.height - node.right?.height

        if balance > 1 {
            if key < node.left?.key {
                => right_rotate node
            } else {
                node.left = left_rotate node.left
                => right_rotate node
            }
        }

        if balance < -1 {
            if key > node.right?.key {
                => left_rotate node
            } else {
                node.right = right_rotate node.right
                => left_rotate node
            }
        }

        node
    }

    fn insert(key: int) {
        root = insert_node root key
    }

    fn left_rotate(y: Node) => Node {
        let x = y.right
        let T2 = x.left

        x.left = y
        y.right = T2
        
        y.height = 1 + max_height y.left y.right
        x.height = 1 + max_height x.left x.right

        x
    }

    fn right_rotate(y: Node) => Node {
        let x = y.left
        let T3 = x.right

        x.right = y
        y.left = T3

        y.height = 1 + max_height y.left y.right
        x.height = 1 + max_height y.left y.right

        x
    }
}

fn main {
    let!? tree = AVLTree{root: nil}
    tree.insert 1
    tree.insert 2
    tree.insert 3
}
