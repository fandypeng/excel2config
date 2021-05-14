"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.hashCode = exports.unifyPath = exports.escapeForRegExp = void 0;
var specials = [
    '-',
    '[',
    ']',
    '/',
    '{',
    '}',
    '(',
    ')',
    '*',
    '+',
    '?',
    '.',
    '\\',
    '^',
    '$',
    '|',
];
var regex = RegExp('[' + specials.join('\\') + ']', 'g');
function escapeForRegExp(str) {
    return str.replace(regex, '\\$&');
}
exports.escapeForRegExp = escapeForRegExp;
function unifyPath(path) {
    return path ? path.replace(/\\/g, '/') : path;
}
exports.unifyPath = unifyPath;
function hashCode(text) {
    var i, chr, len;
    var hash = 0;
    for (i = 0, len = text.length; i < len; i++) {
        chr = text.charCodeAt(i);
        hash = (hash << 5) - hash + chr;
        hash |= 0;
    }
    return hash;
}
exports.hashCode = hashCode;
//# sourceMappingURL=utils.js.map