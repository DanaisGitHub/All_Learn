# Codebase Analysis Notes

## Project Structure
This is a Go HTTP server learning project with the following components:

### Main Applications
- `main.go` - UDP client that continuously sends "Hello UDP" to port 42069
- `cmd/tcplistener/main.go` - TCP server listening on port 42069 that parses HTTP requests

### Internal Packages
- `internal/request/` - HTTP request parsing logic with state machine
- `internal/headers/` - HTTP header parsing and validation

## Key Components

### HTTP Request Parser (`internal/request/request.go`)
- Implements a state machine with states: `initialized`, `requestLine`, `header`, `body`, `done`
- Parses HTTP requests incrementally in 8-byte chunks
- Handles request line parsing (method, target, HTTP version)
- Only supports HTTP/1.1
- Includes panic recovery

### Headers Package (`internal/headers/headers.go`)
- Custom header implementation using `map[string]string`
- Validates header field names using regex for token compliance
- Supports multi-value headers with comma separation
- Case-insensitive header field access

### TCP Server (`cmd/tcplistener/main.go`)
- Listens for TCP connections on port 42069
- Accepts connections and parses HTTP requests using the request package
- Prints parsed request details

### UDP Client (`main.go`)
- Simple UDP client that sends "Hello UDP" messages to port 42069
- Runs in infinite loop

## Observations
1. The project appears to be learning HTTP protocol implementation from scratch
2. Uses 8-byte chunks for incremental parsing
3. Has some commented-out code in `tounderstand.txt` from earlier iterations
4. Includes basic tests using testify
5. The UDP client and TCP server use the same port (42069) which would cause conflicts

## Potential Issues
1. Port conflict between UDP client and TCP server (both use 42069)
2. Some error messages have typos ("couidn't" instead of "couldn't")
3. The request parser has complex state management that could be simplified
4. No actual HTTP response handling - only request parsing

## Recent Changes to Request Parser
The student has refactored the request parser with improvements:
- Created `appendAccumulator()` function for buffer management
- Added `moveState()` function for standardized state transitions
- Replaced `tempStr` with systematic `accumulator` pattern
- Still has mixed responsibilities in main loop (reading + parsing)
- EOF handling is commented out, indicating ongoing work
- Complex accumulator lifecycle management remains