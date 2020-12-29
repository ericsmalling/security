package security

import (
	"fmt"
)

type Config struct {
	// DefaultMode sets the default execution policy for all other commands. It is recommended to set this to "disable"
	// if for restricted setups to avoid accidentally allowing new features coming in with version upgrades.
	DefaultMode ExecutionPolicy `json:"defaultMode" yaml:"defaultMode" default:"allow"`

	// ForceCommand behaves similar to the OpenSSH ForceCommand option. When set this command overrides any command
	// requested by the client and executes this command instead. The original command supplied by the client will be
	// set in the `SSH_ORIGINAL_COMMAND` environment variable.
	//
	// Setting ForceCommand changes subsystem requests into exec requests for the backends.
	ForceCommand string `json:"forceCommand" yaml:"forceCommand"`

	// Env controls whether to allow or block setting environment variables.
	Env EnvConfig `json:"env" yaml:"env"`
	// Command controls whether to allow or block command ("exec") requests via SSh.
	Command CommandConfig `json:"command" yaml:"command"`
	// Shell controls whether to allow or block shell requests via SSh.
	Shell ShellConfig `json:"shell" yaml:"shell"`
	// Subsystem controls whether to allow or block subsystem requests via SSH.
	Subsystem SubsystemConfig `json:"subsystem" yaml:"subsystem"`

	// TTY controls how to treat TTY/PTY requests by clients.
	TTY TTYConfig `json:"tty" yaml:"tty"`

	// Signal configures how to handle signal requests to running programs.
	Signal SignalConfig `json:"signal" yaml:"signal"`

	// MaxSessions drives how many session channels can be open at the same time for a single network connection.
	MaxSessions uint `json:"maxSessions" yaml:"maxSessions"`
}

// Validate validates a shell configuration
func (c Config) Validate() error {
	if err := c.DefaultMode.Validate(); err != nil {
		return fmt.Errorf("invalid defaultMode configuration (%w)", err)
	}
	if err := c.Env.Validate(); err != nil {
		return fmt.Errorf("invalid env configuration (%w)", err)
	}
	if err := c.Command.Validate(); err != nil {
		return fmt.Errorf("invalid command configuration (%w)", err)
	}
	if err := c.Shell.Validate(); err != nil {
		return fmt.Errorf("invalid shell configuration (%w)", err)
	}
	if err := c.Subsystem.Validate(); err != nil {
		return fmt.Errorf("invalid subsystem configuration (%w)", err)
	}
	if err := c.TTY.Validate(); err != nil {
		return fmt.Errorf("invalid TTY configuration (%w)", err)
	}
	if err := c.Signal.Validate(); err != nil {
		return fmt.Errorf("invalid signal configuration (%w)", err)
	}
	return nil
}

// EnvConfig configures setting environment variables.
type EnvConfig struct {
	// Mode configures how to treat environment variable requests by SSH clients.
	Mode ExecutionPolicy `json:"mode" yaml:"mode" default:""`
	// Allow takes effect when Mode is ExecutionPolicyFilter and only allows the specified environment variables to be
	// set.
	Allow []string
	// Allow takes effect when Mode is not ExecutionPolicyDisable and disallows the specified environment variables to
	// be set.
	Deny []string
}

// Validate validates a shell configuration
func (e EnvConfig) Validate() error {
	if err := e.Mode.Validate(); err != nil {
		return fmt.Errorf("invalid mode (%w)", err)
	}
	return nil
}

// CommandConfig controls command executions via SSH (exec requests).
type CommandConfig struct {
	// Mode configures how to treat command execution (exec) requests by SSH clients.
	Mode ExecutionPolicy `json:"mode" yaml:"mode" default:""`
	// Allow takes effect when Mode is ExecutionPolicyFilter and only allows the specified commands to be
	// executed. Note that the match an exact match is performed to avoid shell injections, etc.
	Allow []string
}

// Validate validates a shell configuration
func (c CommandConfig) Validate() error {
	if err := c.Mode.Validate(); err != nil {
		return fmt.Errorf("invalid mode (%w)", err)
	}
	return nil
}

// ShellConfig controls shell executions via SSH.
type ShellConfig struct {
	// Mode configures how to treat shell requests by SSH clients.
	Mode ExecutionPolicy `json:"mode" yaml:"mode" default:""`
}

// Validate validates a shell configuration
func (s ShellConfig) Validate() error {
	if err := s.Mode.Validate(); err != nil {
		return fmt.Errorf("invalid mode (%w)", err)
	}
	return nil
}

// SubsystemConfig controls shell executions via SSH.
type SubsystemConfig struct {
	// Mode configures how to treat subsystem requests by SSH clients.
	Mode ExecutionPolicy `json:"mode" yaml:"mode" default:""`
	// Allow takes effect when Mode is ExecutionPolicyFilter and only allows the specified subsystems to be
	// executed.
	Allow []string
	// Allow takes effect when Mode is not ExecutionPolicyDisable and disallows the specified subsystems to be executed.
	Deny []string
}

// Validate validates a subsystem configuration
func (s SubsystemConfig) Validate() error {
	if err := s.Mode.Validate(); err != nil {
		return fmt.Errorf("invalid mode (%w)", err)
	}
	return nil
}

// TTYConfig controls how to treat TTY/PTY requests by clients.
type TTYConfig struct {
	// Mode configures how to treat TTY/PTY requests by SSH clients.
	Mode ExecutionPolicy `json:"mode" yaml:"mode" default:""`
}

// Validate validates the TTY configuration
func (t TTYConfig) Validate() error  {
	if err := t.Mode.Validate(); err != nil {
		return fmt.Errorf("invalid mode (%w)", err)
	}
	return nil
}

// SignalConfig configures how signal forwarding requests are treated.
type SignalConfig struct {
	// Mode configures how to treat signal requests to running programs
	Mode ExecutionPolicy `json:"mode" yaml:"mode" default:""`
	// Allow takes effect when Mode is ExecutionPolicyFilter and only allows the specified signals to be forwarded.
	Allow []string
	// Allow takes effect when Mode is not ExecutionPolicyDisable and disallows the specified signals to be forwarded.
	Deny []string
}

// Validate validates the signal configuration
func (s SignalConfig) Validate() error  {
	if err := s.Mode.Validate(); err != nil {
		return fmt.Errorf("invalid mode (%w)", err)
	}
	return nil
}

// ExecutionPolicy drives how to treat a certain request.
type ExecutionPolicy string

const (
	// ExecutionPolicyUnconfigured falls back to the default mode. If unconfigured on a global level the default is to
	// "allow".
	ExecutionPolicyUnconfigured ExecutionPolicy = ""

	// ExecutionPolicyEnable allows the execution of the specified method unless the specified option matches the
	// "deny" list.
	ExecutionPolicyEnable ExecutionPolicy = "enable"

	// ExecutionPolicyFilter filters the execution against a specified allow list. If the allow list is empty or not
	// supported this ootion behaves like "disable".
	ExecutionPolicyFilter ExecutionPolicy = "filter"

	// ExecutionPolicyDisable disables the specified method and does not take the allow or deny lists into account.
	ExecutionPolicyDisable ExecutionPolicy = "disable"
)

// Validate validates the execution policy.
func (e ExecutionPolicy) Validate() error {
	switch e {
	case ExecutionPolicyUnconfigured:
	case ExecutionPolicyEnable:
	case ExecutionPolicyFilter:
	case ExecutionPolicyDisable:
	default:
		return fmt.Errorf("invalid mode: %s", e)
	}
	return nil
}
