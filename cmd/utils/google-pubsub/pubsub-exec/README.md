# pubsub-exec

A Google Cloud Pub/Sub-based remote command execution service that allows secure execution of commands and scripts through Pub/Sub messaging.

## Overview

`pubsub-exec` is a Go application that subscribes to a Google Cloud Pub/Sub topic and executes commands or scripts received as messages. It provides a secure way to remotely execute commands by:

1. **Simple Command Execution**: Execute predefined commands with arguments via JSON message
2. **Script Execution**: Execute shell scripts delivered as tarball payloads
3. **Status Reporting**: Publish execution results back to a status topic

## Architecture

The service operates on two Pub/Sub topics:

- **Execution Topic**: `{hostname}_exec` - Receives command execution requests
- **Status Topic**: `{hostname}_status` - Publishes execution results

## Message Formats

### 1. Simple Command Execution

Send a JSON array with the command name and arguments:

```json
["command_name", "arg1", "arg2", "arg3"]
```

**Example:**

```json
["vpn_start", "192.168.1.100"]
```

### 2. Script Execution (Tarball Payload)

Send a tarball containing a `script.sh` file. The tarball is extracted to a temporary directory and the script is executed from within the temporary directory.

**Tarball Structure:**

```
script.sh          # Main script (executable)
other_files/       # Supporting files
config.yaml
```

## Response Format

All executions publish a JSON response to the status topic:

```json
{
  "status": 0,
  "output": "Command output here..."
}
```

- `status`: Exit code (0 = success, non-zero = error)
- `output`: Standard output from the command

## Security Features

- **Command Validation**: Only allows alphanumeric characters, underscores, hyphens, and dots in command names
- **Path Sanitization**: Prevents directory traversal attacks in tarball extraction
- **Scoped Execution**: Limited to a predefined command specified via `-cmd` flag
- **Temporary Isolation**: Scripts execute in temporary directories that are automatically cleaned up

## Usage

### Installation

```bash
go build -o pubsub-exec ./cmd/utils/google-pubsub/pubsub-exec
```

### Configuration

The application uses the standard Go application configuration framework and requires Google Cloud Pub/Sub credentials.

### Running

```bash
./pubsub-exec -cmd=/path/to/your/command
```

**Required Flags:**

- `-cmd`: Path to the command executable to run

**Example:**

```bash
./pubsub-exec -cmd=/usr/local/bin/vpn_manager
```

### Environment Variables

Configure through the standard application configuration framework or environment variables:

- Google Cloud credentials
- Pub/Sub project settings
- Logging configuration

## How It Works

1. **Startup**:
   - Initializes Pub/Sub connection
   - Creates topics based on hostname (`{hostname}_exec`, `{hostname}_status`)
   - Subscribes to execution topic

2. **Message Processing**:
   - Attempts to parse message as JSON array (simple command)
   - If parsing fails, treats message as tarball payload
   - Validates command names for security
   - Executes command in specified working directory

3. **Script Execution Flow**:
   - Creates temporary directory
   - Extracts tarball contents
   - Makes `script.sh` executable (0700)
   - Executes with `script_exec` command
   - Cleans up temporary directory

4. **Response Publishing**:
   - Captures exit code and output
   - Publishes JSON response to status topic

## Example Use Cases

### VPN Management

```json
["vpn_start", "192.168.1.100"]
["vpn_stop", "192.168.1.100"]
```

### System Maintenance

```json
["system_update"]
["backup_database", "/path/to/backup"]
```

### Custom Script Execution

Send a tarball containing:

```
script.sh:
#!/bin/bash
echo "Hello from remote script!"
date
ls -la
```

## Error Handling

- **Invalid Commands**: Rejected with warning log
- **Tarball Errors**: Logged with temporary directory cleanup
- **Execution Failures**: Exit codes and errors published to status topic
- **Network Issues**: Handled by Pub/Sub client with retries

## Logging

Uses structured logging with zerolog:

- Info: Successful executions
- Warn: Non-fatal errors and validation failures
- Error: Critical failures
- Debug: Detailed message processing

## Dependencies

- Google Cloud Pub/Sub Go client
- zerolog for structured logging
- Standard Go libraries for archive/tar handling

## Security Considerations

- Commands are validated against strict regex patterns
- Tarball extraction prevents path traversal attacks
- All script execution occurs in isolated temporary directories
- Command scope is limited to the specified executable via `-cmd` flag
- No shell injection - commands are executed directly
