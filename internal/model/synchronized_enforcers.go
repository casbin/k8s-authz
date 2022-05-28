package model

import (
	"fmt"
	"sync"

	casbin "github.com/casbin/casbin/v2"
	"github.com/casbin/k8s-authz/pkg/casbinhelper"
	admission "k8s.io/api/admission/v1"
)

type EnforcerWrapper struct {
	Enforcer  *casbin.Enforcer
	ModelName string
}

type SynchronizedEnforcerList struct {
	sync.Mutex
	Enforcers []*EnforcerWrapper
}

var EnforcerList *SynchronizedEnforcerList

func init() {
	EnforcerList = NewSynchronizedEnforcerList()
	//test code
	e, err := casbin.NewEnforcer("example/model.conf", "example/policy.csv")
	if err != nil {
		panic(err)
	}
	e.AddFunction("access", casbinhelper.Access)
	EnforcerList.Enforcers = append(EnforcerList.Enforcers, &EnforcerWrapper{Enforcer: e, ModelName: "aaa"})
}

func NewSynchronizedEnforcerList() *SynchronizedEnforcerList {
	return &SynchronizedEnforcerList{
		Enforcers: make([]*EnforcerWrapper, 0),
	}
}

func (s *SynchronizedEnforcerList) Enforce(admission *admission.AdmissionReview) error {
	s.Lock()
	defer s.Unlock()

	for _, wrapper := range s.Enforcers {
		ok, err := wrapper.Enforcer.Enforce(admission)
		if err != nil {
			return fmt.Errorf("%s rejected the request: %s", wrapper.ModelName, err.Error())
		} else if !ok {
			return fmt.Errorf("%s rejected the request", wrapper.ModelName)
		}
	}
	return nil
}
