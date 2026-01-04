# tqtemplate

## Overview

TQTemplate engine provides a simple yet powerful templating system with variable
interpolation, control structures, filters, and expression evaluation. Designed 
specifically for HTML output, all variable interpolation is automatically escaped 
for security.

## BNF Syntax

```bnf
<template>        ::= <content>*

<content>         ::= <literal> | <variable> | <control> | <comment>

<literal>         ::= any text not matching other patterns

<variable>        ::= "{{" <ws>? <expression> <filter-chain>? <ws>? "}}"

<control>         ::= <if-block> | <for-block> | <block> | <extends> | <include>

<comment>         ::= "{#" <any-text> "#}"

<extends>         ::= "{%" <ws>? "extends" <ws> <string> <ws>? "%}"

<include>         ::= "{%" <ws>? "include" <ws> <string> <ws>? "%}"

<block>           ::= <block-tag> <content>* <endblock-tag>

<block-tag>       ::= "{%" <ws>? "block" <ws> <identifier> <ws>? "%}"

<endblock-tag>    ::= "{%" <ws>? "endblock" <ws>? "%}"

<if-block>        ::= <if-tag> <content>* <elseif-tag>* <else-tag>? <endif-tag>

<if-tag>          ::= "{%" <ws>? "if" <ws> <expression> <filter-chain>? <ws>? "%}"

<elseif-tag>      ::= "{%" <ws>? "elseif" <ws> <expression> <filter-chain>? <ws>? "%}" <content>*

<else-tag>        ::= "{%" <ws>? "else" <ws>? "%}" <content>*

<endif-tag>       ::= "{%" <ws>? "endif" <ws>? "%}"

<for-block>       ::= <for-tag> <content>* <endfor-tag>

<for-tag>         ::= "{%" <ws>? "for" <ws> <for-vars> <ws> "in" <ws> <expression> <filter-chain>? <ws>? "%}"

<for-vars>        ::= <identifier> | <identifier> <ws>? "," <ws>? <identifier>

<endfor-tag>      ::= "{%" <ws>? "endfor" <ws>? "%}"

<expression>      ::= <logical-or>

<logical-or>      ::= <logical-and> (("or" | "||") <logical-and>)*

<logical-and>     ::= <equality> (("and" | "&&") <equality>)*

<equality>        ::= <comparison> (("==" | "!=") <comparison>)*

<comparison>      ::= <additive> (("<" | ">" | "<=" | ">=") <additive>)*

<additive>        ::= <multiplicative> (("+" | "-") <multiplicative>)*

<multiplicative>  ::= <unary> (("*" | "/" | "%") <unary>)*

<unary>           ::= "not" <unary> | <primary>

<primary>         ::= <number> | <string> | <path> | "(" <expression> ")"

<filter-chain>    ::= ("|" <filter>)+

<filter>          ::= <identifier> ("(" <filter-args> ")")?

<filter-args>     ::= <filter-arg> ("," <ws>? <filter-arg>)*

<filter-arg>      ::= <string> | <number> | <path>

<path>            ::= <identifier> ("." <identifier>)*

<identifier>      ::= [a-zA-Z_][a-zA-Z0-9_]*

<number>          ::= [0-9]+ ("." [0-9]+)?

<string>          ::= '"' (<char> | <escape-seq>)* '"'

<escape-seq>      ::= "\\" <any-char>

<ws>              ::= [ \t\n\r]+
```

## Operators

### Arithmetic Operators

- `+` Addition (also string concatenation)
- `-` Subtraction
- `*` Multiplication
- `/` Division
- `%` Modulo

### Comparison Operators

- `==` Equal
- `!=` Not equal
- `<` Less than
- `>` Greater than
- `<=` Less than or equal
- `>=` Greater than or equal

### Logical Operators

- `and`, `&&` Logical AND
- `or`, `||` Logical OR
- `not` Logical NOT (unary)

### Operator Precedence (highest to lowest)

1. `not` (unary)
2. `*`, `/`, `%`
3. `+`, `-`
4. `<`, `>`, `<=`, `>=`
5. `==`, `!=`
6. `and`, `&&`
7. `or`, `||`

## Features

- **Variable interpolation** with `{{ }}` syntax
- **Control structures** with `{% %}` syntax (if/elseif/else, for loops)
- **Template inheritance** with `{% extends %}` and `{% block %}`
- **Template inclusion** with `{% include %}`
- **Comments** with `{# #}` syntax
- **Expression evaluation** with full operator support
- **Filters** with pipe syntax `|`
- **Nested data access** with dot notation
- **HTML escaping** by default
- **Raw output** with `raw` filter
- **Custom functions** as filters

---

## Examples

### Example 1: Basic Variable Interpolation

**Data (JSON):**

