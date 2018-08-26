#!/usr/bin/env node
let Daphne = require('./src/daphne.js');


const directory = process.cwd();
const argv = process.argv.slice(2);
const argc = argv.length;


let staticWebsite = new Daphne(directory);

// Enable debugging or not
switch (argv[argc - 1])
{
    case "--debug":
        let debug = require('./src/debugger');
        debug.enableAll();
        break;
}


switch (argv[0])
{
    case "build":
        staticWebsite.buildSite();
        break;

    case "watch":
        staticWebsite.watchSite();
        break;

    case "serve":
        staticWebsite.serveSite();
        break;


    // Help, unkown argument, or no argument
    case "help":
    default:
        if (argv[0] != "help" && argv[0] != undefined)
            console.log(`Unkown argument '${argv[0]}'`);

        console.log("help");
        break;
}

