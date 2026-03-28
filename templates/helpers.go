package templates

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

// LayoutWithContent wraps content inside the Layout component.
func LayoutWithContent(data LayoutData, content templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return Layout(data).Render(templ.WithChildren(ctx, content), w)
	})
}
