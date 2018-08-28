let path = require('path');
let debugFactory = require('../debugger');
let fileUtils = require('../utils/fileUtils');
let projectConfig = require('../managers/projectConfig');


function __getPreview() {
    let parts = this.__content.split(projectConfig.site.read_more);

    this.__preview = parts[0];
}

/**
 * Gets page metadata
 */
function __getPageInformation() {
    let lines = this.__file.content.split('\n');
    let header = [];

    if (lines[0].trim() != projectConfig.compiler.tags.delimeter) {
        this.__debug(`Warning: ${this.__file.relative} doesn't start with header--Won't parse metadata`);
        this.__content = this.__file.content;
        return;
    }

    lines.shift(); // Remove first line (that was already read)

    while (lines.length > 0) {
        let line = lines.shift().trim();
        if (line == projectConfig.compiler.tags.delimeter)
            break; // Done reading metadata

        let parts = line.split(':');
        if (parts.length < 2) {
            throw new Error(`Unable to parse metadata from '${this.__file.relative}', invalid line: ${line}`);
        }

        let key = parts.shift().trim();
        let value = parts.join(':').trim();

        this.__metadata[key] = value;

        header.push(line);
    }

    this.__header = header.join('\n');
    this.__content = lines.join('\n');
}



class FileToCompile {

    constructor(filePath, content) {
        filePath = path.normalize(filePath);
        content = (content == undefined || content == null) ? '' : content;

        let root = path.normalize(projectConfig.site.source_absolute);


        // Information about the file and raw content
        this.__file = {
            absolute: filePath,
            relative: path.relative(root, filePath),
            absoluteDirname: path.dirname(filePath),
            relativeDirname: null,
            name: path.basename(filePath),
            content: content.trim()
        };
        this.__file.relativeDirname = path.dirname(this.__file.relative);

        this.__header = ''; // File header
        this.__metadata = {}; // Parsed file header

        this.__content = this.__file.content; // Content with header removed
        this.__fullPage = this.__file.content; // Content after template and all that stuff
        this.__assets = []; // For posts
        this.__preview = '';

        this.__debug = debugFactory('file:' + this.__file.name);

        // Read header and get preview if possible
        if (this.shouldBeParsed) {
            __getPageInformation.call(this);

            this.__preview = this.__metadata['excerpt'];

            if (this.__preview == undefined)
                __getPreview.call(this);
        }
    }

    /**
     * Gets object for use when rendering the file
     */
    get context() {
        let data = {};

        for (let key in this.__metadata) {
            if (!this.__metadata.hasOwnProperty(key))
                continue;

            data[key] = this.__metadata[key];
        }

        // TODO: Add url and stuff
        data.file_name = this.__file.name;
        data.content = null;
        data.preview = null;

        if (this.shouldBeParsed) {
            data.content = this.content;
            data.preview = this.preview;
        }
        data.isBinary = !this.shouldBeParsed;

        data.assets = [];

        for (let i = 0; i < this.__assets.length; i++) {
            data.assets.push(this.__assets[i].context);
        }

        return data;
    }

    /**
     * Returns the content
     */
    get content() {
        return this.__content;
    }

    /**
     * Returns the expanded content
     */
    get page() {
        return this.__fullPage;
    }

    /**
     * Post preview
     */
    get preview() {
        return this.__preview;
    }

    // File name
    get name() {
        return this.__file.name;
    }

    // Absolute path
    get path() {
        return this.__file.absolute;
    }


    /**
     * If the file should be parsed
     */
    get shouldBeParsed() {
        return fileUtils.canFileBeParsed(this.__file.absolute);
    }
}



module.exports = FileToCompile;
