package template

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	cases := []struct {
		File   string
		Result *Template
		Err    bool
	}{
		/*
		 * Builders
		 */
		{
			"parse-basic.json",
			&Template{
				Builders: map[string]*Builder{
					"something": &Builder{
						Name: "something",
						Type: "something",
					},
				},
			},
			false,
		},
		{
			"parse-builder-no-type.json",
			nil,
			true,
		},
		{
			"parse-builder-repeat.json",
			nil,
			true,
		},

		/*
		 * Provisioners
		 */
		{
			"parse-provisioner-basic.json",
			&Template{
				Provisioners: []*Provisioner{
					&Provisioner{
						Type: "something",
					},
				},
			},
			false,
		},

		{
			"parse-provisioner-pause-before.json",
			&Template{
				Provisioners: []*Provisioner{
					&Provisioner{
						Type:        "something",
						PauseBefore: 1 * time.Second,
					},
				},
			},
			false,
		},

		{
			"parse-provisioner-only.json",
			&Template{
				Provisioners: []*Provisioner{
					&Provisioner{
						Type: "something",
						OnlyExcept: OnlyExcept{
							Only: []string{"foo"},
						},
					},
				},
			},
			false,
		},

		{
			"parse-provisioner-except.json",
			&Template{
				Provisioners: []*Provisioner{
					&Provisioner{
						Type: "something",
						OnlyExcept: OnlyExcept{
							Except: []string{"foo"},
						},
					},
				},
			},
			false,
		},

		{
			"parse-provisioner-override.json",
			&Template{
				Provisioners: []*Provisioner{
					&Provisioner{
						Type: "something",
						Override: map[string]interface{}{
							"foo": map[string]interface{}{},
						},
					},
				},
			},
			false,
		},

		{
			"parse-provisioner-no-type.json",
			nil,
			true,
		},

		{
			"parse-variable-default.json",
			&Template{
				Variables: map[string]*Variable{
					"foo": &Variable{
						Default: "foo",
					},
				},
			},
			false,
		},

		{
			"parse-variable-required.json",
			&Template{
				Variables: map[string]*Variable{
					"foo": &Variable{
						Required: true,
					},
				},
			},
			false,
		},

		{
			"parse-pp-basic.json",
			&Template{
				PostProcessors: [][]*PostProcessor{
					[]*PostProcessor{
						&PostProcessor{
							Type: "foo",
							Config: map[string]interface{}{
								"foo": "bar",
							},
						},
					},
				},
			},
			false,
		},

		{
			"parse-pp-keep.json",
			&Template{
				PostProcessors: [][]*PostProcessor{
					[]*PostProcessor{
						&PostProcessor{
							Type:              "foo",
							KeepInputArtifact: true,
						},
					},
				},
			},
			false,
		},

		{
			"parse-pp-string.json",
			&Template{
				PostProcessors: [][]*PostProcessor{
					[]*PostProcessor{
						&PostProcessor{
							Type: "foo",
						},
					},
				},
			},
			false,
		},

		{
			"parse-pp-map.json",
			&Template{
				PostProcessors: [][]*PostProcessor{
					[]*PostProcessor{
						&PostProcessor{
							Type: "foo",
						},
					},
				},
			},
			false,
		},

		{
			"parse-pp-slice.json",
			&Template{
				PostProcessors: [][]*PostProcessor{
					[]*PostProcessor{
						&PostProcessor{
							Type: "foo",
						},
					},
					[]*PostProcessor{
						&PostProcessor{
							Type: "bar",
						},
					},
				},
			},
			false,
		},

		{
			"parse-pp-multi.json",
			&Template{
				PostProcessors: [][]*PostProcessor{
					[]*PostProcessor{
						&PostProcessor{
							Type: "foo",
						},
						&PostProcessor{
							Type: "bar",
						},
					},
				},
			},
			false,
		},

		{
			"parse-pp-no-type.json",
			nil,
			true,
		},
	}

	for _, tc := range cases {
		f, err := os.Open(fixtureDir(tc.File))
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		tpl, err := Parse(f)
		f.Close()
		if (err != nil) != tc.Err {
			t.Fatalf("err: %s", err)
		}

		if !reflect.DeepEqual(tpl, tc.Result) {
			t.Fatalf("bad: %s\n\n%#v\n\n%#v", tc.File, tpl, tc.Result)
		}
	}
}