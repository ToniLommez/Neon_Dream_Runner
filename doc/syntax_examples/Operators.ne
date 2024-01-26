fn main {
    let! x, y: int // inicializacao padrao em zero

    x = y // atribuição

    x += y // adicao
    x -= y // subtracao
    x *= y // multiplicacao
    x /= y // divisao
    x %= y // modulo

    x <<= y  // bitshift esquerda
    x >>= y  // bitshift direita
    x <<<= y // roundshift direita
    x >>>= y // roundshift esquerda

    x &= y  // logic and
    x ~&= y // logic nand
    x |= y  // logic or
    x ~|= y // logic nor
    x ^= y  // logic xor
    x ~^= y // logic xnor
}
