package usecases

type AddChannelRequest struct {
	Name string
}

type DeleteChannelRequest struct {
	ChannelId uint64
}

type JoinChannelRequest struct {
	ChannelId uint64
}

type LeaveChannelRequest struct {
	ChannelId uint64
}
