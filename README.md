# Swagger-Gen

Generate a swagger 2.0 file directly from comments in sources files. 

### Status
[![Build Status](https://travis-ci.org/macinnir/swagger-gen.svg?branch=master)](https://travis-ci.org/macinnir/swagger-gen)

# About

The goal behind this package is to recursively search for comments/tags within a project and build a swagger file based off of what it finds. It purposely does not use language-specific reflection (e.g. golang ast) so as to leave open the option of multi-language support (albeit a long-term goal).  

For the first couple of releases, it finds model fields based on go structs.  To the point of cross-language support, this will obviously need to change. 

The meta data at the top of the generated swagger file is provided by a file called `swagger-meta.json` which will need to exist in the root of your project. You can read more about this file below.

# Installation
```bash
go get github.com/macinnir/swagger-gen 
```

# Usage

## Initializing A `swagger-meta.json` File
swagger-gen uses a meta file called `swagger-meta.json` that it looks for in your root dir (specified by the `-s` flag, or the CWD). To quickly generate this file with default values, run the following:

```bash
# Create a `swagger-meta.json` file in the current directory
./swagger-gen -i
```

See: [Swagger-meta.json](#swagger-meta)

## Generating A Swagger File
Very helpful if you add your `path/to/go/bin` to your environment path so you can simply run `swagger-gen` from the cli. 

E.g. Windows: %GOPATH%/bin
*nix: $GOPATH/bin

```bash
# Build json swagger file from `src/dir` and write it to `dest/dir`
./swagger-gen -s src/dir -o dest/dir -f json
```

# CLI Args
Flag | Description | Values | Default 
---- | ----------- | ------ | -------
-i | __Input__ <br> Initializes a swagger-meta.json file with default values. It does not prompt for information (yet) so this is just a convenience method to build a placeholder file for you to put in your own information. <br>*Note: Does not work with other commands and will quit after the `swagger-meta.json` file has been generated.* | *none* | n/a
-s | __Source__ <br> The source directory of your code you want scanned. | *string* <br> filepath | `.` (Current directory)
-o | __Output__ <br> The output directory where you want the swagger spec (e.g. `swagger.json`) written to. | *string* <br> file path | `.` (Current Directory)
-f | __Format__ <br> The format of the output file. | *string* <br> `json` or `yaml` | `json` 

<a name="swagger-meta"></a>
# Swagger-meta.json

The swagger-meta.json file contains the top-most information found in a swagger.json file, just without including any of the paths or defintions that this application will generate. 

Described above in the **Options** section, there is an `-i` flag that creates a boilerplate file for convenience that you can update according to your needs.


An example of this file is like so: 

```json
{
    "swagger": "2.0",
    "info": {
        "description": "My API Description",
        "title": "My api title",
        "version": "0.1.0",
        "termsOfService": "TOS",
        "contact": {
            "email": "example@email.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        }
    },
    "host": "myhost.com",
    "basePath": "/v1",
    "tags": null,
    "schemes": [
        "http"
    ],
    "paths": null,
    "definitions": null
}
```

# Comments

## Routes 

### @route 

Positional Arguments for the `@route` tag:

- **OperationID** - String global name of the operation (e.g. `GetUsers`)
- **Method** - String Method (e.g. `GET`|`POST`|`PUT`|`DELETE`, etc.)
- **Route** - String route (e.g. `/users`)
- **Description** - String description of the route 

```go
// @route GetFoo GET /foo Returns a foo object 
```

### @param

Positional Arguments for the `@param` tag:

- **name*** - String name of the parameter
- **type*** - The data type of the parameter. See below on input data models
- **required** - Defaults to `required`. Should be one of two strings: *required* or *optional*
- **in:[transport]** - Defaults to `query`. `transport` should be any of the following:
    - path
    - query
    - form
    - header
    - body
- **description** - The description of the parameter.

Because only the first three arguments are required, parameters, therefore, can take any of the following forms:

```go
// Required path param  
// @param foo int in:path This is the foo param

// Required body param 
// @param foo int in:body This is the foo param

// Optional params
// @param foo int in:path This is the foo param
// @param foo int in:path optional This is the foo param
```

### @return 

Positional parameters for the `@return` tag:

- **ResponseCode** The numeric response code (e.g. `200`)
- **ResponseContent**
    - `empty` is for empty responses (e.g. for 204 no content)
    - Name of one or more models. An array of models is represented with a prefix of `[]`
- **Description** Description of the response 

```go
// @return 200 Foo Returns a Foo object 
// @return 200 []Foo Returns a collection of Foo objects 
// @return 204 empty The Foo object was created
// @return 400 ErrorObj There was an error when creating the Foo object
// @return 404 empty The foo object was not found 
```

## Models 

### @model

Positional parameters for the `@model` tag:
- **ModelName** Global identifier for the ModelName to be referenced when a route specifies an input and/or return model.


## @tag

Tags are comma separated 

```go
// @tag foo,bar,baz
```