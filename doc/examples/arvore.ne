obj node {
    let value: int
    let left: node
    let right: node

    fn init(val int) {
        value = val
    }

    fn insert(val int) {
        if val < value {
            if left is nil {
                left = node
                left.init val
            } else {
                left.insert val
            }
        } else {
            if right is nil {
                right = node
                right.init val
            } else {
                right.insert val
            }
        }
    }

    fn contains(val int) => bool {
        if val == value {
            true
        } else if val < value {
            if left is nil then false else left.contains val
        } else {
            if right is nil then false else right.contains val
        }
    }

    fn print_in_order {
        if left != nil then left.print_in_order
        print "{}", value
        if right != nil then right.print_in_order
    }
}

fn main {
    let! root node
    root.init 5
    root.insert 2
    root.insert 7
    root.insert 1
    root.insert 3
    root.insert 6
    root.insert 8
    root.print_in_order
}