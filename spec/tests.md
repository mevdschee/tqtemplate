# 15. List of Builtin Tests — Jinja2 API

## Implementation Status

✅ **All tests implemented in tqtemplate**

This Go implementation includes all the tests listed below, following Jinja2 semantics.

## Usage

Tests are used with the `is` keyword in conditional expressions:

```
{% if variable is defined %}
    value of variable: {{ variable }}
{% else %}
    variable is not defined
{% endif %}
```

Tests can be negated with `is not`:

```
{% if variable is not null %}
    Variable has a value
{% endif %}
```

## Go Implementation Notes

- All tests work as builtin filters that are automatically registered
- The `is` keyword is preprocessed and converted to filter calls internally
- Missing variables are handled gracefully by `defined` and `undefined` tests
- Tests can be used in `{% if %}`, `{% elseif %}`, and variable expressions (`{{ }}`)
- Type checking is flexible: numeric strings are treated as numbers where appropriate

---

## Test List

`defined`(*value*)

Return true if the variable is defined:

{% if variable is defined %}
    value of variable: {{ variable }}
{% else %}
    variable is not defined
{% endif %}

See the `attr()` (and `default`) filters for a simple way to set undefined variables.

`divisibleby`(*value*, *num*)

Check if a variable is divisible by a number.

`even`(*value*)

Return true if the variable is even.

`iterable`(*value*)

Check if it’s possible to iterate over an object.

`lower`(*value*)

Return true if the variable is lowercased.

`null`(*value*)

Return true if the variable is nil.

**Note:** In tqtemplate, `none` is also available as an alias for `null`.

`number`(*value*)

Return true if the variable is a number.

**Note:** In tqtemplate, this includes integers, floats, and numeric strings (e.g., "123").

`odd`(*value*)

Return true if the variable is odd.

`string`(*value*)

Return true if the object is a string.

`undefined`(*value*)

Like `defined()` but the other way round.

`upper`(*value*)

Return true if the variable is uppercased.

---
Source: [15. List of Builtin Tests — Jinja2 API](https://tedboy.github.io/jinja2/templ15.html)