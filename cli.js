#!/usr/bin/env node
var Daphne = require('./src/daphne.js');

const directory = process.cwd();
const argv = process.argv.slice(2);
const argc = argv.length;


let siteGenerator = new Daphne();

siteGenerator.buildSite(directory);

//console.log(siteGenerator._ProjectPath);