```json
{
    "title": "Welcome",
    "username": "Alice",
    "message": "Hello, World!"
}
```

**Template:**

```html
<!DOCTYPE html>
<html>
    <head>
        <title>{{ title }}</title>
    </head>
    <body>
        <h1>{{ message }}</h1>
        <p>Logged in as: {{ username }}</p>
    </body>
</html>
```

**Output:**

```html
<!DOCTYPE html>
<html>
    <head>
        <title>Welcome</title>
    </head>
    <body>
        <h1>Hello, World!</h1>
        <p>Logged in as: Alice</p>
    </body>
</html>
```

---

### Example 2: HTML Escaping

**Data (JSON):**

```json
{
    "user_input": "<script>alert('XSS')</script>",
    "safe_html": "<strong>Bold Text</strong>"
}
```

**Template:**

```html
<div>
    <p>User input (escaped): {{ user_input }}</p>
    <p>Raw HTML: {{ safe_html|raw }}</p>
</div>
```

**Output:**

```html
<div>
    <p>
        User input (escaped):
        &lt;script&gt;alert(&#039;XSS&#039;)&lt;/script&gt;
    </p>
    <p>Raw HTML: <strong>Bold Text</strong></p>
</div>
```

---

### Example 3: Conditional Rendering

**Data (JSON):**

```json
{
    "user": {
        "name": "Bob",
        "is_admin": true,
        "age": 25
    }
}
```

**Template:**

```html
<div class="user-profile">
    <h2>{{ user.name }}</h2>

    {% if user.is_admin %}
    <span class="badge">Administrator</span>
    {% endif %}
    {% if user.age >= 18 %}
    <p>Adult user ({{ user.age }} years old)</p>
    {% else %}
    <p>Minor user ({{ user.age }} years old)</p>
    {% endif %}
</div>
```

**Output:**

```html
<div class="user-profile">
    <h2>Bob</h2>

    <span class="badge">Administrator</span>
    <p>Adult user (25 years old)</p>
</div>
```

---

### Example 4: If-ElseIf-Else Chain

**Data (JSON):**

```json
{
    "score": 85
}
```

**Template:**

```html
<div class="grade">
    {% if score >= 90 %}
    <span class="A">Grade: A - Excellent!</span>
    {% elseif score >= 80 %}
    <span class="B">Grade: B - Good Job!</span>
    {% elseif score >= 70 %}
    <span class="C">Grade: C - Fair</span>
    {% elseif score >= 60 %}
    <span class="D">Grade: D - Needs Improvement</span>
    {% else %}
    <span class="F">Grade: F - Failed</span>
    {% endif %}
</div>
```

**Output:**

```html
<div class="grade">
    <span class="B">Grade: B - Good Job!</span>
</div>
```

---

### Example 5: For Loops with Arrays

**Data (JSON):**

```json
{
    "fruits": ["Apple", "Banana", "Cherry", "Date"]
}
```

**Template:**

```html
<ul class="fruit-list">
    {% for fruit in fruits %}
    <li>{{ fruit }}</li>
    {% endfor %}
</ul>
```

**Output:**

```html
<ul class="fruit-list">
    <li>Apple</li>
    <li>Banana</li>
    <li>Cherry</li>
    <li>Date</li>
</ul>
```

---

### Example 6: For Loops with Key-Value Pairs

**Data (JSON):**

```json
{
    "products": {
        "laptop": "999.99",
        "mouse": "29.99",
        "keyboard": "79.99"
    }
}
```

**Template:**

```html
<table>
    <thead>
        <tr>
            <th>Product</th>
            <th>Price</th>
        </tr>
    </thead>
    <tbody>
        {% for product, price in products %}
        <tr>
            <td>{{ product }}</td>
            <td>${{ price }}</td>
        </tr>
        {% endfor %}
    </tbody>
</table>
```

**Output:**

```html
<table>
    <thead>
        <tr>
            <th>Product</th>
            <th>Price</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>laptop</td>
            <td>$999.99</td>
        </tr>
        <tr>
            <td>mouse</td>
            <td>$29.99</td>
        </tr>
        <tr>
            <td>keyboard</td>
            <td>$79.99</td>
        </tr>
    </tbody>
</table>
```

---

### Example 7: Nested For Loops

**Data (JSON):**

```json
{
    "grid": [
        [1, 2, 3],
        [4, 5, 6],
        [7, 8, 9]
    ]
}
```

**Template:**

```html
<table class="grid">
    {% for row in grid %}
    <tr>
        {% for cell in row %}
        <td>{{ cell }}</td>
        {% endfor %}
    </tr>
    {% endfor %}
</table>
```

**Output:**

