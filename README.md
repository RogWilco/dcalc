# Dimensional Calculator

A CLI and TUI tool for performing architectural calculations.

## CLI Interface

- Quick evaluation: run `dcalc "1-1/2 + 4-5/8"` (or `dcalc -- 1-1/2 + 4-5/8`) to parse and solve an expression without specifying a subcommand.
- Explicit evaluation: keep `dcalc eval <expression>` for scripting scenarios where clarity is preferred over brevity.
- Conversion utility: `dcalc convert 128 --from in --to ft-in` outputs `10' 8"`; default input unit is inches, output can be `ft-in`, `in`, or `dec-in`.
- Formatting controls: flags such as `--precision 1/16` or `--format feet-inches` tailor the textual result; these apply to both quick and explicit modes.
- Error ergonomics: on invalid input, print highlighted caret feedback and short hints (e.g., “add a space between whole number and fraction”).

```txt
dcalc [command] [args]

dcalc
├── <expression>               # Quick evaluation when a bare expression is provided
├── eval <expression>          # Explicit evaluation with clearer intent for scripts
└── convert [value]            # Convert between units
    ├── <value>                # Numeric input (defaults to inches)
    ├── --from <unit>          # Override input unit: in | ft-in | dec-in
    └── --to <unit>            # Requested output format: in | ft-in | dec-in

Global Options:
  --precision <fraction>       # Output precision (default 1/16)
  --format <style>             # Output style: fractional | feet-inches | decimal
  --json                       # Emit machine-readable JSON (planned v0.2.0)
  --help                       # Show help information
  --version                    # Print application version
  --debug                      # Verbose logging for troubleshooting
```

## Implementation Plan

### Dependencies

- Use `spf13/cobra` for command structure, argument parsing, and help/usage generation.
- Apply Charmbracelet libraries for presentation: `lipgloss` for styling, `glamour` for rich help text, and `log` for structured debug output.
- Build the future TUI on Charmbracelet’s `bubbletea` core with `bubbles` components, reusing `lipgloss` palettes for consistency between interfaces.

### Project Outline

- Build a shared calculator core that normalizes measurements (mixed numbers, fractional strings, decimals) into a `Measurement` backed by thousandths of an inch, supporting add/subtract/multiply/divide and feet↔inch conversion with metric planned.
- Expose the core through a flexible CLI: support bare expressions (`dcalc "2' 3-1/4\" + 5 7/8"`), specialized subcommands like `convert`, and future helpers such as `table`, with optional JSON output via `--json`.
- Plan a TUI (no CLI args) using Bubble Tea + Lip Gloss in a post-v0.1 release: show a tape history, input form, keyboard shortcuts (`a` add, `c` convert, `t` table), validation feedback, clipboard support, and quick unit switching.
- Keep front-ends thin by sharing formatting helpers (`FormatMeasurement`, `FormatFractionTable`) and structured error messages that tutor users on malformed inputs.
- Test the parser and arithmetic thoroughly; add formatting regression tests for the CLI and Bubble Tea model tests for TUI state transitions, covering edge cases like carries, negatives, and 1/32 precision.

### Core Functionality

- Mixed-number arithmetic: add, subtract, multiply, and divide measurements such as `5-3/8" + 4-11/16"` with automatic reduction to the nearest 1/16 or finer when needed.
- Unit conversion: translate between inches, feet/inches, and decimal inches (e.g., `128" -> 10' 8"`), with hooks for future metric conversion (millimeters, centimeters).
- Expression parsing: accept expressions mixing operators, parentheses, and unit suffixes (`ft`, `"`, `'`) so users can perform multi-step calculations in one command.
- Fraction utilities: normalize improper fractions, simplify to lowest terms, and optionally snap to common job-site fractions (1/2, 1/4, 1/8, 1/16, 1/32).
- Result formatting: present results as fractional inches, feet/inches, or decimals, and allow CLI selection via flags or TUI hotkeys.
- Calculation history: maintain a running tape in the TUI to review past steps, label them, and copy results for reuse.

### Parsing & Formatting Rules

