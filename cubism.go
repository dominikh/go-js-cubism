package cubism

import (
	"strings"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/d3"
	"honnef.co/go/js/d3/scale"
	"honnef.co/go/js/util"
)

func Version() string {
	return js.Global.Get("cubism").Get("version").String()
}

type Context struct {
	js.Object
}

func NewContext() *Context {
	return &Context{js.Global.Get("cubism").Call("context")}
}

func (c *Context) Step() time.Duration {
	return time.Duration(c.Call("step").Int()) * time.Millisecond
}

func (c *Context) SetStep(step time.Duration) {
	c.Call("step", step/time.Millisecond)
}

func (c *Context) Size() int {
	return c.Call("size").Int()
}

func (c *Context) SetSize(size int) {
	c.Call("size", size)
}

func (c *Context) ServerDelay() time.Duration {
	return time.Duration(c.Call("serverDelay").Int()) * time.Millisecond
}

func (c *Context) SetServerDelay(delay time.Duration) {
	c.Call("serverDelay", delay/time.Millisecond)
}

func (c *Context) ClientDelay() time.Duration {
	return time.Duration(c.Call("clientDelay").Int()) * time.Millisecond
}

func (c *Context) SetClientDelay(delay time.Duration) {
	c.Call("clientDelay", delay/time.Millisecond)
}

type Graphite struct {
	js.Object
}

func (c *Context) Graphite(url string) *Graphite {
	return &Graphite{c.Call("graphite", url)}
}

func (g *Graphite) Metric(expr string) *Metric {
	return &Metric{g.Call("metric", expr)}
}

func (g *Graphite) Find(pattern string, cb func(results []string, err error)) {
	fn := func(err js.Object, results []string) {
		cb(results, util.Error(err))
	}
	g.Call("find", pattern, fn)
}

func (g *Graphite) URL() string {
	return g.Call("toString").String()
}

type Cube struct {
	js.Object
}

func (c *Context) Cube(url string) *Cube {
	return &Cube{c.Call("cube", url)}
}

type Librato struct {
	js.Object
}

func (c *Context) Librato(user, token string) *Librato {
	return &Librato{c.Call("librato", user, token)}
}

func (c *Context) Constant(value float64) *Metric {
	return &Metric{c.Call("constant", value)}
}

func (c *Context) NewMetric(req Request, name string) *Metric {
	ret := c.Call("metric", func(start, stop js.Object, step int, cb js.Object) {
		tStart := start.Interface().(time.Time)
		tStop := stop.Interface().(time.Time)
		dStep := time.Duration(step) * time.Millisecond
		req.Request(tStart, tStop, dStep, cb)
	}, name)
	return &Metric{ret}
}

func (c *Context) NewAxis() *Axis {
	return &Axis{c.Call("axis")}
}

func (c *Context) NewRule() *Rule {
	return &Rule{c.Call("rule")}
}

func (c *Context) NewHorizon() *Horizon {
	return &Horizon{c.Call("horizon")}
}

type Comparison struct {
	js.Object
}

func (c *Context) NewComparison() *Comparison {
	return &Comparison{c.Call("comparison")}
}

func (co *Comparison) Apply(sel d3.Selection) {
	co.Invoke(sel.Object)
}

func (co *Comparison) SetHeight(height int) {
	co.Call("height", height)
}

func (co *Comparison) Height() int {
	return co.Call("height").Int()
}

func (co *Comparison) SetPrimary(m *Metric) {
	co.Call("primary", m.Object)
}

func (co *Comparison) Primary() *Metric {
	return &Metric{co.Call("primary")}
}

func (co *Comparison) SetSecondary(m *Metric) {
	co.Call("secondary", m.Object)
}

func (co *Comparison) Secondary() *Metric {
	return &Metric{co.Call("secondary")}
}

// TODO way to use functions for primary/secondary/extent/title

func (co *Comparison) SetScale(scale scale.Scale) {
	co.Call("scale", scale.Underlying())
}

func (co *Comparison) Scale() scale.Scale {
	return scale.Wrap(co.Call("scale"))
}

func (co *Comparison) SetExtent(min, max float64) {
	co.Call("extent", min, max)
}

func (co *Comparison) Extent() (min, max float64) {
	vals := co.Call("extent").Interface().([]interface{})
	return vals[0].(float64), vals[1].(float64)
}

func (co *Comparison) SetTitle(title string) {
	co.Call("title", title)
}