```html
<table class="grid">
    <tr>
        <td>1</td>
        <td>2</td>
        <td>3</td>
    </tr>
    <tr>
        <td>4</td>
        <td>5</td>
        <td>6</td>
    </tr>
    <tr>
        <td>7</td>
        <td>8</td>
        <td>9</td>
    </tr>
</table>
```

---

### Example 8: Nested Data Access

**Data (JSON):**

```json
{
    "company": {
        "name": "Tech Corp",
        "employees": [
            {
                "name": "Alice",
                "position": "Developer",
                "salary": 80000
            },
            {
                "name": "Bob",
                "position": "Designer",
                "salary": 75000
            }
        ]
    }
}
```

**Template:**

```html
<div class="company">
    <h1>{{ company.name }}</h1>
    <h2>Employees</h2>
    <ul>
        {% for employee in company.employees %}
        <li>
            <strong>{{ employee.name }}</strong> - {{ employee.position }} (${{
            employee.salary }})
        </li>
        {% endfor %}
    </ul>
</div>
```

**Output:**

```html
<div class="company">
    <h1>Tech Corp</h1>
    <h2>Employees</h2>
    <ul>
        <li>
            <strong>Alice</strong> - Developer ($80000)
        </li>
        <li>
            <strong>Bob</strong> - Designer ($75000)
        </li>
    </ul>
</div>
```

---

### Example 9: Expressions in Variables

**Data (JSON):**

```json
{
    "price": 100,
    "quantity": 3,
    "tax_rate": 0.08
}
```

**Template:**

```html
<div class="invoice">
    <p>Price per item: ${{ price }}</p>
    <p>Quantity: {{ quantity }}</p>
    <p>Subtotal: ${{ price * quantity }}</p>
    <p>Tax (8%): ${{ price * quantity * tax_rate }}</p>
    <p>Total: ${{ price * quantity * (1 + tax_rate) }}</p>
</div>
```

**Output:**

```html
<div class="invoice">
    <p>Price per item: $100</p>
    <p>Quantity: 3</p>
    <p>Subtotal: $300</p>
    <p>Tax (8%): $24</p>
    <p>Total: $324</p>
</div>
```

---

### Example 10: String Concatenation

**Data (JSON):**

```json
{
    "first_name": "John",
    "last_name": "Doe",
    "title": "Dr."
}
```

**Template:**

```html
<div class="profile">
    <h2>{{ title + " " + first_name + " " + last_name }}</h2>
    <p>Full name: {{ first_name + " " + last_name }}</p>
</div>
```

**Output:**

```html
<div class="profile">
    <h2>Dr. John Doe</h2>
    <p>Full name: John Doe</p>
</div>
```

---

### Example 11: Complex Conditions

**Data (JSON):**

```json
{
    "user": {
        "age": 25,
        "is_premium": true,
        "credits": 150
    }
}
```

**Template:**

```html
<div class="access">
    {% if user.age >= 18 && user.is_premium %}
    <p>✓ Full access granted</p>
    {% endif %}
    {% if user.credits > 100 || user.is_premium %}
    <p>✓ Can download premium content</p>
    {% endif %}
    {% if (user.age >= 21 && user.credits > 50) || user.is_premium %}
    <p>✓ Can access exclusive features</p>
    {% endif %}
</div>
```

**Output:**

```html
<div class="access">
    <p>✓ Full access granted</p>
    <p>✓ Can download premium content</p>
    <p>✓ Can access exclusive features</p>
</div>
```

---

### Example 12: For Loop with Conditionals

**Data (JSON):**

```json
{
    "orders": [
        { "id": 1001, "status": "shipped", "total": 99.99 },
        { "id": 1002, "status": "pending", "total": 149.99 },
        { "id": 1003, "status": "delivered", "total": 79.99 },
        { "id": 1004, "status": "cancelled", "total": 199.99 }
    ]
}
```

**Template:**

```html
<table class="orders">
    <tr>
        <th>Order ID</th>
        <th>Total</th>
        <th>Status</th>
    </tr>
    {% for order in orders %}
    <tr class="{% if order.status == 'cancelled' %}cancelled{% elseif order.status == 'delivered' %}success{% endif %}">
        <td>#{{ order.id }}</td>
        <td>${{ order.total }}</td>
        <td>
            {% if order.status == "shipped" %}
            🚚 Shipped
            {% elseif order.status == "pending" %}
            ⏳ Pending
            {% elseif order.status == "delivered" %}
            ✓ Delivered
            {% else %}
            ✗ Cancelled
            {% endif %}
        </td>
    </tr>
    {% endfor %}
</table>
```

**Output:**

