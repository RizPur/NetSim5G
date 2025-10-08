# NetSim5G - Progress Tracker

## Session 1 - 2025-10-08

### Completed Today ✅
- Created Go project structure with proper layout
  - `cmd/netsim5g/` - entry point
  - `internal/ue/` - user equipment package
  - `internal/ran/` - radio access network package
  - `internal/core/` - core network functions package
- Created UE struct with state management (Disconnected, Connected, Idle)
- Moved main.go to correct location (`cmd/netsim5g/main.go`)
- Designed simplified connection flow:
  1. Check if UE's IMSI is in gNodeB's allowed list
  2. Check if gNodeB has capacity
  3. Establish connection

### Key Concepts Discussed
- **Go project structure**: `cmd/`, `internal/`, `pkg/` conventions
- **Slice vs Map**:
  - Slice = ordered collection (like Python list) - O(n) lookup
  - Map = key-value pairs (like Python dict) - O(1) lookup
  - Decision: Use maps for IMSI lookups (both allowed and connected)
- **Separation of concerns**: Two different maps needed:
  - Allowed IMSIs map (authorization)
  - Connected UEs map (current state)

### Design Decisions
- gNodeB needs file with allowed IMSIs (temporary - will move to core network later)
- gNodeB struct needs:
  - Unique ID
  - Max capacity
  - Map of allowed IMSIs
  - Map of currently connected UEs

### Tomorrow's Tasks 🎯

**1. Create UE Package** (`internal/ue/ue.go`)
   - Move UE struct from main.go to ue package
   - Move UEState type and constants
   - Add any helper methods needed for connection

**2. Create gNodeB Package** (`internal/ran/gnodeb.go`)
   - Create gNodeB struct with fields:
     - ID (int or string?)
     - MaxCapacity (int)
     - AllowedIMSIs (map[string]bool)
     - ConnectedUEs (map[string]???) - decide what to store as value
   - Create methods:
     - Check if IMSI is allowed
     - Check capacity
     - Accept connection
     - (Optional) Read allowed IMSIs from file

**3. Update Main** (`cmd/netsim5g/main.go`)
   - Import both packages
   - Create a gNodeB instance
   - Create a UE instance
   - Demo: UE attempts to connect to gNodeB
   - Show the output of the connection process

**Questions to Answer Tomorrow:**
- What should the ConnectedUEs map value be? Just `bool`, or the actual `*UE` object?
- Should gNodeB ID be int or string?
- Constructor functions: `NewUE()` and `NewGNodeB()` - when to use pointers?

### Next Steps (Future Sessions)
- Read allowed IMSIs from a file
- Add multiple UEs connecting
- Handle connection rejection scenarios
- Begin core network functions (AMF, SMF, UPF, UDM)
