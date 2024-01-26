fn main {
    // Lambda
    let f(x: int, y: int) => int, int {
        let doubleX = x * 2
        let doubleY = y * 2
        doubleX, doubleY
    }

    // Poderá ser declarado como f(3 5) ou f(3, 5), mas isso é açucar sintatico
    let x, y = f 3 5 
    print "resultado = {}, {}" x y

    // Lambda sem amarração
    let z = (x) => int {
        let y = x*2
        if y is 3 
            => 5
        y*2
    } 3 // passagem de argumento

    // Lambda funcional
    let f(x) = x + 3
    print "{} {}" (f 1) (f 2)

    // Lambdas funcionais sao apenas para aproximar a sintaxe
    // de uma linguagem funcional pura, os tipos serao implicitados
    let f(x) = {
        let y = x*2
        if y is 3 
            => 5
        y*2
    }
}
