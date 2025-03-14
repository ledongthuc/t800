# T800 - AI-Powered Combat Robot System

![Image](https://github.com/user-attachments/assets/a9680993-a252-47dd-bd2a-fe31e1bf71fc)

[![asciicast](https://asciinema.org/a/708059.svg)](https://asciinema.org/a/708059)

T800 is an advanced combat robot system that uses AI for decision-making and threat management. The system integrates with Ollama for AI-powered decision making, providing sophisticated combat and defensive capabilities.

## Quick Start

### Prerequisites
- Go 1.21 or later
- Ollama installed and running locally
- Git

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/t800.git
cd t800
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables (optional):
```bash
export OLLAMA_BASE_URL="http://localhost:11434"  # Default Ollama URL
export OLLAMA_MODEL="llama3.2"                   # Default model
```

4. Run the system:
```bash
go run main.go
```

## Technical Details

### System Architecture

The T800 system is built with a modular architecture:

```
t800/
├── internal/
│   ├── ai/          # AI decision-making system
│   ├── anatomy/     # Robot physical structure
│   ├── common/      # Shared types and utilities
│   ├── defense/     # Defensive strategies
│   ├── monitoring/  # System monitoring and logging
│   ├── offense/     # Offensive capabilities
│   ├── processor/   # Main system processor
│   └── scanner/     # Threat detection system
├── pkg/             # Public packages
└── main.go          # Application entry point
```

### Key Components

1. **AI Decision Maker**
   - Integrates with Ollama for AI-powered decision making
   - Handles combat decisions and threat engagement
   - Configurable through environment variables

2. **Robot Anatomy**
   - Manages physical components (head, body, arms, legs)
   - Handles health monitoring and damage calculation
   - Supports critical part identification

3. **Defense System**
   - Implements defensive strategies
   - Manages shield activation and evasive maneuvers
   - Handles critical system protection

4. **Offense System**
   - Controls weapon systems
   - Manages attack strategies
   - Handles weapon selection and targeting

5. **Scanner System**
   - Performs threat detection
   - Manages threat tracking
   - Implements threat prediction

### AI Integration

The system uses Ollama for AI decision-making with the following features:
- Combat decision making
- Threat engagement evaluation
- Strategic planning
- Response confidence scoring

## Development Setup

### Development Environment

1. **IDE Setup**
   - Recommended: GoLand or VS Code with Go extension
   - Required extensions:
     - Go
     - Git
     - JSON

2. **Code Style**
   - Follow Go standard formatting
   - Use `go fmt` for code formatting
   - Follow Go best practices and idioms

3. **Testing**
   ```bash
   # Run all tests
   go test ./...

   # Run tests with coverage
   go test ./... -cover
   ```

### Adding New Features

1. **Adding New Weapons**
   - Add weapon definition in `internal/offense/actions.go`
   - Update weapon damage values in `internal/processor/processor.go`
   - Add weapon to available weapons list

2. **Adding New Defensive Strategies**
   - Create new strategy in `internal/defense/actions.go`
   - Register strategy in `internal/defense/strategy.go`
   - Update strategy priorities as needed

3. **Modifying AI Behavior**
   - Update prompts in `internal/ai/decision.go`
   - Modify decision structures as needed
   - Test with different scenarios

### Debugging

1. **Logging**
   - System logs are handled by `internal/monitoring/logger.go`
   - Log levels: INFO, WARNING, ERROR
   - Includes timestamps and context

2. **Health Monitoring**
   - Monitor system health through logs
   - Check component status
   - Verify threat detection

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Ollama for providing the AI capabilities
- The Go community for excellent tools and libraries
- Contributors and maintainers 
