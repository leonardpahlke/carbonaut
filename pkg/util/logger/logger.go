package logger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/fatih/color"
)

type SourceFileMode int

const (
	// Nop does nothing.
	Nop SourceFileMode = iota

	// ShortFile produces only the filename (for example main.go:69).
	ShortFile

	// LongFile produces the full file path (for example /home/frajer/go/src/myapp/main.go:69).
	LongFile
)

type Handler struct {
	groups []string
	attrs  []slog.Attr

	opts Options

	mu  *sync.Mutex
	out io.Writer
}

var DefaultOptions *Options = &Options{
	Level:       slog.LevelInfo,
	TimeFormat:  time.DateTime,
	SrcFileMode: ShortFile,
}

func GetLogLevel(logLevel string) slog.Level {
	switch logLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type Options struct {
	// Level reports the minimum level to log.
	// Levels with lower levels are discarded.
	// If nil, the Handler uses [slog.LevelInfo].
	Level slog.Leveler

	// The time format.
	TimeFormat string

	// SrcFileMode is the source file mode.
	SrcFileMode SourceFileMode
}

func Prefix(prefix string, msg string) string {
	return color.New(color.BgHiWhite, color.FgBlack).Sprint(prefix) + " " + msg
}

// NewHandler creates a new Handler.
func NewHandler(out io.Writer, opts *Options) *Handler {
	h := &Handler{out: out, mu: &sync.Mutex{}}
	if opts != nil {
		h.opts = *opts
	}
	return h
}

func (h *Handler) clone() *Handler {
	return &Handler{
		groups: h.groups,
		attrs:  h.attrs,
		opts:   h.opts,
		mu:     h.mu,
		out:    h.out,
	}
}

// Enabled implements slog.Handler.Enabled .
func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

// Handle implements slog.Handler.Handle .
//
//nolint:gocritic // This implements the slog.Handler.Handle func, hugeparam r slog.Record
func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	var bf bytes.Buffer

	if !r.Time.IsZero() {
		fmt.Fprint(&bf, color.New(color.Faint).Sprint(r.Time.Format(h.opts.TimeFormat)))
		fmt.Fprint(&bf, " ")
	}

	switch r.Level {
	case slog.LevelDebug:
		fmt.Fprint(&bf, color.New(color.BgCyan, color.FgHiWhite).Sprint("DEBUG"))
	case slog.LevelInfo:
		fmt.Fprint(&bf, color.New(color.BgGreen, color.FgHiWhite).Sprint("INFO "))
	case slog.LevelWarn:
		fmt.Fprint(&bf, color.New(color.BgYellow, color.FgHiWhite).Sprint("WARN "))
	case slog.LevelError:
		fmt.Fprint(&bf, color.New(color.BgRed, color.FgHiWhite).Sprint("ERROR"))
	}
	fmt.Fprint(&bf, " ")

	if r.PC != 0 {
		f, _ := runtime.CallersFrames([]uintptr{r.PC}).Next()

		switch h.opts.SrcFileMode {
		case Nop:
			break
		case ShortFile:
			fmt.Fprintf(&bf, "%s:%d ", filepath.Base(f.File), f.Line)
		case LongFile:
			fmt.Fprintf(&bf, "%s:%d ", f.File, f.Line)
		}
	}

	fmt.Fprint(&bf, color.HiWhiteString("| "))

	fmt.Fprint(&bf, r.Message)

	var attrs []slog.Attr
	attrs = append(attrs, h.attrs...)
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})

	for _, a := range attrs {
		fmt.Fprint(&bf, " ")
		for i, g := range h.groups {
			fmt.Fprint(&bf, color.New(color.FgCyan).Sprint(g))
			if i != len(h.groups) {
				fmt.Fprint(&bf, color.New(color.FgCyan).Sprint("."))
			}
		}
		fmt.Fprint(&bf, color.New(color.FgCyan).Sprintf("%s=", a.Key)+a.Value.String())
	}

	fmt.Fprint(&bf, "\n")

	h.mu.Lock()
	_, err := io.Copy(h.out, &bf)
	h.mu.Unlock()

	return err
}

// WithGroup implements slog.Handler.WithGroup .
func (h *Handler) WithGroup(name string) slog.Handler {
	h2 := h.clone()
	h2.groups = append(h2.groups, name)
	return h2
}

// WithAttrs implements slog.Handler.WithAttrs .
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h2 := h.clone()
	h2.attrs = append(h2.attrs, attrs...)
	return h2
}
