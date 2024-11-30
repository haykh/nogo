package utils

type Configuration struct {
	name        string
	desc        string
	value       interface{}
	def         interface{}
	has_default bool
	issecret    bool
}

func (c Configuration) Name() string {
	return c.name
}

func (c Configuration) Value() interface{} {
	return c.value
}

func (c Configuration) IsSecret() bool {
	return c.issecret
}

func NewParameter(
	name, desc string,
	value, def interface{}) Configuration {
	return Configuration{
		name, desc, value, def, true, false,
	}
}

func NewConfiguration(
	name, desc string,
	value interface{},
) Configuration {
	return Configuration{
		name, desc, value, nil, false, false,
	}
}

func NewSecret(
	name, desc string,
	value interface{}) Configuration {
	return Configuration{
		name, desc, value, nil, false, false,
	}
}

func (c *Configuration) Prompt() error {
	if val, err := Prompt(c.desc, c.has_default, c.def); err != nil {
		return err
	} else {
		c.value = val
		return nil
	}
}