```html
<table class="orders">
    <tr>
        <th>Order ID</th>
        <th>Total</th>
        <th>Status</th>
    </tr>
    <tr class="">
        <td>#1001</td>
        <td>$99.99</td>
        <td>
            🚚 Shipped
        </td>
    </tr>
    <tr class="">
        <td>#1002</td>
        <td>$149.99</td>
        <td>
            ⏳ Pending
        </td>
    </tr>
    <tr class="success">
        <td>#1003</td>
        <td>$79.99</td>
        <td>
            ✓ Delivered
        </td>
    </tr>
    <tr class="cancelled">
        <td>#1004</td>
        <td>$199.99</td>
        <td>
            ✗ Cancelled
        </td>
    </tr>
</table>
```

---

### Example 13: Comments

**Data (JSON):**

```json
{
    "username": "Alice",
    "email": "alice@example.com"
}
```

**Template:**

```html
<div class="user">
    {# This is a comment and won't appear in output #}
    <h2>{{ username }}</h2>

    {# Multi-line comment These can span multiple lines and won't be rendered #}
    <p>Email: {{ email }}</p>
    {# TODO: Add phone number field #}
</div>
```

**Output:**

```html
<div class="user">
    <h2>Alice</h2>

    <p>Email: alice@example.com</p>
</div>
```

---

### Example 14: Blog Post List

**Data (JSON):**

```json
{
    "blog": {
        "title": "My Tech Blog",
        "posts": [
            {
                "id": 1,
                "title": "Getting Started with PHP",
                "author": "Alice",
                "date": "2024-01-15",
                "excerpt": "Learn the basics of PHP programming...",
                "published": true,
                "views": 1234
            },
            {
                "id": 2,
                "title": "Advanced Template Engines",
                "author": "Bob",
                "date": "2024-01-20",
                "excerpt": "Deep dive into template engine design...",
                "published": true,
                "views": 856
            },
            {
                "id": 3,
                "title": "Upcoming Features",
                "author": "Alice",
                "date": "2024-02-01",
                "excerpt": "What's coming next...",
                "published": false,
                "views": 0
            }
        ]
    }
}
```

**Template:**

```html
<!DOCTYPE html>
<html>
    <head>
        <title>{{ blog.title }}</title>
    </head>
    <body>
        <header>
            <h1>{{ blog.title }}</h1>
        </header>

        <main>
            {% for post in blog.posts %} {% if post.published %}
            <article class="post">
                <h2>{{ post.title }}</h2>
                <div class="meta">
                    By {{ post.author }} on {{ post.date }}
                    {% if post.views > 1000 %}
                    <span class="popular">🔥 Popular</span>
                    {% endif %}
                </div>
                <p>{{ post.excerpt }}</p>
                <a href="/post/{{ post.id }}">Read more...</a>
            </article>
            {% endif %} {% endfor %}
        </main>
    </body>
</html>
```

**Output:**

```html
<!DOCTYPE html>
<html>
    <head>
        <title>My Tech Blog</title>
    </head>
    <body>
        <header>
            <h1>My Tech Blog</h1>
        </header>

        <main>
            <article class="post">
                <h2>Getting Started with PHP</h2>
                <div class="meta">
                    By Alice on 2024-01-15
                    <span class="popular">🔥 Popular</span>
                </div>
                <p>Learn the basics of PHP programming...</p>
                <a href="/post/1">Read more...</a>
            </article>
            <article class="post">
                <h2>Advanced Template Engines</h2>
                <div class="meta">
                    By Bob on 2024-01-20
                </div>
                <p>Deep dive into template engine design...</p>
                <a href="/post/2">Read more...</a>
            </article>
        </main>
    </body>
</html>
```

---

### Example 15: Dashboard with Statistics

**Data (JSON):**

```json
{
    "dashboard": {
        "user": "Admin",
        "stats": {
            "total_users": 1523,
            "active_users": 892,
            "total_revenue": 45678.90,
            "pending_orders": 23
        },
        "recent_activities": [
            {
                "user": "Alice",
                "action": "registered",
                "time": "2 minutes ago"
            },
            {
                "user": "Bob",
                "action": "made a purchase",
                "time": "5 minutes ago"
            },
            {
                "user": "Charlie",
                "action": "updated profile",
                "time": "10 minutes ago"
            }
        ]
    }
}
```

**Template:**

```html
<!DOCTYPE html>
<html>
    <head>
        <title>Admin Dashboard</title>
    </head>
    <body>
        <h1>Welcome, {{ dashboard.user }}</h1>

        <div class="stats-grid">
            <div class="stat-card">
                <h3>Total Users</h3>
                <p class="number">{{ dashboard.stats.total_users }}</p>
            </div>

            <div class="stat-card">
                <h3>Active Users</h3>
                <p class="number">{{ dashboard.stats.active_users }}</p>
                <small>{{ dashboard.stats.active_users * 100 /
                    dashboard.stats.total_users }}% active</small>
            </div>

            <div class="stat-card {% if dashboard.stats.total_revenue > 40000 %}success{% endif %}">
                <h3>Revenue</h3>
                <p class="number">${{ dashboard.stats.total_revenue }}</p>
            </div>

            <div class="stat-card {% if dashboard.stats.pending_orders > 20 %}warning{% endif %}">
                <h3>Pending Orders</h3>
                <p class="number">{{ dashboard.stats.pending_orders }}</p>
            </div>
        </div>

        <div class="recent-activity">
            <h2>Recent Activity</h2>
            <ul>
                {% for activity in dashboard.recent_activities %}
                <li>
                    <strong>{{ activity.user }}</strong> {{ activity.action }}
                    <span class="time">{{ activity.time }}</span>
                </li>
                {% endfor %}
            </ul>
        </div>
    </body>
</html>
```

