// Because go:generate lines aren't shell commands they do not support * expansion. Hence, we repeat commands here.

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-wsrpc_out=. --go-wsrpc_opt=paths=source_relative telem.proto
//go:generate protoc --go_out=. --go_opt=paths=source_relative telem_automation_custom.proto
//go:generate protoc --go_out=. --go_opt=paths=source_relative telem_enhanced_ea.proto
//go:generate protoc --go_out=. --go_opt=paths=source_relative telem_enhanced_ea_mercury.proto
//go:generate protoc --go_out=. --go_opt=paths=source_relative telem_functions_request.proto
//go:generate protoc --go_out=. --go_opt=paths=source_relative telem_head_report.proto
package telem
