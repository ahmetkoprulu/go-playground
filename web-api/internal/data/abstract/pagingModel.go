package data_models

type PagingModel[T any] struct {
	TotalPage   int `json:"totalPage"`
	CurrentPage int `json:"currentPage"`
	Take        int `json:"take"`
	TotalCount  int `json:"totalCount"`
	Data        []T `json:"data"`
}

type AggregateResult[T any] []struct {
	Metadata []struct {
		Total int `bson:"total"`
	} `bson:"metadata"`
	Data []T `bson:"data"`
}