**Output:**

```html
<!DOCTYPE html>
<html>
    <head>
        <title>Admin Dashboard</title>
    </head>
    <body>
        <h1>Welcome, Admin</h1>

        <div class="stats-grid">
            <div class="stat-card">
                <h3>Total Users</h3>
                <p class="number">1523</p>
            </div>

            <div class="stat-card">
                <h3>Active Users</h3>
                <p class="number">892</p>
                <small>58.568611293499% active</small>
            </div>

            <div class="stat-card success">
                <h3>Revenue</h3>
                <p class="number">$45678.9</p>
            </div>

            <div class="stat-card warning">
                <h3>Pending Orders</h3>
                <p class="number">23</p>
            </div>
        </div>

        <div class="recent-activity">
            <h2>Recent Activity</h2>
            <ul>
                <li>
                    <strong>Alice</strong> registered
                    <span class="time">2 minutes ago</span>
                </li>
                <li>
                    <strong>Bob</strong> made a purchase
                    <span class="time">5 minutes ago</span>
                </li>
                <li>
                    <strong>Charlie</strong> updated profile
                    <span class="time">10 minutes ago</span>
                </li>
            </ul>
        </div>
    </body>
</html>
```

---

### Example 16: Template Inheritance with Extends and Blocks

Template inheritance allows you to build a base "skeleton" template that
contains common elements of your site and defines **blocks** that child
templates can override.

**Base Template (base.html):**

```html
<!DOCTYPE html>
<html>
    <head>
        <title>{% block title %}Default Title{% endblock %}</title>
        {% block head %}{% endblock %}
    </head>
    <body>
        <header>
            <h1>My Website</h1>
            <nav>
                <a href="/">Home</a>
                <a href="/about">About</a>
            </nav>
        </header>

        <main>
            {% block content %}
            <p>Default content</p>
            {% endblock %}
        </main>

        <footer>
            <p>&copy; 2026 My Website</p>
        </footer>
    </body>
</html>
```

**Child Template:**

```html
{% extends 'base.html' %} {% block title %}Welcome Page{% endblock %} {% block
head %}
<style>
    .highlight {
        color: blue;
    }
</style>
{% endblock %} {% block content %}
<h2>Welcome to Our Site!</h2>
<p class="highlight">This content replaces the default content block.</p>
{% endblock %}
```

**Data (JSON):**

```json
{}
```

**Output:**

```html
<!DOCTYPE html>
<html>
    <head>
        <title>Welcome Page</title>
        <style>
            .highlight {
                color: blue;
            }
        </style>
    </head>
    <body>
        <header>
            <h1>My Website</h1>
            <nav>
                <a href="/">Home</a>
                <a href="/about">About</a>
            </nav>
        </header>

        <main>
            <h2>Welcome to Our Site!</h2>
            <p class="highlight">
                This content replaces the default content block.
            </p>
        </main>

        <footer>
            <p>&copy; 2026 My Website</p>
        </footer>
    </body>
</html>
```

**Notes:**

- The `{% extends %}` directive must be the first non-whitespace element in the
  child template
- Child blocks completely replace parent block content
- Blocks not overridden in the child will use the parent's default content
- Template inheritance requires a `TemplateLoader` function to load parent
  templates

---

### Example 17: Nested Blocks

Blocks can be nested within other blocks, and child templates can override
nested blocks.

**Base Template (layout.html):**

```html
<div class="page">
    {% block outer %}
    <div class="container">
        {% block inner %}
        <p>Default inner content</p>
        {% endblock %}
    </div>
    {% endblock %}
</div>
```

**Child Template:**

```html
{% extends 'layout.html' %} {% block inner %}
<h1>Custom Inner Content</h1>
{% endblock %}
```

**Output:**

```html
<div class="page">
    <div class="container">
        <h1>Custom Inner Content</h1>
    </div>
</div>
```

**Notes:**

- When only the inner block is overridden, the outer block structure from the
  parent is preserved
- Each block is independently overrideable

---

### Example 18: Multiple Block Overrides with Variables

**Base Template (blog-layout.html):**

