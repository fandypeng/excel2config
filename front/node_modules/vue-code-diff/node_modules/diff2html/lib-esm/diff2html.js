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
import * as DiffParser from './diff-parser';
import * as fileListPrinter from './file-list-renderer';
import LineByLineRenderer, { defaultLineByLineRendererConfig } from './line-by-line-renderer';
import SideBySideRenderer, { defaultSideBySideRendererConfig } from './side-by-side-renderer';
import { OutputFormatType } from './types';
import HoganJsUtils from './hoganjs-utils';
export var defaultDiff2HtmlConfig = __assign(__assign(__assign({}, defaultLineByLineRendererConfig), defaultSideBySideRendererConfig), { outputFormat: OutputFormatType.LINE_BY_LINE, drawFileList: true });
export function parse(diffInput, configuration) {
    if (configuration === void 0) { configuration = {}; }
    return DiffParser.parse(diffInput, __assign(__assign({}, defaultDiff2HtmlConfig), configuration));
}
export function html(diffInput, configuration) {
    if (configuration === void 0) { configuration = {}; }
    var config = __assign(__assign({}, defaultDiff2HtmlConfig), configuration);
    var diffJson = typeof diffInput === 'string' ? DiffParser.parse(diffInput, config) : diffInput;
    var hoganUtils = new HoganJsUtils(config);
    var fileList = config.drawFileList ? fileListPrinter.render(diffJson, hoganUtils) : '';
    var diffOutput = config.outputFormat === 'side-by-side'
        ? new SideBySideRenderer(hoganUtils, config).render(diffJson)
        : new LineByLineRenderer(hoganUtils, config).render(diffJson);
    return fileList + diffOutput;
}
//# sourceMappingURL=diff2html.js.map