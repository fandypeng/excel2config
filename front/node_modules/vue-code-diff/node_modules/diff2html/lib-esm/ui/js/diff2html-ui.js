var __extends = (this && this.__extends) || (function () {
    var extendStatics = function (d, b) {
        extendStatics = Object.setPrototypeOf ||
            ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
            function (d, b) { for (var p in b) if (Object.prototype.hasOwnProperty.call(b, p)) d[p] = b[p]; };
        return extendStatics(d, b);
    };
    return function (d, b) {
        if (typeof b !== "function" && b !== null)
            throw new TypeError("Class extends value " + String(b) + " is not a constructor or null");
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();
import hljs from 'highlight.js';
import { Diff2HtmlUI as Diff2HtmlUIBase, defaultDiff2HtmlUIConfig } from './diff2html-ui-base';
var Diff2HtmlUI = (function (_super) {
    __extends(Diff2HtmlUI, _super);
    function Diff2HtmlUI(target, diffInput, config) {
        if (config === void 0) { config = {}; }
        return _super.call(this, target, diffInput, config, hljs) || this;
    }
    return Diff2HtmlUI;
}(Diff2HtmlUIBase));
export { Diff2HtmlUI };
export { defaultDiff2HtmlUIConfig };
//# sourceMappingURL=diff2html-ui.js.map