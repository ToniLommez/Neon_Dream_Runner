obj ProcessNode {
    process: Process
    next: ProcessNode?

    get process, next
    set next
}

obj CircularQueue {
    head: ProcessNode?
    tail: ProcessNode?

    fn enqueue(process: Process) {
        let new_node = ProcessNode{process: process, next: nil}

        if head == nil {
            head = tail = new_node
            tail.next = head
        } else {
            tail.next = new_node
            tail = new_node
            tail.next = head
        }
    }

    fn dequeue => Process? {
        if head is nil then => nil

        if head is tail {
            let to_return = head.process
            head = tail = nil
            => to_return
        }

        let to_return = head.process
        head = head.next
        tail.next = head
        
        to_return
    }

    fn show_all {
        let start = head
        if start == nil {
            println "Queue is empty!"
            =>
        }

        let current = start
        do {
            println current.process
            current = current.next
        } while current != start
    }
}

fn main {
    let queue = CircularQueue{}

    let p1 = Process{pid: 1, name: "TextEditor"}
    let p2 = Process{pid: 2, name: "WebBrowser"}
    let p3 = Process{pid: 3, name: "FileExplorer"}

    queue.enqueue p1
    queue.enqueue p2
    queue.enqueue p3

    queue.show_all  // Mostra todos os processos na fila

    queue.dequeue  // Dequeue e executa o processo "TextEditor"
    queue.show_all  // Mostra os processos restantes na fila

    queue.enqueue p1  // Reinserindo "TextEditor" após a execução
    queue.show_all  // Mostra todos os processos na fila novamente
}