func (co *Comparison) Title() string {
	return co.Call("title").String()
}

func (co *Comparison) SetPrimaryFormat(format d3.Formatter) {
	co.Call("formatPrimary", format)
}

func (co *Comparison) PrimaryFormat() d3.Formatter {
	// TODO document that calling this will give us a new function
	f := co.Call("formatPrimary")
	return func(o js.Object) string { return f.Invoke(o).String() }
}

func (co *Comparison) SetChangeFormat(format d3.Formatter) {
	co.Call("formatChange", format)
}

func (co *Comparison) ChangeFormat() d3.Formatter {
	// TODO document that calling this will give us a new function
	f := co.Call("formatChange")
	return func(o js.Object) string { return f.Invoke(o).String() }
}

func (co *Comparison) SetColors(colors []string) {
	co.Call("colors", colors)
}

func (co *Comparison) Colors() []string {
	return util.StringSlice(co.Call("colors"))
}

func (co *Comparison) SetStrokeWidth(width float64) {
	co.Call("strokeWidth", width)
}

func (co *Comparison) StrokeWidth() float64 {
	return co.Call("strokeWidth").Float()
}

func (co *Comparison) Remove(sel d3.Selection) {
	co.Call("remove", sel.Object)
}

func (c *Context) Scale() scale.Scale {
	return scale.Wrap(c.Call("scale"))
}

func (c *Context) Start() {
	c.Call("start")
}

func (c *Context) Stop() {
	c.Call("stop")
}

func (c *Context) AddListener(typ string, fn func(...interface{})) {
	parts := strings.Split(typ, ".")
	actualType := parts[0]
	var cb interface{}
	switch actualType {
	case "change", "beforechange", "prepare":
		cb = func(start, stop js.Object) {
			fn(start.Interface().(time.Time), stop.Interface().(time.Time))
		}
	case "focus":
		// FIXME what about the case of a null index?
		cb = func(idx int) { fn(idx) }
	default:
		panic("unsupported event type " + typ)
	}
	c.Call("on", typ, cb)
}

func (c *Context) RemoveListener(typ string) {
	c.Call("on", typ, nil)
}

func (c *Context) GetListener(typ string) interface{} {
	return c.Call("on", typ).Interface()
}

func (c *Context) Focus(idx int) {
	c.Call("focus", idx)
}

func (c *Context) NoFocus() {
	c.Call("focus", nil)
}

func (cu *Cube) Metric(expr string) *Metric {
	return &Metric{cu.Call("metric", expr)}
}

func (cu *Cube) URL() string {
	return cu.Call("toString").String()
}

func (l *Librato) Metric(name, source string) *Metric {
	return &Metric{l.Call("metric", name, source)}
}

type Request interface {
	Request(start, stop time.Time, step time.Duration, cb js.Object)
}

type RequestFunc func(start, stop time.Time, step time.Duration, cb js.Object)

func (f RequestFunc) Request(start, stop time.Time, step time.Duration, cb js.Object) {
	f(start, stop, step, cb)
}

type Metric struct {
	js.Object
}

func (m *Metric) Add(other Metric) *Metric {
	return &Metric{m.Call("add", other.Object)}
}

func (m *Metric) Subtract(other Metric) *Metric {
	return &Metric{m.Call("subtract", other.Object)}
}

func (m *Metric) Multiplty(other Metric) *Metric {
	return &Metric{m.Call("multiplty", other.Object)}
}

func (m *Metric) Divide(other Metric) *Metric {
	return &Metric{m.Call("divice", other.Object)}
}

func (m *Metric) Shift(offset int) *Metric {
	return &Metric{m.Call("shift", offset)}
}

func (m *Metric) ValueAt(index int) float64 {
	return m.Call("valueAt", index).Float()
}

func (m *Metric) Extent() (float64, float64) {
	vals := m.Call("extent").Interface().([]interface{})
	return vals[0].(float64), vals[1].(float64)
}

func (m *Metric) AddListener(typ string, fn func(...interface{})) {
	parts := strings.Split(typ, ".")
	actualType := parts[0]
	var cb interface{}
	switch actualType {
	case "change":
		cb = func(start, stop js.Object) {
			fn(start.Interface().(time.Time), stop.Interface().(time.Time))
		}
	default:
		panic("unsupported event type " + typ)
	}
	m.Call("on", typ, cb)
}

func (m *Metric) RemoveListener(typ string) {
	m.Call("on", typ, nil)
}

