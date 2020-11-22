package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/labstack/echo/v4"
)

type DependencyReducerContext struct {
	Items []DependencyReducerItem `json:"items"`
}
type DependencyReducerItem struct {
	Key   string `json:"key"`
	Value []byte `json:"Payload"`
}

var contextKey = "x-dependency-reducer-context"

type DependencySetter struct {
	echoContext echo.Context
	reducerName string
}

type DependencyResolver struct {
	echoContext         echo.Context
	resolveContextNames []string
}

func (c *BuildContext) SetReducer(name string) *BuildContext {
	if c.echo == nil {
		panic("Call BuildRequestWithEcho method first")
	}
	c.dependencySetter = &DependencySetter{
		reducerName: name,
		echoContext: c.echo,
	}
	return c
}

func (c *BuildContext) SetResolver(resolveContextNames ...string) *BuildContext {
	if c.echo == nil {
		panic("Call BuildRequestWithEcho method first")
	}
	c.dependencyResolver = &DependencyResolver{
		echoContext:         c.echo,
		resolveContextNames: resolveContextNames,
	}
	return c
}

func (c *DependencySetter) SetContent(content interface{}) {
	if c.echoContext != nil && content != nil {
		c.echoContext.Set(c.reducerName, content)
	}
}

func (c *DependencyResolver) ResolveContent() (string, string) {
	var context DependencyReducerContext
	for _, name := range c.resolveContextNames {
		resolvedContext := c.echoContext.Get(name)
		if resolvedContext != nil {
			b, err := json.Marshal(resolvedContext)
			if err == nil {
				compressed, err := compressWithBrotli(b)
				if err == nil {
					context.Items = append(context.Items, DependencyReducerItem{
						Key:   name,
						Value: compressed,
					})
				}
			}
		}
	}
	contextByteArray, err := json.Marshal(context)
	if err != nil {
		return "", ""
	}
	return contextKey, string(contextByteArray)
}

func compressWithBrotli(value []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	compressor := brotli.NewWriterLevel(buf, 11)
	_, err := compressor.Write(value)
	compressor.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	return buf.Bytes(), nil
}
