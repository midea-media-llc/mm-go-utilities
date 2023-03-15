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
}

type INames[T comparable] interface {
	InternalFindItemNameByListID(ids []*T, claims IClaims) []INameItem[T]
}

type FindNameModel[T comparable] struct {
	Id           *T
	Id2          []*T
	Module       IName[T]
	Module2      INames[T]
	OnCompleted  func(INameItem[T])
	OnCompleted2 func([]INameItem[T])
}

func FindNames[T comparable](claims IClaims, models ...*FindNameModel[T]) {
	callBacks := Select(models, func(model *FindNameModel[T]) Function {
		return func() {
			if model == nil {
				return
			}

			if model.Id != nil {
				result := model.Module.InternalFindItemName(*model.Id, claims)
				model.OnCompleted(result)
			}

			if model.Id2 != nil {
				results := model.Module2.InternalFindItemNameByListID(model.Id2, claims)
				model.OnCompleted2(results)
			}
		}
	})

	RunFuncThreads(callBacks, len(callBacks))
}
