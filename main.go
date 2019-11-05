package main

import (
	pgs "github.com/lyft/protoc-gen-star"
	"github.com/sailmetrix/protoc-gen-twirp_typescript/json"
)

func main() {
	pgs.Init(
		pgs.DebugEnv("DEBUG"),
	).RegisterModule(
		json.Module(),
	).Render()
}

type typescript struct {
	*pgs.ModuleBase
}
