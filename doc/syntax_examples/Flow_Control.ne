fn main {
    loop {
        print "Infinite loop\n"
        break // sai do loop
    }

    loop 50 {
        print "hehe"
    }

    loop 10..50 { // ou loop x..y
        print "Hehe"
    }

    let! x = 10
    while x > 0 {
        print "Countdown: {}\n", x
        x--
    }

    let! x = 10
    until x < 0 {
        print "Countdown: {}\n", x
        x--
    }

    let! x = 0
    do {
        print "Count: {}\n", x
        x++
    } while x != 10

    let! x = 0
    do {
        print "Count: {}\n", x
        x++
    } until x == 10

    let numbers = [1, 2, 3, 4, 5]
    for x: int, n in numbers { // uso de let opcional
        print "[{}]: {}\n", x, n
    }

    for let i = 10; i > 0; i /= 2 {
        print "{}\n", i
    }

    // Injecao de codigo

    pulse before hello_world { print "Hello" }
    pulse inside hello_world @space { print ", " }
    pulse after hello_world { println "World!" }

    fn hello_world = { @space }

    hello_world
}
