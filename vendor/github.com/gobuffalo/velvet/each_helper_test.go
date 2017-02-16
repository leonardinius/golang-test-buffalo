package velvet_test

import (
	"fmt"
	"testing"

	"github.com/gobuffalo/velvet"
	"github.com/stretchr/testify/require"
)

func Test_Each_Helper_NoArgs(t *testing.T) {
	r := require.New(t)
	ctx := velvet.NewContext()
	input := `{{#each }}{{@value}}{{/each}}`

	_, err := velvet.Render(input, ctx)
	r.Error(err)
}

func Test_Each_Helper(t *testing.T) {
	r := require.New(t)
	ctx := velvet.NewContext()
	ctx.Set("names", []string{"mark", "bates"})
	input := `{{#each names }}<p>{{@value}}</p>{{/each}}`

	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	r.Equal("<p>mark</p><p>bates</p>", s)
}

func Test_Each_Helper_Index(t *testing.T) {
	r := require.New(t)
	ctx := velvet.NewContext()
	ctx.Set("names", []string{"mark", "bates"})
	input := `{{#each names }}<p>{{@index}}</p>{{/each}}`

	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	r.Equal("<p>0</p><p>1</p>", s)
}

func Test_Each_Helper_As(t *testing.T) {
	r := require.New(t)
	ctx := velvet.NewContext()
	ctx.Set("names", []string{"mark", "bates"})
	input := `{{#each names as |ind name| }}<p>{{ind}}-{{name}}</p>{{/each}}`

	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	r.Equal("<p>0-mark</p><p>1-bates</p>", s)
}

func Test_Each_Helper_As_Nested(t *testing.T) {
	r := require.New(t)
	ctx := velvet.NewContext()
	users := []struct {
		Name     string
		Initials []string
	}{
		{Name: "Mark", Initials: []string{"M", "F", "B"}},
		{Name: "Rachel", Initials: []string{"R", "A", "B"}},
	}
	ctx.Set("users", users)
	input := `
{{#each users as |user|}}
	<h1>{{user.Name}}</h1>
	{{#each user.Initials as |i|}}
		{{user.Name}}: {{i}}
	{{/each}}
{{/each}}
	`

	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	r.Contains(s, "<h1>Mark</h1>")
	r.Contains(s, "Mark: M")
	r.Contains(s, "Mark: F")
	r.Contains(s, "Mark: B")
	r.Contains(s, "<h1>Rachel</h1>")
	r.Contains(s, "Rachel: R")
	r.Contains(s, "Rachel: A")
	r.Contains(s, "Rachel: B")
}

func Test_Each_Helper_SlicePtr(t *testing.T) {
	r := require.New(t)
	type user struct {
		Name string
	}
	type users []user

	us := &users{
		{Name: "Mark"},
		{Name: "Rachel"},
	}

	ctx := velvet.NewContext()
	ctx.Set("users", us)

	input := `
	{{#each users as |user|}}
		{{user.Name}}
	{{/each}}
	`
	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	r.Contains(s, "Mark")
	r.Contains(s, "Rachel")
}

func Test_Each_Helper_Map(t *testing.T) {
	r := require.New(t)
	ctx := velvet.NewContext()
	data := map[string]string{
		"a": "A",
		"b": "B",
	}
	ctx.Set("letters", data)
	input := `
	{{#each letters}}
		{{@key}}:{{@value}}
	{{/each}}
	`

	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	for k, v := range data {
		r.Contains(s, fmt.Sprintf("%s:%s", k, v))
	}
}

func Test_Each_Helper_Map_As(t *testing.T) {
	r := require.New(t)
	ctx := velvet.NewContext()
	data := map[string]string{
		"a": "A",
		"b": "B",
	}
	ctx.Set("letters", data)
	input := `
	{{#each letters as |k v|}}
		{{k}}:{{v}}
	{{/each}}
	`

	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	for k, v := range data {
		r.Contains(s, fmt.Sprintf("%s:%s", k, v))
	}
}

func Test_Each_Helper_Else(t *testing.T) {
	r := require.New(t)
	ctx := velvet.NewContext()
	data := map[string]string{}
	ctx.Set("letters", data)
	input := `
	{{#each letters as |k v|}}
		{{k}}:{{v}}
	{{else}}
		no letters
	{{/each}}
	`

	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	r.Contains(s, "no letters")
}

func Test_Each_Helper_Else_Collection(t *testing.T) {
	r := require.New(t)
	ctx := velvet.NewContext()
	data := map[string][]string{}
	ctx.Set("collection", data)

	input := `
	{{#each collection.emptykey as |k v|}}
		{{k}}:{{v}}
	{{else}}
		no letters
	{{/each}}
	`

	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	r.Contains(s, "no letters")
}

func Test_Each_Helper_Else_CollectionMap(t *testing.T) {
	r := require.New(t)
	ctx := velvet.NewContext()
	data := map[string]map[string]string{
		"emptykey": map[string]string{},
	}

	ctx.Set("collection", data)

	input := `
	{{#each collection.emptykey.something as |k v|}}
		{{k}}:{{v}}
	{{else}}
		no letters
	{{/each}}
	`

	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	r.Contains(s, "no letters")
}
