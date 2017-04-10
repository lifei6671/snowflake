package snowflake

type Generator interface {
	NextId() (int64,error)
	ToString(int64) string
}
