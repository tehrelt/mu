package sender

type SenderOption func(*manager)

func WithSender(sender Sender) SenderOption {
	return func(manger *manager) {
		manger.appendSender(sender)
	}
}
