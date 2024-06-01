package usecases

type AddChannelRequest struct {
	Name string
}

type DeleteChannelRequest struct {
	ChannelId string
}

type JoinChannelRequest struct {
	ChannelId string
}

type LeaveChannelRequest struct {
	ChannelId string
}
