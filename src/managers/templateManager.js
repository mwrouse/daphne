let path = require('path');
let debug = require('../debugger')('template_manager');
let fileUtils = require('../utils/fileUtils');
let FileToCompile = require('../models/fileToCompile');
let config = require('./projectConfig');

var __templateCache = {}; // Dictionary string => fileToCompile

/**
 * Class for managing templates
 */
class TemplateManager {

    constructor() {
    }

    /**
     * Loads all templates in a project
     */
    loadTemplates() {
        debug('Loading Templates');
        let root = config.compiler.templates_folder_absolute;

        return fileUtils.globFiles(root, '*.*')
            .then((templates) => {
                let waiting = [];

                for (let i = 0; i < templates.length; i++) {
                    let template = templates[i];
                    let name = (template.name).replace(path.extname(template.name), '');

                    if (config.compiler.ignore_absolute.indexOf(template.absolute) != -1)
                        continue; // Ignore the file

                    waiting.push(
                        fileUtils.readEntireFile(template.absolute)
                            .then((content) => {
                                debug(`Found template '${name}'`);
                                return new FileToCompile(template.absolute, content);
                            })
                    );
                }

                return Promise.all(waiting);
            });
    }


    /**
     * Returns true if a template name is valid and exists
     * @param {string} name Template name
     * @returns {boolean} whether template exists or not
     */
    isValidTemplate(name) {
        return (__templateCache.hasOwnProperty(name) || __templateCache[name] != undefined);
    }

    /**
     * Adds a template to the manager
     * @param {string} name Name of the template
     * @param {fileToCompile} file file
     */
    addTemplate(name, file) {
        __templateCache[name] = JSON.parse(JSON.stringify(file)); // Deep copy

        debug(`Template '${name}' was added`);
    }

    /**
     * Retrieves a template
     * @param {string} name template name
     */
    getTemplate(name) {
        if (__templateCache.hasOwnProperty(name) || __templateCache[name] != undefined)
            return JSON.parse(JSON.stringify(__templateCache[name])); // Deep copy

        return null; // No template found
    }

    /**
     * Updates a template
     * @param {string} name template name
     * @param {fileToCompile} file updated file
     */
    updateTemplate(name, file) {
        if (__templateCache.hasOwnProperty(name) || __templateCache[name] != undefined) {
            __templateCache[name] = JSON.parse(JSON.stringify(file)); // Deep copy
            debug(`Template '${name}' was updated`);
        }
    }
}

let instance = new TemplateManager();
module.exports = instance;