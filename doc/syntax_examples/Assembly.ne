fn main() {
    let x = 5
    let y = 10
    println "{}", sum(x, y)
}

#assembly_type = x86
asm sum(x int, y int) => int {
    mov eax, x
    add eax, y
    ret
} x, y
