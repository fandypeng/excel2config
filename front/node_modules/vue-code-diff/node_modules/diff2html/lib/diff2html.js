"use strict";
var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    Object.defineProperty(o, k2, { enumerable: true, get: function() { return m[k]; } });
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.html = exports.parse = exports.defaultDiff2HtmlConfig = void 0;
var DiffParser = __importStar(require("./diff-parser"));
var fileListPrinter = __importStar(require("./file-list-renderer"));
var line_by_line_renderer_1 = __importStar(require("./line-by-line-renderer"));
var side_by_side_renderer_1 = __importStar(require("./side-by-side-renderer"));
var types_1 = require("./types");
var hoganjs_utils_1 = __importDefault(require("./hoganjs-utils"));
exports.defaultDiff2HtmlConfig = __assign(__assign(__assign({}, line_by_line_renderer_1.defaultLineByLineRendererConfig), side_by_side_renderer_1.defaultSideBySideRendererConfig), { outputFormat: types_1.OutputFormatType.LINE_BY_LINE, drawFileList: true });
function parse(diffInput, configuration) {
    if (configuration === void 0) { configuration = {}; }
    return DiffParser.parse(diffInput, __assign(__assign({}, exports.defaultDiff2HtmlConfig), configuration));
}
exports.parse = parse;
function html(diffInput, configuration) {
    if (configuration === void 0) { configuration = {}; }
    var config = __assign(__assign({}, exports.defaultDiff2HtmlConfig), configuration);
    var diffJson = typeof diffInput === 'string' ? DiffParser.parse(diffInput, config) : diffInput;
    var hoganUtils = new hoganjs_utils_1.default(config);
    var fileList = config.drawFileList ? fileListPrinter.render(diffJson, hoganUtils) : '';
    var diffOutput = config.outputFormat === 'side-by-side'
        ? new side_by_side_renderer_1.default(hoganUtils, config).render(diffJson)
        : new line_by_line_renderer_1.default(hoganUtils, config).render(diffJson);
    return fileList + diffOutput;
}
exports.html = html;
//# sourceMappingURL=diff2html.js.map