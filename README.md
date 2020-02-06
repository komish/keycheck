# KeyCheck

## Description
Check a YAML or JSON document for the presence or absence of values at a specific path. Useful for checking a file
for the existence of a path that has been deprecated, or the absence of a path that is required/recommended.

## Specification File Format

This format can be either JSON or YAML and must contain two keys [`path`, `msg`] and an optional key [`required`].

### Example Spec File

```YAML
# search for the presence of key `bar` inside of `foo` and print a message if it is found
- path: "foo.bar"
  msg: "Foos can no longer have bars. Sorry for the inconvenience. Read more at https://example.com/more-info"
# search for the absence of key `that` inside of `this` and print a message if it is not found.
- path: "this.that"
  msg: "this must have a that"
  required: true
```

### Example Execution

The demonstrated spec and target files can be found in the [samples](samples/) directory.

Parsing simple yaml structures. Some comments in-line as the syntax can be a bit confusing.

```
# keycheck --specfile samples/01-specfile.yaml samples/01-target.yaml 
============> Results for file: samples/01-target.yaml

    # Here we check array "fruits" for a "grapefruit" item.
   Item  groceries.fruits.#.grapefruit
Message  Why is grapefruit on the list! I don't like grapefruit.

    # Here we check for specific key "celery" in the veg array.
   Item  groceries.veg.#(celery)
Message  Celery not on the list! I'm trying to eat more celery.

    # Here we check for the existence of "soda" in the "drinks" map.
   Item  groceries.drinks.soda
Message  Skip buying sodas, drink tea instead.
```

You can also pass in multiple files to be checked against a given spec file. This example checks multiple files and the help message is more informative.

```
# keycheck --specfile samples/02-specfile.yaml samples/02-target-*.yaml
============> Results for file: samples/02-target-01.yaml

   Item  config.hostPort
Message  FOUND config.hostPort has been deprecated. Please use separate `config.host` and `config.port`

   Item  config.host
Message  MISSING config.host is required in upcoming releases and backwards compatible with current stable.

   Item  config.port
Message  MISSING config.port is required in upcoming releases and backwards compatible with current stable.

============> Results for file: samples/02-target-02.json

   Item  config.appName
Message  MISSING config.appName is required in upcoming releases and backwards compatible with current stable.
```

## Parsing Syntax

The parsing logic is passing things directly to [gjson](https://github.com/tidwall/gjson#path-syntax) for dot notation parsing. Refer to this project's path syntax documentation for more information.

## Known Issues

* Given an array, this currently does not check to see that all items in the array have a given key. 
* Currently there is no way for this to return the value in the found key.
* If a spec item has a typographical error, keycheck will fail silently. It's impossible to confirm if the key should or should not exist and a typo is effectively an unexpected or missing key.