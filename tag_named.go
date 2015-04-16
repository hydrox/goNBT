package goNBT

type TAG_Named struct {
	Name string
}

func (tag *TAG_Named) GetName() string {
	return tag.Name
}
func (tag *TAG_Named) SetName(name string) {
	tag.Name = name
}
