package multiplexity

type ClientCommandHandler struct {
	config       *Config
	clientCmd    CommandChan
	clientMsgOut MessageChan
	serverMsgOut MessageChan
}
