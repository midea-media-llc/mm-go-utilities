package utils

type IName[T comparable] interface {
	InternalFindItemName(id T, claims IClaims) *string
}

type FindNameModel[T comparable] struct {
	Id          *T
	Module      IName[T]
	OnCompleted func(*string)
}

func FindNames[T comparable](claims IClaims, models ...*FindNameModel[T]) {
	callBacks := Select(models, func(model *FindNameModel[T]) Function {
		return func() {
			if model == nil || model.Id == nil {
				return
			}

			model.OnCompleted(model.Module.InternalFindItemName(*model.Id, claims))
		}
	})

	RunFuncThreads(callBacks, len(callBacks))
}
