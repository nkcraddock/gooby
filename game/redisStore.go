package game

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.google.com/p/go-uuid/uuid"

	"gopkg.in/redis.v3"
)

type RedisStore struct {
	opts *redis.Options
	conn *redis.Client
}

func NewRedisStore(addr, password string, dbindex int64) (*RedisStore, error) {
	store := &RedisStore{
		opts: &redis.Options{
			Addr:     addr,
			Password: password,
			DB:       dbindex,
			PoolSize: 10,
		},
	}

	// Make sure we can ping the redis server
	if err := store.Open(); err != nil {
		return nil, err
	}
	store.Close()

	return store, nil
}

// SaveEvent will add a new event
func (s *RedisStore) SaveEvent(e *Event) error {
	e.Id = uuid.New()

	if err := s.save("events", eventKey(e.Id), e); err != nil {
		return err
	}

	indexKey := playerEventsKey(e.Player)
	return s.push(indexKey, e.Id)
}

// GetEvent retrieves an existing event
func (s *RedisStore) GetEvent(id string) (*Event, error) {
	event := new(Event)
	if err := s.get("events", eventKey(id), event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *RedisStore) GetPlayerEvents(name string, count int64) ([]*Event, error) {
	values, err := s.conn.LRange(playerEventsKey(name), 0, count-1).Result()
	if err != nil {
		return nil, err
	}

	events := make([]*Event, len(values))
	for i, id := range values {
		evt, err := s.GetEvent(id)
		if err != nil {
			return nil, err
		}
		events[i] = evt
	}

	return events, nil
}

// SavePlayer will add a new player or update an existing player
func (s *RedisStore) SavePlayer(p *Player) error {
	return s.save("players", playerKey(p.Name), p)
}

// GetPlayer will retrieve an existing player or return ErrorNotFound
func (s *RedisStore) GetPlayer(name string) (*Player, error) {
	player := new(Player)
	if err := s.get("players", playerKey(name), player); err != nil {
		return nil, err
	}
	return player, nil
}

// ListPlayers retrieves a list of all players
func (s *RedisStore) ListPlayers() ([]*Player, error) {
	results, err := s.list("players")
	if err != nil {
		return nil, err
	}

	players := make([]*Player, len(results))

	i := 0
	for _, playerJson := range results {
		p := new(Player)
		if err = json.Unmarshal([]byte(playerJson), p); err != nil {
			return nil, err
		}

		players[i] = p
		i += 1
	}

	return players, nil
}

// SaveRule stores a new rule or updates an existing rule
func (s *RedisStore) SaveRule(r *Rule) error {
	return s.save("rules", ruleKey(r.Code), r)
}

// GetRule retrieves an existing rule or returns ErrorNotFound
func (s *RedisStore) GetRule(code string) (*Rule, error) {
	rule := new(Rule)
	if err := s.get("rules", ruleKey(code), rule); err != nil {
		return nil, err
	}

	return rule, nil
}

// ListRules retrieves a list of all rules
func (s *RedisStore) ListRules() ([]*Rule, error) {
	results, err := s.list("rules")
	if err != nil {
		return nil, err
	}

	rules := make([]*Rule, len(results))

	i := 0
	for _, ruleJson := range results {
		r := new(Rule)
		if err = json.Unmarshal([]byte(ruleJson), r); err != nil {
			return nil, err
		}

		rules[i] = r
		i += 1
	}

	return rules, nil
}

func (s *RedisStore) Open() error {
	if s.conn != nil {
		return nil
	}

	client := redis.NewClient(s.opts)
	if _, err := client.Ping().Result(); err != nil {
		return err
	}

	s.conn = client
	return nil
}

func (s *RedisStore) Close() {
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}
}

func (s *RedisStore) FlushDb() error {
	return s.conn.FlushDb().Err()
}

func (s *RedisStore) save(col, id string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return s.conn.HSet(col, id, string(jsonData)).Err()
}

func (s *RedisStore) get(col, id string, data interface{}) error {
	jsonData, err := s.conn.HGet(col, id).Result()

	if err == redis.Nil {
		return ErrorNotFound
	} else if err != nil {
		return err
	}

	return json.Unmarshal([]byte(jsonData), data)
}

func (s *RedisStore) list(col string) (map[string]string, error) {
	return s.conn.HGetAllMap(col).Result()
}

func (s *RedisStore) push(col, val string) error {
	return s.conn.LPush(col, val).Err()
}

func playerKey(id string) string {
	return strings.ToLower(id)
}

func ruleKey(id string) string {
	return strings.ToLower(id)
}

func eventKey(id string) string {
	return strings.ToLower(id)
}

func playerEventsKey(id string) string {
	return fmt.Sprintf("players:%s:events", playerKey(id))
}