func (m *Metric) GetListener(typ string) interface{} {
	return m.Call("on", typ).Interface()
}

func (m *Metric) Context() *Context {
	return &Context{m.Get("context")}
}

func (m *Metric) String() string {
	return m.Call("toString").String()
}

func (m *Metric) Underlying() js.Object {
	return m.Object
}

type Axis struct {
	js.Object
}

func (a *Axis) SetOrientation(orientation string) {
	a.Call("orient", orientation)
}

func (a *Axis) Orientation() string {
	return a.Call("orient").String()
}

func (a *Axis) SetTicks(args ...interface{}) {
	// TODO find a way of doing this without empty interfaces
	a.Call("ticks", args...)
}

func (a *Axis) Ticks() interface{} {
	// TODO find a way of doing this without empty interfaces
	return a.Call("ticks").Interface()
}

func (a *Axis) SetTickSubdivide(count int) {
	a.Call("tickSubdivide", count)
}

func (a *Axis) TickSubdivide() int {
	return a.Call("tickSubdivide").Int()
}

func (a *Axis) SetTickSize(major, minor, end int) {
	a.Call("tickSize", major, minor, end)
}

func (a *Axis) TickSize() (major, minor, end int) {
	vals := util.IntSlice(a.Call("tickSize"))
	return vals[0], vals[1], vals[2]
}

func (a *Axis) SetTickPadding(padding int) {
	a.Call("tickPadding", padding)
}

func (a *Axis) TickPadding() int {
	return a.Call("tickPadding").Int()
}

func (a *Axis) SetTickFormat(f d3.Formatter) {
	a.Call("tickFormat", f)
}

func (a *Axis) TickFormat() d3.Formatter {
	// TODO document that calling this will give us a new function
	f := a.Call("tickFormat")
	return func(o js.Object) string { return f.Invoke(o).String() }
}

func (a *Axis) Remove(sel d3.Selection) {
	a.Call("remove", sel.Object)
}

func (a *Axis) Apply(sel d3.Selection) {
	a.Invoke(sel.Object)
}

type Rule struct {
	js.Object
}

func (r *Rule) Apply(sel d3.Selection) {
	r.Invoke(sel.Object)
}

func (r *Rule) Remove(sel d3.Selection) {
	r.Call("remove", sel.Object)
}

func (r *Rule) SetMetric(metric *Metric) {
	r.Call("metric", metric.Object)
}

func (r *Rule) Metric() *Metric {
	return &Metric{r.Call("metric")}
}

type Horizon struct {
	js.Object
}

func (h *Horizon) Apply(sel d3.Selection) {
	h.Invoke(sel.Object)
}

func (h *Horizon) SetMode(mode string) {
	h.Call("mode", mode)
}

func (h *Horizon) Mode() string {
	return h.Call("mode").String()
}

func (h *Horizon) SetHeight(height int) {
	h.Call("height", height)
}

func (h *Horizon) Height() int {
	return h.Call("height").Int()
}

func (h *Horizon) SetMetric(metric *Metric) {
	// TODO support using a function for the metric
	h.Call("metric", metric.Object)
}

func (h *Horizon) Metric() *Metric {
	return &Metric{h.Call("metric")}
}

func (h *Horizon) SetScale(scale scale.Scale) {
	h.Call("scale", scale.Underlying())
}

func (h *Horizon) Scale() scale.Scale {
	return scale.Wrap(h.Call("scale"))
}

func (h *Horizon) SetExtent(min, max float64) {
	h.Call("extent", min, max)
}

func (h *Horizon) Extent() (min, max float64) {
	vals := h.Call("extent").Interface().([]interface{})
	return vals[0].(float64), vals[1].(float64)
}

func (h *Horizon) SetTitle(title string) {
	// TODO support using a function
	h.Call("title", title)
}

func (h *Horizon) Title() string {
	return h.Call("title").String()
}

func (h *Horizon) SetFormat(format d3.Formatter) {
	h.Call("format", format)
}

func (h *Horizon) Format() d3.Formatter {
	f := h.Call("tickFormat")
	return func(o js.Object) string { return f.Invoke(o).String() }
}

func (h *Horizon) SetColors(colors []string) {
	h.Call("colors", colors)
}

func (h *Horizon) Colors() []string {
	return util.StringSlice(h.Call("colors"))
}

func (h *Horizon) Remove(sel d3.Selection) {
	h.Call("remove", sel.Object)
}
