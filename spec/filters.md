# 14. List of Builtin Filters — Jinja2 API

`abs`(*number*)

Return the absolute value of the argument.

`attr`(*obj*, *name*)

Get an attribute of an object. `foo|attr("bar")` works like `foo.bar` just that always an attribute is returned (nil if not found/defined).

`capitalize`(*s*)

Capitalize a value. The first character will be uppercased.

`default`(*value*, *default\_value=u''*, *boolean=False*)

If the value is nil it will return the passed default value, otherwise the value of the variable:

{{ my\_variable|default('my\_variable is not defined') }}

This will output the value of `my_variable` if the variable was not nil, otherwise `'my_variable is not defined'`. If you want to use default with variables that evaluate to false you have to set the second parameter to true:

{{ ''|default('the string was empty', true) }}

`filesizeformat`(*value*, *binary=False*)

Format the value like a ‘human-readable’ file size (i.e. 13 kB, 4.1 MB, 102 Bytes, etc). Per default decimal prefixes are used (Mega, Giga, etc.), if the second parameter is set to True the binary prefixes are used (Mebi, Gibi).

`first`(*seq*, *number=1*)

Return the first (or first number of) items of a slice.

`sprintf`(*value*, *\*args*, *\*\*kwargs*)

Apply sprintf string formatting on an object:

{{ 1/3|format("%0.2f") }}
    -> 0.33

`join`(*value*, *d=u''*, *attribute=None*)

Return a string which is the concatenation of the strings in the slice. The separator between elements is an empty string per default, you can define it with the optional parameter:

{{ [1, 2, 3]|join('|') }}
    -> 1|2|3

{{ [1, 2, 3]|join }}
    -> 123

It is also possible to join certain attributes of an object:

{{ users|join(', ', 'username') }}

New in version 2.6: The attribute parameter was added.

`split`(*str*, *d=u''*)

Return a slice which consists of the 1 or more elements. The separator between elements is an empty string per default, you can define it with the optional parameter:

{{ '1,2,3'|split(',')|join('|') }}
    -> 1|2|3

A string split with empty string splits the characters:

{{ '123'|split|join('|') }}
    -> 1|2|3

`last`(*seq*, *number=1*)

Return the last (or the last number of) items of a slice.

`length`(*object*)

Return the number of items of a slice or collection.

Aliases: `count`

`lower`(*s*)

Convert a value to lowercase.

`debug`(*value*)

Pretty print a variable. Useful for debugging.

Aliases: `d`

`replace`(*s*, *old*, *new*, *count=None*)

Return a copy of the value with all occurrences of a substring replaced with a new one. The first argument is the substring that should be replaced, the second is the replacement string. If the optional third argument `count` is given, only the first `count` occurrences are replaced:

{{ "Hello World"|replace("Hello", "Goodbye") }}
    -> Goodbye World

{{ "aaaaargh"|replace("a", "d'oh, ", 2) }}
    -> d'oh, d'oh, aaargh

`reverse`(*value*)

Reverse the object or return an iterator that iterates over it the other way round.

`round`(*value*, *precision=0*, *method='common'*)

Round the number to a given precision. The first parameter specifies the precision (default is `0`), the second the rounding method:

-   `'common'` (default, alias for up)
-   `'banker'` (alias for even)

-   `'ceil'` always rounds up
-   `'floor'` always rounds down
-   `'up'` rounds half up (common)
-   `'down'` rounds half down
-   `'odd'` rounds half to odd
-   `'even'` rounds half to even (banker)
-   `'awayzero'` round half away from zero
-   `'tozero'` round half towards zero

If you don’t specify a method `'common'` is used.

{{ 42.55|round }}
    -> 43.0
{{ 42.55|round(1, 'floor') }}
    -> 42.5

Note that even if rounded to 0 precision

`sum`(*iterable*, *attribute=None*)

Returns the sum of a slice of numbers plus the value of parameter ‘start’ (which defaults to 0). When the slice is empty it returns start.

{{ [1,2,3]|sum }}
    -> 6

It is also possible to sum up only certain attributes:

Total: {{ items|sum('price') }}

Changed in version 2.6: The attribute parameter was added to allow suming up over attributes. Also the start parameter was moved on to the right.

`title`(*s*)

Return a titlecased version of the value. I.e. words will start with uppercase letters, all remaining characters are left as is.

`trim`(*value*)

Strip leading and trailing whitespace.

`truncate`(*s*, *length=255*, *end='...'*)

Return a truncated copy of the string. The length is specified with the first parameter which defaults to `255`. If the text was in fact truncated it will append an ellipsis sign (`"..."`). If you want a different ellipsis sign than `"..."` you can specify it using the second parameter.

{{ "foo bar baz"|truncate(9) }}
    -> "foo ..."
{{ "foo bar baz"|truncate(9, '......') }}
    -> "foo ba......"

`upper`(*s*)

Convert a value to uppercase.

`urlencode`(*value*)

Escape strings for use in URLs (uses UTF-8 encoding). It accepts both dictionaries and regular strings as well as pairwise iterables.

---
Source: [14. List of Builtin Filters — Jinja2 API](https://tedboy.github.io/jinja2/templ14.html)