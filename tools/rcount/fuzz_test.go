package main

import "testing"

// FuzzParse checks that the readers never panic and that their counts stay
// consistent with each other.
func FuzzParse(f *testing.F) {
	f.Add(namesC)
	f.Add(gramY)
	f.Add(evalC)
	f.Add("FUNTAB R_FunTab[] =\n/* open")
	f.Add(".Call(\"x\" # '\\")
	f.Add(cranIndex)
	f.Fuzz(func(t *testing.T, text string) {
		ft := parseFunTab(text)
		if ft.PrimBuiltin+ft.PrimSpecial+ft.IntBuiltin+ft.IntSpecial != ft.Entries || ft.CFuns > ft.Entries {
			t.Fatalf("inconsistent %+v", ft)
		}
		g := parseGram(text)
		if g.RuleLines+g.Epilogue > g.Lines || g.Nonterminals > g.RuleLines {
			t.Fatalf("inconsistent %+v", g)
		}
		for _, op := range parseOpcodes(text) {
			if op == "" {
				t.Fatal("empty opcode")
			}
		}
		if c := cranStats(text); c.Compiled > c.Packages || c.PureClosure > c.Packages-c.Compiled {
			t.Fatalf("inconsistent %+v", c)
		}
		n := newNative()
		n.add(text)
		for k, s := range n.Sites {
			if len(n.Names[k]) > s {
				t.Fatalf("%s: %d names for %d sites", k, len(n.Names[k]), s)
			}
		}
	})
}
