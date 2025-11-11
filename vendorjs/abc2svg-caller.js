const RenderSvgFromAbc = (code) => {
    let svg = "";
    let error = "";
    const user = {
        read_file: function(fn) {
            console.log("read_file", fn);
        }, // read_file()
        errmsg: function(msg, l, c) {	// get the errors
            error = msg;
        },
        img_out: function(p) {		// image output
            svg += p
        }
    }
    const abcInstance = new abc2svg.Abc(user);
    abcInstance.tosvg("", code, 0, code.length);

    // Mega hack to remove the undefined coming from the title, or something...
    const prefix = "undefined";
    if (svg.startsWith(prefix)) {
        svg = svg.slice(prefix.length);
    }
    if (error != "") {
        const div = document.createElement('div');
        div.textContent = error;
        div.className = "bg-red-500";
        return div.outerHTML;
    }
    return svg;
};
export { abc2svg, RenderSvgFromAbc };
