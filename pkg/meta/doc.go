/*
Package meta can be used to extract metadata from Go structs, such as field names and tags.
To create a meta object, use the New function, passing in a reflect.Value and a tag name
to use for field tags. The resulting object can be used to extract metadata information for
the struct.

The following field tags have special meaning:

  - tag:"-" - The field should be ignored
  - tag:"name" - Use the specified name for the field
  - tag:"name,primarykey" - The field is part of the primary key
  - tag:"name,unique" - The field value should be unique across the collection
  - tag:"name,foreignkey" - The field value references another collection
  - tag:"name,type:text" - The field type is TEXT

Tags can be combined, for example "name,primarykey,unique" would indicate that the field is
both a primary key and unique.

The CollectionName function can be used to extract the name of the collection from the struct.
For example,

	  type MyStruct struct {
		...
	  }

	  func (MyStruct) CollectionName() string {
		return "my_collection"
	  }

The CollectionRef function can be used to reference another collection. For example,

	  type MyStruct struct {
		...
	  }

	  // CREATE TABLE (..., FOREIGN KEY name REFERENCES other_table (name), ...)
	  func (MyStruct) CollectionRef(field string) any {
		switch field {
		case "name":
			return OtherStruct{}
		}
		return nil
	  }
*/
package meta
