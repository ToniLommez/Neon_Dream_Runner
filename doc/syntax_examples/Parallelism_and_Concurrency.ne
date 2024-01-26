fn compute_a(val: int) {
    print "Number: {}\n", val
}

// in:  chan <!  any
// out: chan  !> any
// bi:  chan <!> any
fn compute_b(val: string, ch: chan <!> string) {
    ch <! "{} World", val
}

fn main {
    // Paralelismo sem canal
    !> compute_a 10
    print "Number: 11\n"

  	// Paralelismo com canal
  	let ch
	!> compute_b "hello", ch
	let result <! ch

    print "Result: {}", result
}