```html
<article>
    <header>
        <h1>{% block title %}Untitled{% endblock %}</h1>
        <div class="meta">
            {% block meta %}
            <span>No metadata</span>
            {% endblock %}
        </div>
    </header>
    <div class="content">
        {% block content %}No content{% endblock %}
    </div>
</article>
```

**Child Template:**

```html
{% extends 'blog-layout.html' %} {% block title %}{{ post.title }}{% endblock %}
{% block meta %}
<span class="author">By {{ post.author }}</span>
<span class="date">{{ post.date }}</span>
{% endblock %} {% block content %}
<p>{{ post.body }}</p>
<p>Tags: {{ post.tags }}</p>
{% endblock %}
```

**Data (JSON):**

```json
{
    "post": {
        "title": "Introduction to Templates",
        "author": "Jane Doe",
        "date": "2026-01-04",
        "body": "This is the post body with important information.",
        "tags": "templates, documentation"
    }
}
```

**Output:**

```html
<article>
    <header>
        <h1>Introduction to Templates</h1>
        <div class="meta">
            <span class="author">By Jane Doe</span>
            <span class="date">2026-01-04</span>
        </div>
    </header>
    <div class="content">
        <p>This is the post body with important information.</p>
        <p>Tags: templates, documentation</p>
    </div>
</article>
```

**Notes:**

- Variables can be used within block content
- Multiple blocks can be overridden in a single child template
- Parent template defines the overall structure; child provides specific content

---

### Example 19: Template Inclusion with Include

The `{% include %}` directive inserts another template's content at that point
in the current template.

**Header Template (header.html):**

```html
<header>
    <h1>My Website</h1>
    <nav>
        <a href="/">Home</a>
        <a href="/about">About</a>
    </nav>
</header>
```

**Footer Template (footer.html):**

```html
<footer>
    <p>&copy; 2026 My Website. All rights reserved.</p>
</footer>
```

**Main Template:**

```html
<!DOCTYPE html>
<html>
    <head>
        <title>{{ title }}</title>
    </head>
    <body>
        {% include 'header.html' %}

        <main>
            <h2>{{ heading }}</h2>
            <p>{{ content }}</p>
        </main>

        {% include 'footer.html' %}
    </body>
</html>
```

**Data (JSON):**

```json
{
    "title": "Welcome Page",
    "heading": "Welcome!",
    "content": "This is the main content of the page."
}
```

**Output:**

```html
<!DOCTYPE html>
<html>
    <head>
        <title>Welcome Page</title>
    </head>
    <body>
        <header>
            <h1>My Website</h1>
            <nav>
                <a href="/">Home</a>
                <a href="/about">About</a>
            </nav>
        </header>

        <main>
            <h2>Welcome!</h2>
            <p>This is the main content of the page.</p>
        </main>

        <footer>
            <p>&copy; 2026 My Website. All rights reserved.</p>
        </footer>
    </body>
</html>
```

**Notes:**

- Included templates have access to the same data and functions as the parent
- Multiple includes can be used in a single template
- Includes are useful for reusable components like headers, footers, and
  sidebars

---

### Example 20: Include with Control Structures

Included templates can contain any template features, including loops and
conditionals.

**Item List Template (items.html):**

```html
{% for item in items %}
<div class="item">
    <h3>{{ item.name }}</h3>
    <p>{{ item.description }}</p>
    {% if item.price %}
    <span class="price">${{ item.price }}</span>
    {% endif %}
</div>
{% endfor %}
```

**Main Template:**

```html
<div class="catalog">
    <h1>Product Catalog</h1>
    {% include 'items.html' %}
</div>
```

**Data (JSON):**

```json
{
    "items": [
        {
            "name": "Widget",
            "description": "A useful widget",
            "price": 29.99
        },
        {
            "name": "Gadget",
            "description": "An amazing gadget",
            "price": 49.99
        },
        {
            "name": "Free Sample",
            "description": "Try it free"
        }
    ]
}
```

**Output:**

```html
<div class="catalog">
    <h1>Product Catalog</h1>
    <div class="item">
        <h3>Widget</h3>
        <p>A useful widget</p>
        <span class="price">$29.99</span>
    </div>
    <div class="item">
        <h3>Gadget</h3>
        <p>An amazing gadget</p>
        <span class="price">$49.99</span>
    </div>
    <div class="item">
        <h3>Free Sample</h3>
        <p>Try it free</p>
    </div>
</div>
```

**Notes:**

- Included templates can use loops, conditionals, and all other template
  features
- Data is shared between the main template and included templates
- This is useful for creating reusable list/card components

---

### Example 21: Nested Includes

Included templates can themselves include other templates.

**Icon Template (icon.html):**

```html
<i class="icon-{{ type }}"></i>
```

**Button Template (button.html):**

