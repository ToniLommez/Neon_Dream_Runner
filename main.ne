// Definindo a função Fibonacci
fn fibonacci(n: int) => int {
    case n
        of 0 => 0
        of 1 => 1
    else => fibonacci (n - 1) + fibonacci (n - 2)
}

fn main {
    let n = 10
    let fib = fibonacci n
    print "Fibonacci of {} is {}", n', fib'
}

fn main {
    println "hehe"
}
