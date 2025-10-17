const path = require('path');

const platform = process.platform;
const ext = platform === 'win32' ? '.exe' : '';
const binaryName = `tunny${ext}`;
const binaryPath = path.join(__dirname, 'bin', binaryName);

module.exports = {
  binaryPath,
  version: require('./package.json').version
};

