/*
mongodb package provides a slightly higher level API for MongoDB than the official driver,
although of course it builds on the official driver.

# Documents

Documents are mapped to collections through their type. The type must be a struct, and
the struct must have a field which is a primitive.ObjectID or a string. Use struct tags
to define which field is the key field. For example,

	type Source struct {
		Id         string        `bson:"_id,omitempty"`
		Url        string        `bson:"url,unique"`
	}

The bson tags should either contain metadata or '-' to skip the field in the MongoDB database.
Use the following flags:

  - unique - The field is unique in the collection, and an index is generated for the field (This
    is not part of the underlying MongoDB driver)
  - index  - An index is generated for the field (This is not part of the underlying MongoDB driver)
  - omitempty - The field is omitted from the document if it is empty
*/
package mongodb
