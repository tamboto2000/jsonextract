# JSONExtract

[![Go Reference](https://pkg.go.dev/badge/github.com/tamboto2000/jsonextract.svg)](https://pkg.go.dev/github.com/tamboto2000/jsonextract/v3)

Package jsonextract is a library for extracting JSON from a given source, such as string, bytes, file, and io.Reader, providing methods for editing and evaluating json values. One of the cases to use this library is to extract all json strings in a scraped HTML page

### Installation
JSONExtract require Go v1.14 or higher

```sh
$ GO111MODULE go get github.com/tamboto2000/jsonextract/v3
```

### Sources

You can extract jsons from string, []byte, file, and an io.Reader

```go
// from string
jsons, err := jsonextract.FromString(str)

// from []byte
jsons, err := jsonextract.FromBytes(byts)

// from file path
jsons, err := jsonextract.FromFile(path)

// from io.Reader
jsons, err := jsonextract.FromReader(reader)
```

### JSON Kinds

There's 7 different kinds of JSON:
 - Object
 - Array
 - String
 - Integer
 - Float
 - Boolean
 - Null

You can get JSON kind by calling ```JSON.Kind()```

```go
jsons, err := jsonextract.FromBytes(byts)
if err != nil {
    panic(err.Error())
}

for _, json := range jsons {
    // string
    if json.Kind() == jsonextract.String {
        // print string value
        fmt.Println(json.String())
    }
    
    // int
    if json.Kind() == jsonextract.Integer {
        // print integer value, returned integer value is int64
        fmt.Println(json.Integer())
    }
    
    // float
    if json.Kind() == jsonextract.Float {
        // print float value, returned float is float64
        fmt.Println(json.Float())
    }
    
    // and so on...
}
```

Getter methods, such as ```JSON.String()```, ```JSON.Integer()```, ```JSON.Float()```, etc., will panic if ```JSON``` kind did not match the getter methods, for example, when trying to get string value from ```JSON``` with kind of ```Integer```. ```JSON``` with kind of ```Null``` didn't have getter method

### Modifying Value

Value inside JSON can be modified by setter methods, such as ```JSON.SetStr()```, ```JSON.SetInt()```, etc.

```go
// integer
if json.Kind() == jsonextract.Integer {
    json.EditInt(23)
}

// string
if json.Kind() == jsonextract.String {
    json.EditStr("Hello World!")
}

// and so on...
```

Setter methods will panic if JSON kind did not match with setter method, for example, trying to call ```JSON.SetInt()``` to ```JSON``` with kind of ```String```. ```JSON``` with kind of ```Null``` didn't have setter method

### JSON Object

Json object represented as ```JSON``` with kind of ```Object```. For adding a new value, call ```JSON.AddField()``` with param ```key``` as json field name, and ```val``` for the value

```go
json.AddField("name", "Franklin Collin Tamboto")

json.AddField("email", "tamboto2000@gmail.com")

json.AddField("id", 1)

i := 1
json.AddField("id", &i)


json.AddField("qty", int32(23))


json.AddField("userId", uint(2))


json.AddField("height", float32(3.2))
```

Just like package ```encoding/json```, ```map``` and ```struct``` will be marshaled as json object. ```map``` keys need to be type ```string``` or an integer type, otherwise panic will occur

```go
type User struct {
    ID int `json:"id"`
    Name string `json:"name"`
    // field name will be used as json field name if json tag is not present, 
    // so this field will be marshaled as "Email":"value"
    Email string
}

// add struct
json.AddField("user", User{
    ID: 1,
    Name: "Franklin Collin Tamboto",
    Email: "tamboto2000@gmail.com",
})

// add map
json.AddField("item", map[interface{}]interface{}{
    "Name": "ROG Zephyrus M",
    "Qty": 23,
    123: 456,
})

// add value with int as key
json.AddField(123, 456)

// add value with uint8 as key
json.AddField(uint8(432), "123")
```

To delete an item from object, call ```JSON.DeleteField()``` with param ```key``` as json field name for deletion. ```key``` must be string or an integer type, otherwise panic will occur. ```JSON.DeleteField()``` return ```true``` if item is exist, otherwise ```false``` will returned

```go
// delete with string key
if json.DeleteField("user") {
    fmt.Println("item deleted")
} else {
    fmt.Println("item not exist")
}

// delete with integer key
if json.DeleteField(int32(123)) {
    fmt.Println("item deleted")
} else {
    fmt.Println("item not exist")
}
```

To access all contained items, call ```JSON.Object()```, this method will return ```map[interface{}]*JSON```

```go
fields := json.Object()
for _, f := range fields  {
    // do something...
}
```

Get contained items count by calling ```JSON.Len()```

```go
fmt.Println("len:", json.Len())
```

### JSON Array

Json array represented as ```JSON``` with kind of ```Array```. For adding a new value, call ```JSON.AddItem()``` with param ```val```. Will panic if ```val``` is not valid json value

```go
json.AddItem("Hello world")
```

Get contained items count by calling ```JSON.Len()```

```go
fmt.Println("len:", json.Len())
```

To delete an item, call ```json.DeleteItem()``` with param ```i``` as array index

```go
if json.DeleteItem(2) {
    fmt.Println("item deleted")
} else {
    fmt.Println("item is not exist")
}
```

Get all items by calling ```JSON.Array()```

```go
items := json.Array()
for _, item := range items {
    // do something
}
```

### Additional Note
When an item is edited inside ```JSON```, the contained raw json will automatically re-generated.
*DO NOT PERFORM ANY SETTER OR GETTER METHODS CONCURRENTLY*

### Documentation and Examples
See [Documentation](https://pkg.go.dev/github.com/tamboto2000/jsonextract/v3) for more information. There's also some [examples](https://github.com/tamboto2000/jsonextract/tree/v3/examples) you can look up to 

License
----

MIT