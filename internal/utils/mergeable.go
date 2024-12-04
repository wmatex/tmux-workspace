package utils

type Mergeable[T any] interface {
	Merge(other T)
	Id() string
}

func Merge[T Mergeable[T]](list []T) []T {
	var merged []T

	indexes := make(map[string]int)
	for _, entity := range list {
		if index, ok := indexes[entity.Id()]; ok {
			merged[index].Merge(entity)
		} else {
			merged = append(merged, entity)
			indexes[entity.Id()] = len(merged) - 1
		}
	}

	return merged
}
