package sdk

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/deref/exo/internal/util/errutil"
)

var ErrResourceGone = errutil.NewHTTPError(http.StatusGone, "resource gone")

// TODO: Improve validation and error reporting for reflective calls.

type Controller struct {
	v reflect.Value
}

type Model interface {
	UnmarshalModel(ctx context.Context, s string) error
	MarshalModel(ctx context.Context) (string, error)
}

func NewController(v interface{}) *Controller {
	return &Controller{
		v: reflect.ValueOf(v),
	}
}

func unmarshalModel(ctx context.Context, label string, typ reflect.Type, s string) (Model, error) {
	m := reflect.New(typ.Elem()).Interface().(Model)
	err := m.UnmarshalModel(ctx, s)
	if err != nil {
		err = fmt.Errorf("unmarshaling %s: %w", label, err)
	}
	return m, err
}

func marshalModel(ctx context.Context, typ reflect.Type, m Model) (string, error) {
	s, err := m.MarshalModel(ctx)
	if err != nil {
		err = fmt.Errorf("mashaling updated model: %w", err)
	}
	return s, err
}

func (c *Controller) Identify(ctx context.Context, model string) (iri string, err error) {
	defer errutil.RecoverTo(&err)
	method := c.v.MethodByName("Identify")
	unmarshaledModel, err := unmarshalModel(ctx, "model", method.Type().In(1), model)
	if err != nil {
		return "", err
	}
	res := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(unmarshaledModel),
	})
	iri = res[0].Interface().(string)
	err, _ = res[1].Interface().(error)
	return iri, err
}

func (c *Controller) Create(ctx context.Context, model string) (updatedModel string, err error) {
	defer errutil.RecoverTo(&err)
	method := c.v.MethodByName("Create")
	unmarshaledModel, err := unmarshalModel(ctx, "model", method.Type().In(1), model)
	if err != nil {
		return "", err
	}
	res := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(unmarshaledModel),
	})
	err, _ = res[0].Interface().(error)
	if err != nil {
		return
	}
	return unmarshaledModel.MarshalModel(ctx)
}

func (c *Controller) Read(ctx context.Context, model string) (updatedModel string, err error) {
	defer errutil.RecoverTo(&err)
	method := c.v.MethodByName("Read")
	unmarshaledModel, err := unmarshalModel(ctx, "model", method.Type().In(1), model)
	if err != nil {
		return "", err
	}
	res := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(unmarshaledModel),
	})
	err, _ = res[0].Interface().(error)
	if err != nil {
		return
	}
	return unmarshaledModel.MarshalModel(ctx)
}

func (c *Controller) Update(ctx context.Context, prev string, cur string) (updatedModel string, err error) {
	defer errutil.RecoverTo(&err)
	method := c.v.MethodByName("Read")
	unmarshaledPrev, err := unmarshalModel(ctx, "previous model", method.Type().In(1), prev)
	if err != nil {
		return "", err
	}
	unmarshaledCur, err := unmarshalModel(ctx, "current model", method.Type().In(2), cur)
	if err != nil {
		return "", err
	}
	res := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(unmarshaledPrev),
		reflect.ValueOf(unmarshaledCur),
	})
	err, _ = res[0].Interface().(error)
	if err != nil {
		return
	}
	return unmarshaledCur.MarshalModel(ctx)
}

func (c *Controller) Delete(ctx context.Context, model string) (err error) {
	defer errutil.RecoverTo(&err)
	method := c.v.MethodByName("Delete")
	unmarshaledModel, err := unmarshalModel(ctx, "model", method.Type().In(1), model)
	if err != nil {
		return err
	}
	res := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(unmarshaledModel),
	})
	err, _ = res[0].Interface().(error)
	return err
}
