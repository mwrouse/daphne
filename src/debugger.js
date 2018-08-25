let debug = require('debug');

module.exports = (function(exp){
    let debuggers = {};
    let isEnabled = false;

    function factory(name) {
        if (debuggers[name] == undefined)
            debuggers[name] = debug(`daphne:${name}`);

        if (isEnabled)
            debuggers[name].enabled = true;

        debuggers[name].new = (subnamespace) => {
            return factory(`${name}:${subnamespace}`);
        };
        return debuggers[name];
    }

    // Enable loggers (at any time)
    factory.enableAll = () => {
        isEnabled = true;
        for (let key in debuggers) {
            debuggers[key].enabled = true;
        }
    };

    return factory;
})();
