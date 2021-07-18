package store

/*
	Store service interface.
	Use for implement different types of store, such as:
      - disk store
      - memory store
      - etc ...
   with predefined methods.
*/
type Store interface {
	Write(chunk []byte) error
	Close() error
}