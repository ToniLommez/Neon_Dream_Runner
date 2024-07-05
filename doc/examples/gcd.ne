// gcd = greatest common divisor of two integer
fn gcd(a: int, b:int) => int {
    if b == 0 {
        a
    } else {
        gcd b (a%b)
    }
}

let num1 = 56
let num2 = 98
let result = gcd num1 num2
put "The GCD of " + num1:string + " and " + num2:string + " is " + result:string + "\n"