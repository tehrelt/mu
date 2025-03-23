package sl

import (
	"log/slog"
)

func Err(e error) slog.Attr {
	return slog.String("err", e.Error())
}

func Module(module string) slog.Attr {
	return slog.Attr{
		Key:   "module",
		Value: slog.StringValue(module),
	}
}

func Method(method string) slog.Attr {
	return slog.Attr{
		Key:   "method",
		Value: slog.StringValue(method),
	}
}