```html
<button class="{{ style }}">
    {% include 'icon.html' %} {{ text }}
</button>
```

**Main Template:**

```html
<div class="actions">
    {% include 'button.html' %}
</div>
```

**Data (JSON):**

```json
{
    "type": "save",
    "style": "btn-primary",
    "text": "Save Changes"
}
```

**Output:**

```html
<div class="actions">
    <button class="btn-primary">
        <i class="icon-save"></i>
        Save Changes
    </button>
</div>
```

**Notes:**

- Nested includes allow building complex components from simple parts
- All included templates share the same data context
- This enables component-based template design

---

## Builtin Filters

TQTemplate includes a comprehensive set of builtin filters for common data transformations. Filters are applied using the pipe (`|`) syntax.

### String Filters

#### `lower`

Convert a string to lowercase.

```
{{ "HELLO"|lower }}
→ hello
```

#### `upper`

Convert a string to uppercase.

```
{{ "hello"|upper }}
→ HELLO
```

#### `capitalize`

Capitalize the first character of a string.

```
{{ "hello world"|capitalize }}
→ Hello world
```

#### `title`

Convert a string to title case (capitalize first letter of each word).

```
{{ "hello world"|title }}
→ Hello World
```

#### `trim`

Remove leading and trailing whitespace.

```
{{ "  hello  "|trim }}
→ hello
```

#### `truncate(length, end)`

Truncate a string to a maximum length. Default length is 255, default end is "...".

```
{{ "Hello World"|truncate(8) }}
→ Hello...

{{ "Hello World"|truncate(8, ">>") }}
→ Hello W>>
```

#### `replace(old, new, count)`

Replace occurrences of a substring. If count is omitted, replaces all occurrences.

```
{{ "Hello World"|replace("Hello", "Goodbye") }}
→ Goodbye World

{{ "aaaaargh"|replace("a", "d'oh, ", 2) }}
→ d'oh, d'oh, aaargh
```

#### `split(separator)`

Split a string into a slice. Empty separator splits into characters.

```
{{ "1,2,3"|split(",")|join("|") }}
→ 1|2|3

{{ "123"|split|join("|") }}
→ 1|2|3
```

#### `urlencode`

Encode a string for use in URLs.

```
{{ "hello world"|urlencode }}
→ hello+world

{{ "hello&world=test"|urlencode }}
→ hello%26world%3Dtest
```

### Numeric Filters

#### `abs`

Return the absolute value of a number.

```
{{ -42|abs }}
→ 42

{{ 3.14|abs }}
→ 3.14
```

#### `round(precision, method)`

Round a number to a given precision. Default precision is 0, default method is "common".

Available methods:
- `common` or `up` - Round half up (default)
- `ceil` - Always round up
- `floor` - Always round down
- `down` - Round half down
- `even` or `banker` - Round half to even
- `odd` - Round half to odd
- `awayzero` - Round half away from zero
- `tozero` - Round half towards zero

```
{{ 42.55|round }}
→ 43

{{ 42.55|round(1, "floor") }}
→ 42.5

{{ 2.5|round(0, "even") }}
→ 2
```

#### `sprintf(format)`

Apply sprintf-style formatting to a value.

```
{{ 3.14159|sprintf("%.2f") }}
→ 3.14

{{ 42|sprintf("%05d") }}
→ 00042
```

#### `filesizeformat(binary)`

Format a number as a human-readable file size. Use binary=true for binary prefixes (KiB, MiB).

```
{{ 13000|filesizeformat }}
→ 13.0 kB

{{ 1024|filesizeformat(true) }}
→ 1.0 KiB

{{ 1500000|filesizeformat }}
→ 1.5 MB
```

### Array/Collection Filters

#### `length` / `count`

Return the number of items in a collection.

```
{{ [1, 2, 3]|length }}
→ 3

{{ "hello"|length }}
→ 5
```

#### `first(n)`

Return the first item or first n items of a slice.

```
{{ [1, 2, 3, 4]|first }}
→ 1

{{ [1, 2, 3, 4]|first(2) }}
→ [1, 2]
```

#### `last(n)`

Return the last item or last n items of a slice.

```
{{ [1, 2, 3, 4]|last }}
→ 4

{{ [1, 2, 3, 4]|last(2) }}
→ [3, 4]
```

#### `join(separator, attribute)`

Concatenate items in a slice with a separator. Can optionally join by attribute.

```
{{ [1, 2, 3]|join("|") }}
→ 1|2|3

{{ [1, 2, 3]|join }}
→ 123

{{ users|join(", ", "username") }}
→ alice, bob, charlie
```

#### `reverse`

Reverse a slice or string.

```
{{ [1, 2, 3]|reverse }}
→ [3, 2, 1]

{{ "hello"|reverse }}
→ olleh
```

