// fillVersion saves protoc version number as string
func (g *Generator) fillVersion() {
	if v := g.Plugin.Request.GetCompilerVersion(); v != nil {
		g.protocVersion = fmt.Sprintf("v%v.%v.%v", v.GetMajor(), v.GetMinor(), v.GetPatch())
	}
}
