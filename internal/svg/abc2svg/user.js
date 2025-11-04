var abc2svg = {
    print: print,
    printErr: function(str) {
        std.err.printf("%s\n", str)
    },
    quit: function() {
        std.exit(1)
    },
    readFile: std.loadFile,
    get_mtime: function(fn) {
        return new Date(os.stat(fn)[0].mtime)
    },
    loadjs: function(fn, relay, onerror) {
        try {
            load(fn[0] == "/" ? fn : (path + fn))
            if (relay)
                relay()
        } catch (e) {
            if (onerror)
                onerror()
            else
                abc2svg.printErr("loadjs: Cannot read file " + fn +
                    "\n  " + e.name + ": " + e.message)
            return
        }
    } // loadjs()
} // abc2svg
var user = {
    read_file: function(fn) {	// read a file (main or included)
        var i,
            p = fn,
            file = abc2svg.readFile(p)

        if (!file && fn[0] != '/') {
            for (i = 0; i < abc2svg.path.length; i++) {
                p = abc2svg.path[i] + '/' + fn
                file = abc2svg.readFile(p)
                if (file)
                    break
            }
        }

        if (!file)
            return file

        // memorize the file path
        i = p.lastIndexOf('/')
        if (i > 0) {
            p = p.slice(0, i)
            if (abc2svg.path.indexOf(p) < 0)
                abc2svg.path.unshift(p)
        }

        // convert the file content into a Unix string
        i = file.indexOf('\r')
        if (i >= 0) {
            if (file[i + 1] == '\n')
                file = file.replace(/\r\n/g, '\n')	// M$
            else
                file = file.replace(/\r/g, '\n')	// Mac
        }

        // load the required modules (synchronous)
        abc2svg.modules.load(file)

        return file
    },
    errtxt: '',
    errmsg:			// print or store the error messages
        typeof abc2svg.printErr == 'function'
            ? function(msg, l, c) { abc2svg.printErr(msg) }
            : function(msg, l, c) { user.errtxt += msg + '\n' }
} // user
