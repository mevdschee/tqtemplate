package tqtemplate

import (
	"fmt"
	"strings"
	"testing"
)

// Test basic block definition and rendering
func TestBlockBasic(t *testing.T) {
	tmpl := `<html>
{% block title %}Default Title{% endblock %}
{% block content %}Default Content{% endblock %}
</html>`

	expected := `<html>
Default Title
Default Content
</html>`

	template := NewTemplate()
	result, err := template.Render(tmpl, map[string]any{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Test extends with block overrides
func TestExtendsWithBlockOverride(t *testing.T) {
	// Create a simple in-memory template loader
	templates := map[string]string{
		"base.html": `<html>
<head>
  <title>{% block title %}My Website{% endblock %}</title>
</head>
<body>
  {% block content %}{% endblock %}
</body>
</html>`,
	}

	loader := func(name string) (string, error) {
		if tmpl, exists := templates[name]; exists {
			return tmpl, nil
		}
		return "", fmt.Errorf("template not found: %s", name)
	}

	childTmpl := `{% extends 'base.html' %}

{% block title %}Home Page{% endblock %}

{% block content %}
<h1>Welcome to the home page!</h1>
{% endblock %}`

	expected := `<html>
<head>
  <title>Home Page</title>
</head>
<body>
<h1>Welcome to the home page!</h1>
</body>
</html>`

	template := NewTemplateWithLoader(loader)
	result, err := template.Render(childTmpl, map[string]any{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Test extends with partial block override (some blocks keep default)
func TestExtendsWithPartialOverride(t *testing.T) {
	templates := map[string]string{
		"base.html": `<html>
<head>
  <title>{% block title %}Default Title{% endblock %}</title>
</head>
<body>
  <header>{% block header %}Default Header{% endblock %}</header>
  <main>{% block content %}Default Content{% endblock %}</main>
  <footer>{% block footer %}Default Footer{% endblock %}</footer>
</body>
</html>`,
	}

	loader := func(name string) (string, error) {
		if tmpl, exists := templates[name]; exists {
			return tmpl, nil
		}
		return "", fmt.Errorf("template not found: %s", name)
	}

	childTmpl := `{% extends 'base.html' %}

{% block title %}Custom Title{% endblock %}

{% block content %}<p>Custom content here</p>
{% endblock %}`

	expected := `<html>
<head>
  <title>Custom Title</title>
</head>
<body>
  <header>Default Header</header>
  <main><p>Custom content here</p>
</main>
  <footer>Default Footer</footer>
</body>
</html>`

	template := NewTemplateWithLoader(loader)
	result, err := template.Render(childTmpl, map[string]any{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Test extends with variables in blocks
func TestExtendsWithVariables(t *testing.T) {
	templates := map[string]string{
		"base.html": `<html>
<head>
  <title>{% block title %}{{ site_name }}{% endblock %}</title>
</head>
<body>
  {% block content %}{% endblock %}
</body>
</html>`,
	}

	loader := func(name string) (string, error) {
		if tmpl, exists := templates[name]; exists {
			return tmpl, nil
		}
		return "", fmt.Errorf("template not found: %s", name)
	}

	childTmpl := `{% extends 'base.html' %}

{% block title %}{{ page_title }} - {{ site_name }}{% endblock %}

{% block content %}
<h1>{{ heading }}</h1>
<p>{{ message }}</p>
{% endblock %}`

	data := map[string]any{
		"site_name":  "My Site",
		"page_title": "About",
		"heading":    "About Us",
		"message":    "Welcome to our site!",
	}

	expected := `<html>
<head>
  <title>About - My Site</title>
</head>
<body>
<h1>About Us</h1>
<p>Welcome to our site!</p>
</body>
</html>`

	template := NewTemplateWithLoader(loader)
	result, err := template.Render(childTmpl, data)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Test extends with control structures in blocks
func TestExtendsWithControlStructures(t *testing.T) {
	templates := map[string]string{
		"base.html": `<html>
<body>
  <ul>
  {% block navigation %}
    <li><a href="/">Home</a></li>
  {% endblock %}
  </ul>
  {% block content %}{% endblock %}
</body>
</html>`,
	}

	loader := func(name string) (string, error) {
		if tmpl, exists := templates[name]; exists {
			return tmpl, nil
		}
		return "", fmt.Errorf("template not found: %s", name)
	}

	childTmpl := `{% extends 'base.html' %}

{% block navigation %}
{% for item in menu %}
    <li><a href="{{ item.url }}">{{ item.title }}</a></li>
{% endfor %}
{% endblock %}

{% block content %}
<h1>{{ title }}</h1>
{% if show_list %}
<ul>
{% for item in items %}
  <li>{{ item }}</li>
{% endfor %}
</ul>
{% endif %}
{% endblock %}`

	data := map[string]any{
		"menu": []any{
			map[string]any{"url": "/", "title": "Home"},
			map[string]any{"url": "/about", "title": "About"},
			map[string]any{"url": "/contact", "title": "Contact"},
		},
		"title":     "My Page",
		"show_list": true,
		"items":     []any{"Item 1", "Item 2", "Item 3"},
	}

	expected := `<html>
<body>
  <ul>
    <li><a href="/">Home</a></li>
    <li><a href="/about">About</a></li>
    <li><a href="/contact">Contact</a></li>
  </ul>
<h1>My Page</h1>
<ul>
  <li>Item 1</li>
  <li>Item 2</li>
  <li>Item 3</li>
</ul>
</body>
</html>`

	template := NewTemplateWithLoader(loader)
	result, err := template.Render(childTmpl, data)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Test error when loader not configured
func TestExtendsWithoutLoader(t *testing.T) {
	childTmpl := `{% extends 'base.html' %}
{% block content %}Test{% endblock %}`

	template := NewTemplate()
	_, err := template.Render(childTmpl, map[string]any{})
	if err == nil {
		t.Error("Expected error when loader not configured")
	}
	if !strings.Contains(err.Error(), "loader not configured") {
		t.Errorf("Expected loader error, got: %v", err)
	}
}

// Test error when parent template not found
func TestExtendsTemplateNotFound(t *testing.T) {
	loader := func(name string) (string, error) {
		return "", fmt.Errorf("template not found: %s", name)
	}

	childTmpl := `{% extends 'nonexistent.html' %}
{% block content %}Test{% endblock %}`

	template := NewTemplateWithLoader(loader)
	_, err := template.Render(childTmpl, map[string]any{})
	if err == nil {
		t.Error("Expected error when parent template not found")
	}
	if !strings.Contains(err.Error(), "failed to load parent template") {
		t.Errorf("Expected template not found error, got: %v", err)
	}
}

// Test nested blocks (blocks within blocks) - inherits nested structure
func TestNestedBlocks(t *testing.T) {
	templates := map[string]string{
		"base.html": `<div>
{% block outer %}
  <div class="outer">
  {% block inner %}Inner default{% endblock %}
  </div>
{% endblock %}
</div>`,
	}

	loader := func(name string) (string, error) {
		if tmpl, exists := templates[name]; exists {
			return tmpl, nil
		}
		return "", fmt.Errorf("template not found: %s", name)
	}

	// Override inner block - parent's outer block includes the inner reference
	// so inner will be replaced with child's content
	childTmpl := `{% extends 'base.html' %}

{% block inner %}Custom inner content{% endblock %}`

	expected := `<div>
  <div class="outer">
Custom inner content
  </div>
</div>`

	template := NewTemplateWithLoader(loader)
	result, err := template.Render(childTmpl, map[string]any{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Test empty blocks
func TestEmptyBlocks(t *testing.T) {
	templates := map[string]string{
		"base.html": `<html>
<head>{% block head %}{% endblock %}</head>
<body>{% block body %}Default body{% endblock %}</body>
</html>`,
	}

	loader := func(name string) (string, error) {
		if tmpl, exists := templates[name]; exists {
			return tmpl, nil
		}
		return "", fmt.Errorf("template not found: %s", name)
	}

	childTmpl := `{% extends 'base.html' %}

{% block head %}<title>Page</title>{% endblock %}

{% block body %}{% endblock %}`

	expected := `<html>
<head><title>Page</title></head>
<body></body>
</html>`

	template := NewTemplateWithLoader(loader)
	result, err := template.Render(childTmpl, map[string]any{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Test documenting the current behavior: child blocks do NOT inherit indentation
func TestBlockInheritanceNoIndentationPreservation(t *testing.T) {
	templates := map[string]string{
		"base.html": "<html>\n  <body>\n    <div>\n      {% block content %}Default{% endblock %}\n    </div>\n  </body>\n</html>",
	}

	loader := func(name string) (string, error) {
		if tmpl, exists := templates[name]; exists {
			return tmpl, nil
		}
		return "", fmt.Errorf("template not found: %s", name)
	}

	childTmpl := "{% extends 'base.html' %}\n\n{% block content %}<h1>Title</h1>\n<p>Text</p>{% endblock %}"

	// Expected: child content is NOT indented (replaces block completely)
	expected := "<html>\n  <body>\n    <div>\n<h1>Title</h1>\n<p>Text</p>\n    </div>\n  </body>\n</html>"

	template := NewTemplateWithLoader(loader)
	result, err := template.Render(childTmpl, map[string]any{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// Test basic include functionality
func TestIncludeBasic(t *testing.T) {
	templates := map[string]string{
		"header.html": "<header><h1>Site Header</h1></header>",
		"main.html":   "<div>{% include 'header.html' %}\n<main>Main content</main>\n</div>",
	}

	loader := func(name string) (string, error) {
		if tmpl, exists := templates[name]; exists {
			return tmpl, nil
		}
		return "", fmt.Errorf("template not found: %s", name)
	}

	template := NewTemplateWithLoader(loader)
	result, err := template.Render(templates["main.html"], map[string]any{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "<div><header><h1>Site Header</h1></header>\n<main>Main content</main>\n</div>"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// Test include with variables
func TestIncludeWithVariables(t *testing.T) {
	templates := map[string]string{
		"greeting.html": "<p>Hello, {{ name }}!</p>",
		"main.html":     "<div>{% include 'greeting.html' %}</div>",
	}

	loader := func(name string) (string, error) {
		if tmpl, exists := templates[name]; exists {
			return tmpl, nil
		}
		return "", fmt.Errorf("template not found: %s", name)
	}

	data := map[string]any{
		"name": "Alice",
	}

	template := NewTemplateWithLoader(loader)
	result, err := template.Render(templates["main.html"], data)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "<div><p>Hello, Alice!</p></div>"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// Test multiple includes
func TestMultipleIncludes(t *testing.T) {
	templates := map[string]string{
		"header.html": "<header>Header</header>\n",
		"footer.html": "<footer>Footer</footer>\n",
		"main.html":   "{% include 'header.html' %}\n<main>Content</main>\n{% include 'footer.html' %}\n",
	}

	loader := func(name string) (string, error) {
		if tmpl, exists := templates[name]; exists {
			return tmpl, nil
		}
		return "", fmt.Errorf("template not found: %s", name)
	}

	template := NewTemplateWithLoader(loader)
	result, err := template.Render(templates["main.html"], map[string]any{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "<header>Header</header>\n<main>Content</main>\n<footer>Footer</footer>\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// Test include with control structures
func TestIncludeWithControlStructures(t *testing.T) {
	templates := map[string]string{
		"item.html": "{% for item in items %}<li>{{ item }}</li>\n{% endfor %}",
		"main.html": "<ul>\n{% include 'item.html' %}</ul>",
	}

	loader := func(name string) (string, error) {
		if tmpl, exists := templates[name]; exists {
			return tmpl, nil
		}
		return "", fmt.Errorf("template not found: %s", name)
	}

	data := map[string]any{
		"items": []any{"Apple", "Banana", "Cherry"},
	}

	template := NewTemplateWithLoader(loader)
	result, err := template.Render(templates["main.html"], data)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "<ul>\n<li>Apple</li>\n<li>Banana</li>\n<li>Cherry</li>\n</ul>"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// Test include without loader
func TestIncludeWithoutLoader(t *testing.T) {
	template := NewTemplate()
	_, err := template.Render("{% include 'header.html' %}", map[string]any{})
	if err == nil {
		t.Error("Expected error when loader not configured")
	}
	if !strings.Contains(err.Error(), "template loader not configured") {
		t.Errorf("Expected 'template loader not configured' error, got: %v", err)
	}
}

// Test include template not found
func TestIncludeTemplateNotFound(t *testing.T) {
	loader := func(name string) (string, error) {
		return "", fmt.Errorf("template not found: %s", name)
	}

	template := NewTemplateWithLoader(loader)
	_, err := template.Render("{% include 'missing.html' %}", map[string]any{})
	if err == nil {
		t.Error("Expected error when template not found")
	}
	if !strings.Contains(err.Error(), "failed to load included template") {
		t.Errorf("Expected 'failed to load included template' error, got: %v", err)
	}
}

// Test nested includes
func TestNestedIncludes(t *testing.T) {
	templates := map[string]string{
		"deep.html":   "<span>Deep content</span>",
		"middle.html": "<div>{% include 'deep.html' %}</div>",
		"top.html":    "<section>{% include 'middle.html' %}</section>",
	}

	loader := func(name string) (string, error) {
		if tmpl, exists := templates[name]; exists {
			return tmpl, nil
		}
		return "", fmt.Errorf("template not found: %s", name)
	}

	template := NewTemplateWithLoader(loader)
	result, err := template.Render(templates["top.html"], map[string]any{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "<section><div><span>Deep content</span></div></section>"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
