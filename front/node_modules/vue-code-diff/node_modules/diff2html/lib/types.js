"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DiffStyleType = exports.LineMatchingType = exports.OutputFormatType = exports.LineType = void 0;
var LineType;
(function (LineType) {
    LineType["INSERT"] = "insert";
    LineType["DELETE"] = "delete";
    LineType["CONTEXT"] = "context";
})(LineType = exports.LineType || (exports.LineType = {}));
exports.OutputFormatType = {
    LINE_BY_LINE: 'line-by-line',
    SIDE_BY_SIDE: 'side-by-side',
};
exports.LineMatchingType = {
    LINES: 'lines',
    WORDS: 'words',
    NONE: 'none',
};
exports.DiffStyleType = {
    WORD: 'word',
    CHAR: 'char',
};
//# sourceMappingURL=types.js.map