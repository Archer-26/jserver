package config_go

const (
	ExtLoad_PackageFn = 1
)

var loadfn = make(map[string]func(string) error)

var loadPackageFn = make(map[int32]func())

func RegisterPackageFn(ExtLoad int32, f func()) {
	loadPackageFn[ExtLoad] = f
}

func Load(dir string) error {
	for _, fn := range loadfn {
		if err := fn(dir); err != nil {
			return err
		}
	}
	for _, fn := range loadPackageFn {
		fn()
	}
	return nil
}

type array_int []int64
type array_str []string
type array_float []float32
type int2int map[int64]int64
type int2str map[int64]string
type str2int map[string]int64
type str2str map[string]string

func (t array_int) Len() int   { return len(t) }
func (t array_str) Len() int   { return len(t) }
func (t array_float) Len() int { return len(t) }
func (t int2int) Len() int     { return len(t) }
func (t int2str) Len() int     { return len(t) }
func (t str2int) Len() int     { return len(t) }
func (t str2str) Len() int     { return len(t) }

func (t array_int) Get(index int) int64                 { return t[index] }
func (t array_str) Get(index int) string                { return t[index] }
func (t array_float) Get(index int) float32             { return t[index] }
func (t int2int) Get(key int64) (data int64, ok bool)   { data, ok = t[key]; return }
func (t int2str) Get(key int64) (data string, ok bool)  { data, ok = t[key]; return }
func (t str2int) Get(key string) (data int64, ok bool)  { data, ok = t[key]; return }
func (t str2str) Get(key string) (data string, ok bool) { data, ok = t[key]; return }

func (t array_int) Range(fn func(int, int64) (stop bool)) {
	for k, v := range t {
		if fn(k, v) {
			break
		}
	}
}

func (t array_str) Range(fn func(int, string) (stop bool)) {
	for k, v := range t {
		if fn(k, v) {
			break
		}
	}
}

func (t array_float) Range(fn func(int, float32) (stop bool)) {
	for k, v := range t {
		if fn(k, v) {
			break
		}
	}
}

func (t int2int) Range(fn func(int64, int64) (stop bool)) {
	for k, v := range t {
		if fn(k, v) {
			break
		}
	}
}

func (t int2str) Range(fn func(int64, string) (stop bool)) {
	for k, v := range t {
		if fn(k, v) {
			break
		}
	}
}

func (t str2int) Range(fn func(string, int64) (stop bool)) {
	for k, v := range t {
		if fn(k, v) {
			break
		}
	}
}

func (t str2str) Range(fn func(string, string) (stop bool)) {
	for k, v := range t {
		if fn(k, v) {
			break
		}
	}
}

func (t array_int) Copy() array_int {
	l := len(t)
	if l == 0 {
		return nil
	}
	cp := make(array_int, l)
	for i := range t {
		cp[i] = t[i]
	}
	return cp
}

func (t array_str) Copy() array_str {
	l := len(t)
	if l == 0 {
		return nil
	}
	cp := make(array_str, l)
	for i := range t {
		cp[i] = t[i]
	}
	return cp
}

func (t array_float) Copy() array_float {
	l := len(t)
	if l == 0 {
		return nil
	}
	cp := make(array_float, l)
	for i := range t {
		cp[i] = t[i]
	}
	return cp
}

func (t int2int) Copy() int2int {
	l := len(t)
	if l == 0 {
		return nil
	}
	cp := make(int2int, l)
	for k, v := range t {
		cp[k] = v
	}
	return cp
}

func (t int2str) Copy() int2str {
	l := len(t)
	if l == 0 {
		return nil
	}
	cp := make(int2str, l)
	for k, v := range t {
		cp[k] = v
	}
	return cp
}

func (t str2int) Copy() str2int {
	l := len(t)
	if l == 0 {
		return nil
	}
	cp := make(str2int, l)
	for k, v := range t {
		cp[k] = v
	}
	return cp
}

func (t str2str) Copy() str2str {
	l := len(t)
	if l == 0 {
		return nil
	}
	cp := make(str2str, l)
	for k, v := range t {
		cp[k] = v
	}
	return cp
}
