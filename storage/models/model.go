package models

type IModel interface {
	GetId() string
	SetId(id string)
}

type Model struct {
	Id string
}

func (m *Model) GetId() string {
	return m.Id
}

func (m *Model) SetId(id string) {
	m.Id = id
}
