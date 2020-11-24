package cache

var pool = make(map[string]*RedisConnection)

func Connect(key string, config *RedisConfig) (err error) {

	if conn, exist := pool[key]; exist {

		if conn.TryConnect() == nil {
			return nil
		}
	}

	conn := &RedisConnection{
		Client:      nil,
		RedisConfig: config,
	}

	pool[key] = conn

	return conn.Connect()
}

func Disconnect(key string) (err error) {

	conn, exist := pool[key]
	if !exist {
		return nil
	}

	return conn.Close()
}

func Get(key string) (*RedisConnection, error) {

	conn, exist := pool[key]
	if !exist {
		return nil, errorNotHaveConnection
	}

	return conn, nil
}
