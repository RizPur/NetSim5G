## Project Structure

  netSim5G/              <-- project root
  ├── cmd/
  │   └── netsim5g/      <-- The executable's directory
  │       └── main.go
  ├── internal/
  │   ├── ue/
  │   ├── ran/
  │   └── core/
  └── go.mod

  Why the nesting? Because you might want multiple executables later:
  cmd/
  ├── netsim5g/      <-- Main simulator
  ├── cli/           <-- CLI control tool
  └── benchmarks/    <-- Performance testing tool

Each directory under cmd/ becomes a separate program you can build.


