#set page(paper: "a4", margin: 2cm)
#set text(font: "New Computer Modern", size: 11pt)
#set heading(numbering: "1.1")
#set par(justify: true)

#align(center)[
  #text(size: 24pt, weight: "bold")[Technical Report]
  #v(0.5cm)
  #text(size: 14pt)[A Benchmark Document for go-typst]
  #v(0.3cm)
  #text(size: 11pt, style: "italic")[February 2026]
]

#v(1cm)

#outline(indent: auto)

#pagebreak()

= Introduction

#lorem(150)

== Background

#lorem(200)

== Motivation

#lorem(100)

= Mathematical Foundations

The quadratic formula is given by:

$ x = (-b plus.minus sqrt(b^2 - 4a c)) / (2a) $

And Euler's identity:

$ e^(i pi) + 1 = 0 $

A more complex expression:

$ integral_0^infinity e^(-x^2) dif x = sqrt(pi) / 2 $

The general Stokes' theorem:

$ integral_(partial Omega) omega = integral_Omega dif omega $

== Linear Algebra

A matrix equation:

$ mat(a, b; c, d) vec(x, y) = vec(e, f) $

The determinant:

$ det(A) = sum_(sigma in S_n) "sgn"(sigma) product_(i=1)^n a_(i, sigma(i)) $

= Data Analysis

== Results Table

#table(
  columns: (1fr, 1fr, 1fr, 1fr),
  align: (left, center, center, center),
  table.header(
    [*Method*], [*Time (ms)*], [*Memory (MB)*], [*Accuracy*],
  ),
  [Baseline], [124.5], [48.2], [92.3%],
  [Optimized], [67.8], [32.1], [94.1%],
  [Advanced], [45.2], [28.7], [96.8%],
  [Proposed], [38.1], [24.3], [97.5%],
)

#lorem(80)

== Statistical Summary

#lorem(120)

= Implementation Details

== Architecture

#lorem(180)

== Code Example

```rust
fn fibonacci(n: u64) -> u64 {
    match n {
        0 => 0,
        1 => 1,
        _ => fibonacci(n - 1) + fibonacci(n - 2),
    }
}

fn main() {
    for i in 0..20 {
        println!("F({}) = {}", i, fibonacci(i));
    }
}
```

#lorem(100)

== Performance Characteristics

#lorem(150)

= Discussion

#lorem(200)

== Future Work

#lorem(150)

= Conclusion

#lorem(120)

#pagebreak()

#heading(numbering: none)[References]

+ A. Author, "First Paper Title," _Journal of Examples_, vol. 1, pp. 1--10, 2024.
+ B. Writer, "Second Paper Title," _Conference on Benchmarks_, pp. 100--110, 2025.
+ C. Researcher, "Third Paper Title," _Transactions on Performance_, vol. 42, no. 3, pp. 200--215, 2025.
+ D. Scholar et al., "Fourth Paper Title," _International Review_, vol. 7, pp. 50--65, 2026.
