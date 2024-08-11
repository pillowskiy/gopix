package validator

import (
	"context"

	baseValidator "github.com/go-playground/validator/v10"
)

var v *baseValidator.Validate

func init() {
	v = baseValidator.New(baseValidator.WithRequiredStructEnabled())
}

func ValidateStruct(ctx context.Context, i interface{}) error {
	return v.StructCtx(ctx, i)
}
