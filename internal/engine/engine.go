package engine

import "k8s.io/client-go/rest"

const CONTOUR = "contour"

type (
	ConfigEngine interface {
		GetConfig() (map[string]interface{}, error)
	}

	Engine struct {
		Engines map[string]ConfigEngine
	}
)

func NewEngine(k8s *rest.Config) *Engine {
	return &Engine{
		map[string]ConfigEngine{
			CONTOUR: NewContour(*k8s),
		},
	}
}
