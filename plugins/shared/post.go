package shared

import "github.com/hashicorp/go-plugin"

// PostbuildHandshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var PostbuildHandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BLOX_PLUGIN",
	MagicCookieValue: "postbuild",
}
