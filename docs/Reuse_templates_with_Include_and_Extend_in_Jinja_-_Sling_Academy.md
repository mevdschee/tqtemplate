# Reuse templates with Include and Extend in Jinja - Sling Academy

Jinja is a modern and designer-friendly templating engine for Python. It makes it easy to create dynamic HTML pages by allowing developers to use familiar concepts like loops, conditionals, and flexible inheritance models. Two of the main ways to achieve template reusability in Jinja are through the **"include"** and **"extend"** functionality. In this article, we'll explore how these work and look at some examples to guide us through their usage.

## Table of Contents

1.  [Using the "include" Directive](https://www.slingacademy.com/article/reuse-templates-with-include-and-extend-in-jinja/#using-the-)
    1.  [Example of Using "include"](https://www.slingacademy.com/article/reuse-templates-with-include-and-extend-in-jinja/#example-of-using-)
2.  [Using the "extend" Directive](https://www.slingacademy.com/article/reuse-templates-with-include-and-extend-in-jinja/#using-the-)
    1.  [Example of Using "extend"](https://www.slingacademy.com/article/reuse-templates-with-include-and-extend-in-jinja/#example-of-using-)
3.  [A Note About Performance](https://www.slingacademy.com/article/reuse-templates-with-include-and-extend-in-jinja/#a-note-about-performance)

## Using the "include" Directive

The `include` directive in Jinja allows you to import a template and render it within another template. This is beneficial when you have pieces of HTML that are repeated across different pages, such as headers, footers, or any component that is used multiple times but does not require complex customization.

### Example of Using "include"

Suppose you have a header template *(header.html)* and you want to include it in your main page template *(index.html)*.

```


### header.html ###
<header>
  <h1>Welcome to Our Site!</h1>
  <nav>
    <ul>
      <li><a href="/">Home</a></li>
      <li><a href="/about">About</a></li>
    </ul>
  </nav>
</header>
```

```


### index.html ###
<html>
  <body>
    {% include 'header.html' %}
    <main>
      <p>This is the main content of the page.</p>
    </main>
  </body>
</html>
```

In this example, when *index.html* is rendered, it will include the content from *header.html* at the location specified by `{% include 'header.html' %}`.

## Using the "extend" Directive

The `extend` directive is a foundational concept in template inheritance in Jinja. It allows one template to "inherit" from another template and override certain parts of it. This is useful for creating complex systems with consistent layouts.

### Example of Using "extend"

Let's assume you want all pages to use a common structure defined in *base.html*.

```


### base.html ###
<html>
  <head>
    <title>{% block title %}My Website{% endblock %}</title>
  </head>
  <body>
    <header>
      {% block header %}<h1>Header Title</h1>{% endblock %}
    </header>
    <main>
      {% block content %}{% endblock %}
    </main>
    <footer>
      {% block footer %}<p>Footer Information</p>{% endblock %}
    </footer>
  </body>
</html>
```

```


### page.html ###
{% extends 'base.html' %}

{% block title %}Home Page - My Website{% endblock %}

{% block header %}
  <h1>Customized Header for Home Page</h1>
{% endblock %}

{% block content %}
  <p>This is the home page content.</p>
{% endblock %}
```

Here, *page.html* extends *base.html* and provides its own content for the `title`, `header`, and `content` blocks. When rendered, *page.html* will retain the layout defined in *base.html* but insert the specified blocks at the designated placeholders.

## A Note About Performance

Frequent use of `include` can slow down rendering since each *included* file might require additional overhead. The same care should thus be taken to balance performance with template organization.

In this article, we explored how Jinja's **"include"** and **"extend"** functionality facilitates reusable and manageable templates. By leveraging these directives, you can create more maintainable and DRY (Don't Repeat Yourself) code within your web applications. Templates with complex and shared content structures are much easier to handle, giving both developers and designers a powerful tool for front-end work.

---
Source: [Reuse templates with Include and Extend in Jinja](https://www.slingacademy.com/article/reuse-templates-with-include-and-extend-in-jinja/)
