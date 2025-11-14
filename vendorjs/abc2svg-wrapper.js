import { abc2svg } from './abc2svg-1.cjs';

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
            if (p) {
                svg += p
            }
        }
    }
    const abcInstance = new abc2svg.Abc(user);
    abcInstance.tosvg("", code, 0, code.length);

    if (error != "") {
        const div = document.createElement('div');
        div.textContent = error;
        div.className = "bg-red-500";
        return div.outerHTML;
    }
    return svg;
};
export { RenderSvgFromAbc };
