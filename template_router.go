package mymux

import (
	"regexp"
	"net/http"
	"strings"
	"context"
	"net/url"
	"bytes"
	"fmt"
)

type Vars map[string]string

type key int

const paramsKey key = 0

const (
	RouteSignatureParam = "$route_signature"
)

var routeVariable = regexp.MustCompile(`[^/]*:[^/]*`)

type RouterTemplateHandler struct {
	RouteHandler
	types map[string]string
}

type routePoints struct {
	varStart  int
	varEnd    int
	typeStart int
	typeEnd   int
}

type routeTemplate struct {
	pathSignature string
	template      string
	vars          []routePoints
	method        string
	pattern       *regexp.Regexp
	handler       http.Handler
}

func newVars(ctx context.Context, vars Vars) context.Context {
	return context.WithValue(ctx, paramsKey, vars)
}

func GetVars(r *http.Request) Vars {
	ctx := r.Context()
	vars, ok := ctx.Value(paramsKey).(Vars)
	if !ok {
		return nil
	}
	return vars
}

func NewRouterTemplateHandler() *RouterTemplateHandler {
	h := &RouterTemplateHandler{RouteHandler: *NewRouteHandler(), types: map[string]string{}}
	return h
}


func (handler *RouterTemplateHandler) RegisterType(typeName string, regex string) {
	regexp.MustCompile(regex)
	handler.types[typeName] = regex
}

func (handler *RouterTemplateHandler) AppendRoute(method string, path string, userHandler http.Handler) Route {
	route := handler.newRouteTemplate(method, path, userHandler)
	handler.routes = append(handler.routes, route)
	return route
}

func (handler *RouterTemplateHandler) newRouteTemplate(method string, path string, userHandler http.Handler) *routeTemplate {
	buffer := bytes.Buffer{}

	vars := routeVariable.FindAllStringIndex(path, -1)
	pointList := make([]routePoints, len(vars))
	buffer.WriteString(`^`)
	k := 0
	for i, v := range vars {
		buffer.WriteString(path[k:v[0]])
		t := strings.Split(path[v[0]:v[1]], ":")
		pointList[i] = routePoints{
			varStart:  v[0],
			varEnd:    v[0] + len(t[0]),
			typeStart: v[0] + len(t[0]) + 1,
			typeEnd:   v[1],
		}
		typeName := t[1]

		expr, ok := handler.types[typeName]
		if !ok {
			panic("invalid type " + typeName) //TODO:!! return err
		}
		buffer.WriteString(`(?P<`)
		buffer.WriteString(t[0])
		buffer.WriteString(`>`)
		buffer.WriteString(expr)
		buffer.WriteString(`)`)
		k = v[1]
	}
	buffer.WriteString(path[k:])
	if !strings.HasSuffix(path, "/") {
		buffer.WriteString(`/`)
	}
	buffer.WriteString(`?$`)

	return &routeTemplate{
		pathSignature: pathSignature(pointList, path),
		vars:          pointList,
		template:      path,
		pattern:       regexp.MustCompile(buffer.String()),
		method:        method,
		handler:       userHandler}

}

func pathSignature(points []routePoints, q string) string {
	buffer := bytes.Buffer{}
	i := 0
	if strings.HasPrefix(q, "/") {
		i++
	}
	f := len(q)
	if strings.HasSuffix(q, "/") {
		f--
	}
	for _, v := range points {
		buffer.WriteString(q[i:v.varStart])
		typeName := q[v.typeStart:v.typeEnd]
		buffer.WriteString(typeName)
		i = v.typeEnd
	}
	buffer.WriteString(q[i:f])
	return strings.Replace(buffer.String(), "/", "_", -1)
}

func (handler *RouterTemplateHandler) Dump() {
	for _, r := range handler.routes {
		rt := r.(*routeTemplate)
		fmt.Printf("template:  %s\npattern:   %s\nmethod:    %s\npath sign: %s\n\n", rt.template, rt.pattern, rt.method, rt.pathSignature)
	}
}

func (route *routeTemplate) consume(w http.ResponseWriter, r *http.Request) ConsumeResult {
	url := r.URL.Path

	match := route.pattern.FindStringSubmatch(url)
	lm := len(match)
	if lm == 0 {
		return Unconsumed
	}
	if lm > 1 {
		paramsMap := Vars{}
		for i, name := range route.pattern.SubexpNames() {
			if len(name) > 0 {
				paramsMap[name] = match[i]
			}
		}
		paramsMap[RouteSignatureParam] = route.pathSignature
		r = r.WithContext(newVars(r.Context(), paramsMap))
	}

	if r.Method != route.method {
		return MethodMismatch
	}
	route.handler.ServeHTTP(w, r)
	return Consumed
}

func (route *routeTemplate) URL(m URLVars) string {
	buffer := bytes.Buffer{}
	vars := route.vars
	used := map[string]bool{}
	q := route.template
	i := 0
	for _, v := range vars {
		buffer.WriteString(q[i:v.varStart])
		varName := q[v.varStart:v.varEnd]
		value, ok := m[varName]
		if ok {
			buffer.WriteString(value)
			used[varName] = true
		}
		i = v.typeEnd
	}
	buffer.WriteString(q[i:])
	if len(used) != len(m) {
		i = 0
		for k, v := range m {
			if len(v) > 0 && !used[k] {
				if i == 0 {
					buffer.WriteString("?")
				} else {
					buffer.WriteString("&")
				}
				buffer.WriteString(k)
				buffer.WriteString("=")
				buffer.WriteString(url.PathEscape(v))
				i++
			}
		}
	}
	return buffer.String()
}
