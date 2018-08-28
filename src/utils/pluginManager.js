let debug = require('../debugger')('plugin_manager');
let fileUtils = require('./fileUtils');
let path = require('path');



class PluginManager {
    constructor() {
        this.__plugins = [];
    }


    /**
     * Registers a plugin with the Plugin Manager
     * @param {string} name Folder the plugin is in
     * @param {string} pluginJSONFile Path to the file
     */
    register(name, pluginJSONFile) {
        debug(`Registering Plugin '${name}'`);

        return fileUtils.readEntireFile(pluginJSONFile)
            .then((manifest) => {
                try {
                    manifest = JSON.parse(manifest);
                    if (this.__verifyManifest(manifest)) {
                        debug(`Manifest for ${name} is valid`);

                        this.__plugins.push({
                            root: path.dirname(pluginJSONFile),
                            manifest: manifest
                        });
                    }
                    else {
                        debug(`Manifest for ${name} is invalid`);
                    }
                } catch (e) {
                    debug(`Error parsing manifest for ${name}`);
                    console.error(e);
                }
            });
    }


    /**
     * Verifies that a manifest is valid
     * @param {object} manifest
     */
    __verifyManifest(manifest) {
        let required_properties = [
            { key: 'name', type: 'string' },
            { key: 'functionality', type: 'object' }
        ];

        for (let i = 0; i < required_properties.length; i++) {
            let prop = required_properties[i];

            if (
                (!manifest.hasOwnProperty(prop.key) || manifest[prop.key] == undefined) ||
                (typeof(manifest[prop.key]) != prop.type)
            ) {
                debug(`Manifest is missing '${prop.key}'`);
                return false;
            }
        }

        return true;
    }
}

// Export as a singletone
let instance = new PluginManager();
module.exports = instance;