- Accept inputs with optional whitespace between components (`4 1/2`, `4-1/2`, `4  -  1/2`); normalize hyphenated and spaced forms to a common mixed number.
- Support whole numbers, proper fractions (`5/8`), mixed numbers (`4 5/8` or `4-5/8`), and decimal inches (`3.125`).
- Recognize unit markers: inches (`"` or `in`), feet (`'` or `ft`), and feet/inch combos (`5' 4-1/2"`). Bare numbers default to inches.
- Apply sign to the entire measurement when prefixed (`-4-1/2`, `-(3 3/4)`, or `-2' 1"`); internal fractions must remain positive.
- Allow parentheses and standard operators (`+`, `-`, `*`, `/`) with left-to-right evaluation precedence respecting parentheses.
- Clamp precision by snapping to the configured fraction denominator (1/16 default); preserve exact values until final formatting.
- Reduce fractions to simplest terms and carry extra whole inches into the feet component when formatting in feet/inches.
- Track a unit exponent alongside the numeric value so length÷length collapses to a dimensionless result; dimensionless outputs omit unit glyphs.
- When a dimensionless result comes directly from dividing two measurements, also derive an integer quotient and leftover measurement for layout scenarios (e.g., board counts).
- On parse errors, report the offending segment with a caret indicator and guidance (missing unit, malformed fraction, unknown token).

### Testing Strategy

- Assert parser coverage across mixed numbers, decimals, unit markers, whitespace variants, negatives, and parentheses.
- Validate arithmetic for addition/subtraction/multiplication/division using known carpenter scenarios, ensuring precision snapping behaves as configured.
- Verify dimensional bookkeeping: length÷length yields dimensionless ratios and count+remainder layouts; length÷scalar retains units; scalar÷length errors.
- Exercise formatting outputs for fractional, decimal, and feet/inch styles, including carry/borrow and remainder rendering when applicable.
- Use table-driven tests for parser scenarios where multiple input forms map to shared expectations.

### Module Layout (planned)

```txt
cmd/
└── dcalc/
    └── main.go              # Cobra root; wires commands, global flags, and configuration

internal/
├── calculator/
│   ├── measurement.go       # Core types (Measurement, UnitExponent) and normalization helpers
│   ├── parser.go            # Expression lexer/parser converting strings into AST/Measurements
│   ├── eval.go              # Arithmetic operations, precision snapping, quotient+remainder logic
│   └── calculator_test.go   # Unit tests covering parsing and arithmetic scenarios
├── format/
│   ├── formatter.go         # Rendering helpers for fractional, feet-inches, decimal, JSON
│   └── formatter_test.go    # Output expectation snapshots/regression tests
├── cli/
│   ├── root.go              # Root command setup (bare expression handling, shared flags)
│   ├── eval.go              # Explicit eval command (mainly for scripts/pipelines)
│   ├── convert.go           # Convert command implementation
│   └── cli_test.go          # Lightweight integration tests using Cobra command execution
├── tui/                     # (v0.2+) Bubble Tea implementation when no args are provided
│   ├── model.go             # Bubble Tea model, update loop, and key handling
│   ├── view.go              # Lipgloss-rendered views, tape history, layout output
│   ├── components/          # Reusable bubbles (inputs, tables) and styling definitions
│   └── tui_test.go          # Model/state transition tests
└── meta/
    └── version.go           # Version constant/helper (shared internally across front-ends)
```

## Roadmap

- **v0.1.0**
  - Core measurement parser with add/subtract support and 1/16 precision.
  - CLI bare-expression evaluation plus `convert` subcommand for inches↔feet/inches.
  - Shared formatting helpers to render fractional and feet/inch outputs.
  - Foundational unit tests for parsing, arithmetic, and formatting helpers.
- **v0.2.0**
  - Expand arithmetic to multiply/divide and introduce fraction table generation.
  - Add CLI JSON output, rounding controls, and initial Bubble Tea TUI shell.
  - Flesh out error messaging with recovery hints and short documentation examples.
- **v0.3.0**
  - Introduce project state persistence, named calculations, and export options.
  - Add metric conversions and configurable snapping rules.
  - Begin plugin-style specialty calculators (e.g., stair rise/run).

## Ideas Parking Lot

- Make the precision level configurable (default to the nearest 1/16, allow for 1/32, 1/64, etc.)
- When a calculated result exceeds the precision level, snap to the nearest valid value (e.g., 1/32 -> 1/16, 1/64 -> 1/32).
- When a calculated result is snapped, indicate the reduced precision by coloring the number in a distinct color.
- Support loading a tape from a file (JSON, CSV, etc.) or STDIN.
- Support writing a tape to a file (JSON, CSV, etc.) or STDOUT.
