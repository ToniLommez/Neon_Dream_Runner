fn main {
    // slice
    let! a1 = [1, 2, 3]
    let! a2: [int] = [1, 2, 3]

    // array
    let! b2: [int: 3] = [1, 2, 3] 

    // tupla 
    let! c = (1, 2.0, "3") // imutavel

    // map
    let! x: |(int, int): float|
    let! y: |int: float|

    // obj
    let w = person{a: 1, b: 2, c: 3}
    let z = person{1, 2, 3}
}