#### `sum(attribute)`

Return the sum of numbers in a slice. Can optionally sum by attribute.

```
{{ [1, 2, 3]|sum }}
→ 6

{{ items|sum("price") }}
→ 150.50
```

### Utility Filters

#### `default(value, boolean)`

Return a default value if the input is nil or (with boolean=true) falsy.

```
{{ missing_var|default("N/A") }}
→ N/A

{{ ""|default("empty", true) }}
→ empty

{{ 0|default("zero", true) }}
→ zero
```

#### `attr(name)`

Get an attribute of an object by name.

```
{{ user|attr("email") }}
→ user@example.com

{{ data|attr("nested")|attr("value") }}
→ 42
```

#### `debug` / `d`

Pretty print a value for debugging (outputs JSON).

```
{{ user|debug }}
→ {
  "name": "Alice",
  "email": "alice@example.com"
}
```

#### `raw`

Mark a value as safe HTML that should not be escaped.

```
{{ "<strong>Bold</strong>"|raw }}
→ <strong>Bold</strong>
```

### Filter Chaining

Filters can be chained together:

```
{{ "  hello world  "|trim|upper|replace("WORLD", "FRIEND") }}
→ HELLO FRIEND

{{ items|first(3)|reverse|join(", ") }}
→ 3, 2, 1
```

---

## Custom Filters

Filters can be provided as custom functions when rendering templates. All builtin
filters listed above are available automatically.

**PHP Usage Example:**

```php
$template = new Template('html');

$data = ['name' => 'john doe', 'date' => 'May 13, 1980'];

$functions = [
    'dateFormat' => fn($date, $format) => date($format, strtotime($date))
];

$html = $template->render(
    'Hello {{ name|upper }}, date: {{ date|dateFormat("Y-m-d") }}',
    $data,
    $functions
);
// Output: Hello JOHN DOE, date: 1980-05-13
```

---

## Builtin Tests

TQTemplate supports Jinja2-style tests using the `is` keyword. Tests are used to check properties of values, particularly useful in conditional expressions.

### Syntax

Tests are used with the `is` keyword:

```
{% if variable is testname %}
{% if variable is not testname %}
{% if variable is testname(arg) %}
```

### Available Tests

#### `defined`

Check if a variable is defined (not nil).

```
{% if user is defined %}
    User exists: {{ user }}
{% else %}
    No user found
{% endif %}
```

#### `undefined`

Check if a variable is undefined (nil). Opposite of `defined`.

```
{% if missing is undefined %}
    Variable is not defined
{% endif %}
```

#### `null`

Check if a value is nil.

```
{% if value is null %}
    Value is null
{% endif %}
```

#### `even`

Check if a number is even.

```
{% if count is even %}
    Count is even: {{ count }}
{% endif %}
```

#### `odd`

Check if a number is odd.

```
{% if count is odd %}
    Count is odd: {{ count }}
{% endif %}
```

#### `divisibleby(n)`

Check if a number is divisible by n.

```
{% if total is divisibleby(3) %}
    Total is divisible by 3
{% endif %}
```

#### `number`

Check if a value is a number (int, float, or numeric string).

```
{% if value is number %}
    {{ value }} is a number
{% endif %}
```

#### `string`

Check if a value is a string.

```
{% if name is string %}
    {{ name }} is a string
{% endif %}
```

#### `iterable`

Check if a value can be iterated over (array, map, or string).

```
{% if items is iterable %}
    {% for item in items %}
        {{ item }}
    {% endfor %}
{% endif %}
```

### Test Negation

Tests can be negated using `is not`:

```
{% if value is not null %}
    Value exists
{% endif %}

{% if count is not even %}
    Count is odd
{% endif %}
```

### Tests in Variable Expressions

Tests can also be used in variable expressions and will output `1` for true or empty string for false:

```
{{ num is even }}  {# outputs "1" if num is even, "" otherwise #}
```

---

## Notes

- All output is **HTML-escaped by default** for security
- Use the `raw` filter to output unescaped HTML: `{{ content|raw }}`
- Whitespace in templates is generally preserved
- Lines containing only whitespace and a `{% %}` tag are removed
- Expressions support parentheses for grouping: `{{ (a + b) * c }}`
- Paths use dot notation for nested access: `{{ user.profile.name }}`
- For loops can iterate with values only or with key-value pairs
- Comments are completely removed from output and don't affect whitespace

### Template Inheritance Notes

- The `{% extends %}` directive must be the first non-whitespace element in a
  child template
- Template inheritance requires a `TemplateLoader` function to be configured
- Child block content completely replaces parent block content (no indentation
  inheritance)
- Blocks not overridden in the child will use the parent's default content
- Blocks can be nested, and each can be independently overridden
- Variables, expressions, and all other template features work inside blocks
