#set page(paper: "a4", margin: 1cm)
#set text(font: "New Computer Modern", size: 7pt)

#align(center)[
  #text(size: 16pt, weight: "bold")[Stress Test: ~1000 Page Table]
]

#v(0.3cm)

// ~52 rows per page at 7pt with 1cm margins.
// 52000 rows ≈ 1000 pages.
#table(
  columns: (auto, 1fr, 1fr, 1fr, 1fr, 1fr, auto),
  align: (center, left, right, right, right, right, center),
  table.header(
    [*\#*], [*Name*], [*Col A*], [*Col B*], [*Col C*], [*Ratio*], [*Status*],
  ),
  ..for i in range(52000) {
    let name = if calc.rem(i, 7) == 0 { "Alpha" }
      else if calc.rem(i, 7) == 1 { "Bravo" }
      else if calc.rem(i, 7) == 2 { "Charlie" }
      else if calc.rem(i, 7) == 3 { "Delta" }
      else if calc.rem(i, 7) == 4 { "Echo" }
      else if calc.rem(i, 7) == 5 { "Foxtrot" }
      else { "Golf" }
    let a = calc.rem(i * 37 + 13, 9973)
    let b = calc.rem(i * 53 + 7, 8887)
    let c = calc.rem(i * 71 + 31, 7919)
    let ratio = if b != 0 { calc.round(a / b * 100) / 100 } else { 0 }
    let status = if ratio > 1.5 { emoji.checkmark.heavy }
      else if ratio > 1.0 { "~" }
      else { "—" }
    (
      [#(i + 1)],
      [#name\-#(i + 1)],
      [#a],
      [#b],
      [#c],
      [#ratio],
      [#status],
    )
  }
)
