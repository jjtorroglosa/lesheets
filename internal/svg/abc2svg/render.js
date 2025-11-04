function renderAbcToSvg(filename, abccode) {
    fn,
        out = [],			// output without SVG container
        yo = 0,				// offset of the next SVG
        w = 0;				// max width

    abc2svg.abc_init([])
    var abc = new abc2svg.Abc(user);
    abc.tosvg(filename, abccode);

    var result = abc2svg.abc_end()
    return result;
}
({
    renderAbcToSvg: renderAbcToSvg
})
