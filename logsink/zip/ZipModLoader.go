package zip

type ZipModLoader struct {
	zipImpl ZipMod
	libpath string
}

func NewZipModLoader() *ZipModLoader {
	p := new(ZipModLoader)
	p.zipImpl = NewDefaultZipMod()
	return p

}
