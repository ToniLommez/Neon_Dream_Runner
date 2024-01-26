// Exemplo simples de pattern matching
fn simple_pattern_matching(val: any) {
    case val
        of int => print "É um inteiro."
        of string => print "É uma string."
        of bool => print "É um booleano."
    else => print "Tipo desconhecido."
}

// Exemplo de pattern matching com múltiplas condições para o mesmo caso
fn multiple_condition_pattern_matching(val: int) {
    case val
        of 0 | 1 => print "É zero ou um."
        of 2 | 3 => print "É dois ou três."
    else => print "Outro número."
}

// Exemplo de pattern matching com ranges
fn range_pattern_matching(val: int) {
    case val
        of 0..9 => print "É um dígito."
        of 10..99 => print "É um número de dois dígitos."
    else => print "É um número com três ou mais dígitos."
}

// Exemplo de pattern matching com listas
fn list_pattern_matching(val: any[]) {
    case val
        of [] => print "A lista está vazia."
        of [head] => print "A lista tem um elemento: {}", head
        of [head, tail...] => print "A lista tem vários elementos. O primeiro é: {}", head
    else => print "Não é uma lista."
}

fn seila(x: any?, y: float) => int? {
    case x
        of nil -> print "é null"
        of int if x > 3 | float -> print "é numero e tres"
        of float | string -> print "valor = {}" x
        of int -> print "é so int"
    else => 3
}
