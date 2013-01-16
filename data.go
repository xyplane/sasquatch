package main



var data map[string]interface{} = make(map[string]interface{})


type DataSource interface {
	ReadData(chan error)
	AddData() error
}


type dataFlag struct {
	Key string
	Value string
}

func (df *dataFlag) ReadData(rtn chan error) {
	rtn<- nil
}

func (df *dataFlag) AddData() error {
	data[df.Key] = df.Value
	return nil
}

