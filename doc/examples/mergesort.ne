// Implementação do Merge Sort em Neon
fn merge(left: [int], right: [int]) => [int] {
    let! result = [int]
    let l = 0
    let r = 0

    while l < left.length && r < right.length {
        if left[l] <= right[r] {
            result.Push left[l++]
        } else {
            result.Push right[r++]
        }
    }

    // Acrescentar quaisquer elementos restantes
    for l..left.length { result.Push left[l++] }
    for r..right.length { result.Push right[r++] }

    result
}

fn merge_sort(arr: [int]) => [int] {
    // Caso base
    if arr.length <= 1 then => arr

    let mid = arr.length / 2
    let left = arr[:mid]  // Pega a metade esquerda
    let right = arr[mid:]  // Pega a metade direita

    => merge(merge_sort left, merge_sort right)
}

fn main {
    let arr = [38, 27, 43, 3, 9, 82, 10]
    let sorted = merge_sort arr
    println "Sorted Array: {}" sorted
}
