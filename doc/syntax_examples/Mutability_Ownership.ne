// Definindo uma função
fn add (x int, y int) => int {
    x + y
}

// Função que aceita uma referência mutável
fn inc(x int!) {
    x++
}

// Função que aceita uma referência imutável
fn foo(x int) => int {
    x
}

// Variavel global em função
let greeting = "Hello"
// A variavel sera enviada como copia, utilizando as mesmas regras de mutabilidade
fn greet(name string) => string {
    return "{} {}!", greeting, name
}

fn main {
    // Variavel mutavel
    let! x = 5
    inc x!  // envia x como mutavel
    inc x   // envia x como copia
    
    // Variavel imutavel
    let y = 3
    foo y   // envia y como imutavel
    inc y   // envia y como copia
    inc y!  // erro de compilação
    foo y!  // envia y como mutavel (mas nao faz diferenca)

    // Ternario
    let bigger = if x > y then x else y
    print "Bigger number: {byte}", bigger'

    // Ownership
    let result = add x' y' // x e y estao passando o ownership para result
    print "Result: {bool}" result' // result tambem sera destruido apos a chamada de print

    // Uso de Global
    print "{}", greet "World"
}

