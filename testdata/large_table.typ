#set page(paper: "a4", margin: 1.5cm)
#set text(font: "New Computer Modern", size: 9pt)

#align(center)[
  #text(size: 18pt, weight: "bold")[Large Table Report]
  #v(0.3cm)
  #text(size: 11pt)[1000 Rows of Generated Data]
]

#v(0.5cm)

#table(
  columns: (auto, 1fr, 1fr, 1fr, 1fr, auto),
  align: (center, left, right, right, right, center),
  table.header(
    [*\#*], [*Name*], [*Value A*], [*Value B*], [*Ratio*], [*Status*],
  ),
  // Rows are generated via a loop
  ..for i in range(1000) {
    let name = if calc.rem(i, 5) == 0 { "Alpha" }
      else if calc.rem(i, 5) == 1 { "Bravo" }
      else if calc.rem(i, 5) == 2 { "Charlie" }
      else if calc.rem(i, 5) == 3 { "Delta" }
      else { "Echo" }
    let a = calc.rem(i * 37 + 13, 9973)
    let b = calc.rem(i * 53 + 7, 8887)
    let ratio = if b != 0 { calc.round(a / b * 100) / 100 } else { 0 }
    let status = if ratio > 1.0 { emoji.checkmark } else { "—" }
    (
      [#(i + 1)],
      [#name\-#(i + 1)],
      [#a],
      [#b],
      [#ratio],
      [#status],
    )
  }
)
