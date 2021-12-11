package aqtalk

type SynthFactory struct {
	store map[string]*Synthesizer
}

func (sf *SynthFactory) Get(typ string) (*Synthesizer, error) {
	if sf.store == nil {
		sf.store = make(map[string]*Synthesizer)
	}

	s, ok := sf.store[typ]
	if ok {
		return s, nil
	}

	s, err := NewAqTalk1Synthesizer(typ)
	if err != nil {
		return nil, err
	}
	sf.store[typ] = s
	return s, nil
}