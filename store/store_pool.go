package store

/*
 * Right now, just keeps a collection of store.Interfaces, one for each keyspace.
 * Someday, maybe it'll be a real connection pool, with multiple clients per keyspace.
 */

type Pool struct {
	pool    map[string]Interface
	factory func(keyspace string) Interface
}

/*
 * Maybe this shouldn't be genericized? And should instead be an interface, and then
 * we could have MemPool, RedisPool, MySQLPool, etc. that could handle connections
 * intelligently.
 */
func NewPool(storeFactory func(keyspace string) Interface) *Pool {
	pool := make(map[string]Interface, 0)
	return &Pool{pool, storeFactory}
}

func (p *Pool) Get(keyspace string) (s Interface) {
	// TODO(benjamin) there should probably be a lock here, so we don't
	// have a race between grabbing a nonexistent client and one being created.
	s, ok := p.pool[keyspace]
	if ok {
		return
	}
	p.pool[keyspace] = p.factory(keyspace)
	return p.pool[keyspace]
}
