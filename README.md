# Swagger-Gen

### Status
[![Build Status](https://travis-ci.org/macinnir/swagger-gen.svg?branch=master)](https://travis-ci.org/macinnir/swagger-gen)

# About

> The goal behind this package is to recursively search for tags within a project and build a swagger file based off of what it finds. It purposely does not use language-specific reflection (e.g. golang ast) so as to leave open the option of multi-language support. 

> In this initial commit, it finds model fields based on golang structs.  To the point of cross-language support, this will obviously need to change. Also, there are several hard-coded things in the code that need to be updated, so please treat this project as a PoC/WIP.  

# Installation
```bash
go get github.com/macinnir/swagger-gen 

```

# Usage
Very helpful if you add your `path/to/go/bin` to your environment path so you can simply run `swagger-gen` from the cli. 

E.g. Windows: %GOPATH%/bin
*nix: $GOPATH/bin

```
./swagger-gen -s src/dir -o dest/dir -f json
```

## Options
Name | Flag | Description | Default 
---- | ---- | ----------- | -------
Init | -i | Initializes a swagger-meta.json file with default values. It does not prompt for information (yet) so this is just a convenience method to build a placeholder file for you to put in your own information. *Does not work with other commands and will quit after the `swagger-meta.json` file has been generated.* | n/a
Source | -s | The source directory of your code you want scanned. | `.` (Current directory)
Output | -o | The output directory where you want the swagger spec (e.g. `swagger.json`) written to. | `.` (Current Directory)
Format | -f | The format of the output file. Possible values: `json` or `yaml` | `json` 

# Comments

## Routes 

### @route 

Positional Arguments for the `@route` tag:

- **OperationID** - String global name of the operation (e.g. `GetUsers`)
- **Method** - String Method (e.g. `GET`|`POST`|`PUT`|`DELETE`, etc.)
- **Route** - String route (e.g. `/users`)
- **Description** - String description of the route 

```
    // @route GetFoo /foo Returns a foo object 
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

```
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

```
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

```
// @tag foo,bar,baz
```