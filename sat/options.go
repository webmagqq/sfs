package sat

type Option func(*Options)

//--匿名函数可以理解为不类似类的多态，参数同样的情况下，可以进行各种不同的处理。
//--非常灵活的类型类继续的方法，或说灵活的函数继承操作方法。
type Options struct {
	Path string `json:"path"`
}

func SetPath(path string) Option {
	return func(args *Options) {
		args.Path = path
	}
}
