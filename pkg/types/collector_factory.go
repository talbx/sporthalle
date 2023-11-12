package types

type CollectorFactory struct{}

func (c CollectorFactory) Create(src string) Collector {
	return SporthallenCollector{}
}
