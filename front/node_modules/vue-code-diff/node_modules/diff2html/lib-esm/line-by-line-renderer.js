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
import * as Rematch from './rematch';
import * as renderUtils from './render-utils';
import { LineType, } from './types';
export var defaultLineByLineRendererConfig = __assign(__assign({}, renderUtils.defaultRenderConfig), { renderNothingWhenEmpty: false, matchingMaxComparisons: 2500, maxLineSizeInBlockForComparison: 200 });
var genericTemplatesPath = 'generic';
var baseTemplatesPath = 'line-by-line';
var iconsBaseTemplatesPath = 'icon';
var tagsBaseTemplatesPath = 'tag';
var LineByLineRenderer = (function () {
    function LineByLineRenderer(hoganUtils, config) {
        if (config === void 0) { config = {}; }
        this.hoganUtils = hoganUtils;
        this.config = __assign(__assign({}, defaultLineByLineRendererConfig), config);
    }
    LineByLineRenderer.prototype.render = function (diffFiles) {
        var _this = this;
        var diffsHtml = diffFiles
            .map(function (file) {
            var diffs;
            if (file.blocks.length) {
                diffs = _this.generateFileHtml(file);
            }
            else {
                diffs = _this.generateEmptyDiff();
            }
            return _this.makeFileDiffHtml(file, diffs);
        })
            .join('\n');
        return this.hoganUtils.render(genericTemplatesPath, 'wrapper', { content: diffsHtml });
    };
    LineByLineRenderer.prototype.makeFileDiffHtml = function (file, diffs) {
        if (this.config.renderNothingWhenEmpty && Array.isArray(file.blocks) && file.blocks.length === 0)
            return '';
        var fileDiffTemplate = this.hoganUtils.template(baseTemplatesPath, 'file-diff');
        var filePathTemplate = this.hoganUtils.template(genericTemplatesPath, 'file-path');
        var fileIconTemplate = this.hoganUtils.template(iconsBaseTemplatesPath, 'file');
        var fileTagTemplate = this.hoganUtils.template(tagsBaseTemplatesPath, renderUtils.getFileIcon(file));
        return fileDiffTemplate.render({
            file: file,
            fileHtmlId: renderUtils.getHtmlId(file),
            diffs: diffs,
            filePath: filePathTemplate.render({
                fileDiffName: renderUtils.filenameDiff(file),
            }, {
                fileIcon: fileIconTemplate,
                fileTag: fileTagTemplate,
            }),
        });
    };
    LineByLineRenderer.prototype.generateEmptyDiff = function () {
        return this.hoganUtils.render(genericTemplatesPath, 'empty-diff', {
            contentClass: 'd2h-code-line',
            CSSLineClass: renderUtils.CSSLineClass,
        });
    };
    LineByLineRenderer.prototype.generateFileHtml = function (file) {
        var _this = this;
        var matcher = Rematch.newMatcherFn(Rematch.newDistanceFn(function (e) { return renderUtils.deconstructLine(e.content, file.isCombined).content; }));
        return file.blocks
            .map(function (block) {
            var lines = _this.hoganUtils.render(genericTemplatesPath, 'block-header', {
                CSSLineClass: renderUtils.CSSLineClass,
                blockHeader: file.isTooBig ? block.header : renderUtils.escapeForHtml(block.header),
                lineClass: 'd2h-code-linenumber',
                contentClass: 'd2h-code-line',
            });
            _this.applyLineGroupping(block).forEach(function (_a) {
                var contextLines = _a[0], oldLines = _a[1], newLines = _a[2];
                if (oldLines.length && newLines.length && !contextLines.length) {
                    _this.applyRematchMatching(oldLines, newLines, matcher).map(function (_a) {
                        var oldLines = _a[0], newLines = _a[1];
                        var _b = _this.processChangedLines(file.isCombined, oldLines, newLines), left = _b.left, right = _b.right;
                        lines += left;
                        lines += right;
                    });
                }
                else if (contextLines.length) {
                    contextLines.forEach(function (line) {
                        var _a = renderUtils.deconstructLine(line.content, file.isCombined), prefix = _a.prefix, content = _a.content;
                        lines += _this.generateSingleLineHtml({
                            type: renderUtils.CSSLineClass.CONTEXT,
                            prefix: prefix,
                            content: content,
                            oldNumber: line.oldNumber,
                            newNumber: line.newNumber,
                        });
                    });
                }
                else if (oldLines.length || newLines.length) {
                    var _b = _this.processChangedLines(file.isCombined, oldLines, newLines), left = _b.left, right = _b.right;
                    lines += left;
                    lines += right;
                }
                else {
                    console.error('Unknown state reached while processing groups of lines', contextLines, oldLines, newLines);
                }
            });
            return lines;
        })
            .join('\n');
    };
    LineByLineRenderer.prototype.applyLineGroupping = function (block) {
        var blockLinesGroups = [];
        var oldLines = [];
        var newLines = [];
        for (var i = 0; i < block.lines.length; i++) {
            var diffLine = block.lines[i];
            if ((diffLine.type !== LineType.INSERT && newLines.length) ||
                (diffLine.type === LineType.CONTEXT && oldLines.length > 0)) {
                blockLinesGroups.push([[], oldLines, newLines]);
                oldLines = [];
                newLines = [];
            }
            if (diffLine.type === LineType.CONTEXT) {
                blockLinesGroups.push([[diffLine], [], []]);
            }
            else if (diffLine.type === LineType.INSERT && oldLines.length === 0) {
                blockLinesGroups.push([[], [], [diffLine]]);
            }
            else if (diffLine.type === LineType.INSERT && oldLines.length > 0) {
                newLines.push(diffLine);
            }
            else if (diffLine.type === LineType.DELETE) {
                oldLines.push(diffLine);
            }
        }
        if (oldLines.length || newLines.length) {
            blockLinesGroups.push([[], oldLines, newLines]);
            oldLines = [];
            newLines = [];
        }
        return blockLinesGroups;
    };
    LineByLineRenderer.prototype.applyRematchMatching = function (oldLines, newLines, matcher) {
        var comparisons = oldLines.length * newLines.length;
        var maxLineSizeInBlock = Math.max.apply(null, [0].concat(oldLines.concat(newLines).map(function (elem) { return elem.content.length; })));
        var doMatching = comparisons < this.config.matchingMaxComparisons &&
            maxLineSizeInBlock < this.config.maxLineSizeInBlockForComparison &&
            (this.config.matching === 'lines' || this.config.matching === 'words');
        return doMatching ? matcher(oldLines, newLines) : [[oldLines, newLines]];
    };
    LineByLineRenderer.prototype.processChangedLines = function (isCombined, oldLines, newLines) {
        var fileHtml = {
            right: '',
            left: '',
        };
        var maxLinesNumber = Math.max(oldLines.length, newLines.length);
        for (var i = 0; i < maxLinesNumber; i++) {
            var oldLine = oldLines[i];
            var newLine = newLines[i];
            var diff = oldLine !== undefined && newLine !== undefined
                ? renderUtils.diffHighlight(oldLine.content, newLine.content, isCombined, this.config)
                : undefined;
            var preparedOldLine = oldLine !== undefined && oldLine.oldNumber !== undefined
                ? __assign(__assign({}, (diff !== undefined
                    ? {
                        prefix: diff.oldLine.prefix,
                        content: diff.oldLine.content,
                        type: renderUtils.CSSLineClass.DELETE_CHANGES,
                    }
                    : __assign(__assign({}, renderUtils.deconstructLine(oldLine.content, isCombined)), { type: renderUtils.toCSSClass(oldLine.type) }))), { oldNumber: oldLine.oldNumber, newNumber: oldLine.newNumber }) : undefined;
            var preparedNewLine = newLine !== undefined && newLine.newNumber !== undefined
                ? __assign(__assign({}, (diff !== undefined
                    ? {
                        prefix: diff.newLine.prefix,
                        content: diff.newLine.content,
                        type: renderUtils.CSSLineClass.INSERT_CHANGES,
                    }
                    : __assign(__assign({}, renderUtils.deconstructLine(newLine.content, isCombined)), { type: renderUtils.toCSSClass(newLine.type) }))), { oldNumber: newLine.oldNumber, newNumber: newLine.newNumber }) : undefined;
            var _a = this.generateLineHtml(preparedOldLine, preparedNewLine), left = _a.left, right = _a.right;
            fileHtml.left += left;
            fileHtml.right += right;
        }
        return fileHtml;
    };
    LineByLineRenderer.prototype.generateLineHtml = function (oldLine, newLine) {
        return {
            left: this.generateSingleLineHtml(oldLine),
            right: this.generateSingleLineHtml(newLine),
        };
    };
    LineByLineRenderer.prototype.generateSingleLineHtml = function (line) {
        if (line === undefined)
            return '';
        var lineNumberHtml = this.hoganUtils.render(baseTemplatesPath, 'numbers', {
            oldNumber: line.oldNumber || '',
            newNumber: line.newNumber || '',
        });
        return this.hoganUtils.render(genericTemplatesPath, 'line', {
            type: line.type,
            lineClass: 'd2h-code-linenumber',
            contentClass: 'd2h-code-line',
            prefix: line.prefix === ' ' ? '&nbsp;' : line.prefix,
            content: line.content,
            lineNumber: lineNumberHtml,
        });
    };
    return LineByLineRenderer;
}());
export default LineByLineRenderer;
//# sourceMappingURL=line-by-line-renderer.js.map