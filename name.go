package utils

type INameItem[T comparable] interface {
	GetId() T
	GetName() *string
}

type NameItem[T any] struct {
	Id   T
	Name *string
}

func (e *NameItem[T]) GetId() T {
	return e.Id
}

func (e *NameItem[T]) GetName() *string {
	return e.Name
}

type IName[T comparable] interface {
	InternalFindItemName(id T, claims IClaims) INameItem[T]
	InternalFindItemNameByListID(ids []*T, claims IClaims) []INameItem[T]
}

type FindNameModel[T comparable] struct {
	Id               *T
	Ids              []*T
	OnCompleted      func(INameItem[T])
	OnCompletedSlice func([]INameItem[T])
	Module           IName[T]
}

func FindName[T comparable](claims IClaims, models ...*FindNameModel[T]) {
	callBacks := Select(models, func(model *FindNameModel[T]) Function {
		return func() {
			if model == nil {
				return
			}

			if model.Id != nil {
				result := model.Module.InternalFindItemName(*model.Id, claims)
				model.OnCompleted(result)
			}

			if model.Ids != nil {
				results := model.Module.InternalFindItemNameByListID(model.Ids, claims)
				model.OnCompletedSlice(results)
			}
		}
	})

	RunFuncThreads(callBacks, len(callBacks))
}

func FindNames[T comparable](claims IClaims, modules []IName[T], ids [][]*T, completeds []func([]INameItem[T])) {
	models := make([]*FindNameModel[T], len(modules))
	for i := 0; i < len(modules); i++ {
		models[i] = &FindNameModel[T]{
			Module:           modules[i],
			Ids:              ids[i],
			OnCompletedSlice: completeds[i],
		}
	}
	FindName(claims, models...)
}
