package opa

import (
	"context"
	"opaapp/utils"
	"time"

	"github.com/open-policy-agent/opa/rego"
)

var module = `
package application.authz

# Everyone can see adopted pets
allowed_pets[pet] {
    some i
    pet := input.pet_list[i]
    pet.up_for_adoption == true
}

# Employees can see all pets.
allowed_pets[pet] {
    [header, payload, signature] := io.jwt.decode(input.token)
    payload.employee == true
    some i
    pet := input.pet_list[i]
}
`

func RunRegoQuery(input map[string]interface{}) rego.ResultSet {
	defer utils.TimeTrack(time.Now(), "runRegoQuery")
	ctx := context.TODO()
	query, err := rego.New(
		rego.Query("result = data.application.authz.allowed_pets"),
		rego.Module("example.rego", module),
	).PrepareForEval(ctx)

	if err != nil {
		panic(err)
	}

	results, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		panic(err)
	}
	// else if len(results) == 0 {
	// 	// Handle undefined result.
	// 	panic("Unexpected result")
	// } else if result, ok := results[0].Bindings["x"].(bool); !ok {
	// 	// Handle unexpected result type.
	// 	panic("Unexpected result")
	// } else {
	// 	// Handle result/decision.
	// 	// fmt.Printf("%+v", results) => [{Expressions:[true] Bindings:map[x:true]}]
	// }

	return results
}
