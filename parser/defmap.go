package parser

import (
	"github.com/julelang/jule/pkg/juleio"
	"github.com/julelang/jule/pkg/juletype"
)

func is_accessable(finder, target *File, defIsPub bool) bool {
	return defIsPub || finder == nil || target == nil || finder.Dir == target.Dir
}

// Defmap is definition map.
type Defmap struct {
	Namespaces []*namespace
	Enums      []*Enum
	Structs    []*structure
	Traits     []*trait
	Types      []*TypeAlias
	Funcs      []*Fn
	Globals    []*Var
	side       *Defmap
}

func (dm *Defmap) find_ns_by_id(id string) int {
	for i, t := range dm.Namespaces {
		if t != nil && t.Id == id {
			return i
		}
	}
	return -1
}

func (dm *Defmap) ns_by_id(id string) *namespace {
	i := dm.find_ns_by_id(id)
	if i == -1 {
		return nil
	}
	return dm.Namespaces[i]
}

func (dm *Defmap) find_struct_by_id(id string, f *File) (int, *Defmap, bool) {
	for i, s := range dm.Structs {
		if s != nil && s.Ast.Id == id {
			if is_accessable(f, s.Ast.Token.File, s.Ast.Pub) {
				return i, dm, false
			}
		}
	}
	if dm.side != nil {
		i, m, _ := dm.side.find_struct_by_id(id, f)
		return i, m, true
	}
	return -1, nil, false
}

func (dm *Defmap) struct_by_Id(id string, f *File) (*structure, *Defmap, bool) {
	i, m, canshadow := dm.find_struct_by_id(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Structs[i], m, canshadow
}

func (dm *Defmap) find_trait_vy_id(id string, f *File) (int, *Defmap, bool) {
	for i, t := range dm.Traits {
		if t != nil && t.Ast.Id == id {
			if is_accessable(f, t.Ast.Token.File, t.Ast.Pub) {
				return i, dm, false
			}
		}
	}
	if dm.side != nil {
		i, m, _ := dm.side.find_trait_vy_id(id, f)
		return i, m, true
	}
	return -1, nil, false
}

func (dm *Defmap) trait_by_id(id string, f *File) (*trait, *Defmap, bool) {
	i, m, canshadow := dm.find_trait_vy_id(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Traits[i], m, canshadow
}

func (dm *Defmap) find_enum_by_id(id string, f *File) (int, *Defmap, bool) {
	for i, e := range dm.Enums {
		if e != nil && e.Id == id {
			if is_accessable(f, e.Token.File, e.Pub) {
				return i, dm, false
			}
		}
	}
	if dm.side != nil {
		i, m, _ := dm.side.find_enum_by_id(id, f)
		return i, m, true
	}
	return -1, nil, false
}

func (dm *Defmap) enum_by_id(id string, f *File) (*Enum, *Defmap, bool) {
	i, m, canshadow := dm.find_enum_by_id(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Enums[i], m, canshadow
}

func (dm *Defmap) find_type_by_id(id string, f *File) (int, *Defmap, bool) {
	for i, t := range dm.Types {
		if t != nil && t.Id == id {
			if is_accessable(f, t.Token.File, t.Pub) {
				return i, dm, false
			}
		}
	}
	if dm.side != nil {
		i, m, _ := dm.side.find_type_by_id(id, f)
		return i, m, true
	}
	return -1, nil, false
}

func (dm *Defmap) type_by_id(id string, f *File) (*TypeAlias, *Defmap, bool) {
	i, m, canshadow := dm.find_type_by_id(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Types[i], m, canshadow
}

func (dm *Defmap) find_fn_by_id(id string, f *File) (int, *Defmap, bool) {
	for i, fn := range dm.Funcs {
		if fn != nil && fn.Ast.Id == id {
			if is_accessable(f, fn.Ast.Token.File, fn.Ast.Pub) {
				return i, dm, false
			}
		}
	}
	if dm.side != nil {
		i, m, _ := dm.side.find_fn_by_id(id, f)
		return i, m, true
	}
	return -1, nil, false
}

// fn_by_id returns function by specified id.
//
// Special case:
//  fn_by_id(id) -> nil: if function is not exist.
func (dm *Defmap) fn_by_id(id string, f *File) (*Fn, *Defmap, bool) {
	i, m, canshadow := dm.find_fn_by_id(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Funcs[i], m, canshadow
}

func (dm *Defmap) find_global_by_id(id string, f *File) (int, *Defmap, bool) {
	for i, g := range dm.Globals {
		if g != nil && g.Type.Id != juletype.VOID && g.Id == id {
			if is_accessable(f, g.Token.File, g.Pub) {
				return i, dm, false
			}
		}
	}
	if dm.side != nil {
		i, m, _ := dm.side.find_global_by_id(id, f)
		return i, m, true
	}
	return -1, nil, false
}

func (dm *Defmap) global_by_id(id string, f *File) (*Var, *Defmap, bool) {
	i, m, canshadow := dm.find_global_by_id(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Globals[i], m, canshadow
}

// find_by_id returns index of definition with type if exist.
//
// Special case is;
//  find_by_id(id) -> -1, ' ' if id is not exist
//
// Types;
// 'g' -> global
// 'f' -> function
// 'e' -> enum
// 's' -> struct
// 't' -> type alias
// 'i' -> trait
func (dm *Defmap) find_by_id(id string, f *File) (int, *Defmap, byte) {
	var finders = map[byte]func(string, *juleio.File) (int, *Defmap, bool){
		'g': dm.find_global_by_id,
		'f': dm.find_fn_by_id,
		'e': dm.find_enum_by_id,
		's': dm.find_struct_by_id,
		't': dm.find_type_by_id,
		'i': dm.find_trait_vy_id,
	}
	for code, finder := range finders {
		i, m, _ := finder(id, f)
		if i != -1 {
			return i, m, code
		}
	}
	return -1, nil, ' '
}

func push_defines(dest, src *Defmap) {
	dest.Types = append(dest.Types, src.Types...)
	dest.Traits = append(dest.Traits, src.Traits...)
	dest.Structs = append(dest.Structs, src.Structs...)
	dest.Enums = append(dest.Enums, src.Enums...)
	dest.Globals = append(dest.Globals, src.Globals...)
	dest.Funcs = append(dest.Funcs, src.Funcs...)
}
