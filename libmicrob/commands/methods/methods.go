
package commands_methods

import (
	"github.com/synw/microb/libmicrob/datatypes"
)


func IsValid(command *datatypes.Command) bool {
	valid_commands := []string{"update_routes", "reparse_templates", "ping", "db_status", "routes", "state",
	"set"}
	is_valid := false
	for _, com := range(valid_commands) {
		if (com == command.Name) {
			is_valid = true
			break
		}
	}
	return is_valid
}
