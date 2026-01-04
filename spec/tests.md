# 15. List of Builtin Tests — Jinja2 API

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

`number`(*value*)

Return true if the variable is a number.

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