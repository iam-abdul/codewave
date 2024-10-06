# Bubble-table

<p>
  <a href="https://github.com/Evertras/bubble-table/releases"><img src="https://img.shields.io/github/release/Evertras/bubble-table.svg" alt="Latest Release"></a>
  <a href="https://pkg.go.dev/github.com/evertras/bubble-table/table?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
  <a href='https://coveralls.io/github/Evertras/bubble-table?branch=main'><img src='https://coveralls.io/repos/github/Evertras/bubble-table/badge.svg?branch=main&hash=abc' alt='Coverage Status'/></a>
  <a href='https://goreportcard.com/report/github.com/evertras/bubble-table'><img src='https://goreportcard.com/badge/github.com/evertras/bubble-table' alt='Go Report Card' /></a>
</p>

A customizable, interactive table component for the
[Bubble Tea framework](https://github.com/charmbracelet/bubbletea).

![Styled table](https://user-images.githubusercontent.com/5923958/188168029-0de392c8-dbb0-47da-93a0-d2a6e3d46838.png)

[View above sample source code](./examples/pokemon)

## Contributing

Contributions welcome, please [check the contributions doc](./CONTRIBUTING.md)
for a few helpful tips!

## Features

For a code reference of most available features, please see the [full feature example](./examples/features).
If you want to get started with a simple default table, [check the simplest example](./examples/simplest).

Displays a table with a header, rows, footer, and borders. The header can be
hidden, and the footer can be set to automatically show page information, use
custom text, or be hidden by default.

Columns can be fixed-width [or flexible width](./examples/flex). A maximum
width can be specified which enables [horizontal scrolling](./examples/scrolling),
and left-most columns can be frozen for easier reference.

Border shape is customizable with a basic thick square default. The color can
be modified by applying a base style with `lipgloss.NewStyle().BorderForeground(...)`.

Styles can be applied globally and to columns, rows, and individual cells.
The base style is applied first, then column, then row, then cell when
determining overrides. The default base style is a basic right-alignment.
[See the main feature example](./examples/features) to see styles and
how they override each other.

Styles can also be applied via a style function which can be used to apply
zebra striping, data-specific formatting, etc.

Can be focused to highlight a row and navigate with up/down (and j/k). These
keys can be customized with a KeyMap.

Can make rows selectable, and fetch the current selections.

Events can be checked for user interactions.

Pagination can be set with a given page size, which automatically generates a
simple footer to show the current page and total pages.

Built-in filtering can be enabled by setting any columns as filterable, using
a text box in the footer and `/` (customizable by keybind) to start filtering.

A missing indicator can be supplied to show missing data in rows.

Columns can be sorted in either ascending or descending order. Multiple columns
can be specified in a row. If multiple columns are specified, first the table
is sorted by the first specified column, then each group within that column is
sorted in smaller and smaller groups. [See the sorting example](examples/sorting)
for more information. If a column contains numbers (either ints or floats),
the numbers will be sorted by numeric value. Otherwise rendered string values
will be compared.

If a feature is confusing to use or could use a better example, please feel free
to open an issue.
