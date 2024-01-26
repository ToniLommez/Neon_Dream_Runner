// DEFINIÇÃO INSTAVEL

obj Node {
    value: int
    next: Node?
}

obj Stack {
    head: Node?

    fn! push(value: int) { head = Node { value, next: head } }

    fn! pop => int {
        if head {
            let popped = head.value
            head = head.next
            => popped
        } else {
            => nil
        }
    }

    fn is_empty => bool {
        not head
    }
}

fn main {
    let! stack

    stack.push 1
    stack.push 2
    stack.push 3

    until stack.is_empty {
        print "Popped value: {}", stack.pop
    }
}
