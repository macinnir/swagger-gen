# Swagger-Gen

# About

> The goal behind this package is to recursively search for tags within a project and build a swagger file based off of what it finds. It purposely does not use language-specific reflection (e.g. golang ast) so as to leave open the option of multi-language support. 
> In this initial commit, it finds model fields based on golang structs.  To the point of cross-language support, this will obviously need to change. Also, there are several hard-coded things in the code that need to be updated, so please treat this project as a PoC/WIP.  

# Comment Tags

```
    // @param [fieldName] [dataType] [required|optional] in:[transport] [description]
    // @param foo int required in:query This is the description of the required query int field `foo`
```

## @param

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
// GET, DELETE, PUT
// @param foo int

// Required body param 
// POST
// @param foo int

// @param foo int in:path
// @param foo int in:path optional
optional
required
```

## @tag

Tags are comma separated 

```
// @tag foo,bar,baz
```