"use strict";
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
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.defaultDiff2HtmlUIConfig = exports.Diff2HtmlUI = void 0;
var highlight_js_1 = __importDefault(require("highlight.js"));
var diff2html_ui_base_1 = require("./diff2html-ui-base");
Object.defineProperty(exports, "defaultDiff2HtmlUIConfig", { enumerable: true, get: function () { return diff2html_ui_base_1.defaultDiff2HtmlUIConfig; } });
var Diff2HtmlUI = (function (_super) {
    __extends(Diff2HtmlUI, _super);
    function Diff2HtmlUI(target, diffInput, config) {
        if (config === void 0) { config = {}; }
        return _super.call(this, target, diffInput, config, highlight_js_1.default) || this;
    }
    return Diff2HtmlUI;
}(diff2html_ui_base_1.Diff2HtmlUI));
exports.Diff2HtmlUI = Diff2HtmlUI;
//# sourceMappingURL=diff2html-ui.js.map