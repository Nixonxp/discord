package usecases

type AddChannelRequest struct {
	Name          string
	CurrentUserId string
}

type DeleteChannelRequest struct {
	ChannelId     string
	CurrentUserId string
}

type JoinChannelRequest struct {
	ChannelId     string
	CurrentUserId string
}

type LeaveChannelRequest struct {
	ChannelId     string
	CurrentUserId string
}
