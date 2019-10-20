package zapstackdriver

type validatedField interface {
	validate() error
}
