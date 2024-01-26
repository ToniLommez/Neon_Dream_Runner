// Função de alta ordem recebendo uma função como parametro
fn increment_list(f fn(int, int) => int, lista int[], x int) => int[] {
    let nova int[]
    for e in lista {
        nova.Push f(e x)
    }
    nova
}

// composição de funções
fn double(x int!) => int { x * 2 }
fn increment(x int) => int { x + 1 }

fn main {
    let lista int[] = [1, 2, 3, 4, 5]
    let sum(x y) = x + y
    let result = increment_list sum lista 1
    println "Resultado: {int}" result'

    // Pipeline puro
    let name = " alice "
        |> trim
        |> capitalize
    // Ou inline
    let! name = " alice " |> trim |> capitalize
    name =  capitalize <| trim <| " alice " 

    let f = foo x |> n <0> <1> y <2>
    let g = foo x |> m y // vai ser colocado a direita por padrao

    // Composicao de funcao imutavel e reativa(com amarracao de pipeline)
    fn imutavel = double . increment
    fn reativa  = double | increment

    pulse before double { x++ }
    
    println "imutavel: {}" imutavel 5
    println "reativa:  {}" reativa  5
}
