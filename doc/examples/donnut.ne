let PI = 3.14159265358979323846
let DOUBLE_PI = 6.283185307179586232
let PI_HALF = 1.570796326794896558
let THREE_PI_HALF = 4.712388980384689674

fn fmod(x: float, y: float) => float {
    if y == 0 {
        0
    } else {
        let result:float = x - ((x/y):int):float * y
        if result < 0 { 
            result + y 
        } else { 
            result 
        }
    }
}

fn pow(base: float, exponent: float) => float {
    let! result:float = 1
    let! i
    for i = 0; i < exponent ; i += 1 {
        result = result*base
    }
    result
}

fn sin(x: float) => float {
    let! z:float = 0
    z = fmod x (2*3.14159265358979323846)
    let! result:float = 0
    if z <= 1.570796326794896558 {
        result = 0.028644811879990944 * (pow z 4) - 0.20385110262323103 * (pow z 3) + 0.02090994802843416 * (pow z 2) + 0.9954617512319969 * z + 0.00023060949853418496
    } else {
        if z <= 3.14159265358979323846 {
            result = 0.028643221500985257 * (pow z 4) - 0.15610299471228117 * (pow z 3) - 0.20406096277925495 * (pow z 2) + 1.3562229703585635 * z - 0.1964080629634316
        } else {
            if z <= 4.712388980384689674 {
                result = -0.028646247913692567 * (pow z 4) + 0.5638421123687306 * (pow z 3) - 3.638654926239317 * (pow z 2) + 8.725053225623561 * z - 6.190950354465874
            } else {
                result = -0.028646283653384758 * (pow z 4) + 0.516096855586553827 * (pow z 3) - 2.963668250249625657 * (pow z 2) + 5.536807811523057232 * z - 1.159489388771234974
            }
        }
    }
    result
}

fn cos(x:float) => float {
    sin (PI_HALF - x)
}

let coreString = [string][".", ",", "-", "~", ":", ";", "=", "!", "*", "#", "$", "@"]

let! A:float = 0
let! B:float = 0

let! i: float = 0
let! j: float = 0
let! k: int = 0

let size = 1760
let! z = [float:size]
let! b = [string:size]

let! c:float = 0
let! d:float = 0
let! e:float = 0
let! f:float = 0
let! g:float = 0
let! h:float = 0
let! D:float = 0
let! l:float = 0
let! m:float = 0
let! n:float = 0
let! t:float = 0
let! x:int = 0
let! y:int = 0
let! o:int = 0
let! N:int = 0
let! point:int = 0
let! v:string = ""

clear

while true {
    let! tmp = 0
    for tmp = 0; tmp < size; tmp += 1 {
        z[tmp] = 0
        b[tmp] = " "
    }

    for j = 0.000001; j < 6.28; j += 0.8 {
        for i = 0.000001; i < 6.28; i += 0.6 {
            c = sin i
            d = cos j
            e = sin A
            f = sin j
            g = cos A
            h = d + 2
            D = 1 / (c*h*e + f*g + 5)
            l = cos i
            m = cos B
            n = sin B
            t = c*h*g - f*e

            x = (40 + 30*D*(l*h*m-t*n)):int
            y = (12 + 15*D*(l*h*n+t*m)):int

            o = (x + 80*y):int

            N = (8 * ((f*e-c*d*g)*m - c*d*e - f*g - l*d*n)):int

            if y < 22 && y > 0 && x > 0 && x < 80 && D > z[o] {
                z[o] = D
                point = 0
                if N > 0 {
                    point = N
                }
                b[o] = coreString[point]
            }
        }
    }

    put "\x1b[H"

    for k = 0; k < 1761; k += 1 {
        v = "\n"
        if (k % 80) > 0 {
            v = b[k]
        }
        put v

        A += 0.004
        B += 0.002
    }
}
