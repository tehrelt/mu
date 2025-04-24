package sl

import (
	"log/slog"

	"github.com/google/uuid"
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

func Query(q string) slog.Attr {
	return slog.Attr{
		Key:   "query",
		Value: slog.StringValue(q),
	}
}

func Args(args []interface{}) slog.Attr {
	return slog.Any("args", args)
}

func UUID(key string, id uuid.UUID) slog.Attr {
	return slog.Attr{
		Key:   key,
		Value: slog.StringValue(id.String()),
	}
}
