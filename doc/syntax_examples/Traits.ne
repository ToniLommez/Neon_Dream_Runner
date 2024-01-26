obj Dog {
    name: string

    fn speak => string {
        return "Bark!"
    }
}

trait Speaker {
    fn speak => string
}

fn makeItSpeak(s Speaker) {
    print s.speak
}

fn main {
    let! d Dog
    d.name = "Rover"
    makeItSpeak d
}