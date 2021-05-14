export var LineType;
(function (LineType) {
    LineType["INSERT"] = "insert";
    LineType["DELETE"] = "delete";
    LineType["CONTEXT"] = "context";
})(LineType || (LineType = {}));
export var OutputFormatType = {
    LINE_BY_LINE: 'line-by-line',
    SIDE_BY_SIDE: 'side-by-side',
};
export var LineMatchingType = {
    LINES: 'lines',
    WORDS: 'words',
    NONE: 'none',
};
export var DiffStyleType = {
    WORD: 'word',
    CHAR: 'char',
};
//# sourceMappingURL=types.js.map